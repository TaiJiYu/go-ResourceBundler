package component

import (
	"path"

	"github.com/TaiJiYu/go-ResourceBundler/utils"
)

type unpackClient struct {
	fqmFile *fQMStruct
	option  UnpackOption
}

func unpackInit(o ...UnpackOption) (*unpackClient, error) {
	exePath, err := utils.RunPath()
	if err != nil {
		return nil, err
	}
	fqmPath := path.Join(exePath, defaultOutcomeDir)
	unpackCli := (&unpackClient{
		option: UnpackOption{
			FqmFilePath: fqmPath,
			UseCache:    false,
			SecretKey:   make([]byte, 0),
		},
		fqmFile: newFQMFile(),
	}).checkOption(o...)
	if err := unpackCli.fqmFile.readFqmFromFile(unpackCli.option.FqmFilePath); err != nil {
		return nil, err
	} else {
		return unpackCli, nil
	}
}

func (u *unpackClient) checkOption(o ...UnpackOption) *unpackClient {
	if len(o) == 0 {
		return u
	}
	u.option = o[0]
	u.fqmFile.secretKey = u.option.SecretKey
	return u
}

func (u *unpackClient) key(key string) []byte {
	if !u.option.UseCache {
		return u.fqmFile.key(key)
	}
	if v, ok := defaultCache().get(key); ok {
		return v
	} else {
		data := u.fqmFile.key(key)
		defaultCache().set(key, data)
		return data
	}
}

func (u *unpackClient) close() {
	if u == nil || u.fqmFile == nil {
		return
	}
	u.fqmFile.close()
	if u.option.UseCache {
		defaultCache().cache.Close()
	}
}

func (u *unpackClient) show() {
	u.fqmFile.show()
}
