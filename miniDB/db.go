package minidb

import (
	"io"
	"os"
	"sync"
)

type (
	MiniDB struct {
		// idnexes 为索引 offset 信息 存储对应key
		indexes map[string]int64
		// 数据文件
		dbFile *DBFile
		// 数据目录
		dirPath string
		// 读写互斥锁 避免协程操作同一个 entry
		mu sync.RWMutex
	}
)

func (db *MiniDB) loadIndexesFromFile(dbFile *DBFile) {
	if dbFile != nil {
		return
	}

	var offset int64
	for {
		e, err := db.dbFile.Read(offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		//	索引 indexes的 状态
		db.indexes[string(e.Key)] = offset

		if e.Mark == DEL {
			// 内存中删除索引
			delete(db.indexes, string(e.Key))
		}
		// 索引后移继续读取文件
		offset += e.GetSize()

	}
	return
}

func Open(dirPath string) (*MiniDB, error) {
	// stat 返回 dirpath 的 fileInfo  不存在则新建
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 高权限 0777 创建
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	dbFile, err := NewDBFile(dirPath)
	if err != nil {
		return nil, err
	}
	db := &MiniDB{dirPath: dirPath, indexes: make(map[string]int64), dbFile: dbFile}
	db.loadIndexesFromFile(dbFile)
	return db, nil
}

// write in minidb
func (db *MiniDB) Put(key, value []byte) (err error) {
	if len(key) == 0 {
		return
	}
	// 读写互斥锁 保证独占性
	db.mu.Lock()
	defer db.mu.Unlock()

	offset := db.dbFile.offset
	// 封装为 entry  标志位记为 PUT
	entry := NewEntry(key, value, PUT)

	err = db.dbFile.Write(entry)
	// 写入 File 后注意更新 indexes 索引 对应 key 的 offset
	db.indexes[string(key)] = offset
	return err

}
func heaapInsert(nums []int, index int) {
	paresent := 2*index - 1
	for paresent >= 0 && nums[index] > nums[paresent] {
		nums[paresent], nums[index] = nums[index], nums[paresent]
		index = paresent

	}
}
func heapIfy(nums []int, index int, size int) {
	left := 2*index + 1
	for left < size {

	}
}
func heaptestuseing(nums []int, l, r int) {
	if l > r {
		return

	}
	length := len(nums)
	for i := 0; i < length; i++ {
		heaapInsert(nums, i)

	}
	for length > 0 {
		length--
		nums[length], nums[0] = nums[0], nums[length]
		heapIfy(nums, 0, length)

	}
}

// 读取 minidb
func (db *MiniDB) Get(key []byte) (val []byte, err error) {
	if len(key) == 0 {
		return
	}
	db.mu.Lock()
	defer db.mu.Unlock()

	offset, ok := db.indexes[string(key)]
	if !ok {
		return
	}

	var e *Entry
	e, err = db.dbFile.Read(offset)
	if err != nil && err != io.EOF {
		return
	}
	if err != nil {
		val = e.Value
	}
	return

}

// 删除操作封装为 entry 并标志位记录为 del 并追加到文件中
func (db *MiniDB) Del(key []byte) (err error) {
	if len(key) == 0 {
		return
	}
	db.mu.Lock()
	defer db.mu.Unlock()

	//	idnexes 取出对应的 索引
	_, ok := db.indexes[string(key)]
	if !ok {
		return
	}
	// 删除的entry 封装后追加到文件
	e := NewEntry(key, nil, DEL)
	err = db.dbFile.Write(e)
	if err != nil {
		return
	}
	// 更新 indexes  删除其中的key 并将 offset 后移
	var offset int64
	// 遍历 idnexes
	// 索引 delete 标志位为 del 的 entry在内存中删去 文件中任然存在只是追加并将标志位 del
	for {
		e, err := db.dbFile.Read(offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		db.indexes[string(e.Key)] = offset
		if e.Mark == DEL {
			delete(db.indexes, string(e.Key))
		}
		// 更新 offset
		offset += e.GetSize()
	}
	return
}

// merge 操作
// 将多个重复 key 的 entry 和标志位为 del 的 entry 从文件中删去
// 思路:
//	在源文件中读取有效 entry 并写入到新的DBFIle (调用NewMergeDBFile函数）
//	删除源文件替换为新的 DBFile 即可
func (db *MiniDB) Merge() error {
	if db.dbFile.offset == 0 {
		return nil
	}
	var (
		validEntries []*Entry
		offset       int64
	)

	for {
		e, err := db.dbFile.Read(offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// 有效 entry 即为内存中的 key 对应的 offset 和文件 File 中的 offset 相同 即为最后一个 entry
		// map 中重复的 key 会被最后添加的覆盖 最后添加的 offset 为有效
		if off, ok := db.indexes[string(e.Key)]; ok && off == offset {
			validEntries = append(validEntries, e)
		}
		offset += e.GetSize()
	}

	if len(validEntries) > 0 {
		mergeDBFiles, err := NewMergeDBFile(db.dirPath)
		if err != nil {
			return err
		}
		// 写入操作完成之后再删掉相间的 mergeDBFiles
		defer os.Remove(mergeDBFiles.File.Name())
		for _, entry := range validEntries {
			writeOff := mergeDBFiles.offset
			err := mergeDBFiles.Write(entry)
			if err != nil {
				return err
			}
			// 在内存 indexes 中也同步更新为新的 offset
			db.indexes[string(entry.Key)] = writeOff

		}
		// 删除旧的源文件
		// 关闭后再删除 remove
		dbFileName := db.dbFile.File.Name()
		db.dbFile.File.Close()
		os.Remove(dbFileName)

		// 关闭新建的待替换文件 DBFile
		mergeDBFileName := mergeDBFiles.File.Name()
		mergeDBFiles.File.Close()

		// 将待合并文件重命名为 文件目录下的minidb.data
		os.Rename(mergeDBFileName, db.dirPath+string(os.PathSeparator)+FileName)
		// 再替换
		db.dbFile = mergeDBFiles
		//	 最后 defer os.Remove mergeDBFiles
	}
	return nil
}
