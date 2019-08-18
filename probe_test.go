package probes

import (
	"fmt"
	"testing"
)

func TestCreateProbe(t *testing.T) {
	p := NewProbe()
	fmt.Println(p)
}
