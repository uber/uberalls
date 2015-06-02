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
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const defaultConfig = "config/default.json"

// Config holds application configuration
type Config struct {
	DBType        string
	DBLocation    string
	ListenPort    int
	ListenAddress string
	db            *gorm.DB
}

// ConnectionString returns a TCP string for the HTTP server to bind to
func (c Config) ConnectionString() string {
	return fmt.Sprintf("%s:%d", c.ListenAddress, c.ListenPort)
}

// DB returns a database connection based on configuration
func (c *Config) DB() (*gorm.DB, error) {
	if c.db == nil {
		newdb, err := gorm.Open(c.DBType, c.DBLocation)
		if err != nil {
			return nil, err
		}
		c.db = &newdb
	}
	return c.db, nil
}

// Automigrate runs migrations automatically
func (c Config) Automigrate() error {
	db, err := c.DB()
	if err != nil {
		return err
	}
	model := new(Metric)
	db.AutoMigrate(model)

	db.Model(model).AddIndex(
		"idx_metrics_repository_sha_timestamp",
		"repository",
		"sha",
		"timestamp",
	)
	db.Model(model).AddIndex(
		"idx_metrics_repository_branch_timestamp",
		"repository",
		"branch",
		"timestamp",
	)
	return nil
}

// LoadConfig loads configuration from a file into a Config type
func LoadConfig(c *Config, configPath string) {
	if configPath == "" {
		configPath = defaultConfig
	}

	log.Print("Loading configuration from '", configPath, "'")
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println(err)
		return
	}

	if err = json.Unmarshal(content, c); err != nil {
		log.Fatal(err)
	}

	return
}
