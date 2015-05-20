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
	"github.com/ryanfowler/lines/counter"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func main() {

	c, err := counterFromFlags()
	if err != nil {
		newLine()
		fmt.Println("Error:", err.Error())
		newLine()
		return
	}

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

	fmt.Println("Counting lines...")
	newLine()

	start := time.Now().UTC()
	err = c.ScanDir(dir)
	dur := time.Since(start)

	defer func() {
		fmt.Println("Time:", dur)
		newLine()
	}()

	if err != nil {
		newLine()
		fmt.Println("Error:", err.Error())
		return
	}

	cl := convertAllCounts(c.Cnt)
	sort.Sort(cl)

	cl.printResults()
}

func counterFromFlags() (*counter.Counter, error) {

	var breadth bool
	var filter string
	var exclude string
	var filterDir string
	var excludeDir string

	flag.BoolVar(&breadth, "breadth", false, "use a breadth-first search of directories (default: false)")
	flag.StringVar(&filter, "filter", "", "filter all file and directory names with the provided regex")
	flag.StringVar(&exclude, "exclude", "", "exclude all file and directory names with the provided regex")
	flag.StringVar(&filterDir, "filterDir", "", "filter all directory names with the provided regex")
	flag.StringVar(&excludeDir, "excludeDir", "", "exclude all directory names with the provided regex")

	flag.Parse()

	c := counter.NewCounter()
	c.DepthFirst = !breadth
	if filter != "" {
		rg, err := regexp.Compile(filter)
		if err != nil {
			return nil, err
		}
		c.Filter = rg
	}
	if exclude != "" {
		rg, err := regexp.Compile(exclude)
		if err != nil {
			return nil, err
		}
		c.Exclude = rg
	}
	if filterDir != "" {
		rg, err := regexp.Compile(filterDir)
		if err != nil {
			return nil, err
		}
		c.FilterDir = rg
	}
	if excludeDir != "" {
		rg, err := regexp.Compile(excludeDir)
		if err != nil {
			return nil, err
		}
		c.ExcludeDir = rg
	}

	return c, nil
}
