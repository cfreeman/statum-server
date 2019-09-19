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
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}
var steps map[string]int = map[string]int{"A": 0, "B": 0, "C": 0, "D": 0}

func lerp(srcMin float64, srcMax float64, val float64, dstMin int, dstMax int) int {
	ratio := (math.Min(srcMax, math.Max(srcMin, val)) - srcMin) / (srcMax - srcMin)

	return int(ratio*float64(dstMax-dstMin)) + dstMin
}

func pulse(heartRate chan int, id string) {
	log.Println("Starting a fully hectic Nissan Pulsar (" + id + ").")

	pulseLength := 1000
	start := time.Now()

	for {
		select {
		case hr := <-heartRate:
			if hr > 0 {
				pulseLength = 60000 / hr
			}
		default:
		}

		if time.Now().Sub(start) > (time.Duration(pulseLength) * time.Millisecond) {
			// Broadcast the heartbeat.
			client := osc.NewClient("localhost", 53000)
			msg := osc.NewMessage("/cue/p" + id + "/start")
			log.Println(msg.Address)
			client.Send(msg)

			start = time.Now()
		}

		time.Sleep(50 * time.Millisecond) // Don't chew CPU.
	}
}

func main() {
	log.Println("Starting Statum-Server v0.0.2")

	heartRateA := make(chan int)
	go pulse(heartRateA, "A")

	heartRateB := make(chan int)
	go pulse(heartRateB, "B")

	log.Println("Creating /dat endpoint'")
	http.HandleFunc("/dat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")

		sensorID := r.URL.Query()["id"][0]

		acc, err := strconv.ParseFloat(r.URL.Query()["am"][0], 32)
		if err != nil {
			log.Println("Missing am variable")
		}

		client := osc.NewClient("localhost", 53000)
		msg := osc.NewMessage(fmt.Sprintf("/cue/g%s%d/start", sensorID, int(acc/10.0)))
		log.Println(msg.Address)
		client.Send(msg)

		rot, err := strconv.ParseFloat(r.URL.Query()["rm"][0], 32)
		if err != nil {
			log.Println("Missing rm variable")
		}

		msg = osc.NewMessage(fmt.Sprintf("/cue/r%s%d/start", sensorID, int(rot)))
		log.Println(msg.Address)
		client.Send(msg)
	})

	log.Println("Creating /step endpoint")
	http.HandleFunc("/step", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")

		sensorID := r.URL.Query()["id"][0]
		client := osc.NewClient("localhost", 53000)

		mutex.Lock()
		steps[sensorID] = steps[sensorID] + 1
		msg := osc.NewMessage(fmt.Sprintf("/cue/s%s%d/start", sensorID, steps[sensorID]))
		mutex.Unlock()

		log.Println(msg.Address)
		client.Send(msg)

		msg = osc.NewMessage(fmt.Sprintf("/cue/step%s/start", sensorID))
		log.Println(msg.Address)
		client.Send(msg)
	})

	log.Println("Creating /pulse endpoint")
	http.HandleFunc("/pulse", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")

		sensorID := r.URL.Query()["id"][0]
		bpm, err := strconv.ParseFloat(r.URL.Query()["bpm"][0], 32)
		if err != nil {
			log.Println("Missing bpm variable")
		}

		if strings.Compare(sensorID, "A") == 0 {
			heartRateA <- int(bpm)
		} else {
			heartRateB <- int(bpm)
		}

		client := osc.NewClient("localhost", 53000)
		msg := osc.NewMessage(fmt.Sprintf("/cue/bpm%s%d/start", sensorID, int(bpm)))

		log.Println(msg.Address)
		client.Send(msg)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
