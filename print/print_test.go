package print_test

import (
	"server/chat/print"
	"testing"
)

func TestPrintInvert(t *testing.T) {
	for i := 0; i < 100000; i++ {
		color.PrintInvert("fuck")
	}
}
