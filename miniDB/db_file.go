package minidb

import (
	"os"
)

// 定义字符串常量 在 file创建和 merge 中使用
const FileName = "minidb.data"
const MergeFileName = "minidb.data.merge"

type DBFile struct {
	File *os.File
	// offset偏移量代表读取CRUD的位置与文件开头的偏移
	offset int64
}

// 创建数据文件 通过传入的路径参数
func NewDBFile(path string) (*DBFile, error) {
	// 通过传入的路径拼接 name
	fileName := path + string(os.PathSeparator) + FileName
	return newInternal(fileName)

}

func newInternal(fileName string) (*DBFile, error) {
	//  调用 os.OpenFile 不存在则创建 读写
	// 0644 代表权限为仅所有者读写 其他用户组只有读权限 不可执行
	// 0777 最高权限 读写执行均可
	// 0666 所有人具有读写权限 但均无法执行
	/*
		通过 os.FileMode(0644).String() 返回 -rw-r--r--得到权限具体信息
	*/
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	// stat 为 FileInfo 形容 file的Name
	stat, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}
	return &DBFile{File: file, offset: stat.Size()}, nil

}

// 新建合并时候所需的数据文件
func NewMergeDBFile(path string) (*DBFile, error) {
	fileName := path + string(os.PathSeparator) + MergeFileName
	return newInternal(fileName)
}

// 构建 DBFile 的方法
// Read 注意读取位置的偏移量
func (df *DBFile) Read(offset int64) (e *Entry, err error) {
	buf := make([]byte, entryHeaderSize)
	if _, err := df.File.ReadAt(buf, offset); err != nil {
		return
	}
	if e, err = Decode(buf); err != nil {
		return
	}
	offset += entryHeaderSize
	if e.KeySize > 0 {
		key := make([]byte, e.KeySize)
		if _, err = df.File.ReadAt(key, offset); err != nil {
			return
		}
		e.Key = key
	}
	offset += int64(e.KeySize)
	if e.ValueSize > 0 {
		value := make([]byte, e.ValueSize)
		if _, err = df.File.ReadAt(value, offset); err != nil {
			return
		}
		e.Value = value
	}
	return
}

// write new e to entry

func (df *DBFile) Write(e *Entry) (err error) {
	enc, err := e.Encode()
	if err != nil {
		return err
	}
	_, err = df.File.WriteAt(enc, df.offset)
	// 偏移量对应后移
	df.offset += e.GetSize()
	return
}
