package main

import (
	"flag"

	"github.com/taijiyugame/go-ResourceBundler/api"
	"github.com/taijiyugame/go-ResourceBundler/component"
)

func main() {
	var secretKey, name, resourceDir, outcomeDir string
	flag.StringVar(&secretKey, "s", "", "secretKey")
	flag.StringVar(&resourceDir, "r", "", "resource path")
	flag.StringVar(&outcomeDir, "o", "", "outcome dir")
	flag.StringVar(&name, "n", "", "outcome name")
	flag.Parse()
	packer := api.NewPacker(component.PackOption{
		SecretKey:   []byte(secretKey),
		Name:        name,
		ResourceDir: resourceDir,
		OutcomeDir:  outcomeDir,
	})
	packer.PackResource()
	packer.Save()
}
