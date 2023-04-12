package utils

import "testing"

func TestDeleteFiles(t *testing.T) {
	var path = "../text_to_code/dataset/test_shuffled_with_path_and_id_concode.json"
	var paths []string
	for i := 0; i < 20; i++ {
		tmp1 := AddSuffix(path, i)
		tmp2 := AddSuffix(tmp1, "response")
		paths = append(paths, tmp1)
		paths = append(paths, tmp2)
	}
	DeleteFiles(paths)
}
