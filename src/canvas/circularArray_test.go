package canvas_test

import (
	"fmt"
	"testing"

	"github.com/go-streaming-testing/src/canvas"
)

func TestArrayAppend(t *testing.T) {
	var array canvas.CircularArray = canvas.MakeCircularArray(4)
	array.Append(canvas.CanvasDelta{X: 1})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 2})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 3})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 4})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 5})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 6})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 7})
	if "size 4,start 3,{5 0 0},{6 0 0},{7 0 0},{4 0 0}," != array.Print() {
		t.FailNow()
	}
	//need to check internals to make sure it is properly organized
}

func TestArraySlice(t *testing.T) {
	var array canvas.CircularArray = canvas.MakeCircularArray(4)
	array.Append(canvas.CanvasDelta{X: 1})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 2})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 3})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 4})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 5})
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 6})
	array.Append(canvas.CanvasDelta{X: 7})
	array.Print()
	fmt.Println(array.GetChanges(5))
}
