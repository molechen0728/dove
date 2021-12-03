package dove

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

const (
	_FILE   = 0x01
	_HTTP   = 0x02
	_HTTPS  = 0x04
	_BASE64 = 0x08
)

type tinfo struct {
	no       int
	file     *os.File
	response *http.Response
	bases    []byte
}

func Fly(url string, bsize int, f func(rd *bufio.Reader, objs ...interface{}) error, objs ...interface{}) (err error) {

	reader, err := getReader(url)
	if err != nil {
		return err
	}
	defer reader.Close()

	if bsize <= 0 {
		bsize = 4096
	}

	return f(bufio.NewReaderSize(reader, bsize), objs...)

}

func getReader(url string) (io.ReadCloser, error) {

	info, err := guess(url)
	if err != nil {
		return nil, err
	}

	switch info.no {
	case _HTTP:
		return info.response.Body, nil
	case _HTTPS:
		return info.response.Body, nil
	case _BASE64:
		return io.NopCloser(bytes.NewReader(info.bases)), nil
	case _FILE:
		return info.file, nil
	}

	return nil, errors.New("unsupport type")
}

func guess(url string) (tinfo, error) {

	if strings.HasPrefix(url, "https://") {
		resp, err := http.Get(url)
		if err != nil {
			return tinfo{}, err
		}
		return tinfo{no: _HTTPS, response: resp}, nil
	}

	if strings.HasPrefix(url, "http://") {
		resp, err := http.Get(url)
		if err != nil {
			return tinfo{}, err
		}
		return tinfo{no: _HTTP, response: resp}, nil
	}

	if f, err := os.Open(url); !errors.Is(err, fs.ErrNotExist) {
		return tinfo{no: _FILE, file: f}, nil
	} else if err != nil {
		return tinfo{}, err
	}

	if s, err := base64.StdEncoding.DecodeString(url); err == nil {
		return tinfo{no: _BASE64, bases: s}, nil
	} else {
		return tinfo{}, err
	}

}
