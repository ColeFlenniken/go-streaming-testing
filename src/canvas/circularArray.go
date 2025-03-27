package canvas

import (
	"fmt"
	"strings"
)

type CircularArray struct {
	Deltas        []CanvasDelta
	startChangeid int
	start         int
	end           int
	len           int
}

func MakeCircularArray(cap uint) CircularArray {
	return CircularArray{Deltas: make([]CanvasDelta, cap, cap), start: 0, startChangeid: 0, len: 0}
}

// assume we will never get a change id greater than the greatest change id in the list
func (arr *CircularArray) GetChanges(MRChangeId int) ([]CanvasDelta, error) {
	if arr.startChangeid > MRChangeId+1 {
		return []CanvasDelta{}, fmt.Errorf("last update too old. Need to get the full canvas")
	}
	var leftover = arr.len - (MRChangeId - arr.startChangeid)
	var outputStart int = (arr.start + (MRChangeId - arr.startChangeid)) % len(arr.Deltas)

	if len(arr.Deltas)-outputStart <= arr.len {
		return arr.Deltas[outputStart:arr.end], nil
	}
	var output []CanvasDelta = arr.Deltas[outputStart:]
	leftover -= len(output)
	if leftover > 0 {
		output = append(output, arr.Deltas[0:leftover]...)
	}
	return output, nil
}

// assumes no way to do a pop. Works for our use case
// end is calculatable with start and len so having it is purely for convinience
func (arr *CircularArray) Append(delta CanvasDelta) {
	if arr.len < len(arr.Deltas) {
		arr.Deltas[arr.start+arr.len] = delta
		arr.len++
		arr.end = arr.start + arr.len + 1
		return
	}
	//wraparound case
	arr.Deltas[arr.start] = delta
	arr.startChangeid++
	arr.start = (arr.start + 1) % len(arr.Deltas)
	arr.end = arr.start
}

// for testing only
func (arr *CircularArray) Print() string {
	var strBuilder strings.Builder = strings.Builder{}
	strBuilder.WriteString(fmt.Sprintf("size %v,", arr.len))
	strBuilder.WriteString(fmt.Sprintf("start %v,", arr.start))
	for i := 0; i < len(arr.Deltas); i++ {
		strBuilder.WriteString(fmt.Sprintf("%v,", arr.Deltas[(i)%len(arr.Deltas)]))
	}
	fmt.Println(strBuilder.String())
	return strBuilder.String()
}
