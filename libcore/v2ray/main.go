package main

import (
	"os"
	_ "unsafe"

	"github.com/v2fly/v2ray-core/v5/main/commands"
)

// PC版 适配 qv2ray core & 插件

//go:linkname build github.com/v2fly/v2ray-core/v5.build
var build string

func main() {
	build = "Matsuridayo/Qv2ray"
	commands.CmdRun.Run(commands.CmdRun, os.Args[1:])
}
