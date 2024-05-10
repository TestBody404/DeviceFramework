package Cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

const (
	segmentCount               = 16
	int64One     int64         = 1
	int64Zero    int64         = 0
	negInt64One  int64         = -1
	intTwo                     = 2
	hashInit     uint32        = 2166136261
	prime32      uint32        = 16777619
	twentyYears  time.Duration = 20 * 365 * 24 * time.Hour
)

var (
	notInitErr = errors.New("not initializes")
	paraErr    = errors.New("parameter error")
)

type cacheEle struct {
	key        string
	data       interface{}
	expireTime int64
}

type lruCache struct {
	maxSize    int
	elemIndex  map[string]*list.Element
	*list.List //双向链表可以保持元素的插入顺序，LRU缓存需要频繁地将最近使用的元素移到链表头部，在该策略下可以实现较快的增删查改操作
	mu         sync.Mutex
}

// ConcurrencyLRUCache is a memory-based LRU local cache, default total 16 segment to improve concurrent performance
// LRU is not  real least recently used for the total cache,but just for each buket
// we just need a proper method to clear cache
type ConcurrencyLRUCache struct {
	segment    int
	cacheBuket [segmentCount]*lruCache
}

// Set create or update an element using key
//
//	key:    The identity of an element
//	value:  new value of the element
//	expireTime:    expire time, positive int64 or -1 which means never overdue
func (cl *ConcurrencyLRUCache) Set(key string, value interface{}, expireTime time.Duration) error {
	if cl == nil || cl.cacheBuket[0] == nil {
		return notInitErr
	}
	if expireTime < time.Duration(negInt64One) || expireTime > twentyYears {
		return paraErr
	}
	cacheIndex := cl.index(key)
	if cacheIndex < 0 || cacheIndex >= segmentCount {
		return errors.New("index out of valid value")
	}
	return cl.cacheBuket[cacheIndex].setValue(key, value, expireTime)
}

// Get get the value of a cached element by key. If key do not exist, this function will return nil and an error msg
//
//	key:    The identity of an element
//	return:
//	    value:  the cached value, nil if key do not exist
//	    err:    error info, nil if value is not nil
func (cl *ConcurrencyLRUCache) Get(key string) (interface{}, error) {
	if cl == nil || cl.cacheBuket[0] == nil {
		return nil, notInitErr
	}
	cacheIndex := cl.index(key)
	if cacheIndex < 0 || cacheIndex >= segmentCount {
		return nil, errors.New("index out of valid value")
	}
	return cl.cacheBuket[cacheIndex].getValue(key)
}

// Delete delete the value  by key, no error returned
func (cl *ConcurrencyLRUCache) Delete(key string) {
	if cl == nil || cl.cacheBuket[0] == nil {
		return
	}
	cacheIndex := cl.index(key)
	if cacheIndex < 0 || cacheIndex >= segmentCount {
		return
	}
	cl.cacheBuket[cacheIndex].delValue(key)
}

// SetIfNX if the key not exist or expired, will set the new value to cache and return true ,otherwise return false
func (cl *ConcurrencyLRUCache) SetIfNX(key string, value interface{}, expireTime time.Duration) bool {
	if cl == nil || cl.cacheBuket[0] == nil {
		return false
	}
	if expireTime < time.Duration(negInt64One) || expireTime > twentyYears {
		return false
	}
	cacheIndex := cl.index(key)
	if cacheIndex < 0 || cacheIndex >= segmentCount {
		return false
	}
	return cl.cacheBuket[cacheIndex].setIfNotExist(key, value, expireTime)
}

// index calculate the key hashcode and index the right buket
func (cl *ConcurrencyLRUCache) index(key string) int {
	var hash = hashInit
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return int(hash & (uint32(cl.segment) - 1))
}

// New create an instance of ConcurrencyLRUCache
// maxEntries  the cache size, will to convert to (n/16+n%16>0?1:0)*16
func New(maxEntries int) *ConcurrencyLRUCache {
	if maxEntries <= 0 {
		return nil
	}
	size := maxEntries / segmentCount
	remain := maxEntries % segmentCount
	if remain > 0 {
		size += 1
	}
	var cache [segmentCount]*lruCache
	for i := 0; i < segmentCount; i++ {
		cache[i] = &lruCache{
			maxSize:   size,
			elemIndex: make(map[string]*list.Element, segmentCount),
			List:      list.New(),
			mu:        sync.Mutex{},
		}
	}
	return &ConcurrencyLRUCache{
		segment:    segmentCount,
		cacheBuket: cache,
	}
}

func (c *lruCache) setValue(key string, value interface{}, expireTime time.Duration) error {
	if c == nil || c.elemIndex == nil {
		return errors.New("not initializes")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.elemIndex[key]
	if !ok {
		// if the cache not exist
		c.setInner(key, value, expireTime)
		return nil
	}
	ele, ok := v.Value.(*cacheEle)
	if !ok {
		c.safeDeleteByKey(key, v)
		return errors.New("cacheElement convert failed")
	}
	c.MoveToFront(v)
	pkgElement(ele, value, expireTime)
	return nil
}

func pkgElement(ele *cacheEle, value interface{}, expireTime time.Duration) {
	ele.data = value
	if expireTime == time.Duration(negInt64One) {
		ele.expireTime = negInt64One
		return
	}
	ele.expireTime = time.Now().UnixNano() + int64(expireTime)
}

func (c *lruCache) getValue(key string) (interface{}, error) {
	if c == nil || c.elemIndex == nil {
		return nil, errors.New("not initializes")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.elemIndex[key]
	if !ok {
		return nil, errors.New("no value found")
	}
	c.MoveToFront(v)
	ele, ok := v.Value.(*cacheEle)
	if !ok {
		c.safeDeleteByKey(key, v)
		return nil, errors.New("cacheElement convert failed")
	}
	if ele.expireTime != negInt64One && time.Now().UnixNano() > ele.expireTime {
		// if  cache expired
		c.safeDeleteByKey(key, v)
		return nil, errors.New("the key was expired")
	}
	return ele.data, nil
}

// Delete delete an element
func (c *lruCache) delValue(key string) {
	if c == nil || c.elemIndex == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.elemIndex[key]; ok {
		c.safeDeleteByKey(key, v)
	}
}

func (c *lruCache) setIfNotExist(key string, value interface{}, expireTime time.Duration) bool {
	if c == nil || c.elemIndex == nil {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.elemIndex[key]
	if !ok {
		// if the cache not exist
		c.setInner(key, value, expireTime)
		return true
	}
	ele, ok := v.Value.(*cacheEle)
	if !ok {
		c.safeDeleteByKey(key, v)
		return false
	}
	c.MoveToFront(v)
	if ele.expireTime == negInt64One || time.Now().UnixNano() < ele.expireTime {
		return false
	}
	// if  cache expired
	pkgElement(ele, value, expireTime)
	return true
}

func (c *lruCache) setInner(key string, value interface{}, expireTime time.Duration) {
	if c == nil {
		return
	}
	if c.Len()+1 > c.maxSize {
		c.safeRemoveOldest()
	}
	newElem := &cacheEle{
		key:        key,
		data:       value,
		expireTime: negInt64One,
	}
	if expireTime != time.Duration(negInt64One) {
		newElem.expireTime = time.Now().UnixNano() + int64(expireTime)
	}
	e := c.PushFront(newElem)
	c.elemIndex[key] = e
}

func (c *lruCache) safeDeleteByKey(key string, v *list.Element) {
	if c == nil {
		return
	}
	c.List.Remove(v)
	delete(c.elemIndex, key)
}

func (c *lruCache) safeRemoveOldest() {
	if c == nil {
		return
	}
	v := c.List.Back()
	if v == nil {
		return
	}
	c.List.Remove(v)
	ele, ok := v.Value.(*cacheEle)
	if !ok {
		return
	}
	delete(c.elemIndex, ele.key)
}
