package utils

import "testing"



func AssertNumResult(t *testing.T, expect, real int64) {
	t.Helper()
	if expect != real {
		t.Errorf("Expect: %v, Actual: %v\n", expect, real)
	}
}

func AssertStringResult(t *testing.T, expect, real string) {
	t.Helper()
	if expect != real {
		t.Errorf("Expect: %v, Actual: %v\n", expect, real)
	}
}

func AssertFloatResult(t *testing.T, expect, real float64) {
	t.Helper()
	if expect != real {
		t.Errorf("Expect: %v, Actual: %v\n", expect, real)
	}
}

func AssertTrue(t *testing.T, v bool) {
	t.Helper()
	if v != true {
		t.Error("Assert True but get a false value")
	}
}

func AssertFalse(t *testing.T, v bool) {
	t.Helper()
	if v == true {
		t.Error("Assert True but get a false value")
	}
}

func AssertThread(t *testing.T, threadStat int32) {
	t.Helper()
	if threadStat != 0 {
		t.Error("The clean sessions thread is still alive?")
	}
}
