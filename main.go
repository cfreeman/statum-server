/*
 * Copyright (c) Clinton Freeman 2019
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
 * associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute,
 * sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
 * NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
	"log"
	"math"
	"net/http"
	"strconv"
)

func lerp(srcMin float64, srcMax float64, val float64, dstMin int, dstMax int) int {

	ratio := (math.Min(srcMax, math.Max(srcMin, val)) - srcMin) / (srcMax - srcMin)

	return int(ratio*float64(dstMax-dstMin)) + dstMin
}

func pipeToOSC(r *http.Request, dimension string) {
	sensorID := r.URL.Query()["id"][0]

	am, err := strconv.ParseFloat(r.URL.Query()[dimension+"m"][0], 32)
	if err != nil {
		log.Println("Missing " + dimension + "m variable")
	}
	val := lerp(0, 100, am/100.0, 1, 100)
	client := osc.NewClient("localhost", 53000)
	msg := osc.NewMessage(fmt.Sprintf("/cue/%s%s%d/start", dimension, sensorID, val))
	client.Send(msg)

	log.Println(fmt.Sprintf("%s(%.2f)->/cue/%s%s%d/start", sensorID, am/100.0, dimension, sensorID, val))
}

func main() {
	log.Println("Starting Statum-Server v0.0.1")

	//addr := "localhost:8765"
	//server := &osc.Server{Addr: addr}

	log.Println("Creating /dat endpoint'")
	http.HandleFunc("/dat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")

		pipeToOSC(r, "a")
		pipeToOSC(r, "r")
	})

	log.Println("Creating /step endpoint")
	http.HandleFunc("/step", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")

		sensorID := r.URL.Query()["id"][0]
		steps, err := strconv.ParseFloat(r.URL.Query()["s"][0], 32)
		if err != nil {
			log.Println("Missing s variable")
		}

		client := osc.NewClient("localhost", 53000)
		msg := osc.NewMessage(fmt.Sprintf("/cue/step%s/start", sensorID))
		client.Send(msg)

		log.Println(fmt.Sprintf("%s(%d)->/cue/step%s/start", sensorID, steps, sensorID))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
