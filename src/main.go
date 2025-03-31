package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/go-streaming-testing/src/canvas"
)

//shared state
var mCanvas canvas.ManagedCanvas

func updateData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("issue reading from body")
	}
	var deser []canvas.CanvasDelta = canvas.DeltaDeserialize(body)
	mCanvas.Update(deser)
}

func getData(w http.ResponseWriter, r *http.Request){
	body, err := io.ReadAll(r.)
	if err != nil {
		log.Fatal("issue reading from body")
	}
	
}
func main() {
	canvasData, err := canvas.NewCanvas(500, 500)
	if err != nil {
		log.Fatal(err)
	}
	
	mCanvas = canvas.ManagedCanvas{M : sync.Mutex{}, ChangeLog: canvas.MakeCircularArray(1000), Ts: time.Now(), Id: 1, Canvas: canvasData}

	mux := http.NewServeMux()
	component := Index()
	mux.HandleFunc("/update", updateData)
	mux.Handle("/", templ.Handler(component))

	err = http.ListenAndServe(":3333", mux)
	if err != nil {
		log.Fatal(err)
	}
}
