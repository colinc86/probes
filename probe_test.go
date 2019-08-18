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
	
	if !p.IsActive() {
		t.Error("The probe should be active.")
	}
}

func TestSendValue(t *testing.T) {
	p := NewProbe()
	p.Activate()
	p.C <- 1.0
	p.Flush()	

	s := p.Signal()
	
	if len(s) != 1 {
		t.Errorf("The length of s, %d, should be equal to 1.", len(s))
		return
	}

	v := s[0]
	if v != 1.0 {
		t.Errorf("The value, %f, should be equal to 1.0.", v)
	}
}

func TestPushValue(t *testing.T) {
	p := NewProbe()
	p.Activate()
	p.Push(1.0, false)

	s := p.Signal()
	
	if len(s) != 1 {
		t.Errorf("The length of s, %d, should be equal to 1.", len(s))
		return
	}

	v := s[0]
	if v != 1.0 {
		t.Errorf("The value, %f, should be equal to 1.0.", v)
	}
}

func TestSendPushValue(t *testing.T) {
	p := NewProbe()
	p.InputBufferLength = 10
	p.Activate()
	
	for i := 0; i < 10; i++ {
		p.C <- float64(i)
	}

	p.Push(10.0, true)

	s := p.Signal()
	
	if len(s) != 11 {
		t.Errorf("The length of s, %d, should be equal to 11.", len(s))
		return
	}

	v := s[10]
	if v != 10.0 {
		t.Errorf("The value, %f, should be equal to 10.0.", v)
	}
}
