package cache_simulator

type Cache interface {
	IsCached(p *Packet, update bool) (bool, *int)
	IsCachedWithFiveTuple(f *FiveTuple, update bool) (bool, *int)
	// Cache(p *Packet) []*Packet
	CacheFiveTuple(f *FiveTuple) []*FiveTuple
	InvalidateFiveTuple(f *FiveTuple)
	Clear()
	StatString() string
}

func AccessCache(c Cache, p *Packet) bool {
	hit, _ := c.IsCached(p, true)
	return hit
}
