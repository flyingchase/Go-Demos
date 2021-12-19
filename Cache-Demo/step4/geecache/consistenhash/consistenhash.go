package consistenhash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type (
	Hash func(data []byte) uint32
	Map  struct {
		hash Hash
		// 虚拟节点倍数
		replicas int
		// 哈希环
		keys []int
		// 虚拟节点与真实节点之间的映射表 虚拟节点的哈希值-真实节点名称
		hashMap map[int]string
	}
)

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	// 默认为 crc32.ChecksumIEEE 算法
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 通过添加编号区分不同的虚拟节点
			// 虚拟节点名称 strconv.Itoa(i)+key
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// 环上 hash 值排序
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// 顺时针查找第一个匹配的虚拟节点下标 index
	index := sort.Search(len(m.keys), func(i int) bool {
		// 满足 return true 的最小下标，即最近的匹配节点 index
		return m.keys[i] >= hash
	})
	// 环状 m.keys[]，故取余数
	return m.hashMap[m.keys[index%len(m.keys)]]
}
