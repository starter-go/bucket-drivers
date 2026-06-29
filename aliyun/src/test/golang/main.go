package main

import (
	"os"

	"github.com/starter-go/bucket-drivers/aliyun"
	"github.com/starter-go/units"
)

func main() {

	a := os.Args
	m := aliyun.ModuleForTest()
	c := new(units.Context)

	c.Arguments = a
	c.Module = m
	c.UsePanic = true

	units.Run(c)
}
