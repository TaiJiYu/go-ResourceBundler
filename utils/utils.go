package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path"
	"strings"
)

func RemoveFileType(name string) string {
	return strings.TrimSuffix(name, path.Ext(path.Base(name)))
}

func RunPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return path.Dir(strings.ReplaceAll(exePath, "\\", "/")), nil
}

func CRC(data []byte) []byte {
	return UintToSizeBytes(crc32.ChecksumIEEE(data))
}

// 从右向左检查第index位是否为1,index从0开始
func CheckBit(b byte, index int) bool {
	return ((b >> index) & 1) == 1
}
func BytesToInt(b []byte) int {
	k := append(make([]byte, 8-len(b)), b...)
	return int(binary.BigEndian.Uint64(k))
}

func IntTo2Bytes(i int) []byte {
	return IntToSizeBytes(i, 2)
}

func IntTo4Bytes(i int) []byte {
	return IntToSizeBytes(i, 4)
}
func IntTo6Bytes(i int) []byte {
	return IntToSizeBytes(i, 6)
}
func BoolToByte(b bool) byte {
	if b {
		return 1
	} else {
		return 0
	}
}
func IntToSizeBytes(i int, size int) []byte {
	k := make([]byte, 8)
	binary.BigEndian.PutUint64(k, uint64(i))
	return k[8-size:]
}

func UintToSizeBytes(i uint32) []byte {
	k := make([]byte, 4)
	binary.BigEndian.PutUint32(k, i)
	return k
}

// 对s加密
func AES(secretKey []byte, s []byte) ([]byte, error) {
	secretKey, err := ByteFill(secretKey, 16)
	if err != nil {
		return []byte{}, err
	}
	c, err := aes.NewCipher(secretKey)
	if err != nil {
		return []byte{}, err
	}
	text := make([]byte, aes.BlockSize+len(s))
	iv := text[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, err
	}
	stream := cipher.NewCFBEncrypter(c, iv)
	stream.XORKeyStream(text[aes.BlockSize:], s)
	return text, nil
}
func ByteFill(b []byte, size int) ([]byte, error) {
	if len(b) > size {
		return b, fmt.Errorf("bytes should be less than or equal to size")
	}
	f := make([]byte, size-len(b))
	return append(b, f...), nil
}

func DeAES(secretKey []byte, ciphertext []byte) ([]byte, error) {
	secretKey, err := ByteFill(secretKey, 16)
	if err != nil {
		return []byte{}, err
	}
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return []byte{}, err
	}
	if len(ciphertext) < aes.BlockSize {
		return []byte{}, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}

var UnitsList = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB", "BB"}

// 字节数转文本
func ByteSizeToS(size int) string {
	s := []int{}
	for size >= 0 && len(s) < len(UnitsList)-1 {
		k := size / 1024
		if k == 0 {
			s = append(s, size)
			size = 0
			break
		} else {
			s = append(s, size-k*1024)
			size = k
		}
	}
	s = append(s, size)
	r := ""
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == 0 {
			continue
		}
		p := fmt.Sprintf("%v%v", s[i], UnitsList[i])
		if r != "" {
			r += "+" + p
		} else {
			r = p
		}
	}
	if r == "" {
		return "0B"
	}
	return r

}
