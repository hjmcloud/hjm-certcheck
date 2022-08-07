package main

import (
	_ "hjm-certcheck/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"hjm-certcheck/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
