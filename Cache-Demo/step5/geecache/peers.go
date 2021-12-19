package geecache

type (
	PeerPicker interface {
		// 根据传入的 key 选择相应的结点
		PickPeer(key string) (peer PeerGetter, ok bool)
	}
	PeerGetter interface {
		// 从对应 group 查找缓存
		Get(group string, key string) ([]byte, error)
	}
)
