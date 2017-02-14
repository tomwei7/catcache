# Catcache

simple golang memory cache use lru


### usage

```golang
// NewMultipleCache(maxLength, exprieIn)
cache := NewMultipleCache(1024, 60)
cache.Set(k, v)
cache.Get(k)
```
