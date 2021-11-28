







## 导语

数据的存储模型主要分为 B+树和 LSM 树

![image-20210910153124287](/Users/qlzhou/Library/Application Support/typora-user-images/image-20210910153124287.png)

B+树有二叉查找树演化 增加每层节点的数量 降低树的高度 适配磁盘页 减少磁盘的 IO 操作

- 查询性能稳定 读多写少的场景







![image-20210910153311735](/Users/qlzhou/Library/Application Support/typora-user-images/image-20210910153311735.png)

LSM 树 日志结构合并树 LogStructured Merge Tree 

- 顺序 IO快于随机 IO
- 数据 CRUD 会被记录为日志 追加到磁盘的文件中 
- 适合写多读少的场景



miniDB 采用类似 LSM 的存储结构  为bitcask  也是采用顺序 IO 追加

## 常见操作

### **PUT**



Key 和 Value 被封装为 Entry 类型 内含有两者的值、大小、写入时间

磁盘文件为多个 entry 的集合



内存为 哈希表 key 为下标 value 为值

	- 磁盘记录追加
	- 内存索引更新



### **GET**

- 内存的哈希表中查找 key 对应的索引 即找到 value 在磁盘文件中的位置 取出即可



### **DEL**

- 删除操作封装为 Entry 追加到磁盘文件 标志该 Entry 类型为删除

- 内存中删除哈希表 key 对应的索引



### **Merge**



- 磁盘文件为追加 导致容量的上升
- 同一个 Key 可能存在多条 entry  前述 key 对应的 entry 无效 需要定期合并清理无效的 entry 数据 即为 merge

- 将源文件所有的 entry 取出 将有效 entry 写入临时文件 删除原始数据文件即可



## 知识点：

- 设置 offset 为文件的读取位置
-  采用 indexes 映射存储 key 和 offset   对File的 entry 进行增删改查
- 将删除操作封装为 Entry 标志记为 DEL 追加在文件中  顺序IO 的快捷
- 采用 buf 进行 Decode 和 Encode 这是读写的基础
- merge 时候创建新的 validEntry  依据 indexes 中的 key 对应的 offset 删除源 File 中的重复KeyEntry 和del 标志的 entry 构建新的 entry 并重命名 替换





- os 包内的打开 创建 权限 stat open remove rename close mkdirall modeperm
- RWMutex 读写同步锁  
  - Lock defer Unlock 
  - 读锁占用会阻止写 但不阻止读
  - 写锁时独占



Mutex锁的模式分类 and Mutex 锁底层实现

状态机和信号量

​	state sema

- 使用之后不可copy
- 正常状态和饥饿状态
  - FIFO 模式 唤醒的 goroutine 与新的请求锁的 goroutine 竞争锁的所有
  - 饥饿模式将锁的所有权从 unlock 状态交给等待队列的第一个 

- 







## 数据类型

String List Hash Set ZSet



### List 操作

LPush/Pop RPop/Push Insert Trim



### Hash operations















































































































































