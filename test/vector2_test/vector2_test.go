package vector2_test

import (
	"testing"

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

func TestSub(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	v2 := rvo.NewVector2(1, 1)
	result := rvo.Sub(v1, v2)
	expect := rvo.NewVector2(2, 3)
	if *result != *expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestSub終了")
}

func TestAdd(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	v2 := rvo.NewVector2(1, 1)
	result := rvo.Add(v1, v2)
	expect := rvo.NewVector2(4, 5)
	if *result != *expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestAdd終了")
}

func TestMul(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	v2 := rvo.NewVector2(1, 1)
	result := rvo.Mul(v1, v2)
	expect := float64(7)
	if result != expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestMul終了")
}

func TestMulOne(t *testing.T) {
	v := rvo.NewVector2(3, 4)
	result := rvo.MulOne(v, 0.5)
	expect := rvo.NewVector2(1.5, 2)
	if *result != *expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestMulOne終了")
}

func TestDiv(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	result := rvo.Div(v1, 2)
	expect := rvo.NewVector2(1.5, 2)
	if *result != *expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestDiv終了")
}

func TestEqual(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	v2 := rvo.NewVector2(1, 1)

	//異なるベクトルの場合、falseになる
	if rvo.Equal(v1, v2) != false {
		t.Error("\n実際： ", true, "\n理想： ", false)
	}

	// 等しいベクトルの場合、trueになる
	if rvo.Equal(v1, v1) != true {
		t.Error("\n実際： ", false, "\n理想： ", true)
	}

	t.Log("TestEqual終了")
}

func TestNotEqual(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	v2 := rvo.NewVector2(1, 1)

	//異なるベクトルの場合、trueになる
	if rvo.NotEqual(v1, v2) != true {
		t.Error("\n実際： ", false, "\n理想： ", true)
	}

	// 等しいベクトルの場合、falseになる
	if rvo.NotEqual(v1, v1) != false {
		t.Error("\n実際： ", true, "\n理想： ", false)
	}

	t.Log("TestNotEqual終了")
}

func TestMulSum(t *testing.T) {
	v := rvo.NewVector2(3, 4)
	result := rvo.MulSum(v, 0.5)
	expect := rvo.NewVector2(4.5, 6)
	if *result != *expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestMulSum終了")
}

func TestDivSum(t *testing.T) {
	v := rvo.NewVector2(3, 4)
	result := rvo.DivSum(v, 2)
	expect := rvo.NewVector2(4.5, 6)
	if *result != *expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestDivSum終了")
}

func TestAddSum(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	v2 := rvo.NewVector2(2, 1)
	result := rvo.AddSum(v1, v2)
	expect := rvo.NewVector2(8, 9)
	if *result != *expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestAddSum終了")
}

func TestSubSum(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	v2 := rvo.NewVector2(2, 1)
	result := rvo.SubSum(v1, v2)
	expect := rvo.NewVector2(4, 7)
	if *result != *expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestSubSum終了")
}

func TestSqr(t *testing.T) {
	v := rvo.NewVector2(3, 4)
	result := rvo.Sqr(v)
	expect := float64(25)
	if result != expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestSqr終了")
}

func TestAbs(t *testing.T) {
	v := rvo.NewVector2(-3, 4)
	result := rvo.Abs(v)
	expect := float64(5)
	if result != expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestAbs終了")
}

func TestNormalize(t *testing.T) {
	v := rvo.NewVector2(3, 4)
	result := rvo.Normalize(v)
	expect := rvo.NewVector2(float64(3)/5, float64(4)/5)
	if *result != *expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestNormalize終了")
}

func TestDet(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	v2 := rvo.NewVector2(1, 1)
	result := rvo.Det(v1, v2)
	expect := float64(-1)
	if result != expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestDet終了")
}

func TestLeftOf(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	v2 := rvo.NewVector2(1, 1)
	v3 := rvo.NewVector2(2, 3)
	result := rvo.LeftOf(v1, v2, v3)
	expect := float64(-1)
	if result != expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestLeftOf終了")
}

func TestSqPointLineSegment(t *testing.T) {
	v1 := rvo.NewVector2(3, 4)
	v2 := rvo.NewVector2(1, 1)
	v3 := rvo.NewVector2(2, 3)
	result := rvo.DistSqPointLineSegment(v1, v2, v3)
	expect := 0.07692307692307705
	if result != expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestSqPointLineSegment終了")
}
