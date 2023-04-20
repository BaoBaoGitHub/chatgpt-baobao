package utils

import (
	"log"
	"testing"
)

import (
	"fmt"
)

var path = "../text_to_code/dataset/test_shuffled_with_path_and_id_concode.json"
var successfulContent = "Assuming you have a class called `FunctionNode` that represents a function node, and that it has instance variables for the function name, parameters, and variables associated with it, you could generate mappings for each function node and its associated parameters and variables using the following Java Code:\n\n```\nMap<FunctionNode, Map<String, String>> mappings = new HashMap<>();\n\nfor (FunctionNode node : functionNodes) {\n    Map<String, String> nodeMappings = new HashMap<>();\n    nodeMappings.put(\"functionName\", node.getFunctionName());\n\n    Map<String, String> parameterMappings = new HashMap<>();\n    for (String parameterName : node.getParameters()) {\n        parameterMappings.put(parameterName, \"parameter\");\n\n    }\n    nodeMappings.put(\"parameters\", parameterMappings.toString());\n\n    Map<String, String> variableMappings = new HashMap<>();\n    for (String variableName : node.getVariables()) {\n        variableMappings.put(variableName, \"variable\");\n    }\n    nodeMappings.put(\"variables\", variableMappings.toString());\n\n    mappings.put(node, nodeMappings);\n}\n```\n\nThis Code creates a `HashMap` called `mappings` that maps each `FunctionNode` object to a `HashMap` of mappings for that node. The mappings for each node include the function name, the parameters and their types (in this example, \"parameter\"), and the variables and their types (in this example, \"variable\"). \n\nYou would need to modify this Code to match the implementation details of your `FunctionNode` class."
var respPath = "../text_to_code/dataset/test_file_response.json"
var testFilesPath = []string{
	"../text_to_code/dataset/test_file_1_response.json", "../text_to_code/dataset/test_file_2_response.json",
}
var java_cs_path = "../code_to_code/dataset/references.txt"
var cs_java_path = "../code_to_code/dataset/references.txt"
var finalPath = "../text_to_code/dataset/test_shuffled_with_path_and_id_concode_response.json"
var concodePath = "../text_to_code/dataset/test_shuffled_with_path_and_id_concode.json"
var predictionsPath = "../text_to_code/dataset/evaluator/predictions.txt"
var answersPath = "../text_to_code/dataset/evaluator/answers.json"
var cSharpPath = "../code_to_code/dataset/evaluator/references.txt"
var code2codePredictionPath = "../code_to_code/dataset/evaluator/predictions.txt"
var code2codeResponsePath = "../code_to_code/dataset/test.java-cs.txt_response.json"

func TestGetData(t *testing.T) {
	data := ReadFromJsonFile(path)
	for _, line := range data {
		for k, v := range line {
			if "nl" == k {
				fmt.Println(v)
			}
		}
	}
}

func TestAddISuffix(t *testing.T) {
	fileName := AddSuffix("test.json", 1)
	if fileName != "test_1.json" {
		t.Error()
	}

	fileName = AddSuffix("./test.json", 2)
	if fileName != "./test_2.json" {
		t.Error()
	}

	fileName = AddSuffix("test.json", "test")

	if fileName != "test_test.json" {
		t.Error()
	}
}

func TestLineCounter(t *testing.T) {
	counter, err := LineCounter(path)
	FatalCheck(err)
	if counter != 2000 {
		t.Error()
	}
}

func TestSplitJsonFile(t *testing.T) {
	num := 2
	fileNames := SplitFile(path, num)
	fmt.Println(fileNames)
}

func TestSplitJsonFile2(t *testing.T) {
	num := 2
	fileNames := SplitFile(cs_java_path, num)
	fmt.Println(fileNames)
}

func TestGetCodeFromString(t *testing.T) {
	code := GetCodeFromString(successfulContent)
	fmt.Println(code)
}

func TestConvertResponseStringToStruct(t *testing.T) {
	//response := ConvertStringToResponse(successfulContent)
	//fmt.Println(response)
}

func TestWriteToJSONFileFromString(t *testing.T) {
	//WriteToJSONFileFromString(respPath, successfulContent)
}

func TestGetMergeFileName(t *testing.T) {
	beforeMerge := []string{"test_shuffled_with_path_and_id_concode_0_response.json"}
	log.Println(GetMergeFileName(beforeMerge))
}

func TestMergeJSONFile(t *testing.T) {
	MergeJSONFile(testFilesPath)
}

func TestGetFinalFromJSONFIle(t *testing.T) {
	//GetPredictionFromJSONFIle(finalPath, predictionsPath)
	GetPredictionFromJSONFIle(code2codeResponsePath, code2codePredictionPath)
}

func TestDeleteImport(t *testing.T) {
	s := "import java.util.HashMap;\\nimport java.util.Map;\\n\\npublic class Example {\\n    public static void main(String[] args) {\\n        Map<String, String> mappings = new HashMap<>();\\n        \\n        // function parameters\\n        mappings.put(\\\"length\\\", \\\"int\\\");\\n        mappings.put(\\\"width\\\", \\\"int\\\");\\n        \\n        // function variable\\n        mappings.put(\\\"area\\\", \\\"double\\\");\\n        \\n        // print mappings\\n        System.out.println(\\\"Mappings for calculateArea function:\\\");\\n        for (String key : mappings.keySet()) {\\n            String value = mappings.get(key);\\n            System.out.println(key + \\\" : \\\" + value);\\n        }\\n    }\\n    \\n    public static double calculateArea(int length, int width) {\\n        double area = length * width;\\n        return area;\\n    }\\n}"
	lines := ModifyCodeFormat(s)
	fmt.Println(lines)
}

func TestGenerateAnswersFromJSONFile(t *testing.T) {
	GenerateAnswersFromJSONFile("../text_to_code/dataset/ref/test_shuffled_with_path_and_id_concode.json",
		"../text_to_code/dataset/ref/test.json", "../text_to_code/dataset/ref/answers.json")
}
