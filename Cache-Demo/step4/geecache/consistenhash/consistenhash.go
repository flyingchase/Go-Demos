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
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[index%len(m.keys)]]
}
