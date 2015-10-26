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
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// file extension regexp
var feRegexp = regexp.MustCompile(`\.[^.]*$`)

// The Count struct contains all the line type totals for a single file
type Count struct {
	Total    int64
	Code     int64
	BlockCom int64
	LineCom  int64
	Mix      int64
	Empty    int64
	Files    int64
}

// the Counter struct contains the count map and the regexp to find file
// extensions
type Counter struct {
	// muCnt is the mutex that protects the cnt map
	muCnt sync.Mutex
	// Cnt is a map with language names as keys and Count objects as values
	cnt map[string]*Count
	// Regexp to filter all files and directories to match
	filter *regexp.Regexp
	// Regexp to exclude all files and directories that match
	exclude *regexp.Regexp
	// Regexp to filter all directories t0 match
	filterDir *regexp.Regexp
	// Regexp to exclude all directories that match
	excludeDir *regexp.Regexp
	// chScanFile is the channel that receives files to scan
	chScanFile chan string
	// chDone is the channel that is closed when no more files are left to scan
	chDone chan struct{}
	// wg is the waitgroup for all of the scanFileWorkers
	wg sync.WaitGroup
}

// create and return a pointer to a new Coutner
func NewCounter() *Counter {
	c := &Counter{
		muCnt:      sync.Mutex{},
		cnt:        make(map[string]*Count),
		chScanFile: make(chan string),
		chDone:     make(chan struct{}),
		wg:         sync.WaitGroup{},
	}
	c.startScanFileWorkers()
	return c
}

// Count the lines of each supported file located in the provided directory
func (c *Counter) CountLines(p string) (err error) {
	err = c.scanDir(p)
	close(c.chDone)
	c.wg.Wait()
	return
}

// scandir scans a directory sending file paths on chScanFile
func (c *Counter) scanDir(d string) error {
	// read all from the directory
	ds, err := ioutil.ReadDir(d)
	if err != nil {
		return err
	}

	// obtain slice and map of the child directories and files
	cDirs := make([]string, 0)
	var name string
	for _, dir := range ds {
		name = dir.Name()
		// filter or exclude names based on regexp
		if c.filter != nil && !c.filter.Match([]byte(name)) {
			continue
		}
		if c.exclude != nil && c.exclude.Match([]byte(name)) {
			continue
		}
		if dir.IsDir() {
			// filter or exclude directories based on regexp
			if c.filterDir != nil && !c.filterDir.Match([]byte(name)) {
				continue
			}
			if c.excludeDir != nil && c.excludeDir.Match([]byte(name)) {
				continue
			}
			cDirs = append(cDirs, filepath.Join(d, name))
			continue
		}
		// send on chScanFile
		c.wg.Add(1)
		c.chScanFile <- filepath.Join(d, name)
	}

	if err = c.scanAllDirs(cDirs); err != nil {
		return err
	}

	return nil
}

// start two scanFileWorkers for each cpu
func (c *Counter) startScanFileWorkers() {
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go c.scanFileWorker()
	}
}

// scans a file and writes the line count to the counter;
// if an error is encountered, it is returned
func (c *Counter) scanFileWorker() {
	var path string
	for {
		select {
		case <-c.chDone:
			return
		case path = <-c.chScanFile:
			// filter only recogized file types
			lang := LANGS[strings.ToLower(feRegexp.FindString(path))]
			if lang == nil {
				c.wg.Done()
				continue
			}
			// open the file
			f, err := os.Open(path)
			if err != nil {
				c.wg.Done()
				continue
			}

			fs := newFileScanner(f, lang)
			for fs.scanner.Scan() {
				fs.countLine()
			}
			fs.countLine()
			// check for error here?

			c.addCount(fs.cnt, fs.lang)
			f.Close()
			c.wg.Done()
		}
	}
}

// calls ScanDir on each provided directory
func (c *Counter) scanAllDirs(dirs []string) error {
	for _, dir := range dirs {
		if err := c.scanDir(dir); err != nil {
			return err
		}
	}
	return nil
}

// add a count item to the total count for the specified language
func (c *Counter) addCount(cnt *Count, lang *Language) {
	c.muCnt.Lock()
	defer c.muCnt.Unlock()
	if tCnt := c.cnt[lang.Name]; tCnt != nil {
		tCnt.Total += cnt.Total
		tCnt.Empty += cnt.Empty
		tCnt.Code += cnt.Code
		tCnt.LineCom += cnt.LineCom
		tCnt.BlockCom += cnt.BlockCom
		tCnt.Mix += cnt.Mix
		tCnt.Files += 1
		return
	}
	cnt.Files = 1
	c.cnt[lang.Name] = cnt

}

// fileScanner scans a file and counts each line of code
type fileScanner struct {
	scanner   *bufio.Scanner
	cnt       *Count
	lang      *Language
	inComment bool
	curLine   string
	lineCom   *regexp.Regexp
	blockComS *regexp.Regexp
	blockComE *regexp.Regexp
}

// create a new fileScanner struct
func newFileScanner(f *os.File, l *Language) *fileScanner {
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	return &fileScanner{
		scanner:   s,
		cnt:       new(Count),
		lang:      l,
		inComment: false,
		lineCom:   regexp.MustCompile(l.lCom),
		blockComS: regexp.MustCompile(l.bComS),
		blockComE: regexp.MustCompile(l.bComE),
	}
}

// Count a line under the appropriate type.
func (fs *fileScanner) countLine() {

	// inc total count
	fs.cnt.Total += 1
	// assign current line with leading and trailing spaces removed
	fs.curLine = strings.TrimSpace(fs.scanner.Text())

	// if inside a block comment
	if fs.inComment {
		fs.cnt.BlockCom += 1
		if fs.stillInBlockComment() {
			return
		}
		fs.inComment = false
		return
	}

	// check if empty line
	if fs.curLine == "" {
		fs.cnt.Empty += 1
		return
	}

	// if line comment regexp exists, check if in the current line
	if fs.lang.lCom != "" {
		if idx := fs.lineCom.FindStringIndex(fs.curLine); len(idx) > 0 {
			if idx[0] == 0 {
				// comment is at the start of the line
				fs.cnt.LineCom += 1
				return
			}
			// comment is after some code
			fs.cnt.Mix += 1
			return
		}
	}

	// if block comment regexp exists, check if in the current line
	if fs.lang.bComS != "" && fs.lang.bComE != "" {
		idxs := fs.blockComS.FindAllStringIndex(fs.curLine, -1)
		// check if starting block commment
		if fs.inBlockComment(idxs) {
			if idxs[0][0] != 0 {
				// block comment starts after some code; mixed line
				fs.cnt.Mix += 1
				return
			}
			// block comment at the start of the line
			fs.cnt.BlockCom += 1
			return
		}
		// check if block comment started and ended on the same line
		if len(idxs) > 0 {
			return
		}
	}

	// is a line of code; repeat
	fs.cnt.Code += 1
}

// returns true if still in block comment, false otherwise
func (fs *fileScanner) stillInBlockComment() bool {
	idxe := fs.blockComE.FindAllStringIndex(fs.curLine, -1)
	if len(idxe) == 0 {
		return true
	}
	idxs := fs.blockComS.FindAllStringIndex(fs.curLine, -1)
	if len(idxs) == 0 {
		return false
	}
	return idxe[len(idxe)-1][0] < idxs[len(idxs)-1][0]
}

// Takes slices of the start block comment(s) in the current line as a parameter.
// Returns true if starting a block comment, false otherwise.
func (fs *fileScanner) inBlockComment(idxs [][]int) bool {
	if len(idxs) == 0 {
		return false
	}
	idxe := fs.blockComE.FindAllStringIndex(fs.curLine, -1)
	if len(idxe) == 0 {
		fs.inComment = true
		return true
	}
	if idxe[len(idxe)-1][0] < idxs[len(idxs)-1][0] {
		fs.inComment = true
		return true
	}
	// block comment ends on same line, check if mixed line
	if idxs[0][0] != 0 {
		fs.cnt.Mix += 1
		return false
	}
	// current line is a block comment, no other code
	fs.cnt.BlockCom += 1
	return false
}
