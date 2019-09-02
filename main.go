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
	"log"
	//"math"
	"net/http"
	"strconv"
	//"time"
)

func main() {
	log.Println("Starting Statum-Server v0.0.1")

	//addr := "localhost:8765"
	//server := &osc.Server{Addr: addr}

	log.Println("Creating /dat endpoint'")
	http.HandleFunc("/dat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")

		log.Println("DAT!");

		ax, err := strconv.ParseFloat(r.URL.Query()["ax"][0], 32)
		if err != nil {
			log.Println("Missing ax variable");
		}

		log.Println("DAT: %f", ax/100.0)



		// if err != nil {
		// 	log.Fatal("Unable to parse '/l' argument.")
		// }

		// id := lerp(0, 100, f, 1, 100)
		// client := osc.NewClient("localhost", 53000)
		// msg := osc.NewMessage(fmt.Sprintf("/cue/l%d/start", id))
		// client.Send(msg)
		// log.Println(fmt.Sprintf("%s (%.2f)", msg.Address, f))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

