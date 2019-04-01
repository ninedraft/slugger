package main

import (
	"strconv"
)

type Range []struct{}

func (r Range) Ints() []int {
	var v = make([]int, r.Len())
	for i := range r {
		v[i] = i
	}
	return v
}

func (r Range) Floats() []float64 {
	var v = make([]float64, r.Len())
	for i := range r {
		v[i] = float64(i)
	}
	return v
}

func (r Range) Len() int {
	return len(r)
}

func (r Range) String() string {
	return "Range(0, " + strconv.Itoa(r.Len()) + ")"
}

func R(n int) Range {
	return make(Range, n)
}

func (r Range) StreamTo(stop <-chan struct{}, target chan<- int) {
	for i := range r {
		select {
		case <-stop:
			return
		case target <- i:
			continue
		}
	}
}

func (r Range) Stream(stop <-chan struct{}) <-chan int {
	var stream = make(chan int)
	go func() {
		defer close(stream)
		for i := range r {
			select {
			case <-stop:
				return
			case stream <- i:
				continue
			}
		}
	}()
	return stream
}
