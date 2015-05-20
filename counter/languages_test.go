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
	"fmt"
	"testing"
)

func TestLangGo(t *testing.T) {
	c := NewCounter()
	c.ScanFile("./test_files/go_sample.go", LANGS[".go"])
	cnt := c.Cnt["Go"]
	fmt.Println(cnt.Code)
	if cnt.Total != 16 {
		t.Error("Go - 'Total' lines of code incorrect")
	}
	if cnt.Code != 6 {
		t.Error("Go - 'Code' lines of code incorrect")
	}
	if cnt.Mix != 2 {
		t.Error("Go - 'Mix' lines of code incorrect")
	}
	if cnt.LineCom != 1 {
		t.Error("Go - 'Line comments' lines of code incorrect")
	}
	if cnt.BlockCom != 4 {
		t.Error("Go - 'Block comments' lines of code incorrect")
	}
	if cnt.Empty != 3 {
		t.Error("Go - 'Empty' lines of code incorrect")
	}
}
