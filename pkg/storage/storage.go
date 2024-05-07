package storage

import (
	"cmp"
	"container/heap"
	"my_storage/pkg/avl_tree"
	"sync"
	"time"
)

type heapTTL[Key cmp.Ordered, Value any] []*elementStorage[Key, Value]

func (h *heapTTL[Key, Value]) Push(x any) {
	*h = append(*h, x.(*elementStorage[Key, Value]))
}

func (h *heapTTL[Key, Value]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h heapTTL[Key, Value]) Len() int {
	return len(h)
}

func (h heapTTL[Key, Value]) Less(i, j int) bool {
	return h[i].expireAt.Before(h[j].expireAt)
}

func (h heapTTL[Key, Value]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index, h[j].index = i, j
}

type elementStorage[Key cmp.Ordered, Value any] struct {
	key      Key
	value    Value
	expireAt time.Time
	index    int
}

type Storage[Key cmp.Ordered, Value any] struct {
	root      *avl_tree.AVLTree[Key, *elementStorage[Key, Value]]
	heap      heapTTL[Key, Value]
	m         sync.RWMutex
	updateC   chan struct{}
	finishedC chan struct{}
	once      sync.Once
}

func NewStorage[Key cmp.Ordered, Value any]() *Storage[Key, Value] {
	return &Storage[Key, Value]{
		root:      &avl_tree.AVLTree[Key, *elementStorage[Key, Value]]{},
		heap:      make([]*elementStorage[Key, Value], 0, 1024),
		updateC:   make(chan struct{}, 1),
		finishedC: make(chan struct{}, 1),
	}
}

func (s *Storage[Key, Value]) Start() {
	s.once.Do(func() {
		go func() {
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-ticker.C:
					s.m.Lock()
					for len(s.heap) > 0 && s.heap[0].expireAt.Before(time.Now()) {
						s.root.Remove(s.heap[0].key)
						heap.Pop(&s.heap)
					}

					if len(s.heap) > 0 {
						ticker.Reset(time.Until(s.heap[0].expireAt))
					}
					s.m.Unlock()
				case <-s.updateC:
					if len(s.heap) > 0 {
						ticker.Reset(time.Until(s.heap[0].expireAt))
					}
				case <-s.finishedC:
					break
				}
			}
		}()
	})
}

func (s *Storage[Key, Value]) Stop() {
	s.finishedC <- struct{}{}
}

func (s *Storage[Key, Value]) Set(key Key, value Value, expireAt time.Time) {
	s.m.Lock()
	defer s.m.Unlock()

	temp := &elementStorage[Key, Value]{key: key, value: value, expireAt: expireAt, index: len(s.heap)}

	s.remove(key) // to update expiredAt if something with this key exists

	s.root.Add(key, temp)
	heap.Push(&s.heap, temp)

	s.updateC <- struct{}{}
}

func (s *Storage[Key, Value]) Get(key Key) (Value, bool) {
	s.m.RLock()
	defer s.m.RUnlock()

	result, ok := s.root.Search(key)
	return result.value, ok
}

func (s *Storage[Key, Value]) Delete(key Key) {
	s.m.Lock()
	defer s.m.Unlock()

	s.remove(key)
}

func (s *Storage[Key, Value]) remove(key Key) {
	if find, ok := s.root.Search(key); ok {
		if len(s.heap) <= find.index {
			return
		}

		heap.Remove(&s.heap, find.index)
		s.root.Remove(find.key)
	}
}

func (s *Storage[Key, Value]) GetRoot() *avl_tree.AVLTree[Key, *elementStorage[Key, Value]] {
	return s.root
}
