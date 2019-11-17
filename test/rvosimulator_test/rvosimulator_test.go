package rvosimulator_test

import (
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport

	rvo "../../src/rvosimulator"
)

func TestNewRVOSimulator(t *testing.T) {
	sim := rvo.NewRVOSimulator(1.5, 1, 1, 1, 1, 1, 1, rvo.NewVector2(1, 1))
	if sim.TimeStep != 1.5 {
		t.Error("TimeStep: \n実際： ", sim.TimeStep, "\n理想： ", 1.5)
	}
	if sim.GlobalTime != 0.0 {
		t.Error("\n実際： ", sim.GlobalTime, "\n理想： ", 0.0)
	}

	t.Log("TestFlip終了")
}
