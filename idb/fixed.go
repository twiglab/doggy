package idb

import (
	"sync"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

const DefaultCapacity = 128

type fixed struct {
	size         int
	capacity     int
	callbackEmit func([]*influxdb3.Point)

	points []*influxdb3.Point
	mu     sync.Mutex
}

func NewFixed(size int) *fixed {
	f := &fixed{
		capacity: size,
		points:   make([]*influxdb3.Point, 0, size),
	}
	return f
}

func (f *fixed) SetEmitCallback(callback func([]*influxdb3.Point)) {
	f.callbackEmit = callback
}

func (f *fixed) AddOne(p *influxdb3.Point) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.size == f.capacity {
		f.callbackEmit(f.emitPoints())
	}

	f.points[f.size] = p
	f.size = f.size + 1

	return nil
}

func (f *fixed) emitPoints() []*influxdb3.Point {
	points := f.points[:f.size-1]
	f.size = 0
	return points
}
