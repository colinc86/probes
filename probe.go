package probes

import (
	"fmt"
	"math"
	"sync"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Probe types reprepsent a digital probe which captures a signal and outputs
// that signal to an array and, optionally, an image for inspection.
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

// NewProbe creates and returns a new probe with the given name.
func NewProbe() *Probe {
	return &Probe{
		MaximumSignalLength: math.MaxInt32,
		InputBufferLength:   1,
		C:                   make(chan float64),
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
// Optionally, it also saves a plot of its signal.
func (p *Probe) Deactivate() []float64 {
	if !p.IsActive() {
		return nil
	}

	p.active = false
	close(p.C)

	return p.signal
}

// IsActive returns the state of the probe.
func (p *Probe) IsActive() bool {
	p.activeMutex.Lock()
	defer p.activeMutex.Unlock()
	return p.active
}

// Flush flushes the probe's input buffer if values were given to the buffer via
// its input channel.
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

// ClearSignal removes all elements from the probe's input signal.
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

// RecentValue retrieves the most recent value collected by the probe or 0.0 if
// a value has not been collected.
func (p *Probe) RecentValue() float64 {
	p.signalMutex.Lock()
	defer p.signalMutex.Unlock()

	if len(p.signal) > 0 {
		return p.signal[len(p.signal)-1]
	}

	return 0.0
}

// WriteSignalToPNG draws the probe's input signal to a PNG and saves it with
// the provided filename and appends the ".png" extension.
func (p *Probe) WriteSignalToPNG(filename string) error {
	p.signalMutex.Lock()
	defer p.signalMutex.Unlock()

	pe, err := plot.New()
	if err != nil {
		return err
	}

	pe.Title.Text = "Probe Input"
	pe.X.Label.Text = "Update"
	pe.Y.Label.Text = "Value"

	plotValues := make(plotter.XYs, len(p.signal))
	for i, v := range p.signal {
		plotValues[i].X = float64(i)
		plotValues[i].Y = v
	}

	if err := plotutil.AddLines(pe, "Probe Input", plotValues); err != nil {
		return err
	}

	// Save the plot to a PNG file.
	if err := pe.Save(16*vg.Inch, 8*vg.Inch, fmt.Sprintf("%s.png", filename)); err != nil {
		return err
	}

	return nil
}
