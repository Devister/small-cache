package bucket

type Bucket interface {
	Set(hkey uint64, key []byte, value []byte)
	Get(hkey uint64, key []byte) []byte
	Delete(hkey uint64, key []byte)
}
