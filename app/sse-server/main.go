// Copyright 2023 Michael Graff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/OpsMx/go-app-base/sse"
)

func main() {
	http.HandleFunc("/", makeStream)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func makeStream(w http.ResponseWriter, r *http.Request) {
	s := sse.NewSSE(r.Body)
	defer r.Body.Close()

	w.Header().Set("content-type", "text/event-stream")

	for i := 1; i <= 10; i++ {
		e := sse.Event{"data": fmt.Sprintf("Chunk #%d", i)}
		err := s.Write(w, e)
		if err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
