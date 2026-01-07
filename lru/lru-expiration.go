package main

import (
	"container/list"
	"time"
)

type ExpirableLRUCache[K comparable, V any] interface {
	Get(K) (V, bool)
	Add(K, V)
	Len() int
}

type CacheEntry[K comparable, V any] struct {
	key        K
	value      V
	listRef    *list.Element
	expiration time.Time
}

type ExpirableLruCache[K comparable, V any] struct {
	cache    map[K]*CacheEntry[K, V]
	capacity int
	order    *list.List
	ttl      time.Duration
}

func NewExpirableLRUCache[K comparable, V any](capacity int, ttl time.Duration) ExpirableLRUCache[K, V] {
	if capacity <= 0 {
		panic("capacity must be > 0")
	}
	if ttl <= 0 {
		panic("ttl must be > 0")
	}
	return &ExpirableLruCache[K, V]{
		cache:    make(map[K]*CacheEntry[K, V], capacity),
		order:    list.New(),
		capacity: capacity,
		ttl:      ttl,
	}
}

func (l *ExpirableLruCache[K, V]) Len() int {
	return len(l.cache)
}

func (l *ExpirableLruCache[K, V]) purgeExpired() {
	now := time.Now()
	for _, entry := range l.cache {
		if entry.expiration.Before(now) {
			delete(l.cache, entry.key)
			l.order.Remove(entry.listRef)
		}
	}
}

func (l *ExpirableLruCache[K, V]) Get(k K) (V, bool) {
	if element, found := l.cache[k]; found {
		if element.expiration.Before(time.Now()) {
			//it's expired, remmoved from the list
			delete(l.cache, k)
			l.order.Remove(element.listRef)
			var v V
			return v, false
		} else {
			l.order.MoveToFront(element.listRef)
			return element.value, true
		}
	} else {
		var v V
		return v, false
	}
}

func (l *ExpirableLruCache[K, V]) Add(key K, value V) {
	expiration := time.Now().Add(l.ttl)
	if element, found := l.cache[key]; found {
		element.expiration = expiration
		l.order.MoveToFront(element.listRef)
		element.value = value
	} else {
		entry := &CacheEntry[K, V]{
			key:        key,
			value:      value,
			expiration: expiration,
		}

		elem := l.order.PushFront(entry)
		entry.listRef = elem
		l.cache[key] = entry

		if l.Len() > l.capacity {
			lastElement := l.order.Back()

			if lastElement != nil {
				evictEntry := lastElement.Value.(*CacheEntry[K, V])
				delete(l.cache, evictEntry.key)
				l.order.Remove(lastElement)
			}
		}
	}
	l.purgeExpired()
}

//we should add mutex as well
