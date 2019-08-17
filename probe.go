package probes

import (
	"sync"
)

type ProbeState int32

const (
	ProbeStateInactive ProbeState = 0
	ProbeStateActive   ProbeState = 1
)

type Probe interface {
}

type ProbeFloat64 struct {
	MaximumSignalLength int
	signal              []float64
	state               ProbeState
	stateMutex          sync.Mutex
	c                   chan float64
}

func (p *ProbeFloat64) Activate(bufferSize int) chan float64 {
	p.stateMutex.Lock()
	p.state = ProbeStateActive
	p.stateMutex.Unlock()

	p.c = make(chan float64, bufferSize)
	go func() {
		for f := range p.c {
			p.signal = append(p.signal, f)
			if len(p.signal) > p.MaximumSignalLength {
				p.signal = p.signal[len(p.signal)-p.MaximumSignalLength:]
			}
		}
	}()

	return p.c
}

func (p *ProbeFloat64) Deactivate() {
	close(p.c)

	p.stateMutex.Lock()
	p.state = ProbeStateInactive
	p.stateMutex.Unlock()
}
