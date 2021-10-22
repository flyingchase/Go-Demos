package geecache

// 只读的数据结构表示缓存值
type ByteView struct {
	b []byte
}

// Len 方法实现 lru 内 Value 接口
func (v ByteView) Len() int {
	return len(v.b)
}
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
func (v ByteView) String() string {
	return string(v.b)
}
