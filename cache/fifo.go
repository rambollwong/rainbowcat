package cache

import (
	"container/list"
	"sync"
)

// FIFOCache represents a First-In-First-Out (FIFO) cache with a fixed size.
// It stores key-value pairs and evicts the oldest entry when the maximum number of elements is reached.
type FIFOCache[K, V any] struct {
	mu              sync.RWMutex
	threadSafe      bool
	maxElements     int
	currentElements int
	_list           *list.List
	cache           map[any]*list.Element

	onRemoved func(k K, v V)
}

// cacheEntry represents a single entry in the FIFO cache.
// It contains a key-value pair.
type cacheEntry[K, V any] struct {
	key   K
	value V
}

// NewFIFOCache creates a new FIFOCache with the specified maximum number of elements.
func NewFIFOCache[K, V any](maxElements int, threadSafe bool) *FIFOCache[K, V] {
	return &FIFOCache[K, V]{
		threadSafe:  threadSafe,
		maxElements: maxElements,
		_list:       list.New(),
		cache:       make(map[any]*list.Element),
	}
}

// SetOnRemovedCallBack registers a callback function that will be invoked when any entry is eliminated or removed.
func (c *FIFOCache[K, V]) SetOnRemovedCallBack(callback func(k K, v V)) {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}
	c.onRemoved = callback
}

// putAndOverwriteIfExist puts a new key-value pair into the FIFO cache.
// If the key already exists, it either overwrites the existing value or retains the existing value based on the 'overwrite' parameter.
// It returns a boolean indicating whether the operation was successful.
func (c *FIFOCache[K, V]) putAndOverwriteIfExist(k K, v V, overwrite bool) bool {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	// Check if the key already exists in the cache
	ele, ok := c.cache[k]

	// If the key exists
	if ok {
		if overwrite {
			// Move the existing entry to the front of the list
			c._list.MoveToFront(ele)
			// Update the value of the existing entry
			ele.Value.(*cacheEntry[K, V]).value = v
			return true // Operation successful
		}
		return false // Operation unsuccessful (key exists and overwrite is false)
	}

	// If the key does not exist
	// Create a new cache entry
	newEntry := &cacheEntry[K, V]{k, v}
	// Put the new cache entry at the front of the list
	newEle := c._list.PushFront(newEntry)
	c.cache[k] = newEle
	c.currentElements++

	// Check if we need to eliminate an entry
	if c.currentElements > c.maxElements {
		// Eliminate a cache entry from the back of the list
		eleEliminated := c._list.Back()
		if eleEliminated != nil {
			entryEliminated, _ := eleEliminated.Value.(*cacheEntry[K, V])
			delete(c.cache, entryEliminated.key)
			c._list.Remove(eleEliminated)
			c.currentElements--
			if c.onRemoved != nil {
				c.onRemoved(entryEliminated.key, entryEliminated.value)
			}
		}
	}
	return true // Operation successful
}

// Put puts a new key-value pair into the FIFO cache, overwriting the existing value if the key already exists.
func (c *FIFOCache[K, V]) Put(k K, v V) {
	c.putAndOverwriteIfExist(k, v, true)
}

// PutIfNotExist puts a new key-value pair into the FIFO cache if the key does not already exist.
// It returns a boolean indicating whether the operation was successful (key did not exist in the cache).
func (c *FIFOCache[K, V]) PutIfNotExist(k K, v V) bool {
	return c.putAndOverwriteIfExist(k, v, false)
}

// Get retrieves the value associated with the specified key from the FIFO cache.
// It returns the value and a boolean indicating whether the key was found in the cache.
func (c *FIFOCache[K, V]) Get(k K) (v V, found bool) {
	if c.threadSafe {
		c.mu.RLock()
		defer c.mu.RUnlock()
	}

	// Check if the key exists in the cache
	ele, ok := c.cache[k]
	if !ok {
		return v, false // Key not found
	}

	// Retrieve the value from the cache entry
	return ele.Value.(*cacheEntry[K, V]).value, true // Return the value and indicate key found
}

// Remove removes the entry with the specified key from the FIFO cache.
// It returns a boolean indicating whether the entry was successfully removed.
func (c *FIFOCache[K, V]) Remove(k K) bool {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	// Check if the key exists in the cache
	ele, ok := c.cache[k]
	if ok {
		// Remove the entry from the linked list
		c._list.Remove(ele)

		// Delete the entry from the cache map
		delete(c.cache, k)

		// Decrease the count of current elements in the cache
		c.currentElements--

		// Trigger the onRemoved callback function, if provided
		if c.onRemoved != nil {
			entry, _ := ele.Value.(*cacheEntry[K, V])
			c.onRemoved(entry.key, entry.value)
		}

		return true // Entry successfully removed
	}

	return false // Entry not found in the cache
}

// Exist checks if the specified key exists in the FIFO cache.
// It returns a boolean indicating whether the key exists in the cache.
func (c *FIFOCache[K, V]) Exist(k K) bool {
	if c.threadSafe {
		c.mu.RLock()
		defer c.mu.RUnlock()
	}

	// Check if the key exists in the cache
	_, ok := c.cache[k]
	return ok
}

// Clear clears all entries from the FIFO cache.
func (c *FIFOCache[K, V]) Clear() {
	if c.threadSafe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	// Reset the number of current elements to zero
	c.currentElements = 0

	// Create a new empty cache map
	c.cache = make(map[interface{}]*list.Element)

	// Create a new empty linked list
	c._list = list.New()
}

// Size returns the current number of elements in the FIFO cache.
func (c *FIFOCache[K, V]) Size() int {
	if c.threadSafe {
		c.mu.RLock()
		defer c.mu.RUnlock()
	}

	// Return the current number of elements in the cache
	return c.currentElements
}
