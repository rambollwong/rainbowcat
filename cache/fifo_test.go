package cache

import (
	"sync"
	"testing"
)

func TestFIFOCache_PutAndGet(t *testing.T) {
	// Create a cache with maximum capacity of 2
	cache := NewFIFOCache[string, int](2, false)

	// Test putting and getting elements
	cache.Put("key1", 1)
	cache.Put("key2", 2)

	// Test getting existing elements
	if value, found := cache.Get("key1"); !found || value != 1 {
		t.Errorf("Expected key1 to be 1, got %d, found: %v", value, found)
	}

	if value, found := cache.Get("key2"); !found || value != 2 {
		t.Errorf("Expected key2 to be 2, got %d, found: %v", value, found)
	}

	// Test getting non-existing element
	if _, found := cache.Get("key3"); found {
		t.Error("Expected key3 to not be found")
	}
}

func TestFIFOCache_PutIfNotExist(t *testing.T) {
	cache := NewFIFOCache[string, int](2, false)

	// First put should succeed
	if success := cache.PutIfNotExist("key1", 1); !success {
		t.Error("Expected PutIfNotExist to succeed for key1")
	}

	// Putting the same key again should fail
	if success := cache.PutIfNotExist("key1", 2); success {
		t.Error("Expected PutIfNotExist to fail for existing key1")
	}

	// Check that the value remains the original value
	if value, found := cache.Get("key1"); !found || value != 1 {
		t.Errorf("Expected key1 to be 1, got %d", value)
	}
}

func TestFIFOCache_Overwrite(t *testing.T) {
	cache := NewFIFOCache[string, int](2, false)

	// Put an element
	cache.Put("key1", 1)

	// Overwrite the element
	cache.Put("key1", 10)

	// Check that the value has been updated
	if value, found := cache.Get("key1"); !found || value != 10 {
		t.Errorf("Expected key1 to be 10, got %d", value)
	}
}

func TestFIFOCache_Remove(t *testing.T) {
	cache := NewFIFOCache[string, int](2, false)

	// Put elements
	cache.Put("key1", 1)
	cache.Put("key2", 2)

	// Remove an existing element
	if removed := cache.Remove("key1"); !removed {
		t.Error("Expected key1 to be removed")
	}

	// Confirm the element has been removed
	if _, found := cache.Get("key1"); found {
		t.Error("Expected key1 to not be found after removal")
	}

	// Try to remove a non-existing element
	if removed := cache.Remove("key3"); removed {
		t.Error("Expected removal of non-existent key to return false")
	}
}

func TestFIFOCache_Exist(t *testing.T) {
	cache := NewFIFOCache[string, int](2, false)

	// Put an element
	cache.Put("key1", 1)

	// Check for existing element
	if !cache.Exist("key1") {
		t.Error("Expected key1 to exist")
	}

	// Check for non-existing element
	if cache.Exist("key2") {
		t.Error("Expected key2 to not exist")
	}
}

func TestFIFOCache_FIFO_Eviction(t *testing.T) {
	cache := NewFIFOCache[int, string](3, false)

	// Put elements until maximum capacity is reached
	cache.Put(1, "value1")
	cache.Put(2, "value2")
	cache.Put(3, "value3")

	// Check all elements exist
	if !cache.Exist(1) || !cache.Exist(2) || !cache.Exist(3) {
		t.Error("All elements should exist before eviction")
	}

	// Put a fourth element, eviction should be triggered
	cache.Put(4, "value4")

	// Check that the first element was removed (FIFO)
	if cache.Exist(1) {
		t.Error("Expected key 1 to be evicted (FIFO)")
	}

	// Check that other elements still exist
	if !cache.Exist(2) || !cache.Exist(3) || !cache.Exist(4) {
		t.Error("Keys 2, 3, and 4 should exist after eviction")
	}
}

func TestFIFOCache_Clear(t *testing.T) {
	cache := NewFIFOCache[string, int](2, false)

	// Put elements
	cache.Put("key1", 1)
	cache.Put("key2", 2)

	// Clear the cache
	cache.Clear()

	// Check that the cache is empty
	if cache.Exist("key1") || cache.Exist("key2") {
		t.Error("Expected cache to be empty after Clear")
	}

	if cache.Size() != 0 {
		t.Errorf("Expected cache size to be 0 after Clear, got %d", cache.Size())
	}
}

func TestFIFOCache_Size(t *testing.T) {
	cache := NewFIFOCache[string, int](3, false)

	// Initial size should be 0
	if size := cache.Size(); size != 0 {
		t.Errorf("Expected initial size to be 0, got %d", size)
	}

	// Put elements
	cache.Put("key1", 1)
	if size := cache.Size(); size != 1 {
		t.Errorf("Expected size to be 1 after adding one element, got %d", size)
	}

	cache.Put("key2", 2)
	if size := cache.Size(); size != 2 {
		t.Errorf("Expected size to be 2 after adding two elements, got %d", size)
	}

	// Remove an element
	cache.Remove("key1")
	if size := cache.Size(); size != 1 {
		t.Errorf("Expected size to be 1 after removing one element, got %d", size)
	}
}

func TestFIFOCache_OnRemovedCallback(t *testing.T) {
	var removedKey string
	var removedValue int

	cache := NewFIFOCache[string, int](2, false)

	// Set callback function
	cache.SetOnRemovedCallBack(func(k string, v int) {
		removedKey = k
		removedValue = v
	})

	// Put elements until maximum capacity is reached
	cache.Put("key1", 1)
	cache.Put("key2", 2)

	// Put a third element, eviction should be triggered
	cache.Put("key3", 3)

	// Check that the callback was called correctly
	if removedKey != "key1" {
		t.Errorf("Expected removed key to be key1, got %s", removedKey)
	}

	if removedValue != 1 {
		t.Errorf("Expected removed value to be 1, got %d", removedValue)
	}

	// Manually remove an element to test callback
	cache.Remove("key2")
	if removedKey != "key2" {
		t.Errorf("Expected removed key to be key2, got %s", removedKey)
	}

	if removedValue != 2 {
		t.Errorf("Expected removed value to be 2, got %d", removedValue)
	}
}

func TestFIFOCache_MoveToFront(t *testing.T) {
	cache := NewFIFOCache[string, int](3, false)

	// Put elements
	cache.Put("key1", 1) // Oldest
	cache.Put("key2", 2) // Middle
	cache.Put("key3", 3) // Newest

	// Overwrite the first element, it should be moved to the front
	cache.Put("key1", 10)

	// Put a new element to trigger eviction
	cache.Put("key4", 4)

	// Since key1 was moved to the front, key2 should be the one evicted
	if cache.Exist("key1") && !cache.Exist("key2") {
		// Verify the value was updated
		if value, _ := cache.Get("key1"); value != 10 {
			t.Errorf("Expected key1 value to be 10, got %d", value)
		}
	} else {
		t.Error("Expected key2 to be evicted, not key1")
	}
}

// Test edge case: cache with zero capacity
func TestFIFOCache_ZeroCapacity(t *testing.T) {
	cache := NewFIFOCache[string, int](0, false)

	// Any element added should be immediately evicted
	cache.Put("key1", 1)

	// Check that the element was immediately evicted
	if cache.Exist("key1") {
		t.Error("Expected key1 to be immediately evicted with zero capacity")
	}

	// Size should be 0
	if size := cache.Size(); size != 0 {
		t.Errorf("Expected size to be 0 with zero capacity, got %d", size)
	}
}

// Test edge case: cache with capacity of 1
func TestFIFOCache_OneCapacity(t *testing.T) {
	cache := NewFIFOCache[string, int](1, false)

	// Put one element
	cache.Put("key1", 1)

	// Check that element exists
	if !cache.Exist("key1") {
		t.Error("Expected key1 to exist")
	}

	// Put a second element, the first should be evicted
	cache.Put("key2", 2)

	// Check that the first element was evicted and the second exists
	if cache.Exist("key1") {
		t.Error("Expected key1 to be evicted")
	}

	if !cache.Exist("key2") {
		t.Error("Expected key2 to exist")
	}
}

func TestFIFOCache_ThreadSafe(t *testing.T) {
	cache := NewFIFOCache[int, int](100, true)
	var wg sync.WaitGroup

	// Test concurrent Put operations
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Put(i, i*2)
		}(i)
	}
	wg.Wait()

	// Verify all items were inserted
	for i := 0; i < 100; i++ {
		if value, found := cache.Get(i); !found || value != i*2 {
			t.Errorf("Expected key %d to have value %d, got %d, found: %v", i, i*2, value, found)
		}
	}

	// Test concurrent Get operations
	wg = sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			value, found := cache.Get(i)
			if !found || value != i*2 {
				t.Errorf("Concurrent get failed for key %d: value=%d, found=%v", i, value, found)
			}
		}(i)
	}
	wg.Wait()

	// Test mixed concurrent operations
	wg = sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(2)
		// Concurrent Put
		go func(i int) {
			defer wg.Done()
			cache.Put(i+100, i*3)
		}(i)

		// Concurrent Get
		go func(i int) {
			defer wg.Done()
			cache.Get(i)
		}(i)
	}
	wg.Wait()
}

func TestFIFOCache_NonThreadSafe(t *testing.T) {
	// Non-thread-safe cache should not be used in concurrent scenarios
	// This test just verifies it works in single-threaded context
	cache := NewFIFOCache[int, int](10, false)

	for i := 0; i < 10; i++ {
		cache.Put(i, i*2)
	}

	for i := 0; i < 10; i++ {
		if value, found := cache.Get(i); !found || value != i*2 {
			t.Errorf("Expected key %d to have value %d, got %d, found: %v", i, i*2, value, found)
		}
	}
}
