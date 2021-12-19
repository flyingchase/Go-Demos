package lru

import (
	"container/list"
)

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

// vim 测试
func quicksort(nums []int) {
	if len(nums) == 0 {
		return
	}
	quicksorthelper(nums, 0, len(nums)-1)

}
func quicksorthelper(nums []int, l int, r int) {
	if l < r {
		p := paratition(nums, l, r)
		quicksorthelper(nums, l, p[0]-1)
		quicksorthelper(nums, p[1]+1, r)
	}
}
func paratition(nums []int, l, r int) []int {
	less, more := l-1, r
	for l < more {
		if nums[l] < nums[r] {
			less++
			nums[less], nums[l] = nums[l], nums[less]
			l++
		} else if nums[l] > nums[r] {
			more--
			nums[l], nums[more] = nums[more], nums[l]
		} else {
			l++
		}
	}
	nums[more], nums[r] = nums[r], nums[more]
	return []int{less + 1, more}
}
func mergesort(nums []int) {
	if len(nums) == 0 {
		return
	}
	mergesortHelper(nums, 0, len(nums)-1)
}
func mergesortHelper(nums []int, l, r int) {
	if l >= r {
		return
	}
	mid := l + (r-l)>>1
	mergesortHelper(nums, l, mid)
	mergesortHelper(nums, mid+1, r)
	merge(nums, l, mid, r)
}
func merge(nums []int, l, mid, r int) {
	p1, p2, i, helper := l, mid, 0, make([]int, r-l+1)
	for p1 <= mid && p2 <= r {
		if nums[p1] < nums[p2] {
			helper[i] = nums[p1]
			p1++
		} else {
			helper[i] = nums[p2]
			p2++
		}
		i++
	}
	copy(helper[i:], nums[p1:mid])
	copy(helper[i:], nums[p2:])
	copy(nums, helper)
}

type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
	Value int
}

func zigzaLevelTraversalBT(root *TreeNode) [][]int {
	res := make([][]int, 0)
	cur := root
	queue := make([]*TreeNode, 0)
	queue = append(queue, cur)
	flag := false
	for len(queue) > 0 {
		size := len(queue)
		lists := make([]int, 0)
		for size > 0 {
			size--
			lists = append(lists, cur.Value)
			if cur.Left != nil {
				queue = append(queue, cur.Left)
			}
			if cur.Right != nil {
				queue = append(queue, cur.Right)
			}
		}
		if !flag {
			for i, j := 0, len(lists)-1; i < j; i, j = i+1, j-1 {
				lists[i], lists[j] = lists[j], lists[i]
			}
		}
		flag = !flag
		res = append(res, lists)
	}
	return res
}
