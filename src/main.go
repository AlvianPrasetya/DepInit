package main

import (
	"fmt"

	"github.com/AlvianPrasetya/DepInit/src/depinit"
)

func main() {
	dm := depinit.NewDepMngr(false)

	// a -> b -> c -> d
	addModule(dm, "a")
	addModule(dm, "b", "a")
	addModule(dm, "c", "a", "b")
	addModule(dm, "d", "a", "b", "c")

	dm.Init()
}

func addModule(dm *depinit.DepMngr, name string, deps ...string) {
	dm.AddModule(name, func() error {
		fmt.Printf("init %s\n", name)
		return nil
	}, deps...)
}
