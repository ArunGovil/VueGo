//go:build ignore

package main

import (
	"log"

	"./assets"
	"github.com/shurcooL/vfsgen"
)

func main() {
	var err error

	err = vfsgen.Generate(assets.Assets, vfsgen.Options{
		Filename:     "./assets/appcode.go",
		PackageName:  "assets",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})

	if err != nil {
		log.Fatal(err)
	}
}
