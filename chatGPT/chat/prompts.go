package chat

import (
	"fmt"
	"strings"
)

const (
	FULL_PROMPTS int = 1
)

// GenerateQueryBasedPromts 根据data为代码生成功能制造完全版的prompts
func GenerateQueryBasedPromts(data map[string]any) string {
	className := data["className"].(string)
	memberVariablesMap := data["memberVariables"].(map[string]any)
	memberFunctionsMap := data["memberFunctions"].(map[string]any)
	nl := data["nl"].(string)

	var memberVariablesSlice []string
	for k, v := range memberVariablesMap {
		vStr := v.(string)
		kStr := k
		variable := vStr + " " + kStr
		memberVariablesSlice = append(memberVariablesSlice, variable)
	}
	memberVariablesStr := strings.Join(memberVariablesSlice, ",")

	var memberFunctionsSlice []string
	var tokenslice []string
	for key, value := range memberFunctionsMap {
		tokenslice = nil
		//fmt.Println("Key:", key)
		// 使用类型断言判断内层 value 是否为 []interface{}
		if innerMap, ok := value.([]interface{}); ok {
			// 遍历内层 []interface{}
			for _, innerValue := range innerMap {
				// 使用类型断言判断内层内层 value 是否为 []interface{}
				if innerInnerMap, ok := innerValue.([]interface{}); ok {
					// 遍历内层内层 []interface{}
					for i, innerInnerValue := range innerInnerMap {
						v := innerInnerValue.(string)
						tokenslice = append(tokenslice, v)
						if i == 0 {
							tokenslice = append(tokenslice, key)
							tokenslice = append(tokenslice, "(")
						}
						if i != len(innerInnerMap)-1 && i != 0 {
							tokenslice = append(tokenslice, ",")
						}
						//fmt.Println(innerInnerValue)
					}
					tokenslice = append(tokenslice, ")")
				}
			}
			functionHead := strings.Join(tokenslice, " ")
			memberFunctionsSlice = append(memberFunctionsSlice, functionHead)
		}
	}
	memberFunctionsStr := strings.Join(memberFunctionsSlice, ",")

	role := `As a senior Java developer, you'll be given information about a Java class including its name, member variables, and member function headers.`
	addition := ` Additionally, a natural language description will be provided for a specific member function.`
	task := ` Your task is to implement this member function within the given class.`
	rspFormat := ` Please respond with the complete code inside a single code block, without any explanations.`

	desc := fmt.Sprintf(` The Java class name is %s, member variables are %s, and member functions headers are %s. The natural language description is %s.`, className, memberVariablesStr, memberFunctionsStr, nl)
	ends := ` Please provide the Java member function implementation based on this description.`

	return role + addition + task + rspFormat + desc + ends
}
