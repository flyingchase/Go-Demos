package consistenhash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	// 自定义 hash 算法
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	hash.Add("6", "4", "2")
	// 真实节点 2 4 6, replicates 在 New 中为 3 则
	// 02/12/22 04/14/24 06/16/26的虚拟节点 hash
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s,shoud have yielded %s ", k, v)
		}
	}
	hash.Add("8")
	testCases["27"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s,shoud have yielded %s ", k, v)
		}
	}
}
