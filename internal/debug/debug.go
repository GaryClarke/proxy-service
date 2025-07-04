package debug

import (
	"os"

	"github.com/davecgh/go-spew/spew"
)

// DD dumps whatever you pass in to stdout and then panics to stop execution.
func DD(vals ...interface{}) {
	spew.Fdump(os.Stdout, vals...)
	panic("💥 debug dump – halting request")
}
