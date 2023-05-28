package main

import (
	"clifile/internal/util/file"
	"fmt"
)

func main() {
	ch, err := file.ReadRunes(".gitignore")
	if err == nil {
		for run := range ch {
			fmt.Print(string(run))
		}
	} else {
		fmt.Printf("err: %v\n", err)
	}
}
