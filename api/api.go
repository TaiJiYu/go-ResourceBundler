package api

import "github.com/TaiJiYu/go-ResourceBundler/component"

type IPacker interface {
	Add(key string, data []byte)
	AddResourceDir(dir string)
	AddFile(key string, path string)
	PackResource()
	Save() error
}

// packer
func NewPacker(o ...component.PackOption) IPacker { return component.NewPacker(o...) }

type IUnPacker interface {
	Key(key string) []byte
	Close()
	Show()
}

// unpacker
func UnpackerInit(o ...component.UnpackOption) (IUnPacker, error) {
	return component.UnpackerInit(o...)
}
