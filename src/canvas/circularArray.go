package canvas

import (
	"fmt"
	"strings"
)

// start refers to the start index while startChangeid refers to the change id of the element at the start index. Since
// changeIds increment by 1, all change ids in the list can be inferred by knowing the start
type CircularArray struct {
	Deltas        []CanvasDelta
	startChangeid int
	start         int
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
	var outputLen = arr.len - (MRChangeId - arr.startChangeid + 1)
	fmt.Printf("outlen: %v\n", outputLen)
	var outputStart int = (arr.start + (MRChangeId - arr.startChangeid + 1)) % len(arr.Deltas)
	var output []CanvasDelta = make([]CanvasDelta, outputLen)
	for i := 0; i < outputLen; i++ {
		output[i] = arr.Deltas[(outputStart+i)%len(arr.Deltas)]
	}
	return output, nil
}

// assumes no way to do a pop. Works for our use case
// end is calculatable with start and len so having it is purely for convinience
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

func (arr *CircularArray) GetLatestChangeId() int {
	return arr.start + arr.len - 1
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
