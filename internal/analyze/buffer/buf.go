package buffer

import (
	"fmt"
	"sync"
)

const (
	toImportsBuffer = iota << 1
	toStructsBuffer
)

type Buffer interface {
	Len() int
	Size() int
	HandleEvent(event Event, errChan chan<- error)
}

type Event interface {
	ToBuffer() int
	Execute(buffer Buffer, errChan chan<- error)
}

type ManagerBuffer struct {
	StructBuffer *StructBuffer
	ImportBuffer *ImportBuffer

	WaitGroup sync.WaitGroup

	eventChan chan Event
	stopChan  chan struct{}
	errChan   chan error
}

func (buf *ManagerBuffer) SendEvent(event ...Event) {
	for i := 0; i < len(event); i++ {
		buf.eventChan <- event[i]
	}
}

func (buf *ManagerBuffer) handleEvent(event Event) {
	switch event.ToBuffer() {
	case toImportsBuffer:
		buf.ImportBuffer.HandleEvent(event, buf.errChan)
	case toStructsBuffer:
		buf.StructBuffer.HandleEvent(event, buf.errChan)
	default:
		buf.errChan <- fmt.Errorf("unknown buffer type: %d", event.ToBuffer())
	}
}

func (buf *ManagerBuffer) Start() {
	for {
		select {
		case event, ok := <-buf.eventChan:
			if !ok {
				return
			}
			buf.handleEvent(event)
		case <-buf.stopChan:
			return
		}
	}
}

func (mb *ManagerBuffer) Stop() {
	close(mb.stopChan)
}

func NewManagerBuffer(errChan chan error) *ManagerBuffer {
	return &ManagerBuffer{
		StructBuffer: newStructBuffer(),
		ImportBuffer: newImportBuffer(),
		WaitGroup:    sync.WaitGroup{},
		eventChan:    make(chan Event),
		stopChan:     make(chan struct{}),
		errChan:      errChan,
	}
}
