package parse

import (
	"testing"
)

func TestGenResFromJSONFile(t *testing.T) {
	GenResFromJSONFile("../dataset/code_generation/references_api.json", "../dataset/code_generation/references_api.txt", APIMode)
}
