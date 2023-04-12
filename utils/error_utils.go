package utils

import "fmt"

func Check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func FatalCheck(e error) {
	if e != nil {
		panic(e)
	}
}
