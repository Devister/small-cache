package entry

type Entry struct {
	Len int
}

func (e *Entry) Key() []byte {
	return []byte{}
}

func (e *Entry) Value() []byte {
	return []byte{}
}

func (e *Entry) Cap() int {
	return 0
}

func (e *Entry) Set(key, value []byte) error {
	return nil
}
