package probes

import (
	"fmt"
	"testing"
)

func TestCreateProbe(t *testing.T) {
	p := NewProbeFloat64(100)
	fmt.Println(p)
}
