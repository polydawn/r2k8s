package main

import (
	"fmt"
	"os"

	"go.polydawn.net/r2k8s/pkg/kuddle"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Please specify one arg: a k8s yaml file to munge.\n")
		os.Exit(1)
	}
	err := kuddle.InterpolateFile(
		os.Args[1],
		kuddle.FormulaLoaderForPath("."),
		"/dev/stdout",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(3)
	}
}
