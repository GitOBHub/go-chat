package color_test

import (
	"go-chat/color"
	"testing"
)

func TestPrintInvert(t *testing.T) {
	for i := 0; i < 100000; i++ {
		color.PrintInvert("fuck")
	}
}
