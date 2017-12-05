package pipeline

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func NewLogger() *log.Logger {
	return log.New(ioutil.Discard, "", log.LstdFlags)
}

func BenchmarkPipeline_Run(b *testing.B) {
	ch := make(chan interface{}, 100)

	//定义一个流水线
	pip := New().Listen(ch).
		Process(10, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job^^^^^^^^^1:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(5, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job========2:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(3, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job*********3:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Run()

	//给流水线增加作业
	for i := 0; i < 100; i++ {
		ch <- i
	}
	pip.JobSendEnd()
}

func BenchmarkPipeline_Wait(b *testing.B) {
	ch := make(chan interface{}, 100)

	//定义一个流水线
	pip := New().Listen(ch).
		Process(10, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job^^^^^^^^^1:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(5, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job========2:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(3, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job*********3:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Run()

	//给流水线增加作业
	for i := 0; i < 100; i++ {
		ch <- i
	}
	pip.JobSendEnd()

	pip.Wait()
}

func BenchmarkPipeline_Output(b *testing.B) {
	ch := make(chan interface{}, 100)
	op := make(chan interface{}, 101)

	//定义一个流水线
	pip := New().Listen(ch).Output(op).
		Process(10, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job^^^^^^^^^1:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(5, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job========2:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(3, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job*********3:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Run()

	//给流水线增加作业
	for i := 0; i < 100; i++ {
		ch <- i
	}
	pip.JobSendEnd()

	pip.Wait()

	go func() {
		for num := range op {
			//如果收到结束任务，不用再等待了
			_, ok := num.(EndJob)
			if ok {
				break
			}
			// fmt.Println(num)
		}
	}()
}

func BenchmarkPipeline_OutputSmallThanEntry(b *testing.B) {
	defer func() {
		recover()
	}()

	ch := make(chan interface{}, 100)
	op := make(chan interface{}, 100)

	//定义一个流水线
	pip := New().SetLogger(NewLogger()).Listen(ch).Output(op).
		Process(10, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job^^^^^^^^^1:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(5, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job========2:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(3, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job*********3:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Run()

	//给流水线增加作业
	for i := 0; i < 100; i++ {
		ch <- i
	}
	pip.JobSendEnd()

	pip.Wait()
}

func BenchmarkPipeline_NotListen(b *testing.B) {
	ch := make(chan interface{}, 100)
	defer func() {
		recover()
	}()

	//定义一个流水线
	pip := New().SetLogger(NewLogger()).Run()

	//给流水线增加作业
	for i := 0; i < 100; i++ {
		ch <- i
	}
	pip.JobSendEnd()

	pip.Wait()
}

func BenchmarkPipeline_NotProcess(b *testing.B) {
	ch := make(chan interface{}, 100)
	defer func() {
		recover()
	}()

	//定义一个流水线
	pip := New().SetLogger(NewLogger()).Listen(ch).Run()

	//给流水线增加作业
	for i := 0; i < 100; i++ {
		ch <- i
	}
	pip.JobSendEnd()

	pip.Wait()
}

func BenchmarkPipeline_Buffer(b *testing.B) {
	ch := make(chan interface{}, 100)

	//定义一个流水线
	pip := New().Buffer(10).Listen(ch).
		Process(10, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job^^^^^^^^^1:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(5, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job========2:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(3, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job*********3:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Run()

	//给流水线增加作业
	for i := 0; i < 100; i++ {
		ch <- i
	}
	pip.JobSendEnd()

	pip.Wait()
}

func BenchmarkPipeline_Process(b *testing.B) {
	ch := make(chan interface{}, 100)

	//定义一个流水线
	pip := New().Listen(ch).
		Process(10, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job^^^^^^^^^1:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(5, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job========2:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Process(3, func(num interface{}) (interface{}, error) {
			val, ok := num.(int)
			if !ok {
				return nil, fmt.Errorf("xxxx")
			}

			//fmt.Printf("recive job*********3:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
			//time.Sleep(2 * time.Second)

			return val, nil
		}).
		Run()

	//给流水线增加作业
	for i := 0; i < 100; i++ {
		ch <- i
	}
	pip.JobSendEnd()

	pip.Wait()
}

func BenchmarkPipeline_ProcessError(b *testing.B) {
	ch := make(chan interface{}, 100)

	//定义一个流水线
	pip := New().SetLogger(NewLogger()).Listen(ch).
		Process(2, func(num interface{}) (interface{}, error) {
			return nil, errors.New("error testing")
		}).
		Process(5, func(num interface{}) (interface{}, error) {
			return nil, errors.New("error testing")
		}).
		Process(3, func(num interface{}) (interface{}, error) {
			return nil, errors.New("error testing")
		}).
		Run()

	//给流水线增加作业
	for i := 0; i < 100; i++ {
		ch <- i
	}
	pip.JobSendEnd()

	pip.Wait()
}
