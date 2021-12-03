package dove

import (
	"bufio"
	"os"
	"testing"
)

func TestDove(t *testing.T) {

	err := Fly(
		"https://pkg.go.dev/fmt",
		1024,

		func(rd *bufio.Reader, objs ...interface{}) error {
			_, err := rd.WriteTo(os.Stdout)
			return err
		},
	)

	if err != nil {
		t.Error(err)
	}
}
