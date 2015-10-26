//
// lines
//
// Copyright (c) 2015 Ryan Fowler
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
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func main() {
	// get initialized Counter
	c, err := counterFromFlags()
	if err != nil {
		newLine()
		fmt.Println("Error :", err.Error())
		newLine()
		return
	}
	// get directory
	args := os.Args
	if len(args) == 1 || strings.HasPrefix(args[len(args)-1], "-") {
		newLine()
		fmt.Println("Error: No directory provided")
		newLine()
		return
	}
	dir := args[len(args)-1]
	if dir[len(dir)-1] == '/' {
		dir = dir[:len(dir)-1]
	}
	// begin counting
	fmt.Println("Counting lines...")
	newLine()
	start := time.Now().UTC()
	err = c.CountLines(dir)
	dur := time.Since(start)
	// print time when exiting
	defer func() {
		fmt.Println("Time:", dur)
		newLine()
	}()
	// if encountered, print error
	if err != nil {
		newLine()
		fmt.Println("Error:", err.Error())
		return
	}
	// sort and print results
	cl := convertAllCounts(c.cnt)
	sort.Sort(cl)
	cl.printResults()
}

// set up Counter
func counterFromFlags() (*Counter, error) {
	// flag variables
	var filter string
	var exclude string
	var filterDir string
	var excludeDir string
	// parse flags
	flag.StringVar(&filter, "filter", "", "filter all file and directory names with the provided regex")
	flag.StringVar(&exclude, "exclude", "", "exclude all file and directory names with the provided regex")
	flag.StringVar(&filterDir, "filterDir", "", "filter all directory names with the provided regex")
	flag.StringVar(&excludeDir, "excludeDir", "", "exclude all directory names with the provided regex")
	flag.Parse()
	// create and initialize Counter
	c := NewCounter()
	if filter != "" {
		rg, err := regexp.Compile(filter)
		if err != nil {
			return nil, err
		}
		c.filter = rg
	}
	if exclude != "" {
		rg, err := regexp.Compile(exclude)
		if err != nil {
			return nil, err
		}
		c.exclude = rg
	}
	if filterDir != "" {
		rg, err := regexp.Compile(filterDir)
		if err != nil {
			return nil, err
		}
		c.filterDir = rg
	}
	if excludeDir != "" {
		rg, err := regexp.Compile(excludeDir)
		if err != nil {
			return nil, err
		}
		c.excludeDir = rg
	}

	return c, nil
}
