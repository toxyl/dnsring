package dnsring

import (
	"container/ring"
	"time"
)

type Ring struct {
	ring *ring.Ring
}

func (rs *Ring) set(servers ...*Resolver) {
	rs.ring = ring.New(len(servers))
	for _, s := range servers {
		rs.ring.Value = s
		rs.ring = rs.ring.Next()
	}
}

func (rs *Ring) Next() *Resolver {
	rs.ring = rs.ring.Next()
	return rs.ring.Value.(*Resolver)
}

func New(servers ...Server) *Ring {
	res := []*Resolver{}
	for _, r := range servers {
		res = append(res, NewResolver(r.Host, uint(r.Port), time.Duration(r.Timeout)*time.Millisecond))
	}
	rs := &Ring{}
	rs.set(res...)
	return rs
}
