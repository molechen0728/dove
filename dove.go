package dove

import (
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
	_FTP    = 0x08
	_BASE64 = 0x10
)

type Dover interface {
	Throw(error) Dover
	Drop() Dover
	Return() Dover
	Next(int) Dover
	Util(string) Dover
	Buffer(int) Dover
	Get(int) (interface{}, error)
	Ok() bool
}

type ecmd struct {
	e   error
	cmd string
}

type dove struct {
	ecmd   ecmd
	next   int
	params []interface{}
	buffer []byte
	ok     bool
	memory bool
	reader io.ReadCloser
}

type tinfo struct {
	no       int
	file     *os.File
	response *http.Response
	bases    []byte
}

func (m *dove) Throw(err error) Dover {
	m.ecmd.e = err
	return m
}

func (m *dove) Drop() Dover {
	m.next += len(m.buffer)
	return m
}

func (m *dove) Return() Dover {
	m.ecmd.cmd = "return"
	return m
}

func (m *dove) Next(n int) Dover {
	m.next = n
	return m
}

func (m *dove) Util(end string) Dover {
	return m
}

func (m *dove) Buffer(n int) Dover {

	tmp := m.buffer[:0]
	copy(tmp, m.buffer)
	m.buffer = make([]byte, n)
	copy(m.buffer, tmp)

	if len(tmp) > 0 {
		_, _ = io.CopyBuffer(bytes.NewBuffer(m.buffer), m.reader, tmp)
	}

	return m
}

func (m *dove) Ok() bool {
	return m.ok
}

func (m *dove) Get(n int) (interface{}, error) {

	if len(m.params) == 0 {
		return nil, errors.New("has no object")
	}

	if n < 0 || n >= len(m.params) {
		return nil, errors.New("invalid index")
	}

	return m.params[n], nil
}

func (m *dove) lazyRead(b []byte) (n int, err error) {
	return m.reader.Read(b)
}

func Fly(uri string, f func(i Dover, n int, b []byte) (o Dover), objs ...interface{}) (err error) {

	i, n, b := &dove{params: objs}, 0, []byte{}

	defer catch(i, &err, n, b)

	var o = f(i, n, b)
	o.(*dove).ok = true
	b = make([]byte, len(o.(*dove).buffer))

	reader, err := o.(*dove).getReader(uri)
	if err != nil {
		return err
	}
	defer reader.Close()

	i.reader = reader

	for {

		i = o.(*dove)
		if i.ecmd.e != nil {
			err = i.ecmd.e
			return
		}

		if i.ecmd.cmd == "return" {
			return
		}

		n, err = i.lazyRead(b)

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		o = f(i, n, b)

	}
}

func (m *dove) getReader(uri string) (io.ReadCloser, error) {

	switch info := guess(uri); info.no {
	case _HTTP:
		return info.response.Body, nil
	case _HTTPS:
		return info.response.Body, nil
	case _FTP:
		return nil, errors.New("unsupport type")
	case _BASE64:
		m.memory = true
		return io.NopCloser(bytes.NewReader(info.bases)), nil
	case _FILE:
		return info.file, nil
	}

	return nil, errors.New("unsupport type")
}

func guess(uri string) tinfo {

	if strings.HasPrefix(uri, "https://") {
		resp, err := http.Get(uri)
		if err != nil {
			return tinfo{}
		}
		return tinfo{no: _HTTPS, response: resp}
	}

	if strings.HasPrefix(uri, "http://") {
		resp, err := http.Get(uri)
		if err != nil {
			return tinfo{}
		}
		return tinfo{no: _HTTP, response: resp}
	}

	if strings.HasPrefix(uri, "ftp://") {
		return tinfo{no: _FTP}
	}

	if f, err := os.Open(uri); !errors.Is(err, fs.ErrNotExist) {
		return tinfo{no: _FILE, file: f}
	}

	if s, err := base64.StdEncoding.DecodeString(uri); err == nil {
		return tinfo{no: _BASE64, bases: s}
	}

	return tinfo{}
}

func catch(d *dove, exp *error, n int, b []byte) {

	if obj := recover(); obj != nil {

		if e, ok := obj.(error); ok && !d.ok &&
			strings.Contains(e.Error(), "runtime error: index out of range") {

			*exp = errors.New("call Dover's Ok function first ")
			return
		}

		panic(obj)
	}
}
