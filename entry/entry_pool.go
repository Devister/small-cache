package entry

type Pool struct {
}

func (p *Pool) GetEntry(key, value []byte) *Entry {
	return &Entry{}
}

func (p *Pool) RecycleEntry(e *Entry) {

}
