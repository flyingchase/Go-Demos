package geecache

type (
	PeerPicker interface {
		PickerPeer(key string) (peer PeerGetter, ok bool)
	}
	PeerGetter interface {
		Get(group string, key string) ([]byte, error)
	}
)


