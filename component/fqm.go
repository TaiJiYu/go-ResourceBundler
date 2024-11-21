package component

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/TaiJiYu/go-ResourceBundler/utils"
)

const (
	defaultName        string = "my_resource"
	defaultResourceDir string = "resource"
	defaultDir         string = "./"
	defaultExtensions  string = "fqm"
	defaultOutcomeDir  string = "outcome"
	headerSize         int    = 30
)

var (
	fQMHeader      = [3]byte{0x66, 0x71, 0x6D}
	sdkVersionMain = 1
	sdkVersionSub  = 0
)

type fQMStruct struct {
	file           *os.File
	baseInfo       *fQMBaseInfo
	indexInfo      []*fQMIndexInfo
	secretKey      []byte
	keyMap         map[string]fQMDataInfo
	keyIndexData   []byte
	keyNameData    []byte
	data           []byte
	dataBeginIndex int64
}

type fQMBaseInfo struct {
	header         [3]byte
	crc_32         [4]byte
	sdkVersionMain int  // 1
	sdkVersionSub  int  // 1
	ifEncrypt      bool // 是否加密
	createTime     time.Time
	indexInfoSize  int
	indexNameSize  int
	dataSize       int
}

type fQMIndexInfo struct {
	keyBeginIndex int
	keySize       int
	dataInfo      fQMDataInfo
}

type fQMDataInfo struct {
	dataBeginIndex int
	dataSize       int
}

func newFQMFile() *fQMStruct {
	return &fQMStruct{
		baseInfo: &fQMBaseInfo{
			header:         fQMHeader,
			crc_32:         [4]byte{},
			sdkVersionMain: sdkVersionMain,
			sdkVersionSub:  sdkVersionSub,
			ifEncrypt:      false,
			createTime:     time.Now(),
		},
		indexInfo:    make([]*fQMIndexInfo, 0),
		keyIndexData: make([]byte, 0),
		keyNameData:  make([]byte, 0),
		data:         make([]byte, 0),
		secretKey:    make([]byte, 0),
		keyMap:       make(map[string]fQMDataInfo, 0),
	}
}

func (f *fQMStruct) setSecret(secretKey []byte) error {
	if len(secretKey) > 16 {
		return fmt.Errorf("secretKey is too large and must be less than or equal to 16")
	}
	if len(secretKey) == 0 {
		return nil
	}
	f.baseInfo.ifEncrypt = len(secretKey) != 0
	f.secretKey = secretKey
	return nil
}

func (f *fQMStruct) save(path string) error {
	if err := f.aesBaseInfo(); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	ret := make([]byte, 0)
	ret = append(ret, f.baseInfo.header[:]...)
	ret = append(ret, []byte{0, 0, 0, 0}...)
	ret = append(ret, byte(f.baseInfo.sdkVersionMain), byte(f.baseInfo.sdkVersionSub))
	ret = append(ret, utils.BoolToByte(f.baseInfo.ifEncrypt))
	ret = append(ret, utils.IntToSizeBytes(f.baseInfo.createTime.Year(), 2)...)
	ret = append(ret, byte(f.baseInfo.createTime.Month()), byte(f.baseInfo.createTime.Day()))
	ret = append(ret, utils.IntToSizeBytes(len(f.keyIndexData), 4)...)
	ret = append(ret, utils.IntToSizeBytes(len(f.keyNameData), 6)...)
	ret = append(ret, utils.IntToSizeBytes(len(f.data), 6)...)
	ret = append(ret, f.keyIndexData...)
	ret = append(ret, f.keyNameData...)
	ret = append(ret, f.data...)
	crc := utils.CRC(ret)
	ret[3], ret[4], ret[5], ret[6] = crc[0], crc[1], crc[2], crc[3]
	_, err = file.Write(ret)
	return err
}

func (f *fQMStruct) aesBaseInfo() error {
	keyIndexText, err := utils.AES(f.secretKey, f.keyIndexData)
	if err != nil {
		return err
	}
	keyNameText, err := utils.AES(f.secretKey, f.keyNameData)
	if err != nil {
		return err
	}
	f.keyIndexData = keyIndexText
	f.baseInfo.indexInfoSize = len(f.keyIndexData)
	f.keyNameData = keyNameText
	f.baseInfo.indexNameSize = len(f.keyNameData)
	return nil
}

func (f *fQMStruct) addByte(key string, data []byte) error {
	data, err := f.aesData(data)
	if err != nil {
		return err
	}
	f.keyIndexData = append(f.keyIndexData, f.getKeyInfo(key, data).toBytes()...)
	f.keyNameData = append(f.keyNameData, []byte(key)...)
	f.data = append(f.data, data...)
	return nil
}

func (f *fQMStruct) aesData(data []byte) ([]byte, error) {
	if len(f.secretKey) == 0 {
		return data, nil
	}
	if aesData, err := utils.AES(f.secretKey, data); err != nil {
		return data, err
	} else {
		return aesData, nil
	}
}

func (f *fQMStruct) deAesData(data []byte) ([]byte, error) {
	if len(f.secretKey) == 0 {
		return data, nil
	}
	if realData, err := utils.DeAES(f.secretKey, data); err != nil {
		return data, err
	} else {
		return realData, nil
	}
}

func (f *fQMStruct) getKeyInfo(key string, data []byte) fQMIndexInfo {
	return fQMIndexInfo{
		keyBeginIndex: len(f.keyNameData),
		keySize:       len([]byte(key)),
		dataInfo: fQMDataInfo{
			dataBeginIndex: len(f.data),
			dataSize:       len(data),
		},
	}
}

func (i fQMIndexInfo) toBytes() []byte {
	ret := make([]byte, 0)
	ret = append(ret, utils.IntTo4Bytes(i.keyBeginIndex)...)
	ret = append(ret, utils.IntTo2Bytes(i.keySize)...)
	ret = append(ret, utils.IntTo6Bytes(i.dataInfo.dataBeginIndex)...)
	ret = append(ret, utils.IntTo6Bytes(i.dataInfo.dataSize)...)
	return ret
}

func (f *fQMStruct) readFqmFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	f.file = file
	if buf, err := f.readBytes(headerSize); err != nil {
		return err
	} else {
		f.readBaseInfoFromBytes(buf)
	}
	if buf, err := f.readBytes(f.baseInfo.indexInfoSize); err != nil {
		return err
	} else if data, err := f.deAesData(buf); err != nil {
		return err
	} else {
		f.readIndexInfo(data)
	}
	if buf, err := f.readBytes(f.baseInfo.indexNameSize); err != nil {
		return err
	} else if data, err := f.deAesData(buf); err != nil {
		return err
	} else {
		f.readIndexName(data)
	}
	f.dataBeginIndex = int64(headerSize + f.baseInfo.indexInfoSize + f.baseInfo.indexNameSize)
	return nil
}

func (f *fQMStruct) close() {
	if f.file == nil {
		return
	}
	f.file.Close()
}

func (f *fQMStruct) readBytes(bufSize int) ([]byte, error) {
	buf := make([]byte, bufSize)
	if size, err := f.file.Read(buf); err != nil || size != bufSize {
		return buf, fmt.Errorf("read err:%v, Files may be corrupted", err)
	}
	return buf, nil
}

func (f *fQMStruct) readBaseInfoFromBytes(data []byte) {
	f.baseInfo.header = [3]byte(data[:3])
	f.baseInfo.crc_32 = [4]byte(data[3:7])
	f.baseInfo.sdkVersionMain = int(data[7])
	f.baseInfo.sdkVersionSub = int(data[8])
	f.baseInfo.ifEncrypt = utils.CheckBit(data[9], 0)
	f.baseInfo.createTime = time.Date(utils.BytesToInt(data[10:12]), time.Month(int(data[12])), int(data[13]), 0, 0, 0, 0, time.Local)
	f.baseInfo.indexInfoSize = utils.BytesToInt(data[14:18])
	f.baseInfo.indexNameSize = utils.BytesToInt(data[18:24])
	f.baseInfo.dataSize = utils.BytesToInt(data[24:30])
}

func (f *fQMStruct) readIndexInfo(indexInfo []byte) {
	if f.baseInfo.ifEncrypt {

	}
	for i := 0; i < len(indexInfo)/18; i++ {
		chunkIndex := i * 18
		f.indexInfo = append(f.indexInfo, &fQMIndexInfo{
			keyBeginIndex: utils.BytesToInt(indexInfo[chunkIndex : chunkIndex+4]),
			keySize:       utils.BytesToInt(indexInfo[chunkIndex+4 : chunkIndex+6]),
			dataInfo: fQMDataInfo{
				dataBeginIndex: utils.BytesToInt(indexInfo[chunkIndex+6 : chunkIndex+12]),
				dataSize:       utils.BytesToInt(indexInfo[chunkIndex+12 : chunkIndex+18]),
			},
		})
	}
}

func (f *fQMStruct) readIndexName(nameData []byte) {
	for _, keyInfo := range f.indexInfo {
		name := nameData[keyInfo.keyBeginIndex : keyInfo.keyBeginIndex+keyInfo.keySize]
		f.keyMap[string(name)] = keyInfo.dataInfo
	}
}

func (f *fQMStruct) show() {
	fmt.Printf("Pack sdk version:%v.%v\n", f.baseInfo.sdkVersionMain, f.baseInfo.sdkVersionSub)
	fmt.Printf("Encryption status:%v\n", f.baseInfo.ifEncrypt)
	fmt.Printf("CRC-32 check digit:%v\n", hex.EncodeToString(f.baseInfo.crc_32[:]))
	fmt.Printf("Create time:%v\n", f.baseInfo.createTime.Format(time.DateOnly))
	fmt.Printf("Key num:%v Data size:%v\n", len(f.indexInfo), utils.ByteSizeToS(f.baseInfo.dataSize))
	for k, v := range f.keyMap {
		fmt.Printf("key:%v data_index:%v data_size:%v\n", k, v.dataBeginIndex, v.dataSize)
	}
}

// 一定成功，否则panic
func (f *fQMStruct) key(key string) []byte {
	v, ok := f.keyMap[key]
	if !ok {
		panic("key is not exist")
	}
	if _, err := f.file.Seek(f.dataBeginIndex+int64(v.dataBeginIndex), 0); err != nil {
		panic(err)
	}
	if buf, err := f.readBytes(v.dataSize); err != nil {
		panic(err)
	} else if data, err := f.deAesData(buf); err != nil {
		panic(err)
	} else {
		return data
	}

}
