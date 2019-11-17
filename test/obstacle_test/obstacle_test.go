package obstacle_test

import (
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport

	rvo "../../src/rvosimulator"
)

func TestNewObstacle(t *testing.T) {
	result := rvo.NewObstacle()
	expect := &rvo.Obstacle{
		ID:           0,
		IsConvex:     false,
		NextObstacle: nil,
		PrevObstacle: nil,
	}
	if result == expect {
		t.Error("\n実際： ", result, "\n理想： ", expect)
	}

	t.Log("TestObstacle終了")
}
