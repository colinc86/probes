package probes

import (
	"fmt"
	"testing"
)

func TestCreateProbe(t *testing.T) {
	p := NewProbe(100)
	fmt.Println(p)
}
