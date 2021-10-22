package lru

import "container/list"

type (
	Cache struct {
		// 允许使用的最大内存
		maxBytes int64
		// 当前使用的内存
		nbytes int64
		ll     *list.List
		cache  map[string]*list.Element
		// 某条记录删除时候的回调函数 可为 nil
		OnEvicted func(key string, value Value)
	}
	entrty struct {
		key   string
		value Value
	}
	// 值实现 Len() 方法返回所占的内存大小 则均可视为 value
	Value interface {
		Len() int
	}
)

// 实例化，实现 New()函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 从字典从查找对应的双向链表的结点，再移动到队首
/*
 */
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 移动到队首（约定 front 为队首）
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entrty)
		return kv.value, true
	}
	return
}

// delete 即移出队尾
/*
	-
*/
func (c *Cache) RemoveOldest() {
	// 取出队尾结点，删除
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entrty)
		delete(c.cache, kv.key)
		// update 目前所使用的内存
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 回调函数非空则使用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 更新结点
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entrty)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 新增结点 先在队首添加新节点，再在map 存储 key 和对应的 value 的映射
		ele := c.ll.PushFront(&entrty{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 新增结点超出设定值 移出 oldest
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 添加数据的数量
func (c *Cache) Len() int {
	return c.ll.Len()
}
