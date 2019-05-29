package resolver

import (
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

type ResolverBuilder struct {
	WatchInterval time.Duration
	ClientConfig  *consulapi.Config
}

func (r *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	consul, err := consulapi.NewClient(r.ClientConfig)
	if err != nil {
		return nil, err
	}

	res := Resolver{
		target:        target,
		cc:            cc,
		consul:        consul,
		addr:          make(chan []resolver.Address, 1),
		done:          make(chan struct{}, 1),
		watchInterval: r.WatchInterval,
	}

	go res.updater()
	go res.watcher()
	res.resolve()

	return &res, nil
}

func (r *ResolverBuilder) Scheme() string {
	return Scheme
}
