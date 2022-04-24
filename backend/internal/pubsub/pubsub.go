package pubsub

import "sync"

type Pubsub struct {
	mu     sync.RWMutex
	subs   map[string][]chan string
	closed bool
}

func NewPubsub() *Pubsub {
	return &Pubsub{
		mu:   sync.RWMutex{},
		subs: make(map[string][]chan string),
	}
}

func (ps *Pubsub) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 10)
	ps.subs[topic] = append(ps.subs[topic], ch)

	return ch
}

func (ps *Pubsub) Publish(topic string, msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		go func(ch chan string) {
			ch <- msg
		}(ch)
	}
}

func (ps *Pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}
