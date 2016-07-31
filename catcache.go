package catcache

import (
	"sync"
	"time"
)

const (
	Expired = iota
	Nonexist
)

type CacheError struct {
	Code int
	Msg  string
}

func (err CacheError) Error() string {
	return err.Msg
}

//数据节点
type CacheData struct {
	Key      string
	Previous *CacheData
	Next     *CacheData
	Expired  int64       //过期时间
	Data     interface{} //数据
}

//多键值缓存
type MultipleCache struct {
	MaxLength     int
	Length        int
	ExpireIn      int64
	cacheListHead *CacheData
	cacheMap      map[string]*CacheData
	rwm           sync.RWMutex
}

//单个变量缓存
type SingleCache struct {
	ExpireIn  int64
	Expired   int64
	cacheData interface{}
}

func NewSingleCache(expireIn int64) *SingleCache {
	return &SingleCache{
		ExpireIn: Expired,
	}
}

func (p *SingleCache) Set(v interface{}) {
	p.cacheData = v
	p.Expired = time.Now().Unix() + p.ExpireIn
}

func (p *SingleCache) Get() (interface{}, error) {
	timeStamp := time.Now().Unix()
	if timeStamp > p.Expired {
		return p.cacheData, CacheError{Expired, "Cache Expired"}
	}
	return p.cacheData, nil
}

func NewMultipleCache(maxLength int, expireIn int64) *MultipleCache {
	if maxLength < 1 {
		panic("Length can't less 1")
	}
	return &MultipleCache{
		MaxLength:     maxLength,
		Length:        0,
		ExpireIn:      expireIn,
		cacheListHead: nil,
		cacheMap:      make(map[string]*CacheData),
	}
}

//设置一个key-value
func (p *MultipleCache) Set(key string, v interface{}) {
	if p.Length >= p.MaxLength {
		p.del(p.cacheListHead.Previous)
	}
	p.rwm.Lock()
	defer p.rwm.Unlock()
	cacheData := &CacheData{
		Key:     key,
		Expired: time.Now().Unix() + p.ExpireIn,
		Data:    v,
	}
	if p.cacheListHead == nil {
		cacheData.Next = cacheData
		cacheData.Previous = cacheData
		p.cacheListHead = cacheData
	} else {
		cacheData.Previous = p.cacheListHead.Previous
		p.cacheListHead.Previous = cacheData
		cacheData.Next = p.cacheListHead
		p.cacheListHead = cacheData
	}
	p.cacheMap[key] = cacheData
	p.Length++
	return
}

//获取一个key-value
func (p *MultipleCache) Get(key string) (interface{}, error) {
	p.rwm.RLock()
	defer p.rwm.RUnlock()
	timeStamp := time.Now().Unix()
	if val, ok := p.cacheMap[key]; ok {
		if val.Expired < timeStamp {
			go p.del(val)
			return val.Data, CacheError{Expired, "Key Expired"}
		}
		return val.Data, nil
	}
	return nil, CacheError{Nonexist, "key Nonexist"}
}

//删除一个cachedata
func (p *MultipleCache) del(cacheData *CacheData) {
	p.rwm.Lock()
	defer p.rwm.Unlock()
	delete(p.cacheMap, cacheData.Key)
	p.Length--
	if cacheData == p.cacheListHead {
		if p.Length == 1 {
			p.cacheListHead = nil
			return
		}
		p.cacheListHead = cacheData.Next
	}
	cacheData.Next.Previous = cacheData.Previous
	cacheData.Previous.Next = cacheData.Next
	return
}
