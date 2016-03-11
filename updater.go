package updater

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	copyWriteUrl = "http://update.cz88.net/ip/copywrite.rar"
	qqwryUrl     = "http://update.cz88.net/ip/qqwry.rar"
)

func getContent(url string) (b []byte, err error) {
	if strings.Contains(url, "copywrite.rar") {
		return ioutil.ReadFile("copywrite.rar")
	}
	if strings.Contains(url, "qqwry.rar") {
		return ioutil.ReadFile("qqwry.rar")
	}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	return
}

func getKey(b []byte) (key uint32, err error) {
	if len(b) != 280 {
		return 0, errors.New("copywrite.rar is corrupt")
	}
	key = binary.LittleEndian.Uint32(b[20:])
	return
}

func decrypt(b []byte, key uint32) (_ []byte, err error) {
	for i := 0; i < 0x200; i++ {
		key *= uint32(0x805)
		key++
		key &= uint32(0xff)
		b[i] = b[i] ^ byte(key)
	}
	rc, err := zlib.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

func Fetch() (b []byte, err error) {
	copyWriteData, err := getContent(copyWriteUrl)
	if err != nil {
		return
	}

	qqwryData, err := getContent(qqwryUrl)
	if err != nil {
		return
	}

	key, err := getKey(copyWriteData)
	if err != nil {
		return
	}

	return decrypt(qqwryData, key)
}
