package canvas_test

import (
	"fmt"
	"testing"

	"github.com/go-streaming-testing/src/canvas"
)

func TestAdd(t *testing.T) {
	var cv canvas.Canvas = canvas.Canvas{
		Width:  2,
		Height: 2,
		Pixels: []byte{1, 2, 3, 4},
	}
	ser := canvas.Serialize(cv)
	fmt.Println(ser)
	deser := canvas.Deserialize(ser)
	fmt.Println(deser)

}
