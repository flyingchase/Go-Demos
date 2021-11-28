package minidb

import "encoding/binary"

//entry为数据结构的封装 内含有 key val 两者的大小 创建的时间 标志位

const entryHeaderSize = 10

const (
	// 两个常量 0 和 1 分别代表 entry 的 mark 为追加还是删除
	PUT uint16 = iota
	DEL
)

type Entry struct {
	Key       []byte
	Value     []byte
	KeySize   uint32
	ValueSize uint32
	// 标志位 代表 entry 在磁盘中是追加还是删除 0=put 1=del
	Mark uint16
}

func NewEntry(key, value []byte, mark uint16) *Entry {
	return &Entry{Key: key, Value: value, Mark: mark, KeySize: uint32(len(key)), ValueSize: uint32(len(value))}
}

// 创建 Entry 的方法包括 GetSize Encode（返回字节数组） Decode（解码）
func (e *Entry) GetSize() int64 {
	// add the constvalue cap entryHeaderSize
	return int64(entryHeaderSize + e.KeySize + e.ValueSize)

}
func (e *Entry) Encode() ([]byte, error) {
	buf := make([]byte, e.GetSize())
	// 将 keysize转化为字节数组 存储在 buf 前四位
	binary.BigEndian.PutUint32(buf[0:4], e.KeySize)
	binary.BigEndian.PutUint32(buf[4:8], e.ValueSize)
	// mark 字节存储在 buf 的最后两位 8 9
	binary.BigEndian.PutUint16(buf[8:10], e.Mark)
	// 将 buf HeaderSize 后存储为 key
	copy(buf[entryHeaderSize:entryHeaderSize+e.KeySize], e.Key)
	// buf剩下部分存储 value
	copy(buf[entryHeaderSize+e.KeySize:], e.Value)
	return buf, nil
}

// 从 buf 字节数组中解码得到 entry HeaderSize中解码
// 注意非 entry 的方法 而是函数
func  Decode(buf []byte) (*Entry, error) {
	ks := binary.BigEndian.Uint32(buf[0:4])
	vs := binary.BigEndian.Uint32(buf[4:8])
	mark := binary.BigEndian.Uint16(buf[8:10])

	return &Entry{
		KeySize:   ks,
		ValueSize: vs,
		Mark:      mark,
	}, nil
}
