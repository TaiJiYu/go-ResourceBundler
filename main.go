package main

import (
	"flag"
	"strings"

	"github.com/TaiJiYu/go-ResourceBundler/api"
	"github.com/TaiJiYu/go-ResourceBundler/component"
)

func main() {
	var secretKey, name, resourceDir, outcomeDir, ignoreDirS string
	flag.StringVar(&secretKey, "s", "", "secretKey")
	flag.StringVar(&resourceDir, "r", "", "resource path")
	flag.StringVar(&outcomeDir, "o", "", "outcome dir")
	flag.StringVar(&name, "n", "", "outcome name")
	flag.StringVar(&ignoreDirS, "i", "", "ignore dir,ex:ign,ign1")

	flag.Parse()
	packer := api.NewPacker(component.PackOption{
		SecretKey:   []byte(secretKey),
		Name:        name,
		ResourceDir: resourceDir,
		OutcomeDir:  outcomeDir,
		IgnoreDir:   strings.Split(ignoreDirS, ","),
	})
	packer.PackResource()
	packer.Save()
}
