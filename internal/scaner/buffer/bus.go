package buffer

import (
	"fmt"
	"sync"

	"github.com/g10z3r/archx/pkg/bloom"
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

type BufferEventBus struct {
	StructBuffer *StructBuffer
	ImportBuffer *ImportBuffer

	WaitGroup sync.WaitGroup

	eventChan chan bufferEvent
	stopChan  chan struct{}
	errChan   chan error
}

func (buf *BufferEventBus) SendEvent(event bufferEvent) {
	buf.eventChan <- event
}

func (buf *BufferEventBus) handleEvent(event bufferEvent) {
	switch event.ToBuffer() {
	case toImportsBuffer:
		buf.ImportBuffer.HandleEvent(event, buf.errChan)
	case toStructsBuffer:
		buf.StructBuffer.HandleEvent(event, buf.errChan)
	default:
		buf.errChan <- fmt.Errorf("unknown buffer type: %d", event.ToBuffer())
	}
}

func (buf *BufferEventBus) Open() {
	for {
		select {
		case event, ok := <-buf.eventChan:
			if !ok {
				return
			}
			buf.handleEvent(event)
		case <-buf.stopChan:
			close(buf.eventChan)
			return
		}
	}
}

func (mb *BufferEventBus) Close() {
	close(mb.stopChan)
}

func NewBufferEventBus(mod string, impTotal int, errChan chan error) *BufferEventBus {
	return &BufferEventBus{
		StructBuffer: newStructBuffer(),
		ImportBuffer: newImportBuffer(
			mod,
			bloom.FilterConfig{
				ExpectedItemCount:        uint64(impTotal),
				DesiredFalsePositiveRate: 0.01,
			},
		),
		WaitGroup: sync.WaitGroup{},
		eventChan: make(chan bufferEvent),
		stopChan:  make(chan struct{}),
		errChan:   errChan,
	}
}
