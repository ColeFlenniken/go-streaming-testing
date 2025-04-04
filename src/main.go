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
	Kind string `json:"kind"`
	Data []byte `json:"data"`
}

// shared state
var mCanvas canvas.ManagedCanvas

func updateData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	//for testing
	println(body)
	if err != nil {
		log.Fatal("issue reading from body")
	}
	var deser []canvas.CanvasDelta
	json.Unmarshal(body, &deser)
	mCanvas.Update(deser)
}

func getData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	//for testing
	println(body)
	if err != nil {
		log.Fatal("issue reading from body")
	}
	changeId, err := strconv.Atoi(string(body))
	if err != nil {
		log.Fatal(err)
	}
	deltas, err := mCanvas.GetChanges(changeId)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	deltas = append(deltas, canvas.CanvasDelta{X: 30, Y: 20, Color: 3}, canvas.CanvasDelta{X: 31, Y: 20, Color: 3}, canvas.CanvasDelta{X: 32, Y: 20, Color: 3}, canvas.CanvasDelta{X: 33, Y: 20, Color: 3})
	output, err := json.Marshal(deltas)
	if err != nil {
		log.Fatal(err)
	}
	println("data")
	println("writing " + string(output))
	w.Write(output)
}
func main() {
	canvasData, err := canvas.NewCanvas(1000, 1000)
	if err != nil {
		log.Fatal(err)
	}

	mCanvas = canvas.ManagedCanvas{M: sync.Mutex{}, ChangeLog: canvas.MakeCircularArray(1000), Ts: time.Now(), Id: 1, Canvas: canvasData}

	mux := http.NewServeMux()
	component := Index()
	mux.HandleFunc("/update", updateData)
	mux.HandleFunc("/getData", getData)
	mux.Handle("/", templ.Handler(component))
	mux.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))
	err = http.ListenAndServe(":3333", mux)
	if err != nil {
		log.Fatal(err)
	}
}
