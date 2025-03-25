package canvas

import (
	"fmt"
	"strings"
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
func (arr *CircularArray) GetChanges(MRChangeId int) ([]CanvasDelta, error) {
	if arr.startChangeid > MRChangeId {
		return []CanvasDelta{}, fmt.Errorf("last update too old. Need to get the full canvas")
	}
	var outputStart int = (arr.start + (MRChangeId - arr.startChangeid)) % len(arr.Deltas)
	if len(arr.Deltas)-outputStart >= arr.len {
		return arr.Deltas[outputStart : outputStart+arr.len+1], nil
	}
	fmt.Printf("\n\n%v\n\n", len(arr.Deltas)-outputStart)
	var output []CanvasDelta = arr.Deltas[outputStart:]
	//this calc is wrong
	var leftover int = arr.len - outputStart - (len(arr.Deltas) - outputStart)
	output = append(output, arr.Deltas[0:leftover]...)
	return output, nil
}

// assumes no way to do a pop. Works for our use case
func (arr *CircularArray) Append(delta CanvasDelta) {
	if arr.len < len(arr.Deltas) {
		arr.Deltas[arr.start+arr.len] = delta
		arr.len++
		return
	}
	//wraparound case
	arr.Deltas[arr.start] = delta
	arr.startChangeid++
	arr.start = (arr.start + 1) % len(arr.Deltas)
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
