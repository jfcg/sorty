package sorty

import (
	"os"
	"testing"

	"github.com/jfcg/sixb/v2"
)

func BenchmarkSortB(b *testing.B) {
	b.StopTimer()
	env, arg := os.Environ(), os.Args
	slc := make([][]byte, 16*(len(arg)+len(env)))

	for q := 0; q < b.N; q++ {
		// fill slc
		for i, r := 16, 0; i > 0; i-- {
			for k := len(arg) - 1; k >= 0; k-- {
				slc[r] = sixb.Bytes(arg[k])
				r++
			}
			for k := len(env) - 1; k >= 0; k-- {
				slc[r] = sixb.Bytes(env[k])
				r++
			}
		}
		b.StartTimer()
		sortB(slc)
		b.StopTimer()
	}
	if isSortedB(slc) != 0 {
		b.Fatal("sortB error")
	}
}
