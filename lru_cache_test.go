package zippyzap

import (
	"testing"
)

import "github.com/stretchr/testify/require"

func TestLRUCache_CRUD(t *testing.T) {
	cache := CreateAndStartLRUCache(2)

	defer cache.Shutdown()

	require.Equal(t, 0, cache.Size())

	testKey := "testKey"
	testVal := "testVal"
	testUpdatedVal := "updatedVal"

	testKey2 := 69  // Mature humor
	testVal2 := 420 // I live in washington

	testKey3 := "testKey3"
	testVal3 := map[string]int{}

	// ***** Test first Put() and Get() calls work *****

	cache.Put(testKey, testVal)
	require.Equal(t, 1, cache.Size())
	require.True(t, cache.ContainsKey(testKey))

	foundVal, found := cache.Get(testKey)
	require.True(t, found)
	require.Equal(t, testVal, foundVal)

	// ***** Test Put() under same key works *****

	cache.Put(testKey, testUpdatedVal)
	require.Equal(t, 1, cache.Size())
	require.True(t, cache.ContainsKey(testKey))

	foundVal, found = cache.Get(testKey)
	require.True(t, found)
	require.Equal(t, testUpdatedVal, foundVal)

	// ***** Test Put() under different key works *****

	cache.Put(testKey2, testVal2)
	require.Equal(t, 2, cache.Size())
	require.True(t, cache.ContainsKey(testKey))
	require.True(t, cache.ContainsKey(testKey2))

	foundVal, found = cache.Get(testKey2)
	require.True(t, found)
	require.Equal(t, testVal2, foundVal)

	// ***** Test third Put() keeps size at 2 (since that is max size) *****

	cache.Put(testKey3, testVal3)
	require.Equal(t, 2, cache.Size())
	_, found = cache.Get(testKey)
	require.False(t, found)
	require.True(t, cache.ContainsKey(testKey2))
	require.True(t, cache.ContainsKey(testKey3))

	foundVal, found = cache.Get(testKey3)
	require.True(t, found)
	require.Equal(t, testVal3, foundVal)

	// ***** Test Remove() works
	cache.Remove(testKey3)
	require.Equal(t, 1, cache.Size())
	require.False(t, cache.ContainsKey(testKey3))
}

func BenchmarkLRUCache_Put_SameInput(b *testing.B) {
	cache := CreateAndStartLRUCache(1)

	defer cache.Shutdown()

	for i := 0; i < b.N; i++ {
		cache.Put("key", "val")
	}
}

func BenchmarkLRUCache_Put_RandomInput(b *testing.B) {
	cache := CreateAndStartLRUCache(1)

	defer cache.Shutdown()

	for i := 0; i < b.N; i++ {
		cache.Put(i, "val")
	}
}
