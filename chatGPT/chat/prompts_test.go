package chat

import (
	"fmt"
	"testing"
)

func TestTaskPrompts(t *testing.T) {
	nl := "task"
	task := fmt.Sprintf(`Write a method named function that %s.`, nl)
	guidelines :=
		`When writing the method, please follow these guidelines:
- Remove all comments from the code.
- Remove any 'throws' statements.
- Remove any function modifiers (e.g. 'public', 'private', etc.).
- Change the method name to 'function'.
- Change the argument name to arg0, arg1, ...
- Change any local variable names to loc0, loc1, ...`
	res := task + "\n" + guidelines
	println(res)
}
