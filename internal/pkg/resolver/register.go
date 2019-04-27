package resolver

import (
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

func RegisterDefault(addr string, port int, watchInterval time.Duration) {
	config := consulapi.DefaultConfig()
	config.Address = addr + ":" + strconv.Itoa(port)
	resolver.Register(&ResolverBuilder{
		WatchInterval: watchInterval,
		ClientConfig:  config,
	})
}
