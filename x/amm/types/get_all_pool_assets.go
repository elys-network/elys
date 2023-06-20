package types

func (p Pool) GetAllPoolAssets() []PoolAsset {
	copyslice := make([]PoolAsset, len(p.PoolAssets))
	copy(copyslice, p.PoolAssets)
	return copyslice
}
