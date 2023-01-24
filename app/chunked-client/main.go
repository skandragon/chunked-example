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
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", "http://localhost:8090/", nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	for {
		data := make([]byte, 1024)
		l, err := res.Body.Read(data)
		if err == io.EOF {
			fmt.Printf("EOF\n")
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		s := strings.TrimSuffix(string(data[0:l]), "\n")
		fmt.Printf("Got %d bytes: %s\n", l, s)
	}
}
