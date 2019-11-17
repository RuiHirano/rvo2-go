package agent_test

import (
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport

	rvo "../../src/rvosimulator"
)

func TestFlip(t *testing.T) {
	v := rvo.NewVector2(1, 1)
	result := rvo.Flip(v)
	expect := rvo.NewVector2(-1, -1)
	if *result != *expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestFlip終了")
}
