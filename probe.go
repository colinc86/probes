package probes

import (
	"fmt"
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

	// The probe's name.
	Name string

	// The probe's input channel.
	C chan float64

	signal      []float64
	signalMutex sync.Mutex
	active      bool
	activeMutex sync.Mutex
}

// MARK: Initializers

// NewProbe creates and returns a new probe.
func NewProbe(name string, maxSignalLength int) *Probe {
	return &Probe{
		MaximumSignalLength: maxSignalLength,
		Name:                name,
		C:                   make(chan float64),
	}
}

// MARK: Public methods

// Activate activates the probe and begins waiting for signal values over its
// input channel.
func (p *Probe) Activate(bufferSize int) {
	p.activeMutex.Lock()
	defer p.activeMutex.Unlock()

	if p.active {
		return
	}

	p.active = true
	p.signal = nil
	p.C = make(chan float64, bufferSize)

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
func (p *Probe) Deactivate(produceImage bool) []float64 {
	p.activeMutex.Lock()
	defer p.activeMutex.Unlock()

	if !p.active {
		return nil
	}

	p.active = false
	close(p.C)

	if produceImage {
		p.plotSignal(p.signal, "Probe Input", fmt.Sprintf("Probe %s", p.Name), "Value", "Update", fmt.Sprintf("%s.png", p.Name))
	}

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

// MARK: Private methods

func (p *Probe) plotSignal(signal []float64, series string, title string, xAxis string, yAxis string, file string) {
	pe, err := plot.New()
	if err != nil {
		panic(err)
	}

	pe.Title.Text = title
	pe.X.Label.Text = xAxis
	pe.Y.Label.Text = yAxis

	errorValues := make(plotter.XYs, len(signal))

	for i, v := range signal {
		errorValues[i].X = float64(i)
		errorValues[i].Y = v
	}

	err = plotutil.AddLines(pe,
		series, errorValues,
	)

	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := pe.Save(16*vg.Inch, 8*vg.Inch, file); err != nil {
		panic(err)
	}
}
