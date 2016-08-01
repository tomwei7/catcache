package catcache

import (
	"strconv"
	"testing"
)

func TestMultipleCache(t *testing.T) {
	cache := NewMultipleCache(16, 60)
	for i := 0; i < 1024; i++ {
		//cache.Set(strconv.Itoa(i), "data "+strconv.Itoa(i)+"\n")
		cache.Set(strconv.Itoa(i), [1024 * 1024]byte{})
	}
	//cache.Display()
	//t.Log(cache)
	for {
	}
}
