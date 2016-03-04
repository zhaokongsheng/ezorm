package cache

import (
	"github.com/ezbuy/ezorm/codec"
	"github.com/golang/groupcache"
)

type GroupCache struct {
	group *groupcache.Group
	codec codec.Codec
}

func NewGroupCache(group *groupcache.Group, codec codec.Codec) *GroupCache {
	return &GroupCache{
		group: group,
		codec: codec,
	}
}

func (c *GroupCache) Get(key string, dest interface{}) error {
	var data []byte
	if err := c.group.Get(nil, key, groupcache.AllocatingByteSliceSink(&data)); err != nil {
		return err
	}

	return c.codec.Decode(data, dest)
}

var _ Cache = (*GroupCache)(nil)
