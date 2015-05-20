package counter

import (
	"bufio"
	//"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	//"time"
)

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
	// internal variables
	extReg    *regexp.Regexp
	inComment bool
	lcReg     *regexp.Regexp
	bcsReg    *regexp.Regexp
	bceReg    *regexp.Regexp
}

// create and return a pointer to a new Coutner
func NewCounter() *Counter {
	return &Counter{
		Cnt:        make(map[string]*Count),
		DepthFirst: true,
		extReg:     regexp.MustCompile(`\.[^.]*$`),
		inComment:  false,
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

	// obtain slices of the child directories and files
	cDirs := make([]string, 0)
	files := make(map[string]*Language)
	var name string
	for _, dir := range ds {
		name = dir.Name()

		if c.Filter != nil && !c.Filter.Match([]byte(name)) {
			continue
		}
		if c.Exclude != nil && c.Exclude.Match([]byte(name)) {
			continue
		}

		if dir.IsDir() {
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
		if lang := LANGS[strings.ToLower(c.extReg.FindString(name))]; lang != nil {
			files[d+"/"+name] = lang
		}
	}

	// scan files using a depth first algorithm
	if c.DepthFirst {
		err = c.scanAllDirs(cDirs)
		if err != nil {
			return err
		}
		err = c.scanAllFiles(files)
		if err != nil {
			return err
		}
		return nil
	}

	// scan files using a breadth first algorithm
	err = c.scanAllFiles(files)
	if err != nil {
		return err
	}
	err = c.scanAllDirs(cDirs)
	if err != nil {
		return err
	}

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

// calls ScanFile on each provided file
func (c *Counter) scanAllFiles(files map[string]*Language) error {
	for path, lang := range files {
		if err := c.ScanFile(path, lang); err != nil {
			return err
		}
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

	// initialize scanner + variables
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	cnt := new(Count)
	c.inComment = false
	c.lcReg = regexp.MustCompile(lang.lCom)
	c.bcsReg = regexp.MustCompile(lang.bComS)
	c.bceReg = regexp.MustCompile(lang.bComE)

	// scan through each line in the file
	for s.Scan() {
		c.countLine(s, cnt, lang)
	}
	c.countLine(s, cnt, lang)
	// check for error here?

	// add count to main counter
	c.addCount(cnt, lang)

	return nil
}

// count a line under the appropriate count type
// takes pointers to the scanner and Count struct as parameters
func (c *Counter) countLine(s *bufio.Scanner, cnt *Count, lang *Language) {

	// inc total count and remove leading/trailing whitespace
	cnt.Total += 1

	// if inside block comment
	if c.inComment {
		cnt.BlockCom += 1
		if c.stillInBlockComment(s.Text()) {
			return
		}
		c.inComment = false
		return
	}

	line := strings.TrimSpace(s.Text())

	// empty line
	if line == "" {
		cnt.Empty += 1
		return
	}

	// verify that line comment regexp exists for current language
	if lang.lCom != "" {
		if idx := c.lcReg.FindStringIndex(line); len(idx) > 0 {
			if idx[0] == 0 {
				// at the start of the line
				cnt.LineCom += 1
				return
			}
			// comment is after some code
			cnt.Mix += 1
			return
		}
	}

	// verify that block comment regexps exist for current language
	if lang.bComS != "" && lang.bComE != "" {
		idxs := c.bcsReg.FindAllStringIndex(line, -1)
		// check if starting block commment
		if c.inBlockComment(idxs, line, cnt) {
			if idxs[0][0] != 0 {
				// block comment starts after some code; mixed line
				cnt.Mix += 1
				return
			}
			// block comment at the start of the line
			cnt.BlockCom += 1
			return
		}
		if len(idxs) > 0 {
			return
		}
	}

	// is a line of code; repeat
	cnt.Code += 1
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

// takes slices of the end block comment(s) and the start block comment(s) in
// the current line; returns true if still in block comment, false otherwise
func (c *Counter) stillInBlockComment(line string) bool {
	idxe := c.bceReg.FindAllStringIndex(line, -1)
	if len(idxe) == 0 {
		return true
	}
	idxs := c.bcsReg.FindAllStringIndex(line, -1)
	if len(idxs) == 0 {
		return false
	}
	return idxe[len(idxe)-1][0] < idxs[len(idxs)-1][0]
}

// takes slices of the end block comment(s) and the start block comment(s) in
// the current line; returns true if starting a block comment, false otherwise
func (c *Counter) inBlockComment(idxs [][]int, line string, cnt *Count) bool {
	if len(idxs) == 0 {
		return false
	}
	idxe := c.bceReg.FindAllStringIndex(line, -1)
	if len(idxe) == 0 {
		c.inComment = true
		return true
	}
	if idxe[len(idxe)-1][0] < idxs[len(idxs)-1][0] {
		c.inComment = true
		return true
	}
	// block comment ends on same line, check if mixed line
	if idxs[0][0] != 0 {
		cnt.Mix += 1
		return false
	}
	cnt.BlockCom += 1
	return false
}
