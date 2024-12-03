package main

import (
	// Modules of packwiz
	_ "github.com/packwiz/packwiz/basedon"
	"github.com/packwiz/packwiz/cmd"
	_ "github.com/packwiz/packwiz/curseforge"
	_ "github.com/packwiz/packwiz/migrate"
	_ "github.com/packwiz/packwiz/modrinth"
	_ "github.com/packwiz/packwiz/settings"
	_ "github.com/packwiz/packwiz/url"
	_ "github.com/packwiz/packwiz/utils"
)

func main() {
	cmd.Execute()
}
