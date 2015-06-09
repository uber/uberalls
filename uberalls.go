// Copyright (c) 2015 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"log"
	"net/http"
	"os"
)

var configLocations = []string{
	"UBERALLS_CONFIG",
	"UBERALLS_SECRETS",
}

var config Config

func init() {
	Configure()
}

func Configure() error {
	log.Println("Configuring...")
	config = Config{}

	configPaths := make([]string, 0, 2)
	for _, env := range configLocations {
		if envValue := os.Getenv(env); envValue != "" {
			configPaths = append(configPaths, envValue)
		}
	}
	return LoadConfigs(&config, configPaths)
}

// MakeServeMux instantiates an http ServeMux for the server
func MakeServeMux() *http.ServeMux {
	db, err := config.DB()
	if err != nil {
		log.Fatalf("Unable to initialize DB connection: %v", err)
	}

	if err := config.Automigrate(); err != nil {
		log.Fatalf("Could not establish database connection: %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/health", NewHealthHandler(db))
	mux.Handle("/metrics", NewMetricsHandler(db))

	return mux
}

func main() {
	mux := MakeServeMux()
	listenString := config.ConnectionString()
	log.Printf("Listening on %s... ", listenString)

	log.Fatal(http.ListenAndServe(listenString, mux))
}
