package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/go-streaming-testing/src/canvas"
)

// aka if a full canvas or a delta array is returned
type changeData struct {
	kind string
	data []byte
}

// shared state
var mCanvas canvas.ManagedCanvas

func updateData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("issue reading from body")
	}
	var deser []canvas.CanvasDelta
	json.Unmarshal(body, &deser)
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
	deltas, canvas := mCanvas.GetChanges(changeId)
	var kind string = "deltas"
	if deltas == nil {
		kind = "canvas"
	}
	data, err := json.Marshal(canvas)
	if err != nil {
		log.Fatal(err)
	}
	fullcanvas := changeData{kind: kind, data: data}
	output, err := json.Marshal(fullcanvas)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(output)
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
