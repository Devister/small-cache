package bucket

type Bucket interface {
	Set(hkey uint64, key []byte, value []byte) error
	Get(hkey uint64, key []byte) ([]byte, error)
	Delete(hkey uint64, key []byte) error
}
