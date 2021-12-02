package dove

import (
	"errors"
	"io"
	"os"
	"strings"
)

type Dover interface {
	Throw(error) Dover
	Drop() Dover
	Return() Dover
	Next(int) Dover
	Buffer(int) Dover
	Get(int) (interface{}, error)
}

type ecmd struct {
	e   error
	cmd string
}

type dove struct {
	emsg   ecmd
	next   int
	params []interface{}
	buffer []byte
}

func (m *dove) Throw(err error) Dover {
	m.emsg.e = err
	return m
}

func (m *dove) Drop() Dover {
	m.next += len(m.buffer)
	return m
}

func (m *dove) Return() Dover {
	m.emsg.cmd = "return"
	return m
}

func (m *dove) Next(n int) Dover {
	m.next = n
	return m
}

func (m *dove) Buffer(n int) Dover {

	tmp := m.buffer[:0]
	copy(tmp, m.buffer)
	m.buffer = make([]byte, n)
	copy(m.buffer, tmp)

	return m
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

func Fly(uri string, f func(i Dover, n int, b []byte) (o Dover), objs ...interface{}) (err error) {

	i, n, b := &dove{params: objs}, 0, make([]byte, 1024)

	defer catch(i, &err, n, b)

	var o = f(i, n, b)
	b = make([]byte, len(o.(*dove).buffer))

	reader, err := getReader(uri)
	if err != nil {
		return err
	}
	defer reader.Close()

	for {
		i = o.(*dove)
		if i.emsg.e != nil {
			err = i.emsg.e
			return
		}

		if i.emsg.cmd == "return" {
			return
		}
		n, err = reader.Read(b)

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return
		}

		o = f(i, n, b)

	}
}

func getReader(uri string) (io.ReadCloser, error) {

	if uri == "E:\\logs\\nlp-service.log" {
		return os.Open(uri)
	}

	return nil, nil
}

func catch(d *dove, exp *error, n int, b []byte) {
	if obj := recover(); obj != nil {
		if e, ok := obj.(error); ok && n == 0 && len(b) == 0 &&
			strings.Contains(e.Error(), "runtime error: index out of range") {
			*exp = errors.New("check the byte slice for safety, please! ")
			return
		}
		panic(obj)
	}
}
