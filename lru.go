package main

import (
	"container/list"
	"fmt"
)

type LRUCache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
	Len() int
}

type Entry[K comparable, V any] struct {
	Key     K
	Value   V
	Element *list.Element
}

type lruCache[K comparable, V any] struct {
	cache    map[K]*Entry[K, V]
	capacity int
	order    *list.List //double linked list to track LRU order
}

func NewLRUCache[K comparable, V any](capacity int) LRUCache[K, V] {
	return &lruCache[K, V]{
		cache:    make(map[K]*Entry[K, V], capacity),
		capacity: capacity,
		order:    list.New(),
	}
}

func (l *lruCache[K, V]) Get(key K) (V, bool) {
	if entry, found := l.cache[key]; found {
		l.order.MoveToFront(entry.Element)
		return entry.Value, true
	} else {
		var zero V
		return zero, false
	}
}

func (l *lruCache[K, V]) Put(key K, value V) {

	if entry, found := l.cache[key]; found {
		entry.Value = value
		l.order.MoveToFront(entry.Element)
	} else {
		l.cache[key] = &Entry[K, V]{
			Key:   key,
			Value: value,
		}
		elem := l.order.PushFront(l.cache[key])
		l.cache[key].Element = elem
		if len(l.cache) > l.capacity {
			backElem := l.order.Back()
			if backElem != nil {
				evictEntry := backElem.Value.(*Entry[K, V])
				delete(l.cache, evictEntry.Key)
				l.order.Remove(backElem)
			}
		}
	}
}

func (l *lruCache[K, V]) Len() int {
	return len(l.cache)
}

func main() {
	lruCache := NewLRUCache[string, int](2)
	lruCache.Put("one", 1)
	lruCache.Put("two", 2)
	fmt.Println(lruCache.Get("one"))
}
