package buffer

import (
	"fmt"
	"sync"
)

const (
	toImportsBuffer = iota << 1
	toStructsBuffer
)

type bufferBus interface {
	Len() int
	Size() int
	HandleEvent(event bufferEvent, errChan chan<- error)
}

type bufferEvent interface {
	ToBuffer() int
	Execute(buffer bufferBus, errChan chan<- error)
}

type ManagerBuffer struct {
	StructBuffer *StructBuffer
	ImportBuffer *ImportBuffer

	WaitGroup sync.WaitGroup

	eventChan chan bufferEvent
	stopChan  chan struct{}
	errChan   chan error
}

func (buf *ManagerBuffer) SendEvent(event ...bufferEvent) {
	for i := 0; i < len(event); i++ {
		buf.eventChan <- event[i]
	}
}

func (buf *ManagerBuffer) handleEvent(event bufferEvent) {
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
		eventChan:    make(chan bufferEvent),
		stopChan:     make(chan struct{}),
		errChan:      errChan,
	}
}
