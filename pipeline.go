package pipeline

import (
	"log"
	"sync"
)

// Pipeline
type Pipeline struct {
	entry      chan interface{}
	workspaces []Workspace
	buffer     int
	connect    []chan interface{}
	wg         sync.WaitGroup
}

// Workspace
type Workspace struct {
	worker int
	handle func(interface{}) (interface{}, error)
}

// EndJob
type EndJob struct{}

// New new a Pipeline
func New() *Pipeline {
	p := &Pipeline{}
	p.buffer = 0
	return p
}

// Listen listen an channel
func (p *Pipeline) Listen(ch chan interface{}) *Pipeline {
	p.entry = ch
	return p
}

// Buffer how much to send job at the same time
func (p *Pipeline) Buffer(val int) *Pipeline {
	p.buffer = val
	return p
}

// End job send over
func (p *Pipeline) End() *Pipeline {
	p.entry <- EndJob{}
	return p
}

// Wait wait all job is done
func (p *Pipeline) Wait() {
	p.wg.Wait()
}

// Process set a processing
func (p *Pipeline) Process(worker int, handle func(interface{}) (interface{}, error)) *Pipeline {
	p.workspaces = append(p.workspaces, Workspace{worker: worker, handle: handle})
	return p
}

// Run run this pipeline
func (p *Pipeline) Run() *Pipeline {
	if p.entry == nil {
		log.Panic("Missing entry")
	}

	l := len(p.workspaces)
	if l == 0 {
		log.Panic("Workspace at least one")
	}

	p.connect = append(p.connect, p.entry)
	for i := 0; i < l-1; i++ {
		p.connect = append(p.connect, make(chan interface{}, p.buffer))
	}
	p.connect = append(p.connect, nil)

	for i := 0; i < l; i++ {
		workshop := p.workspaces[i]
		entry := p.connect[i]
		next := p.connect[i+1]

		p.work(entry, workshop.worker, workshop.handle, next)
	}

	return p
}

// work
func (p *Pipeline) work(entry chan interface{}, worker int, handle func(interface{}) (interface{}, error), next chan interface{}) {
	p.wg.Add(1)
	go func() {
		//等待所有的任务处理完
		var wg sync.WaitGroup

		workers := make(chan int, worker)
		for num := range entry {
			//如果收到结束任务，不用再等待了
			_, ok := num.(EndJob)
			if ok {
				break
			}

			workers <- 1
			wg.Add(1)
			go func(num interface{}) {
				ret, err := handle(num)
				if err != nil {
					log.Println(err.Error())
					return
				}

				p.writeNext(next, ret)

				wg.Done()
				<-workers
			}(num)
		}

		wg.Wait()
		//这个作业区的任务都处理完了，
		p.writeNext(next, EndJob{})

		p.wg.Done()
	}()
}

// writeNext
func (p *Pipeline) writeNext(next chan interface{}, v interface{}) {
	if next != nil {
		next <- v
	}
}
