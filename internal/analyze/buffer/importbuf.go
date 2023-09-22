package buffer

import "sync"

type ImportBuffer struct {
	mutex  sync.Mutex
	lenght int
	size   int

	Imports map[string]string
}

func (buf *ImportBuffer) HandleEvent(event Event, errChan chan<- error) {
	event.Execute(buf, errChan)
}

func (buf *ImportBuffer) Size() int {
	return buf.size
}

func (buf *ImportBuffer) Len() int {
	return buf.lenght
}

func newImportBuffer() *ImportBuffer {
	return &ImportBuffer{
		mutex:   sync.Mutex{},
		Imports: make(map[string]string),
	}
}
