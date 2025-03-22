package canvas

import (
	"fmt"
	"time"
)

type ManagedCanvas struct {
	canvas Canvas
	ts     time.Time
	id     int
}

type Canvas struct {
	Width  uint
	Height uint
	Pixels []byte
}

type CanvasDelta struct {
	X     uint
	Y     uint
	Color byte
}

func (mCanvas *ManagedCanvas) Update(deltas []CanvasDelta) error {
	var canvas = mCanvas.canvas
	if len(canvas.Pixels) != int(canvas.Height*canvas.Width) {
		fmt.Errorf("Error: number of pixels in canvas is %v, while width*height of canvas is %v", len(canvas.Pixels), canvas.Height*canvas.Width)
	}
	var len = len(deltas)
	for i := 0; i < len; i++ {
		if deltas[i].Y < 0 || deltas[i].Y >= canvas.Height {
			return fmt.Errorf("Error: Y value of delta out of bounds. Height of Canvas is %v, while the Y of the delta of index %v is %v", canvas.Height, i, deltas[i].Y)
		}
		if deltas[i].X < 0 || deltas[i].X >= canvas.Width {
			return fmt.Errorf("Error: x value of delta out of bounds. Width of Canvas is %v, while the Y of the delta of index %v is %v", canvas.Height, i, deltas[i].X)
		}
	}
	//seperate loop is a convinience that prevents rollback logic. Can be changed later for improved efficiency assuming rate if bounds errors is low
	for i := 0; i < len; i++ {
		canvas.Pixels[canvas.Width*deltas[i].Y+deltas[i].X] = deltas[i].Color
	}
	mCanvas.ts = time.Now()
	return nil
}
