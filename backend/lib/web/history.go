package web

import (
	"io"
	"sync"
)

// simple rolling history buffer
type History struct {
	sync.RWMutex
	buffer [][]byte
	max    int
}

func NewHistory() *History {
	return &History{
		buffer: make([][]byte, 0, 1000),
		max:    1000,
	}
}

func (h *History) Write(data []byte) (int, error) {
	h.Lock()
	defer h.Unlock()

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	h.buffer = append(h.buffer, dataCopy)

	if len(h.buffer) > h.max {
		h.buffer = h.buffer[1:]
	}

	return len(data), nil
}

func (h *History) Close() error {
	return nil
}

func (h *History) Replay(w io.Writer) error {
	h.RLock()
	defer h.RUnlock()

	for _, chunk := range h.buffer {
		_, err := w.Write(chunk)
		if err != nil {
			return err
		}
	}
	return nil
}
