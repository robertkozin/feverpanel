package web

import (
	"io"
	"sync"
)

type Broker struct {
	sync.RWMutex
	writers []io.WriteCloser
}

func NewBroker(writers ...io.WriteCloser) *Broker {
	return &Broker{writers: writers}
}

func (b *Broker) Write(data []byte) (int, error) {
	b.RLock()
	defer b.RUnlock()

	for _, w := range b.writers {
		n, err := w.Write(data)
		if err != nil || n != len(data) {
			_ = w.Close()
			defer b.remove(w)
		}
	}

	return len(data), nil
}

func (b *Broker) Add(writer io.WriteCloser) {
	b.Lock()
	b.writers = append(b.writers, writer)
	b.Unlock()
}

func (b *Broker) Remove(writer io.WriteCloser) {
	b.Lock()
	b.remove(writer)
	b.Unlock()
}

// assumes lock is held
func (b *Broker) remove(writer io.WriteCloser) {
	for i, w := range b.writers {
		if w == writer {
			b.writers = append(b.writers[:i], b.writers[i+1:]...)
			return
		}
	}
}
