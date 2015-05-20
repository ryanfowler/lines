package counter

import (
	"fmt"
	"testing"
)

func BenchmarkScanDir(b *testing.B) {
	c := NewCounter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ScanDir("/Users/ryanfowler/Documents/go/src/github.com/ryanfowler")
	}
	b.StopTimer()
	fmt.Println()
	for lang, cnt := range c.cnt {
		fmt.Println(lang + ": ")
		if cnt.Code+cnt.LineCom+cnt.BlockCom+cnt.Mix+cnt.Empty != cnt.Total {
			fmt.Print("Something went wrong, counts don't equal")
		}
		fmt.Println("  Total lines:", cnt.Total)
		fmt.Println("  Code lines:", cnt.Code)
		fmt.Println("  Line comments:", cnt.LineCom)
		fmt.Println("  Block comments:", cnt.BlockCom)
		fmt.Println("  Mix lines:", cnt.Mix)
		fmt.Println("  Empty lines:", cnt.Empty)
	}
	fmt.Println()

}
