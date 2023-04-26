package translation

import (
	"testing"
)

func TestGetJSONFilePathFromOtherFile(t *testing.T) {
	if ModifyFileExtToJSON("test0.txt") != "test0.json" {
		t.Error()
	}
}
