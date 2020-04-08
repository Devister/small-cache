package small_cache

const (
	defaultBucketNum = 1024
)

type CacheConfig struct {
	ItemNum         uint64
	BucketNum       uint64
	HardMemoryLimit int // hard memory limit in MB
}

func parseConfig(cfg *CacheConfig) error {
	if cfg.BucketNum == 0 {
		if cfg.ItemNum == 0 {
			cfg.BucketNum = defaultBucketNum
		} else {
			cfg.BucketNum = powCeiling(cfg.ItemNum / 256)
		}
	}
	return nil
}

func powCeiling(n uint64) uint64 {
	ret := uint64(1)
	m := n
	for n > 1 {
		n = n >> 1
		ret = ret << 1
	}
	if ret == m {
		return ret
	} else {
		return ret << 1
	}
}
