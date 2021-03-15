package zippyzap

import (
	"math/rand"
	"testing"
	"time"
)

import "github.com/stretchr/testify/require"

func init() {
	rand.Seed(time.Now().Unix())
}

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

// The following benchmarks mirror the benchmarks
// written by hashicorp so that the 2 libs can be
// compared. Their benchmarks can be found here:
// https://github.com/hashicorp/golang-lru/blob/80c98217689d6df152309d574ccc682b21dc802c/lru_test.go

func BenchmarkLRU_Rand(b *testing.B) {
	l := CreateAndStartLRUCache(8192)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = rand.Int63() % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.Put(trace[i], trace[i])
		} else {
			_, ok := l.Get(trace[i])
			if ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))

	b.StopTimer()
	l.Shutdown()
	b.StartTimer()
}

func BenchmarkLRU_Freq(b *testing.B) {
	l := CreateAndStartLRUCache(8192)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = rand.Int63() % 16384
		} else {
			trace[i] = rand.Int63() % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Put(trace[i], trace[i])
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		_, ok := l.Get(trace[i])
		if ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))

	b.StopTimer()
	l.Shutdown()
	b.StartTimer()
}
