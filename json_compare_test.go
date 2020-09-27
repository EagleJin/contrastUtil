package main

import (
	"fmt"
	"testing"
)

func TestJsonDiffDict(t *testing.T) {
	fmt.Println(">>>>>>>>>")
	var (
		json1 map[string]interface{}
		json2 map[string]interface{}
		diff  = &JsonDiff{HasDiff: false, Result: ""}
		//filed = flag.String("b","123","dddd")
	)
	if err := loadJson("./1.json", &json1); err != nil {
		fmt.Println(err)
	}
	if err := loadJson("./2.json", &json2); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("diff0:%p\n", &diff)
	jsonDiffDict(json1, json2, 1, diff)
	fmt.Println("===================================")
	fmt.Println(diff.HasDiff, diff.Result)
	fmt.Println("<<<<<<<<<<<")
}
