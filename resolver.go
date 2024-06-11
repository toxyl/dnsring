package dnsring

import (
	"context"
	"fmt"
	"net"
	"time"
)

type Resolver struct {
	Host     string        `yaml:"host"`
	Port     uint          `yaml:"port"`
	Resolver *net.Resolver `yaml:"-"`
}

func (dr *Resolver) HostToIP(host string) string {
	return cache.hostToIP(host, dr)
}

func NewResolver(host string, port uint, timeout time.Duration) *Resolver {
	d := &net.Dialer{
		Timeout: timeout,
	}
	return &Resolver{
		Host: host,
		Port: port,
		Resolver: &net.Resolver{
			PreferGo:     true,
			StrictErrors: false,
			Dial: func(ctx context.Context, network string, address string) (net.Conn, error) {
				return d.DialContext(ctx, network, fmt.Sprintf("%s:%d", host, port))
			},
		},
	}
}
