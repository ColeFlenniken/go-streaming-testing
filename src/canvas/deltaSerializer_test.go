package canvas_test

import (
	"fmt"
	"testing"

	"github.com/go-streaming-testing/src/canvas"
)

func TestSimpleDeltaSerialize(t *testing.T) {
	var dcv []canvas.CanvasDelta = make([]canvas.CanvasDelta, 3)
	dcv[0] = canvas.CanvasDelta{X: 3, Y: 92, Color: 4}
	dcv[1] = canvas.CanvasDelta{X: 55, Y: 932, Color: 7}
	dcv[2] = canvas.CanvasDelta{X: 90, Y: 83, Color: 1}
	ser := canvas.DeltaSerialize(dcv)
	fmt.Println(ser)
	deser := canvas.DeltaDeserialize(ser)
	fmt.Println(deser)

}
