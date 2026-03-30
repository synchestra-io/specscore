package main

import (
	"os"

	cli "github.com/synchestra-io/specscore/internal/cli"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		cli.Fatal(err)
	}
}
