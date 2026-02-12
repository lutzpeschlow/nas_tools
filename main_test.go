package main

import (
	"testing"

	"github.com/lutzpeschlow/nas_tools/objects"
	"github.com/lutzpeschlow/nas_tools/read"
)

func TestReadNasFile(t *testing.T) {
	// test without BEGIN BULK, regression test file is reg_test_01.dat
	mod1 := objects.Model{}
	var len_01 int
	var len_02 int
	len_01 = 0
	len_02 = 0

	len_01, len_02, err := read.ReadNasFile("regression_tests//reg_test_01.dat", &mod1)
	if err != nil {
		t.Errorf("problem reading file: %v", err)
	}
	if len_01 != 47 {
		t.Errorf("len_01 wrong: got %d, want %d", len_01, 47)
	}
	if len_02 != 47 {
		t.Errorf("len_02 wrong: got %d, want %d", len_02, 47)
	}
	// test with BEGIN BULK, regression test file is reg_test_02.dat
	mod2 := objects.Model{}
	len_01 = 0
	len_02 = 0
	len_01, len_02, err = read.ReadNasFile("regression_tests//reg_test_02.dat", &mod2)
	if err != nil {
		t.Errorf("problem reading file: %v", err)
	}
	if len_01 != 47 {
		t.Errorf("len_01 wrong: got %d, want %d", len_01, 47)
	}
	if len_02 != 47 {
		t.Errorf("len_02 wrong: got %d, want %d", len_02, 47)
	}
}

// test
//
//	run function without final error
func TestRun_Success(t *testing.T) {
	err := run()
	if err != nil {
		t.Errorf("run() should not deliver any error: %v", err)
	}
}
