package main

import (
	"os"

	"github.com/starter-go/bucket-drivers/aliyun"
	"github.com/starter-go/starter"
)

func main() {

	a := os.Args
	m := aliyun.ModuleForTest()
	i := starter.Init(a)

	i.MainModule(m)

	i.WithPanic(true).Run()
}
