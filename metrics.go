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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
)

// Metric represents code coverage
type Metric struct {
	ID                  int64   `gorm:"primary_key:yes" json:"id"`
	Repository          string  `sql:"not null" json:"repository"`
	Sha                 string  `sql:"not null" json:"sha"`
	Branch              string  `json:"branch"`
	PackageCoverage     float64 `sql:"not null" json:"packageCoverage"`
	FilesCoverage       float64 `sql:"not null" json:"filesCoverage"`
	ClassesCoverage     float64 `sql:"not null" json:"classesCoverage"`
	MethodCoverage      float64 `sql:"not null" json:"methodCoverage"`
	LineCoverage        float64 `sql:"not null" json:"lineCoverage"`
	ConditionalCoverage float64 `sql:"not null" json:"conditionalCoverage"`
	Timestamp           int64   `sql:"not null" json:"timestamp"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// MetricsHandler represents a metrics handler
type MetricsHandler struct {
	db *gorm.DB
}

const defaultBranch = "master"

func writeError(w io.Writer, message string, err error) {
	formattedMessage := fmt.Sprintf("%s: %v", message, err)

	log.Println(formattedMessage)

	errorMsg := errorResponse{
		Error: formattedMessage,
	}

	errorString, encodingError := json.Marshal(errorMsg)
	if encodingError != nil {
		encodingErrorMessage := fmt.Sprintf("Unable to encode response message %v", encodingError)
		log.Printf(encodingErrorMessage)
	}

	w.Write(errorString)
}

func respondWithMetric(w http.ResponseWriter, m Metric) {
	bodyString, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, "unable to encode response", err)
		return
	}

	w.Write([]byte(bodyString))
}

// ExtractMetricQuery extracts a query from the request
func ExtractMetricQuery(form url.Values) Metric {
	repository := form["repository"][0]
	query := Metric{
		Repository: repository,
	}

	if len(form["sha"]) < 1 {
		if len(form["branch"]) < 1 {
			query.Branch = defaultBranch
		} else {
			query.Branch = form["branch"][0]
		}
	} else {
		query.Sha = form["sha"][0]
	}
	return query
}

func (mh MetricsHandler) handleMetricsQuery(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, "error parsing params", err)
		return
	}

	if len(r.Form["repository"]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, "missing 'repository'", errors.New("need repository"))
		return
	}

	query := ExtractMetricQuery(r.Form)

	m := new(Metric)
	mh.db.Where(&query).Order("timestamp desc").First(m)

	if m.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		writeError(w, "no rows found", errors.New("nope"))
		return
	}

	respondWithMetric(w, *m)
}

func (mh MetricsHandler) handleMetricsSave(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, "no response body", errors.New("nil body"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	m := new(Metric)
	if err := decoder.Decode(m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, "unable to decode body", err)
		return
	}

	if err := mh.RecordMetric(m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, "error recording metric", err)
	} else {
		respondWithMetric(w, *m)
	}
}

// RecordMetric saves a Metric to the database
func (mh MetricsHandler) RecordMetric(m *Metric) error {
	if m.Repository == "" || m.Sha == "" {
		return errors.New("missing required field")
	}

	if m.Timestamp == 0 {
		m.Timestamp = time.Now().Unix()
	}

	mh.db.Create(m)
	return nil
}

type handler func(w http.ResponseWriter, r *http.Request)

// NewMetricsHandler creates a new MetricsHandler
func NewMetricsHandler(db *gorm.DB) MetricsHandler {
	return MetricsHandler{
		db: db,
	}
}

// ServeHTTP handles an HTTP request for metrics
func (mh MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		mh.handleMetricsQuery(w, r)
	} else {
		mh.handleMetricsSave(w, r)
	}
	return
}
