package memcachedpool

import "github.com/vdrpkv/kvstore/internal/pkg/memcached"

type MemcachedClient struct {
	ClientPool
}

var _ memcached.Client = MemcachedClient{}

func (c MemcachedClient) Set(key string, flags int16, exptime int, val string) (clientErr error) {
	use := func(c memcached.Client) { clientErr = c.Set(key, flags, exptime, val) }

	if acqErr := c.aur(use); acqErr != nil {
		return acqErr
	}

	return

}

func (c MemcachedClient) Gat(exptime int, keys ...string) (items []memcached.Item, clientErr error) {
	use := func(c memcached.Client) { items, clientErr = c.Gat(exptime, keys...) }

	if acqErr := c.aur(use); acqErr != nil {
		return nil, acqErr
	}

	return
}

func (c MemcachedClient) Get(keys ...string) (items []memcached.Item, clientErr error) {
	use := func(c memcached.Client) { items, clientErr = c.Get(keys...) }

	if acqErr := c.aur(use); acqErr != nil {
		return nil, acqErr
	}

	return
}

func (c MemcachedClient) Delete(key string) (clientErr error) {
	use := func(c memcached.Client) { clientErr = c.Delete(key) }

	if acqErr := c.aur(use); acqErr != nil {
		return acqErr
	}

	return
}

func (c MemcachedClient) aur(use func(c memcached.Client)) error {
	return AcquireUseReleaseClient(c, use)
}
