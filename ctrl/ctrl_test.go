package ctrl

import (
	"strings"
	"testing"

	"github.com/lutzpeschlow/nas_tools/objects"
)

// TEST
//
// valid json file
func TestReadControlJsonFile_Success(t *testing.T) {
	ctrlObj := objects.Control{}

	err := ReadControlJsonFile("testdata/control_valid.json", &ctrlObj, "linux")
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	// Prüfe Ergebnisse
	if ctrlObj.Action != "READ" {
		t.Errorf("Expected action 'myaction', got '%s'", ctrlObj.Action)
	}
	if ctrlObj.InputFile != "reg_test_01.dat" {
		t.Errorf("Expected InputFile 'reg_test_01.dat', got '%s'", ctrlObj.InputFile)
	}
	if ctrlObj.OutputFile != "result.txt" {
		t.Errorf("Expected OutputFile 'result.txt', got '%s'", ctrlObj.OutputFile)
	}
}

// TEST
//
//	error should output in case of non-existent json file
func TestReadControlJsonFile_MissingFile(t *testing.T) {
	// instance of object
	ctrlObj := objects.Control{}
	// try to read a non-existent json file
	err := ReadControlJsonFile("testdata/nonexistent.json", &ctrlObj, "linux")
	// in case of success, direct fatal
	if err == nil {
		t.Fatal("Expected error for missing file")
	}
	// check for PASS
	//    error contains string - read control file
	if !strings.Contains(err.Error(), "read control file") {
		t.Errorf("Expected file read error, got: %v", err)
	}
}

// TEST
//
//	invalid action in json file
func TestReadControlJsonFile_InvalidAction(t *testing.T) {
	// instance of object
	ctrlObj := objects.Control{}
	// read invalid json fil
	err := ReadControlJsonFile("testdata/control_invalid.json", &ctrlObj, "linux")
	// error should be exist, if nil go out with fatal
	if err == nil {
		t.Fatal("Expected action not found error")
	}
	// no action found, error message should contain string ...
	if !strings.Contains(err.Error(), "action myaction not found") {
		t.Errorf("Expected action error, got: %v", err)
	}
}
