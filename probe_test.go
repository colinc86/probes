package probes

import (
	"math"
	"testing"
)

func TestCreateProbe(t *testing.T) {
	p := NewProbe()

	if p.MaximumSignalLength != math.MaxInt32 {
		t.Errorf("Maximum signal length, %d, should equal %d.", p.MaximumSignalLength, math.MaxInt32)
	}

	if p.InputBufferLength != 1 {
		t.Errorf("Input buffer length, %d, should equal %d.", p.InputBufferLength, 1)
	}
}

func TestProbeActivation(t *testing.T) {
	p := NewProbe()
	p.Activate()
	p.C <- 1.0

	s := p.Signal()
	if len(s) != 1 {
		t.Error("The length of s should be equal to 1.")
	}

	v := s[0]
	if v != 1.0 {
		t.Errorf("The value, %f, should be equal to 1.0.", v)
	}
}
