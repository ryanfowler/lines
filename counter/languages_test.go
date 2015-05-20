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
