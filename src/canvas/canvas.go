package canvas

import (
	"fmt"
	"sync"
	"time"
)

// TODO merge canvas and Managed Canvas. Build perhaps using an interface. This will allow managed canvas to has a array of pixels
// with associated timestamps that helps on sync. Note this should not get serialized (unless specs change) as the times are just server side
type ManagedCanvas struct {
	Canvas    Canvas
	M         sync.Mutex
	Ts        time.Time
	Id        int
	ChangeLog CircularArray
}

type Canvas struct {
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Pixels []byte `json:"pixels"`
}

type CanvasDelta struct {
	X     uint `json:"x"`
	Y     uint `json:"y"`
	Color byte `json:"color"`
}

// need to use something other than errors as  control flow.Alsoo need to make sure to add a type to response to let client know what type of data they are receiving
func (mCanvas *ManagedCanvas) GetChanges(MRChangeId int) ([]CanvasDelta, Canvas) {
	mCanvas.M.Lock()
	defer mCanvas.M.Unlock()
	output, err := mCanvas.ChangeLog.GetChanges(MRChangeId)
	if err != nil {
		return nil, mCanvas.Canvas
	}
	return output, Canvas{}
}

func NewCanvas(height uint, width uint) (Canvas, error) {
	if height > 1<<12 || width > 1<<12 {
		return Canvas{}, fmt.Errorf("invalid creation dimensions: height and width may be 4096 at maximum")
	}
	return Canvas{Height: height, Width: width, Pixels: make([]byte, height*width)}, nil
}

func (mCanvas *ManagedCanvas) Update(deltas []CanvasDelta) error {
	mCanvas.M.Lock()
	defer mCanvas.M.Unlock()

	if deltas == nil {
		return fmt.Errorf("deltas cannot be nil")
	}
	var canvas = mCanvas.Canvas
	expectedPixels := int(canvas.Height * canvas.Width)
	if len(canvas.Pixels) != expectedPixels {
		return fmt.Errorf("invalid canvas dimensions: got %d pixels, expected %d",
			len(canvas.Pixels), expectedPixels)
	}
	var length = len(deltas)
	for i := 0; i < length; i++ {
		if deltas[i].Y >= canvas.Height {
			return fmt.Errorf("y value of delta out of bounds. Height of Canvas is %v, while the Y of the delta of index %v is %v", canvas.Height, i, deltas[i].Y)
		}
		if deltas[i].X >= canvas.Width {
			return fmt.Errorf("x value of delta out of bounds. Width of Canvas is %v, while the Y of the delta of index %v is %v", canvas.Height, i, deltas[i].X)
		}
		if deltas[i].Color >= 8 {
			return fmt.Errorf("color value must be in the range 0-7")
		}
	}
	//seperate loop is a convinience that prevents rollback logic. Can be changed later for improved efficiency assuming rate if bounds errors is low
	for i := 0; i < length; i++ {
		canvas.Pixels[canvas.Width*deltas[i].Y+deltas[i].X] = deltas[i].Color
		mCanvas.ChangeLog.Append(deltas[i])
	}
	mCanvas.Ts = time.Now()
	return nil
}

func (mCanvas *ManagedCanvas) GetCanvas() Canvas {
	mCanvas.M.Lock()
	defer mCanvas.M.Unlock()
	return mCanvas.Canvas
}
