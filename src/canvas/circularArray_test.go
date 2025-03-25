package canvas_test

import (
	"testing"

	"github.com/go-streaming-testing/src/canvas"
)

func TestArrayAppend(t *testing.T) {
	var array canvas.CircularArray = canvas.MakeCircularArray(4)
	array.Append(canvas.CanvasDelta{X: 1})

	array.Append(canvas.CanvasDelta{X: 2})

	array.Append(canvas.CanvasDelta{X: 3})

	array.Append(canvas.CanvasDelta{X: 4})

	array.Append(canvas.CanvasDelta{X: 5})

	array.Append(canvas.CanvasDelta{X: 6})

	array.Append(canvas.CanvasDelta{X: 7})
	//need to check internals to make sure it is properly organized
}
