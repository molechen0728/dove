package dove

import (
	"testing"
)

func TestDove(t *testing.T) {

	err := Fly(
		"E:\\logs\\nlp-service.log",

		func(i Dover, n int, b []byte) (o Dover) {
			// v1, _ := i.Get(0)
			// t.Log(v1.(int))
			// v2, _ := i.Get(1)
			// t.Log(v2.(int))
			// v3, _ := i.Get(2)
			// t.Log(v3.(int))

			if len(b) > 0 {
				t.Log("==", string(b))
				return i.Next(1).Buffer(32)
			}
			return i.Next(1).Buffer(32)
		},
		// 1, 2, 3,
	)
	if err != nil {
		t.Error(err)
	}
}
