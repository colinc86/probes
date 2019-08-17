package probes

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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
	Identifier          uuid.UUID
	C                   chan float64
	signal              []float64
	state               ProbeState
	stateMutex          sync.Mutex
}

func NewProbeFloat64(maxSignalLength int) *ProbeFloat64 {
	return &ProbeFloat64{
		MaximumSignalLength: maxSignalLength,
		Identifier:          uuid.New(),
		C:                   make(chan float64),
	}
}

func (p *ProbeFloat64) Activate(bufferSize int) chan float64 {
	p.stateMutex.Lock()
	p.state = ProbeStateActive
	p.stateMutex.Unlock()

	p.C = make(chan float64, bufferSize)
	go func() {
		for f := range p.C {
			p.signal = append(p.signal, f)
			if len(p.signal) > p.MaximumSignalLength {
				p.signal = p.signal[len(p.signal)-p.MaximumSignalLength:]
			}
		}
	}()

	return p.C
}

func (p *ProbeFloat64) Deactivate(produceImage bool) []float64 {
	close(p.C)

	p.stateMutex.Lock()
	p.state = ProbeStateInactive
	p.stateMutex.Unlock()

	if produceImage {
		plotSignal(p.signal, "Probe Input", fmt.Sprintf("Probe %s", p.Identifier), "Value", "Update", fmt.Sprintf("%s.png", p.Identifier))
	}

	return p.signal
}

func plotSignal(signal []float64, series string, title string, xAxis string, yAxis string, file string) {
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
