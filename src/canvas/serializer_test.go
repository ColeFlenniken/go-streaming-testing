package canvas_test

import (
	"testing"

	"github.com/go-streaming-testing/src/canvas"
	"github.com/google/go-cmp/cmp"
)

func TestSimpleSerialize(t *testing.T) {
	var cv canvas.Canvas = canvas.Canvas{
		Width:  2,
		Height: 2,
		Pixels: []byte{1, 2, 3, 4},
	}
	ser := canvas.Serialize(cv)

	deser := canvas.Deserialize(ser)
	if !cmp.Equal(deser, cv) {
		t.FailNow()
	}

}

func TestSerialize(t *testing.T) {
	var cv canvas.Canvas = canvas.Canvas{
		Width:  4,
		Height: 4,
		Pixels: []byte{
			1, 2, 3, 4,
			5, 6, 7, 3,
			6, 0, 1, 2,
			3, 4, 5, 6,
		},
	}
	ser := canvas.Serialize(cv)
	deser := canvas.Deserialize(ser)
	if !cmp.Equal(deser, cv) {
		t.FailNow()
	}
}
