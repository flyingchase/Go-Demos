package geecache

import (
	"fmt"
	"log"
	"sync"
)

type (
	Getter interface {
		Get(key string) ([]byte, error)
	}
	// GetterFunc 通过函数实现 Getter接口的 Get 方法，接口型函数
	GetterFunc func(key string) ([]byte, error)
)

// GetterFunc 含有 Get 方法，在方法内调用自身，实现了接口 Getter
// GetterFunc 是实现接口的函数
// 接口型函数只应用于接口内部之定义一个方法的接口
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	// GetGroup 只涉及特定名称 Group 读取，不涉及写操作 故使用读写锁
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// 创建 Group的一个实例 g
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	// g 纳入 groups 存储
	groups[name] = g
	return g
}

// GetGroup 返回所创建的 group 的实例 g 不存在则 nil
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is empty and required")
	}
	// 缓存命中
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCacha] hit")
		return v, nil
	}
	return g.load(key)

}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
