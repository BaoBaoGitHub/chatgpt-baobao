package utils

import (
	"fmt"
	"path/filepath"
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
	fmt.Println(CalcClassNumFromPath("D:\\Code\\Go\\src\\github.com\\BaoBaoGitHub\\chatgpt-for-se-tasks\\text_to_code\\dataset\\guided_prompts_api_exception\\test0\\predictions.txt"))
}

func TestAddSpace(t *testing.T) {
	//src := "D:\\Code\\Go\\src\\github.com\\BaoBaoGitHub\\chatgpt-for-se-tasks\\text_to_code\\dataset\\guided_prompts_api_exception_conciseness\\test0\\predictions.txt"
	//dst := AddSuffix(src, "space")
	src := "D:\\学习\\研一\\Guiding ChatGPT for SE tasks\\guiding实验数据\\text2code\\guided_prompts_api_exception_conciseness\\test0\\predictions_space.txt"
	dst := "D:\\学习\\研一\\Guiding ChatGPT for SE tasks\\guiding实验数据\\text2code\\guided_prompts_api_exception_conciseness\\test0\\predictions.txt"
	AddSpace(src, dst)
}

func TestDeleteOverride(t *testing.T) {
	src := "D:\\Code\\Go\\src\\github.com\\BaoBaoGitHub\\chatgpt-for-se-tasks\\code_to_code\\dataset\\guided_prompts_api_exception\\predictions.txt"
	dst := "D:\\Code\\Go\\src\\github.com\\BaoBaoGitHub\\chatgpt-for-se-tasks\\code_to_code\\dataset\\guided_prompts_api_exception\\predictions_annotation.txt"
	deleteOverride(src, dst)
}

func TestGetPredictionWithoutCommentsWithSpaceFromJSONFile(t *testing.T) {
	srcJSONPath := "D:\\Code\\Go\\src\\github.com\\BaoBaoGitHub\\chatgpt-for-se-tasks\\text_to_code\\dataset\\guided_prompts_api_exception_conciseness\\round0\\test_shuffled_with_path_and_id_concode_response.json"
	tgtDir := filepath.Dir(srcJSONPath)
	GetPredictionWithoutCommentsWithSpaceFromJSONFile(srcJSONPath, tgtDir)
}
