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

package counter

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var extRegexp = regexp.MustCompile(`\.[^.]*$`)

// the Count struct contains all the total counts for each line type
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

	// Cnt is a map with language names as keys and Count objects as values
	Cnt map[string]*Count
	// DepthFirst indicates whether the file system should be visited using
	// a depthfirst algorithm (default) or breadth first
	DepthFirst bool
	// Regexp to filter all files and directories to match
	Filter *regexp.Regexp
	// Regexp to exclude all files and directories that match
	Exclude *regexp.Regexp
	// Regexp to filter all directories t0 match
	FilterDir *regexp.Regexp
	// Regexp to exclude all directories that match
	ExcludeDir *regexp.Regexp
	// Indicates if async
	Async bool
}

// create and return a pointer to a new Coutner
func NewCounter() *Counter {
	return &Counter{
		Cnt:        make(map[string]*Count),
		DepthFirst: true,
		Async:      false,
	}
}

// scans a directory using a recursive depth-first algorithm and adds calls the
// ScanFile function for each file with a valid extension
func (c *Counter) ScanDir(d string) error {

	// read all from the directory
	ds, err := ioutil.ReadDir(d)
	if err != nil {
		return err
	}

	// obtain slice and map of the child directories and files
	cDirs := make([]string, 0)
	files := make(map[string]*Language)
	var name string
	for _, dir := range ds {
		name = dir.Name()

		// filter or exclude names based on regexp
		if c.Filter != nil && !c.Filter.Match([]byte(name)) {
			continue
		}
		if c.Exclude != nil && c.Exclude.Match([]byte(name)) {
			continue
		}

		if dir.IsDir() {
			// filter or exclude directories based on regexp
			if c.FilterDir != nil && !c.FilterDir.Match([]byte(name)) {
				continue
			}
			if c.ExcludeDir != nil && c.ExcludeDir.Match([]byte(name)) {
				continue
			}
			cDirs = append(cDirs, d+"/"+name)
			continue
		}
		// filter only recogized file types
		if lang := LANGS[strings.ToLower(extRegexp.FindString(name))]; lang != nil {
			files[d+"/"+name] = lang
		}
	}

	// scan files using a depth first algorithm
	if c.DepthFirst {
		if err = c.scanAllDirs(cDirs); err != nil {
			return err
		}
		if c.Async {
			var err error
			ch := make(chan error, len(files))
			c.scanAllFilesAsync(files, ch)
			for i := 0; i < len(files); i++ {
				err = <-ch
				if err != nil {
					return err
				}
			}
		} else {
			if err = c.scanAllFiles(files); err != nil {
				return err
			}
		}
		return nil
	}

	// scan files using a breadth first algorithm
	if c.Async {
		var err error
		ch := make(chan error, len(files))
		c.scanAllFilesAsync(files, ch)
		for i := 0; i < len(files); i++ {
			err = <-ch
			if err != nil {
				return err
			}
		}
	} else {
		if err = c.scanAllFiles(files); err != nil {
			return err
		}
	}
	if err = c.scanAllDirs(cDirs); err != nil {
		return err
	}

	return nil
}

// scans a file and writes the line count to the counter;
// if an error is encountered, it is returned
func (c *Counter) ScanFile(path string, lang *Language) error {

	// open the file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// initialize fileScanner
	fs := newFileScanner(f, lang)

	// scan through each line in the file
	for fs.scanner.Scan() {
		fs.countLine()
	}
	fs.countLine()
	// check for error here?

	// add count to main counter
	c.addCount(fs.cnt, fs.lang)

	return nil
}

// calls ScanDir on each provided directory
func (c *Counter) scanAllDirs(dirs []string) error {
	for _, dir := range dirs {
		if err := c.ScanDir(dir); err != nil {
			return err
		}
	}
	return nil
}

func (c *Counter) scanAllFilesAsync(files map[string]*Language, ch chan error) {
	for path, lang := range files {
		go func(p string, l *Language) {
			if err := c.ScanFile(p, l); err != nil {
				ch <- err
				return
			}
			ch <- nil
		}(path, lang)
	}
}

// calls ScanFile on each provided file
func (c *Counter) scanAllFiles(files map[string]*Language) error {
	for path, lang := range files {
		if err := c.ScanFile(path, lang); err != nil {
			return err
		}
	}
	return nil
}

// add a count item to the total count for the specified language
func (c *Counter) addCount(cnt *Count, lang *Language) {
	if tCnt := c.Cnt[lang.Name]; tCnt != nil {
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
	c.Cnt[lang.Name] = cnt

}

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

// count a line under the appropriate count type
// takes pointers to the scanner and Count struct as parameters
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

// takes the current line as a string
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

// takes slices of the start block comment(s) in the current line, and the
// current line as a string;
// returns true if starting a block comment, false otherwise
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
