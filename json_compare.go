/**
 * 比较两个未知结构的json
 * 参考：https://www.cnblogs.com/wangzhao765/p/9662331.html
 * json作为一个类map的结构体，它的value可能分为3类：
 * 1. json。json的值可能还是json。这就意味着，遇到了值为json的情况，我们需要进行嵌套的比较。另外一点需要注意的，是json结构体本身是无序的，所以比较过程中，要处理好这一点。
 * 2. jsonArray。json的值也有可能是jsonArray。这不仅带来了嵌套比较，还要注意，jsonArray跟json相比，它是有序的。
 * 3. 简单值。这里的简单值包括字符串，实数和布尔值。简单值只需要比较类型和值是否相同即可，也不存在嵌套的情况。
 * 思路：对于两个json结构体json1和json2，我们首先要遍历json1的键值对，检查json2是否存在对应的键值对，然后根据值的类型分别进行处理。
 */
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

/**
保存Json差异内容
*/
type JsonDiff struct {
	HasDiff bool
	Result  string
}

func marshal(j interface{}) string {
	value, _ := json.Marshal(j)
	return string(value)
}

func jsonDiffDict(json1, json2 map[string]interface{}, depth int, diff *JsonDiff) {
	blank := strings.Repeat(" ", (2 * (depth - 1)))
	for key, value := range json1 {
		quotedKey := fmt.Sprintf("\"%s\"", key)
		if strings.Contains(Settings.ignoreFields, key) {
			fmt.Println("ignoreFiled is: ", key)
			continue
		}
		if _, ok := json2[key]; ok {
			switch value.(type) {
			case map[string]interface{}:
				if _, ok2 := json2[key].(map[string]interface{}); !ok2 {
					//fmt.Println("jsonDiffDict>>map[string]interface{}..key:", key)
					diff.HasDiff = true
					diff.Result = diff.Result + "\n-" + blank + quotedKey + ": " + marshal(value) + ","
					diff.Result = diff.Result + "\n+" + blank + quotedKey + ": " + marshal(json2[key])
				} else {
					//fmt.Println("jsonDiffDict>>map[string]interface{}.else.key:",key)
					//diff.Result = diff.Result + "\n" + longBlank + quotedKey + ": "
					jsonDiffDict(value.(map[string]interface{}), json2[key].(map[string]interface{}), depth+1, diff)
				}
			case []interface{}:
				//diff.Result = diff.Result + "\n" + longBlank + quotedKey + ": "
				if _, ok2 := json2[key].([]interface{}); !ok2 {
					//fmt.Println("jsonDiffDict>>[]interface{}..key:", key)
					diff.HasDiff = true
					diff.Result = diff.Result + "\n-" + blank + quotedKey + ": " + marshal(value) + ","
					diff.Result = diff.Result + "\n+" + blank + quotedKey + ": " + marshal(json2[key])
				} else {
					//fmt.Println("jsonDiffDict>>[]interface{}.else.key:",key)
					jsonDiffList(value.([]interface{}), json2[key].([]interface{}), depth+1, diff)
				}
			default:
				if !reflect.DeepEqual(value, json2[key]) {
					//fmt.Println("jsonDiffDict>>default..key:", key)
					diff.HasDiff = true
					diff.Result = diff.Result + "\n-" + blank + quotedKey + ": " + marshal(value) + ","
					diff.Result = diff.Result + "\n+" + blank + quotedKey + ": " + marshal(json2[key])
				} else {
					//fmt.Println("jsonDiffDict>>default.else. equeal:",key)
					//diff.Result = diff.Result + "\n" + longBlank + quotedKey + ": " + marshal(value)
				}
			}
		} else {
			//fmt.Println("jsonDiffDict>>if json2[key]..else..key:", key)
			diff.HasDiff = true
			diff.Result = diff.Result + "\n-" + blank + quotedKey + ": " + marshal(value)
		}
		//diff.Result = diff.Result + ","
	}
	for key, value := range json2 {
		if _, ok := json1[key]; !ok {
			//fmt.Println("jsonDiffDict>>range json2..json1[key]..key:", key)
			diff.HasDiff = true
			diff.Result = diff.Result + "\n+" + blank + "\"" + key + "\"" + ": " + marshal(value) + ","
		}
	}
	//diff.Result = diff.Result + "\n" + blank + "}"
}

func jsonDiffList(json1, json2 []interface{}, depth int, diff *JsonDiff) {
	blank := strings.Repeat(" ", (2 * (depth - 1)))
	size := len(json1)
	if size > len(json2) {
		size = len(json2)
	}
	for i := 0; i < size; i++ {
		switch json1[i].(type) {
		case map[string]interface{}:
			if _, ok := json2[i].(map[string]interface{}); ok {
				jsonDiffDict(json1[i].(map[string]interface{}), json2[i].(map[string]interface{}), depth+1, diff)
			} else {
				//fmt.Println("jsonDiffList>>map[string]interface{}..else..")
				diff.HasDiff = true
				diff.Result = diff.Result + "\n-" + blank + marshal(json1[i]) + ","
				diff.Result = diff.Result + "\n+" + blank + marshal(json2[i])
			}
		case []interface{}:
			if _, ok2 := json2[i].([]interface{}); !ok2 {
				//fmt.Println("jsonDiffList>>[]interface{}..")
				diff.HasDiff = true
				diff.Result = diff.Result + "\n-" + blank + marshal(json1[i]) + ","
				diff.Result = diff.Result + "\n+" + blank + marshal(json2[i])
			} else {
				jsonDiffList(json1[i].([]interface{}), json2[i].([]interface{}), depth+1, diff)
			}
		default:
			if !reflect.DeepEqual(json1[i], json2[i]) {
				//fmt.Println("jsonDiffList>>default..")
				diff.HasDiff = true
				diff.Result = diff.Result + "\n-" + blank + marshal(json1[i]) + ","
				diff.Result = diff.Result + "\n+" + blank + marshal(json2[i])
			} else {
				//fmt.Println("jsonDiffList>>default.else...")
				//diff.Result = diff.Result + "\n" + longBlank + marshal(json1[i])
			}
		}
		//diff.Result = diff.Result + ","
	}
	for i := size; i < len(json1); i++ {
		//fmt.Println("jsonDiffList>>i < len(json1)...")
		diff.HasDiff = true
		diff.Result = diff.Result + "\n-" + blank + marshal(json1[i])
		diff.Result = diff.Result + ","
	}
	for i := size; i < len(json2); i++ {
		//fmt.Println("jsonDiffList>>i < len(json2)...")
		diff.HasDiff = true
		diff.Result = diff.Result + "\n+" + blank + marshal(json2[i])
		diff.Result = diff.Result + ","
	}
	//diff.Result = diff.Result + "\n" + blank + "]"
}

func loadJson(path string, dist interface{}) (err error) {
	var content []byte
	if content, err = ioutil.ReadFile(path); err == nil {
		err = json.Unmarshal(content, dist)
	}
	return err
}
