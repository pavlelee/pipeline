# Pipeline
A task must be handled by multiple steps, like a pipeline, the first processing will go into the next processing stage. You can control the buffer of the task and the number of workers at each stage.

### Useage
processing a task flow

```
ch := make(chan interface{}, 10)
//ouput must big than entry, becouse pipeline will send EndJob when task over.
op := make(chan interface{}, 11)

//defind a task flow
pip := pipeline.New().Buffer(10).Listen(ch).Output(op).
  Process(7, func(num interface{}) (interface{}, error) {
    val, ok := num.(int)
    if !ok {
      return nil, fmt.Errorf("xxxx")
    }

    fmt.Printf("recive job^^^^^^^^^1:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
    //time.Sleep(2 * time.Second)

    return val, nil
  }).
  Process(5, func(num interface{}) (interface{}, error) {
    val, ok := num.(int)
    if !ok {
      return nil, fmt.Errorf("xxxx")
    }

    fmt.Printf("recive job========2:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
    //time.Sleep(2 * time.Second)

    return val, nil
  }).
  Process(3, func(num interface{}) (interface{}, error) {
    val, ok := num.(int)
    if !ok {
      return nil, fmt.Errorf("xxxx")
    }

    fmt.Printf("recive job*********3:%d，time: %s\n", val, time.Now().Format("2006-01-02T15:04:05.999999-07:00"))
    //time.Sleep(2 * time.Second)

    return val, nil
  }).
  Run()
  
//Add your jobs to channel
for i := 0; i < 100; i++ {
  ch <- i
}
pip.End()

//wait all done
pip.Wait()

//get results form op(channel)
go func() {
  for num := range op {
    //如果收到结束任务，那就可以break了
    _, ok := num.(pipeline.EndJob)
    if ok {
      break
    }
    fmt.Println(num)
  }
}()
```
if your pipeline use once, you can call ```pip.End()```
if you want wait pipeline done, you can call ```pip.Wait```

### API
1. ```func (*Pipeline) Listen(ch chan interface{}) *Pipeline```
Listen an channel

2. ```func (*Pipeline) Buffer(val int) *Pipeline```
How much to send job at the same time

3. ```func (*Pipeline) End() *Pipeline```
Job send over

4. ```func (*Pipeline) Wait() *Pipeline```
Wait all job is done

5. ```func (*Pipeline) Process(worker int, handle func(interface{}) (interface{}, error)) *Pipeline```
Set a processing

6. ```func (*Pipeline) Run() *Pipeline```
Run this pipeline

6. ```func (*Pipeline) Output(ch chan interface{}) *Pipeline```
Accept the final result
