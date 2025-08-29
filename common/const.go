package common

import "fmt"

func Recover() {
	if r := recover(); r != nil {
		fmt.Println("Recovered: ", r)
	}
}
