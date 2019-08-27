// Package probes provides types to inspect varying numerical values safely
// across many goroutines.
package probes

import (
	"math"
	"sync"
)

// Probe types collect floating point values and expose those values safely.
type Probe struct {
	// The maximum length of the probe's internal signal history.
	MaximumSignalLength int

	// The length of the input channel's buffer.
	InputBufferLength int

	// The probe's input channel.
	C chan float64

	signal      []float64
	signalMutex sync.Mutex
	active      bool
	activeMutex sync.Mutex
}

// MARK: Initializers

// NewProbe creates and returns a new probe.
func NewProbe() *Probe {
	return &Probe{
		MaximumSignalLength: math.MaxInt32,
		InputBufferLength:   1,
	}
}

// MARK: Public methods

// Activate activates the probe and begins waiting for signal values over its
// input channel.
func (p *Probe) Activate() {
	if p.IsActive() {
		return
	}

	p.active = true
	p.C = make(chan float64, p.InputBufferLength)

	go func() {
		for f := range p.C {
			p.signalMutex.Lock()
			p.signal = append(p.signal, f)
			if len(p.signal) > p.MaximumSignalLength {
				p.signal = p.signal[len(p.signal)-p.MaximumSignalLength:]
			}
			p.signalMutex.Unlock()
		}
	}()
}

// Deactivate deactivates the probe and returns the signal it collected.
func (p *Probe) Deactivate() {
	if !p.IsActive() {
		return
	}

	p.active = false
	close(p.C)
}

// IsActive returns the state of the probe.
func (p *Probe) IsActive() bool {
	p.activeMutex.Lock()
	defer p.activeMutex.Unlock()
	return p.active
}

// Flush flushes the probe's input channel buffer.
func (p *Probe) Flush() {
	p.signalMutex.Lock()
	defer p.signalMutex.Unlock()

	length := len(p.C)
	for i := 0; i < length; i++ {
		p.signal = append(p.signal, <-p.C)
	}

	if len(p.signal) > p.MaximumSignalLength {
		p.signal = p.signal[len(p.signal)-p.MaximumSignalLength:]
	}
}

// Push appends a value to the probe's signal after, optionally, flushing any
// values in the probe's input channel buffer.
func (p *Probe) Push(value float64, flush bool) {
	if !p.IsActive() {
		return
	}

	p.signalMutex.Lock()
	defer p.signalMutex.Unlock()

	if flush {
		length := len(p.C)
		for i := 0; i < length; i++ {
			p.signal = append(p.signal, <-p.C)
		}
	}

	p.signal = append(p.signal, value)

	if len(p.signal) > p.MaximumSignalLength {
		p.signal = p.signal[len(p.signal)-p.MaximumSignalLength:]
	}
}

// ClearSignal removes all elements from the probe's signal.
func (p *Probe) ClearSignal() {
	p.signalMutex.Lock()
	p.signal = nil
	p.signalMutex.Unlock()
}

// Signal returns the probe's accumulated signal.
func (p *Probe) Signal() []float64 {
	p.signalMutex.Lock()
	defer p.signalMutex.Unlock()
	return p.signal
}

// RecentValue retrieves the most recent value collected by the probe, or 0.0 if
// a value has not been collected.
func (p *Probe) RecentValue() float64 {
	p.signalMutex.Lock()
	defer p.signalMutex.Unlock()

	if len(p.signal) > 0 {
		return p.signal[len(p.signal)-1]
	}

	return 0.0
}
