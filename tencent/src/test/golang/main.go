package main

import (
	"os"

	"github.com/starter-go/bucket-drivers/tencent"
	"github.com/starter-go/units"
)

func main() {

	a := os.Args
	m := tencent.ModuleForTest()

	// i := starter.Init(a)
	// i.MainModule(m)
	// i.WithPanic(true).Run()

	ctx := &units.Context{
		Arguments: a,
		Module:    m,
		UsePanic:  true,
	}

	units.Run(ctx)

}
