package dove

import (
	"testing"
)

func TestDove(t *testing.T) {

	err := Fly(
		"E:\\logs\\nlp-service.log",
		// "https://pkg.go.dev/fmt",
		// "http://192.168.0.206/arm/vpr_onnx_4.1.4_Linux_aarch64_20210715.tar.gz",

		func(i Dover, n int, b []byte) (o Dover) {
			// v1, _ := i.Get(0)
			// t.Log(v1.(int))
			// v2, _ := i.Get(1)
			// t.Log(v2.(int))
			// v3, _ := i.Get(2)
			// t.Log(v3.(int))
			if i.Ok() {
				t.Log("\n", string(b))
				return i.Next(1).Buffer(32)
			}
			return i.Next(1).Buffer(128)
		},

		// 1, 2, 3,
	)

	if err != nil {
		t.Error(err)
	}
}
