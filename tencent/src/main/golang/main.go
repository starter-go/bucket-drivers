package main

import (
	"os"

	"github.com/starter-go/bucket-drivers/tencent"
	"github.com/starter-go/starter"
)

func main() {

	a := os.Args
	m := tencent.Module()
	i := starter.Init(a)

	i.MainModule(m)

	i.WithPanic(true).Run()
}
