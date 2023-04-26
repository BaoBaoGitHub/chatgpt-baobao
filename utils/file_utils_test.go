package utils

import (
	"fmt"
	"testing"
)

var javaCodePath = "../code_to_code/dataset/test0.java-cs.txt.java"
var referencesPath = "../code_to_code/dataset/evaluator/references.txt"

func TestDeleteFiles(t *testing.T) {
	//var path = "../text_to_code/dataset/test_shuffled_with_path_and_id_concode.json"
	path := "../code_to_code/dataset/references.txt"
	var paths []string
	for i := 0; i < 20; i++ {
		tmp1 := AddSuffix(path, i)
		tmp2 := AddSuffix(tmp1, "response")
		paths = append(paths, tmp1)
		paths = append(paths, tmp2)
	}
	DeleteFiles(paths)
}

func TestGenerateReferencesFromPath(t *testing.T) {
	GenerateReferencesFromPath(javaCodePath, referencesPath)
}

func TestCalcClassNumFromPath(t *testing.T) {
	fmt.Println(CalcClassNumFromPath("D:\\Code\\Go\\src\\github.com\\BaoBaoGitHub\\chatgpt-for-se-tasks\\text_to_code\\dataset\\task_prompts\\predictions.txt"))
}
