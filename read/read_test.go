package read

import "testing"

func TestRead(t *testing.T) {
	if Add(2, 3) != 5 {
		t.Errorf("Error")
	}
}
