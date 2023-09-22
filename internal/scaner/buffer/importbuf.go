package buffer

import "sync"

type ImportBuffer struct {
	mutex  sync.Mutex
	lenght int
	size   int

	Imports map[string]string
}

func (buf *ImportBuffer) HandleEvent(event bufferEvent, errChan chan<- error) {
	event.Execute(buf, errChan)
}

func (buf *ImportBuffer) Size() int {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	return buf.size
}

func (buf *ImportBuffer) Len() int {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	return buf.lenght
}

func (buf *ImportBuffer) IsPresent(key string) (string, bool) {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	path, exists := buf.Imports[key]
	return path, exists
}

func newImportBuffer() *ImportBuffer {
	return &ImportBuffer{
		mutex:   sync.Mutex{},
		Imports: make(map[string]string),
	}
}