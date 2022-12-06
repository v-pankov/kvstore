package memcachedpool

import (
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached"
)

type ClientPool interface {
	AcquireClient() (memcached.Client, error)
	ReleaseClient(memcached.Client)
}

func AcquireUseReleaseClient(p ClientPool, useFn func(memcached.Client)) error {
	c, err := p.AcquireClient()
	if err != nil {
		return fmt.Errorf("acquire client: %w", err)
	}
	useFn(c)
	p.ReleaseClient(c)
	return nil
}
