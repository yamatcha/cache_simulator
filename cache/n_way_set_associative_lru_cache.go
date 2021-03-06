package cache

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"hash/crc32"
)

type NWaySetAssociativeLRUCache struct {
	Sets []FullAssociativeLRUCache // len(Sets) = Size / Way, each size == Way
	Way  uint
	Size uint
}

func fiveTupleToBigEndianByteArray(f *FiveTuple) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, *f)
	return buf.Bytes()
}

func (cache *NWaySetAssociativeLRUCache) StatString() string {
	return ""
}

func (cache *NWaySetAssociativeLRUCache) IsCached(p *Packet, update bool) (bool, *int) {
	return cache.IsCachedWithFiveTuple(p.FiveTuple(), update)
}

func (cache *NWaySetAssociativeLRUCache) setIdxFromFiveTuple(f *FiveTuple) uint {
	maxSetIdx := cache.Size / cache.Way
	crc := crc32.ChecksumIEEE(fiveTupleToBigEndianByteArray(f))
	return uint(crc) % maxSetIdx
}

func (cache *NWaySetAssociativeLRUCache) IsCachedWithFiveTuple(f *FiveTuple, update bool) (bool, *int) {
	setIdx := cache.setIdxFromFiveTuple(f)
	return cache.Sets[setIdx].IsCachedWithFiveTuple(f, update) // TODO: return meaningful value
}

func (cache *NWaySetAssociativeLRUCache) CacheFiveTuple(f *FiveTuple) []*FiveTuple {
	setIdx := cache.setIdxFromFiveTuple(f)
	return cache.Sets[setIdx].CacheFiveTuple(f)
}

func (cache *NWaySetAssociativeLRUCache) InvalidateFiveTuple(f *FiveTuple) {
	setIdx := cache.setIdxFromFiveTuple(f)
	cache.Sets[setIdx].InvalidateFiveTuple(f)
}

func (cache *NWaySetAssociativeLRUCache) Clear() {
	panic("Not implemented")
}

func (cache *NWaySetAssociativeLRUCache) Description() string {
	return "NWaySetAssociativeLRUCache"
}

func (cache *NWaySetAssociativeLRUCache) ParameterString() string {
	return fmt.Sprintf("{\"Type\": \"%s\", \"Way\": %d, \"Size\": %d}", cache.Description(), cache.Way, cache.Size)
}

func NewNWaySetAssociativeLRUCache(size, way uint) *NWaySetAssociativeLRUCache {
	if size%way != 0 {
		panic("Size must be multiplier of way")
	}

	sets_size := size / way
	sets := make([]FullAssociativeLRUCache, sets_size)

	for i := uint(0); i < sets_size; i++ {
		sets[i] = *NewFullAssociativeLRUCache(way)
	}

	return &NWaySetAssociativeLRUCache{
		Sets: sets,
		Way:  way,
		Size: size,
	}
}
