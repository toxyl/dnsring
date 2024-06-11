package dnsring

import (
	"context"
	"net"
	"strings"
	"sync"
)

type Cache struct {
	mu         *sync.Mutex
	cache      map[string]string
	exclusions map[string]string // host -> IP
}

func (c *Cache) addExclusions(exclusions map[string]string) {
	for k, v := range exclusions {
		c.exclusions[k] = v
	}
}

func (c *Cache) has(k string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.cache[k]
	return ok
}

func (c *Cache) get(k string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.cache[k]
	if !ok {
		return k
	}
	return v
}

func (c *Cache) set(k, v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[k] = v
}

func (c *Cache) hostToIP(address string, resolver *Resolver) string {
	host, port, err := net.SplitHostPort(address)
	if err != nil && strings.Contains(err.Error(), "missing port in address") {
		host = address
		port = ""
	}
	if !cache.has(host) {
		ips, err := resolver.Resolver.LookupIP(context.Background(), "ip", host)
		if err != nil {
			cache.set(host, "")
			return ""
		}
		var resolvedIP net.IP
		for _, ip := range ips {
			if i := ip.To4(); i != nil {
				resolvedIP = i
				break
			}
			if i := ip.To16(); i != nil {
				resolvedIP = i
				break
			}
		}
		if resolvedIP == nil || resolvedIP.String() == "0.0.0.0" {
			cache.set(host, "") // the host doesn't resolve to an IP address
			return ""
		}
		ip := resolvedIP.String()

		for h, i := range c.exclusions {
			ipMatches := ip == i
			hostMatches := strings.EqualFold(host, h)
			if h[0] == '!' {
				// user wants to exclude everything matching the IP but NOT the host
				h = h[1:]
				hostMatches = !strings.EqualFold(host, h)
			}
			if ipMatches && hostMatches {
				cache.set(host, "")
				return ""
			}
		}
		cache.set(host, resolvedIP.String())
	}
	res := cache.get(host)
	if port != "" {
		return res + ":" + port
	}

	return res
}

var (
	cache = &Cache{
		mu:         &sync.Mutex{},
		cache:      map[string]string{},
		exclusions: map[string]string{},
	}
)

func AddCacheExclusions(exclusions map[string]string) {
	cache.addExclusions(exclusions)
}
