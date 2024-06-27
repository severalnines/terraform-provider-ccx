package resources

import (
	"fmt"
	"os"
)

func debug(name string, val any) {
	f, err := os.Create("/home/fx/Projects/Severalnines/aardvark/CCX-4516/debug.log")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if val == nil {
		fmt.Fprintf(f, "%s is nil\n", name)
	} else {
		fmt.Fprintf(f, "%s = %T\n%+v\n", name, val, val)
	}

	os.Exit(1)
}
