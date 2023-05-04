package chat

import (
	"fmt"
	"strings"
)

const (
	TaskPrompts                      = "task_prompts"
	GuidedPromptsWithAPIAndException = "guided_prompts_api_exception"
	//上面是大家都用的

	//下面是generation单独用的
	FullPrompts                                    = "full_prompts"
	TestPrompts                                    = "test_prompts"
	DetailedPrompts                                = "detailed_prompts"
	GuidedPromptsWithAPI                           = "guided_prompts_api"
	DetailedPromptsWithoutRemove                   = "detaileed_prompts_without_remove_statement"
	GuidedPromptsWithAPIAndExceptionAndConciseness = "guided_prompts_api_exception_conciseness"

	//下面是translation单独用的
	TaskPromptsWithBackticks                                = "task_prompts_backticks"
	TaskPromptsWithBackticksAndConciseness                  = "task_prompts_backticks_conciseness"
	TaskPromptsWithAnnotation                               = "task_prompts_annotation"
	TaskPromptsWithBackticksAndAnnotationAndAPI             = "task_prompts_backtick_annotation_api"
	TaskPromptsWithBackticksAndAnnotationAndException       = "task_prompts_backtick_annotation_exception"
	TaskPromptsWithBackticksAndAnnotationAndAPIAndException = "task_prompts_backtick_annotation_api_exception"
)

// GenerateQueryBasedPromts 根据data为代码生成功能制造完全版的prompts
func GenerateQueryBasedPromts(data map[string]any, promptMode string, line ...string) string {
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
	memberVariablesStr := strings.Join(memberVariablesSlice, ", ")

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
							tokenslice = append(tokenslice, ", ")
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
	memberFunctionsStr := strings.Join(memberFunctionsSlice, ", ")

	var res string
	//	guidelines := `When writing the method, please follow these guidelines:
	//- Remove all comments from the code.
	//- Remove all 'throws' statements.
	//- Remove all function modifiers (e.g. 'public', 'private', etc.).
	//- Change the method name to 'function'.
	//- Change the argument name to arg0, arg1, ...
	//- Change any local variable names to loc0, loc1, ...`

	guidelines1 := `When writing the method, please follow these guidelines:
- Remove all comments from the code.
- Remove all 'throws' statements.
- Remove all function modifiers (e.g. 'public', 'private', etc.).
- Change the method name to 'function'.
- Change the argument name to arg0, arg1, ...
- Change any local variable names to loc0, loc1, ...
- Return a Java method instead of a class`

	switch promptMode {
	case FullPrompts:
		{
			role := `As a senior Java developer, you'll be given information about a Java class including its name, member variables, and member function signatures.`
			addition := ` Additionally, a natural language description will be provided for a specific member function.`
			task := ` Your task is to implement this member function according to natural description within the given class.`
			rspFormat := ` Please respond with the complete member function code inside a single code block, without any explanations.`
			desc := fmt.Sprintf(` The Java class name is %s, member variables are %s, and member functions signatures are %s. The natural language description is %s.`, className, memberVariablesStr, memberFunctionsStr, nl)
			ends := ` Please provide the Java member function implementation based on this description.`

			res = role + addition + task + rspFormat + desc + ends
		}
	case TestPrompts:
		{
			//context := fmt.Sprintf(`Remember that you have a Java class named "%s", member variables "%s", member functions "%s".`, className, memberVariablesStr, memberFunctionsStr)
			//requirement := fmt.Sprintf(` Write a method named function  to "%s". `, nl)
			//res = context + requirement + guidelines1

			task := fmt.Sprintf(`Write a method named function within the %s class that %s. The class %s has member variables %s and member functions %s`, className, nl, className, memberVariablesStr, memberFunctionsStr)
			res = task + "\n" + guidelines1
		}
	case TaskPrompts:
		{
			task := fmt.Sprintf(`Write a Java method that %s`, nl)
			res = task
		}
	case DetailedPrompts:
		{
			context := fmt.Sprintf(`Remember you have a Java class named "%s", member variables "%s", member functions "%s".`, className, memberVariablesStr, memberFunctionsStr) + "\n"
			requirement := fmt.Sprintf(` Write a method named function  to "%s" `, nl)
			guidelines := `remove comments; remove summary; remove throws; remove function modifiers; change method name to "function"; change argument names to "arg0", "arg1"...; change local variable names to "loc0", "loc1"...`
			res = context + requirement + guidelines
		}
	case DetailedPromptsWithoutRemove:
		{
			context := fmt.Sprintf(`Remember you have a Java class named "%s", member variables "%s", member functions "%s".`, className, memberVariablesStr, memberFunctionsStr) + "\n"
			requirement := fmt.Sprintf(` Write a method named function  to "%s" `, nl)
			res = context + requirement
		}
	case GuidedPromptsWithAPIAndException:
		{
			context := fmt.Sprintf(`Remember you have a Java class named "%s", member variables "%s", member functions "%s".`, className, memberVariablesStr, memberFunctionsStr) + "\n"
			if strings.TrimSpace(line[0]) == "" {
				line[0] = ""
			} else {
				line[0] = fmt.Sprintf("that calls %s ", line[0])
			}
			if strings.TrimSpace(line[1]) == "true" {
				line[1] = ""
			} else {
				line[1] = "out"
			}

			requirement := fmt.Sprintf(` Write a method named function %swith%s exception handling to "%s" `, line[0], line[1], nl)
			guidelines := `remove comments; remove summary; remove throws; remove function modifiers; change method name to "function"; change argument names to "arg0", "arg1"...; change local variable names to "loc0", "loc1"...`
			res = context + requirement + guidelines
		}
	case GuidedPromptsWithAPI:
		{
			context := fmt.Sprintf(`Remember you have a Java class named "%s", member variables "%s", member functions "%s".`, className, memberVariablesStr, memberFunctionsStr) + "\n"
			if strings.TrimSpace(line[0]) == "" {
				line[0] = ""
			} else {
				line[0] = fmt.Sprintf("that calls %s ", line[0])
			}

			requirement := fmt.Sprintf(` Write a method named function %s to "%s" `, line[0], nl)
			guidelines := `remove comments; remove summary; remove throws; remove function modifiers; change method name to "function"; change argument names to "arg0", "arg1"...; change local variable names to "loc0", "loc1"...`
			res = context + requirement + guidelines
		}
	case GuidedPromptsWithAPIAndExceptionAndConciseness:
		{
			context := fmt.Sprintf(`Remember you have a Java class named "%s", member variables "%s", member functions "%s".`, className, memberVariablesStr, memberFunctionsStr) + "\n"
			if strings.TrimSpace(line[0]) == "" {
				line[0] = ""
			} else {
				line[0] = fmt.Sprintf("that calls %s ", line[0])
			}
			if strings.TrimSpace(line[1]) == "true" {
				line[1] = ""
			} else {
				line[1] = "out"
			}

			requirement := fmt.Sprintf(` Write a concise method named function %swith%s exception handling to "%s" `, line[0], line[1], nl)
			guidelines := `remove comments; remove summary; remove throws; remove function modifiers; change method name to "function"; change argument names to "arg0", "arg1"...; change local variable names to "loc0", "loc1"...`
			res = context + requirement + guidelines
		}
	default:
		{
			res = ""
		}
	}

	//return "Use temperature=0. " + res
	return res
}
