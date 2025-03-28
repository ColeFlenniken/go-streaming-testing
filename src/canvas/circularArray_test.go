package canvas_test

import (
	"fmt"
	"testing"

	"github.com/go-streaming-testing/src/canvas"
	"github.com/google/go-cmp/cmp"
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
	array.Append(canvas.CanvasDelta{X: 2})
	changes, err := array.GetChanges(0)
	if err != nil {
		fmt.Printf("fail on error from array.getChanges(0)")
		t.FailNow()
	}
	if !cmp.Equal(changes, []canvas.CanvasDelta{canvas.CanvasDelta{X: 2}}) {
		fmt.Printf("fail on wrong output from array.getChanges(0)")
		t.FailNow()
	}
	array.Append(canvas.CanvasDelta{X: 3})
	array.Append(canvas.CanvasDelta{X: 4})
	array.Append(canvas.CanvasDelta{X: 5})
	changes, err = array.GetChanges(0)
	if err != nil {
		t.FailNow()
	}
	if !cmp.Equal(changes, []canvas.CanvasDelta{{X: 2}, {X: 3}, {X: 4}, {X: 5}}) {
		fmt.Print(changes)
		t.FailNow()
	}
	//array.Print()
	array.Append(canvas.CanvasDelta{X: 6})
	array.Append(canvas.CanvasDelta{X: 7})
	changes, err = array.GetChanges(5)
	if err != nil {
		t.FailNow()
	}
	if !cmp.Equal(changes, []canvas.CanvasDelta{canvas.CanvasDelta{X: 7}}) {
		t.FailNow()
	}
}
