package canvas_test

import (
	"testing"

	"github.com/go-streaming-testing/src/canvas"
	"github.com/google/go-cmp/cmp"
)

func TestSimpleDeltaSerialize(t *testing.T) {
	var dcv []canvas.CanvasDelta = make([]canvas.CanvasDelta, 3)
	dcv[0] = canvas.CanvasDelta{X: 3, Y: 92, Color: 4}
	dcv[1] = canvas.CanvasDelta{X: 55, Y: 932, Color: 7}
	dcv[2] = canvas.CanvasDelta{X: 90, Y: 83, Color: 1}
	ser := canvas.DeltaSerialize(dcv)

	deser := canvas.DeltaDeserialize(ser)
	if !cmp.Equal(deser, dcv) {
		t.FailNow()
	}

}

func TestDeltaSerialize(t *testing.T) {
	var dcv []canvas.CanvasDelta = make([]canvas.CanvasDelta, 10)
	dcv[0] = canvas.CanvasDelta{X: 3, Y: 92, Color: 4}
	dcv[1] = canvas.CanvasDelta{X: 55, Y: 932, Color: 7}
	dcv[2] = canvas.CanvasDelta{X: 90, Y: 83, Color: 1}
	dcv[3] = canvas.CanvasDelta{X: 15, Y: 200, Color: 6}
	dcv[4] = canvas.CanvasDelta{X: 0, Y: 10, Color: 2}
	dcv[5] = canvas.CanvasDelta{X: 100, Y: 500, Color: 3}
	dcv[6] = canvas.CanvasDelta{X: 78, Y: 300, Color: 7}
	dcv[7] = canvas.CanvasDelta{X: 42, Y: 90, Color: 5}
	dcv[8] = canvas.CanvasDelta{X: 19, Y: 75, Color: 0}
	dcv[9] = canvas.CanvasDelta{X: 64, Y: 280, Color: 1}

	ser := canvas.DeltaSerialize(dcv)
	deser := canvas.DeltaDeserialize(ser)
	if !cmp.Equal(deser, dcv) {
		t.FailNow()
	}
}
