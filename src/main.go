package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/go-streaming-testing/src/canvas"
)

// shared state
var mCanvas canvas.ManagedCanvas

func updateData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("issue reading from body")
	}
	var deser []canvas.CanvasDelta = canvas.DeltaDeserialize(body)
	mCanvas.Update(deser)
}

func getData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("issue reading from body")
	}
	changeId, err := strconv.Atoi(string(body))
	if err != nil {
		log.Fatal(err)
	}
	deltas, err := mCanvas.GetChanges(changeId)
	if err != nil {
		canvasFull := canvas.Serialize(mCanvas.GetCanvas())
		w.Write(append([]byte{byte(0)}, canvasFull...))
	}
	canvasDeltas := canvas.DeltaSerialize(deltas)
	w.Write(append([]byte{1}, canvasDeltas...))
}
func main() {
	canvasData, err := canvas.NewCanvas(500, 500)
	if err != nil {
		log.Fatal(err)
	}

	mCanvas = canvas.ManagedCanvas{M: sync.Mutex{}, ChangeLog: canvas.MakeCircularArray(1000), Ts: time.Now(), Id: 1, Canvas: canvasData}

	mux := http.NewServeMux()
	component := Index()
	mux.HandleFunc("/update", updateData)
	mux.HandleFunc("/getData", getData)
	mux.Handle("/", templ.Handler(component))

	err = http.ListenAndServe(":3333", mux)
	if err != nil {
		log.Fatal(err)
	}
}
