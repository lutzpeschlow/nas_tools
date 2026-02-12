package ctrl

import (
	"strings"
	"testing"

	"github.com/lutzpeschlow/nas_tools/objects"
)

func TestReadControlJsonFile_Success(t *testing.T) {
	ctrlObj := objects.Control{}

	err := ReadControlJsonFile("testdata/control_valid.json", &ctrlObj, "linux")
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	// Prüfe Ergebnisse
	if ctrlObj.Action != "myaction" {
		t.Errorf("Expected action 'myaction', got '%s'", ctrlObj.Action)
	}
	if ctrlObj.InputFile != "test.nas" {
		t.Errorf("Expected InputFile 'test.nas', got '%s'", ctrlObj.InputFile)
	}
	if ctrlObj.FullInputPath != "input/test.nas" {
		t.Errorf("Expected FullInputPath 'input/test.nas', got '%s'", ctrlObj.FullInputPath)
	}
}

func TestReadControlJsonFile_MissingFile(t *testing.T) {
	ctrlObj := objects.Control{}

	err := ReadControlJsonFile("testdata/nonexistent.json", &ctrlObj, "linux")
	if err == nil {
		t.Fatal("Expected error for missing file")
	}
	if !strings.Contains(err.Error(), "read control file") {
		t.Errorf("Expected file read error, got: %v", err)
	}
}

func TestReadControlJsonFile_InvalidAction(t *testing.T) {
	// testdata/control_invalid.json ohne "myaction"
	ctrlObj := objects.Control{}

	err := ReadControlJsonFile("testdata/control_invalid.json", &ctrlObj, "linux")
	if err == nil {
		t.Fatal("Expected action not found error")
	}
	if !strings.Contains(err.Error(), "action myaction not found") {
		t.Errorf("Expected action error, got: %v", err)
	}
}
