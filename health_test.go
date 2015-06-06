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
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/uber/uberalls"
)

var _ = Describe("Health handler", func() {
	var response *httptest.ResponseRecorder
	var c Config

	JustBeforeEach(func() {
		request, _ := http.NewRequest("GET", "/health", nil)
		response = httptest.NewRecorder()

		db, err := c.DB()
		Expect(err).ToNot(HaveOccurred())

		handler := NewHealthHandler(db)
		handler.ServeHTTP(response, request)
	})

	Context("With a valid DB connection", func() {
		BeforeEach(func() {
			c = Config{
				DBType:     "sqlite3",
				DBLocation: "test.sqlite",
			}
		})

		It("Should be HTTP OK", func() {
			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})

	Context("With an invalid DB connection", func() {
		BeforeEach(func() {
			c = Config{
				DBType:     "mysql",
				DBLocation: ":-1",
			}
		})

		It("should not be OK", func() {
			Expect(response.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})
