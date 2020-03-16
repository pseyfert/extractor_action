package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Printf("%s:\t%s\n", pair[0], pair[1])
	}
}
