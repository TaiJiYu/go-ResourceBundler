package component

import (
	"fmt"
	"go-resource-bundler/utils"
	"os"
	"path"
)

type packageByte struct {
	fqmFile *fQMStruct
	option  PackOption
}

func newPackage(options ...PackOption) *packageByte {
	exePath, err := utils.RunPath()
	if err != nil {
		panic(err)
	}
	return (&packageByte{
		fqmFile: newFQMFile(),
		option: PackOption{
			SecretKey:   make([]byte, 0),
			Name:        defaultName,
			ResourceDir: path.Join(exePath, defaultResourceDir),
			OutcomeDir:  path.Join(exePath, defaultOutcomeDir),
		},
	}).checkOption(options...)
}

func (p *packageByte) checkOption(o ...PackOption) *packageByte {
	if len(o) == 0 {
		return p
	}
	optionU := o[0]
	if len(optionU.SecretKey) != 0 {
		if len(optionU.SecretKey) > 16 {
			panic("secretKey is too large and must be less than or equal to 16")
		}
		p.option.SecretKey = optionU.SecretKey
		if err := p.fqmFile.setSecret(optionU.SecretKey); err != nil {
			panic(err)
		}
	}
	if optionU.Name != "" {
		p.option.Name = utils.RemoveFileType(optionU.Name)
	}
	if optionU.ResourceDir != "" {
		p.option.ResourceDir = optionU.ResourceDir
	}
	if optionU.OutcomeDir != "" {
		p.option.OutcomeDir = optionU.OutcomeDir
	}
	return p
}

func (p *packageByte) add(key string, data []byte) {
	p.fqmFile.addByte(key, data)
}
func (p *packageByte) addFile(key string, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	p.add(key, data)
}

func (p *packageByte) save() error {
	name := fmt.Sprintf("%s.%s", path.Join(p.option.OutcomeDir, p.option.Name), defaultExtensions)
	return p.fqmFile.save(name)
}

func (p *packageByte) addResourceDir(dir string) {
	p.addDir(dir, "")
}

func (p *packageByte) packResource() {
	p.addDir(p.option.ResourceDir, "")
}

func (p *packageByte) addDir(baseDir, dirName string) {
	dir := path.Join(baseDir, dirName)
	resous, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, r := range resous {
		if r.IsDir() {
			p.addDir(baseDir, path.Join(dirName, r.Name()))
		} else {
			filePath := path.Join(dirName, r.Name())
			p.addFile(filePath, path.Join(baseDir, filePath))
		}
	}
}
