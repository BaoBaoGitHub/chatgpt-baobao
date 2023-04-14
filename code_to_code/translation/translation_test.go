package translation

import (
	"testing"
)

func TestGetJSONFilePathFromOtherFile(t *testing.T) {
	if ModifyFileExtToJSON("test.txt") != "test.json" {
		t.Error()
	}
}
