package main

import (
	"github.com/suzuki-shunsuke/ghatm/pkg/cli"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

var version = ""

func main() {
	urfave.Main("ghatm", version, cli.Run)
}
