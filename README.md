# Go demo collection

## goroutine Pool Ants

将`goroutine` 池化，复用 goroutine，减轻 runtime 的调度压力

启动服务之前初始化 goroutine pool 池，Pool 维护类似栈的 LIFO 队列，存放处理任务的 Worker，在 client 提交 task 到 Pool 后，在 Pool 内部接收 task:

1. 检查 worker 队列中是否有可用的 worker，并取出 task
2. 没用可用 worker，判断是否超容
   1. 是则判断工作池是否非阻塞模式，是则返回 nil，否则阻塞等待直到 worker 被释放回到 pool
   2. 否 则新开一个 worker 处理 task
3. 每个 worker 处理完后放回 pool 队列中等待；
