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
	. "github.com/uber/uberalls"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("File loading", func() {
	var c *Config

	BeforeEach(func() {
		c = &Config{}
	})

	It("Should load a default file", func() {
		LoadConfig(c, "")
		Expect(c.DBType).ToNot(BeEmpty())
	})

	It("Should error on non-existent files", func() {
		Expect(LoadConfig(c, "non-existent")).To(HaveOccurred())
		Expect(c.DBType).To(BeEmpty())
	})

	It("Should error on bad paths", func() {
		Expect(LoadConfigs(c, []string{"non-existent"})).To(HaveOccurred())
	})

	It("Should work with good paths", func() {
		Expect(LoadConfigs(c, []string{DefaultConfig})).ToNot(HaveOccurred())
	})
})

var _ = Describe("Database connections", func() {
	It("Should have a connection string", func() {
		c := Config{
			ListenAddress: "somehost",
			ListenPort:    1,
		}

		Expect(c.ConnectionString()).To(Equal("somehost:1"))
	})

	Context("With an invalid connection", func() {
		var c *Config

		BeforeEach(func() {
			c = &Config{
				DBType:     "unknown",
				DBLocation: "mars",
			}
		})

		It("Should throw an error", func() {
			conn, err := c.DB()
			Expect(conn).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("Should error trying to automigrate", func() {
			Expect(c.Automigrate()).ToNot(Succeed())
		})
	})
})
