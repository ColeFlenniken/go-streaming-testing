package canvas

import (
	"fmt"
)

type CircularArray struct {
	Deltas        []CanvasDelta
	startChangeid int
	start         int
	len           int
}

func MakeCircularArray(cap uint) CircularArray {
	return CircularArray{Deltas: make([]CanvasDelta, cap), start: 0, startChangeid: 0, len: 0}
}

// assume we will never get a change id greater than the greatest change id in the list
func (arr *CircularArray) getChanges(MRChangeId int) ([]CanvasDelta, error) {
	if arr.startChangeid > MRChangeId {
		return []CanvasDelta{}, fmt.Errorf("last update too old. Need to get the full canvas")
	}
	var outputStart int = arr.start + (MRChangeId - arr.startChangeid)
	if len(arr.Deltas)-outputStart >= arr.len {
		return arr.Deltas[outputStart : outputStart+arr.len+1], nil
	}
	var output []CanvasDelta = arr.Deltas[outputStart:]
	var leftover int = arr.len - (len(arr.Deltas) - outputStart)
	output = append(output, arr.Deltas[0:leftover]...)
	return output, nil
}

// assumes no way to do a pop. Works for our use case
func (arr *CircularArray) append(delta CanvasDelta) {
	if arr.len < len(arr.Deltas) {
		arr.Deltas[arr.start+arr.len] = delta
		arr.len++
	}
	//wraparound case
	arr.Deltas[arr.start] = delta
	arr.startChangeid++
	arr.start = (arr.start + 1) % len(arr.Deltas)
}
