package idb

import (
	"sync"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3/batching"
)

const DefaultCapacity = 128

type fixed struct {
	size          int
	capacity      int
	callbackReady func()
	callbackEmit  func([]*influxdb3.Point)

	points []*influxdb3.Point
	mu     sync.Mutex
}

func NewFixed(opts ...batching.Option) *fixed {
	f := &fixed{
		capacity:      DefaultCapacity,
		callbackReady: func() {},
		callbackEmit:  func([]*influxdb3.Point) {},
	}

	for _, o := range opts {
		o(f)
	}

	f.points = make([]*influxdb3.Point, 0, f.capacity)

	return f

}

func (f *fixed) SetSize(s int) {

}

func (f *fixed) SetCapacity(s int) {
	f.capacity = s
}

func (f *fixed) SetReadyCallback(rct func()) {
	f.callbackReady = rct
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

func main() {
}
