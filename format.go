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
	"bytes"
	"fmt"
	"github.com/ryanfowler/lines/counter"
	"math"
	"strconv"
)

func addCommasToItem(orig string, buf *bytes.Buffer) string {
	if len(orig) <= 3 {
		return orig
	}
	defer buf.Reset()
	i := int(math.Floor(math.Mod(float64(len(orig)), 3) + 0.5))
	if i == 0 {
		i = 3
	}
	prev := 0
	for ; i < len(orig); i += 3 {
		buf.WriteString(orig[prev:i])
		buf.WriteString(",")
		prev = i
	}
	buf.WriteString(orig[prev:])
	return buf.String()
}

func printCount(cnt count, w widths, allLeftAligned bool) {
	if allLeftAligned {
		printVal(cnt.lang, w.lang, true)
		printVal(cnt.files, w.files, true)
		printVal(cnt.code, w.code, true)
		printVal(cnt.mix, w.mix, true)
		printVal(cnt.com, w.com, true)
		printVal(cnt.empty, w.empty, true)
		printVal(cnt.total, w.total, true)
		printDivider()
		return
	}
	printVal(cnt.lang, w.lang, true)
	printVal(cnt.files, w.files, false)
	printVal(cnt.code, w.code, false)
	printVal(cnt.mix, w.mix, false)
	printVal(cnt.com, w.com, false)
	printVal(cnt.empty, w.empty, false)
	printVal(cnt.total, w.total, false)
	printDivider()
}

func printVal(val string, l int, leftAligned bool) {
	printDivider()
	printSpace()
	if leftAligned {
		fmt.Print(val)
		for i := 0; i < l-len(val); i++ {
			printSpace()
		}
	} else {
		for i := 0; i < l-len(val); i++ {
			printSpace()
		}
		fmt.Print(val)
	}
	printSpace()
}

func printDivider() {
	fmt.Print("|")
}

func newLine() {
	fmt.Print("\n")
}

func printDash() {
	fmt.Print("-")
}

func printSpace() {
	fmt.Print(" ")
}

type count struct {
	lang  string
	nums  *counter.Count
	total string
	code  string
	com   string
	mix   string
	empty string
	files string
}

func convertCount(name string, cnt *counter.Count) count {
	return count{
		lang:  name,
		nums:  cnt,
		total: strconv.FormatInt(cnt.Total, 10),
		code:  strconv.FormatInt(cnt.Code, 10),
		com:   strconv.FormatInt(cnt.BlockCom+cnt.LineCom, 10),
		mix:   strconv.FormatInt(cnt.Mix, 10),
		empty: strconv.FormatInt(cnt.Empty, 10),
		files: strconv.FormatInt(cnt.Files, 10),
	}
}

type countList []count

func convertAllCounts(cnts map[string]*counter.Count) countList {
	cl := make(countList, 0, len(cnts))
	for name, cnt := range cnts {
		cl = append(cl, convertCount(name, cnt))
	}
	return cl
}

func (cl countList) printResults() {
	if len(cl) == 0 {
		fmt.Println("No source files found")
		return
	}
	if len(cl) > 1 {
		cl = cl.addTotal()
	}
	cl.addCommas()
	w := cl.tableWidths()
	wSum := w.sum()
	for i := 0; i < wSum; i++ {
		printDash()
	}
	newLine()
	printCount(count{
		lang:  "Language",
		files: "Files",
		code:  "Code",
		mix:   "Mixed",
		com:   "Comments",
		empty: "Empty",
		total: "Total",
	}, w, true)
	newLine()
	for i := 0; i < wSum; i++ {
		printDash()
	}
	for _, cnt := range cl {
		newLine()
		if cnt.lang == "Totals:" {
			for i := 0; i < wSum; i++ {
				printDash()
			}
			newLine()
		}
		printCount(cnt, w, false)
	}
	newLine()
	for i := 0; i < wSum; i++ {
		printDash()
	}
	newLine()
}

func (cl countList) addTotal() countList {
	var files int64
	var code int64
	var mix int64
	var com int64
	var empty int64
	var total int64
	for _, val := range cl {
		files += val.nums.Files
		code += val.nums.Code
		mix += val.nums.Mix
		com += val.nums.BlockCom + val.nums.LineCom
		empty += val.nums.Empty
		total += val.nums.Total
	}
	return append(cl, count{
		lang:  "Totals:",
		files: strconv.FormatInt(files, 10),
		code:  strconv.FormatInt(code, 10),
		mix:   strconv.FormatInt(mix, 10),
		com:   strconv.FormatInt(com, 10),
		empty: strconv.FormatInt(empty, 10),
		total: strconv.FormatInt(total, 10),
	})
}

func (cl countList) addCommas() {
	buf := new(bytes.Buffer)
	for idx, val := range cl {
		cl[idx].files = addCommasToItem(val.files, buf)
		cl[idx].code = addCommasToItem(val.code, buf)
		cl[idx].mix = addCommasToItem(val.mix, buf)
		cl[idx].com = addCommasToItem(val.com, buf)
		cl[idx].empty = addCommasToItem(val.empty, buf)
		cl[idx].total = addCommasToItem(val.total, buf)
	}
}

func (cl countList) Len() int {
	return len(cl)
}

func (cl countList) Less(i, j int) bool {
	return cl[i].nums.Total > cl[j].nums.Total
}

func (cl countList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

// the type widths holds the integer value of the longest string length for
// each column
type widths struct {
	lang  int
	total int
	code  int
	com   int
	mix   int
	empty int
	files int
}

// gets longest lengths for each column
// len(string) finds the number of bytes, but all characters used here should
// be 1 byte
func (cl countList) tableWidths() widths {
	w := widths{
		lang:  len("Language"),
		total: len("Total"),
		code:  len("Code"),
		com:   len("Comments"),
		mix:   len("Mixed"),
		empty: len("Empty"),
		files: len("Files"),
	}
	for _, cnt := range cl {
		if len(cnt.lang) > w.lang {
			w.lang = len(cnt.lang)
		}
		if len(cnt.total) > w.total {
			w.total = len(cnt.total)
		}
		if len(cnt.code) > w.code {
			w.code = len(cnt.code)
		}
		if len(cnt.com) > w.com {
			w.com = len(cnt.com)
		}
		if len(cnt.mix) > w.mix {
			w.mix = len(cnt.mix)
		}
		if len(cnt.empty) > w.empty {
			w.empty = len(cnt.empty)
		}
		if len(cnt.files) > w.files {
			w.files = len(cnt.files)
		}
	}
	return w
}

// returns the entire table width based on each column width
func (w widths) sum() int {
	return 3*7 + 1 + (w.lang + w.total + w.code + w.com + w.mix + w.empty + w.files)
}
