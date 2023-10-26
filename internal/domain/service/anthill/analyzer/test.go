package analyzer

import (
	"fmt"
	"sync"
	"time"
)

type compassEvent interface {
	Name() string
}
type LogEvent struct {
	Message string
}

func (l *LogEvent) Name() string {
	return "LogEvent"
}

type DataEvent struct {
	Data string
}

func (d *DataEvent) Name() string {
	return "DataEvent"
}

type Root struct {
	eventCh       chan compassEvent
	unsubscribeCh chan struct{}
}

func NewRoot() *Root {
	return &Root{
		eventCh:       make(chan compassEvent, 1),
		unsubscribeCh: make(chan struct{}),
	}
}
func (r *Root) Subscribe() (<-chan compassEvent, chan struct{}) {
	return r.eventCh, r.unsubscribeCh
}
func (r *Root) EmitLog(logMessage string) {
	r.eventCh <- &LogEvent{Message: logMessage}
}
func (r *Root) EmitData(data string) {
	r.eventCh <- &DataEvent{Data: data}
}
func (r *Root) Parse() {
	v := &Visitor2{eventCh: r.eventCh}
	v.DoSomething()
}

// /// VISITOR /////
type Visitor2 struct {
	eventCh chan compassEvent
}

func (v *Visitor2) DoSomething() {
	for i := 0; i < 10; i++ {
		go func() {
			v.eventCh <- &LogEvent{Message: "Visitor did something."}
		}()
	}
}

// /////////////////
func main() {
	root := NewRoot()
	var wg sync.WaitGroup
	eventCh, unsubscribeCh := root.Subscribe()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case e := <-eventCh:
				switch ev := e.(type) {
				case *LogEvent:
					fmt.Printf("Log Subscriber received a log event: Message=%s\n", ev.Message)
				case *DataEvent:
					fmt.Printf("Data Subscriber received a data event: Data=%s\n", ev.Data)
				default:
					fmt.Printf("Unknown event type: %s\n", e.Name())
				}
			case <-unsubscribeCh:
				return
			}
		}
	}()
	root.Parse()
	time.Sleep(time.Second)
	close(unsubscribeCh)
	wg.Wait()
}
