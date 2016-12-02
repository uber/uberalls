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

package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/jinzhu/gorm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/uber/uberalls"
)

func getMetricsResponse(method string, body *strings.Reader, params string, db *gorm.DB) *httptest.ResponseRecorder {
	var request *http.Request
	url := "/metrics"
	if params != "" {
		url = fmt.Sprintf("%s?%s", url, params)
	}
	if body == nil {
		request, _ = http.NewRequest(method, url, nil)
	} else {
		request, _ = http.NewRequest(method, url, body)
	}
	response := httptest.NewRecorder()
	handler := NewMetricsHandler(db)
	handler.ServeHTTP(response, request)
	return response
}

var _ = Describe("/metrics handler", func() {
	var (
		c  *Config
		db *gorm.DB
	)

	BeforeEach(func() {
		c = &Config{
			DBType:     "sqlite3",
			DBLocation: "test.sqlite",
		}
		db, _ = c.DB()
	})

	var response *httptest.ResponseRecorder
	It("Should automigrate", func() {
		Expect(c.Automigrate()).To(Succeed())
	})

	It("Should generate a 404 for non-existent repos", func() {
		response = getMetricsResponse("GET", nil, "repository=foo", db)
		Expect(response.Code).To(Equal(http.StatusNotFound))
	})

	Context("With an empty POST body", func() {
		response = getMetricsResponse("POST", nil, "", db)

		It("Should be non-OK", func() {
			Expect(response.Code).ToNot(Equal(http.StatusOK))
		})
	})

	Context("With valid JSON", func() {
		dummyJSON := `{
				"repository": "test", "packageCoverage": 38, "filesCoverage": 39,
        "classesCoverage": 40, "methodCoverage": 41, "lineCoverage": 42,
				"conditionalCoverage": 43, "sha": "deadbeef"}`

		BeforeEach(func() {
			response = getMetricsResponse("POST", strings.NewReader(dummyJSON), "", db)
		})

		It("Should be HTTP OK", func() {
			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("Should decode a metric", func() {
			metric := new(Metric)
			decoder := json.NewDecoder(response.Body)
			Expect(decoder.Decode(metric)).To(Succeed())
			Expect(metric.ID).To(BeNumerically(">", 0))
		})

		Context("Retrieving the metric", func() {
			var metric *Metric

			BeforeEach(func() {
				response = getMetricsResponse("GET", nil, "repository=test&sha=deadbeef", db)
			})

			It("Should be HTTP OK", func() {
				Expect(response.Code).To(Equal(http.StatusOK))
			})

			It("Should decode the metric", func() {
				metric = new(Metric)
				decoder := json.NewDecoder(response.Body)
				Expect(decoder.Decode(metric)).To(Succeed())
			})

			It("Should have the correct values", func() {
				Expect(metric.PackageCoverage).To(Equal(38.))
				Expect(metric.FilesCoverage).To(Equal(39.))
				Expect(metric.ClassesCoverage).To(Equal(40.))
				Expect(metric.MethodCoverage).To(Equal(41.))
				Expect(metric.LineCoverage).To(Equal(42.))
				Expect(metric.ConditionalCoverage).To(Equal(43.))
			})
		})
	})

	Context("With invalid JSON", func() {
		badJSON := `{}`

		BeforeEach(func() {
			response = getMetricsResponse("POST", strings.NewReader(badJSON), "", db)
		})

		It("Should not be OK posting an empty JSON body", func() {
			Expect(response.Code).ToNot(Equal(http.StatusOK))
		})
	})

	Context("With bad configuration", func() {
		badConfig := Config{
			DBType: "unknown",
		}

		BeforeEach(func() {
			db, _ := badConfig.DB()
			response = getMetricsResponse("GET", nil, "", db)
		})

		It("Should not be OK", func() {
			Expect(response.Code).ToNot(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("ExtractMetricsQuery", func() {
	Context("When branch is specified in the query", func() {
		values := url.Values{
			"repository": []string{"foo"},
			"branch":     []string{"master"},
		}
		query := ExtractMetricQuery(values)

		It("Should extract the master branch", func() {
			Expect(query.Branch).To(Equal("master"))
		})
	})

	Context("When no branch is specified", func() {
		values := url.Values{
			"repository": []string{"foo"},
		}
		query := ExtractMetricQuery(values)

		It("Should extract the default branch", func() {
			Expect(query.Branch).To(Equal("origin/master"))
		})
	})
})
