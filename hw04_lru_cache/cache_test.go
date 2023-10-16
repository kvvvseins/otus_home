package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("проверка выталкивания из-за размера", func(t *testing.T) {
		cachedValues := []string{"aaa", "bbb", "ccc"}

		c := createFullCache(cachedValues)

		c.Set("ddd", 100)
		require.Equal(t, len(cachedValues), c.capacity)

		val, ok := c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 1, val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Equal(t, nil, val)

		require.Equal(t, len(cachedValues), c.queue.Len())
	})

	t.Run("проверка выталкивания давно используемых элементов", func(t *testing.T) {
		cachedValues := []string{"zzzzz", "eeee", "ddd", "ccc", "bbb", "aaa"}

		c := createFullCache(cachedValues)

		c.Set("ddd", 100)
		c.Get("bbb")
		c.Get("zzzzz")
		c.Get("eeee")

		// вытесняем ccc
		c.Set("ccc2", 111)

		val, ok := c.Get("ccc")
		require.False(t, ok)
		require.Equal(t, nil, val)

		require.Equal(t, len(cachedValues), c.queue.Len())
	})

	t.Run("тестируем clear", func(t *testing.T) {
		cachedValues := []string{"zzzzz", "eeee", "ddd", "ccc", "bbb", "aaa"}

		c := createFullCache(cachedValues)
		c.Clear()

		require.Equal(t, len(c.items), c.queue.Len())
		require.Equal(t, 0, c.queue.Len())
		require.Equal(t, 0, len(c.items))
	})
}

func createFullCache(cachedValues []string) *lruCache {
	c := &lruCache{
		capacity: len(cachedValues),
		queue:    NewList(),
		items:    make(map[Key]*ListItem, len(cachedValues)),
		mu:       &sync.Mutex{},
	}

	// заполняем кеш
	for i, v := range cachedValues {
		_ = c.Set(Key(v), i)
	}

	return c
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
