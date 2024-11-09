package component

import (
	"go-resource-bundler/utils"
	"path"
)

type unpackClient struct {
	fqmFile *fQMStruct
	option  UnpackOption
}

var unpackCli *unpackClient

func unpackInit(o ...UnpackOption) error {
	exePath, err := utils.RunPath()
	if err != nil {
		return err
	}
	fqmPath := path.Join(exePath, defaultOutcomeDir)
	unpackCli = (&unpackClient{
		option: UnpackOption{
			FqmFilePath: fqmPath,
			SecretKey:   make([]byte, 0),
		},
		fqmFile: newFQMFile(),
	}).checkOption(o...)
	return unpackCli.fqmFile.readFqmFromFile(unpackCli.option.FqmFilePath)
}

func (u *unpackClient) checkOption(o ...UnpackOption) *unpackClient {
	if len(o) == 0 {
		return u
	}
	u.option = o[0]
	u.fqmFile.secretKey = u.option.SecretKey
	return u
}

func key(key string) []byte {
	if unpackCli == nil || unpackCli.fqmFile == nil {
		return []byte{}
	}
	if v, ok := defaultCache().get(key); ok {
		return v
	} else {
		data := unpackCli.fqmFile.key(key)
		defaultCache().set(key, data)
		return data
	}
}

func close() {
	if unpackCli == nil || unpackCli.fqmFile == nil {
		return
	}
	unpackCli.fqmFile.close()
	defaultCache().cache.Close()
}

func show() {
	if unpackCli == nil || unpackCli.fqmFile == nil {
		return
	}
	unpackCli.fqmFile.show()
}