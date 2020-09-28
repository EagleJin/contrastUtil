/**
 * 注意：
 * 1、bufio.ReadLine() 读满缓冲区就返回，剩下的字节不会丢弃，留着下次读取.如果一行太长，会被截断
 * 2、Scanner在初始化的时候有设置一个maxTokenSize，这个值默认是MaxScanTokenSize = 64 * 1024 ，当一行的长度大于64*1024即65536之后，
 * 就会出现ErrTooLong错误,当遇到错误时，scanner就会自动退出
 * 3、通过设置Scanner buffer更改默认最大限制，可以解决超过默认限制，报：error: bufio.Scanner: token too long 错误，导致退出问题
 * 	切记：设置buffer一定在Scan 之前，不然会不工作，Buffer() 方法上注释有说明
 * https://blog.csdn.net/tianlongtc/article/details/80148509
 */
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//var filepath = flag.String("filepath", "defautl", "文件路径")
//var diffresult = flag.String("diffresult", "result.log", "对比结果文件")
//var ignoreFileds = flag.String("ignore-fields","","忽略字段")
const (
	HTTP_STATUS_304 = "304"
	MAX_CAPACITY    = 1024 * 1024
)

func readfile() {
	startTime := time.Now()
	flag.Parse()
	//fmt.Println("-filepath:", *filepath)
	file, err := os.Open(Settings.filePath)
	if err != nil {
		fmt.Println("read file fail..", err)
		return
	}
	defer file.Close()

	//r := bufio.NewReader(file)

	r := bufio.NewScanner(file)
	buf := make([]byte, MAX_CAPACITY)
	r.Buffer(buf, MAX_CAPACITY)

	// 打开文件，重复多次写入内容，defer关闭文件
	f, err := OpenFile()
	if err != nil {
		fmt.Println("open file fail..", err)
		return
	}
	// 创建新的Write对象
	w := bufio.NewWriter(f)

	//count := 0
	var buffer bytes.Buffer
	flagClear := false
	sourceResponse := ""
	replayResponse := ""
	writeFlag := false
	lineNo := 0
	// 记录总请求数量
	sourceRequestCount := 0
	// 记录返回结果差异的请求数量
	diffResponseCount := 0
	responseFlag := false
	replayResFlag := false
	for r.Scan() {
		//line, err := r.ReadString('\n')
		/*if err == io.EOF {
			fmt.Println("EOF break..")
			break
		}*/
		/*if err != nil {
			if err == io.EOF {

			}
			panic(err)
			fmt.Println("readline fail..", err)
		}*/

		/**
		1、记录文件内容
		2、找到两次返回结果进行比较
		3、如果两次结果相同，则清空记录的内容，重新记录
		4、如果两次结果不同，则写入文件
		*/
		// 记录文件行号
		lineNo++
		content := r.Text()
		//fmt.Println("")
		//fmt.Println("line==>>>", content)
		//fmt.Println("flagClear==>>>", flagClear)

		// 1 代表请求，说明是一个新的请求，要把原来设置的 返回和回放的相关变量置空。防止个别请求没有返回，导致对比结果错误
		if strings.HasPrefix(content, "1 ") {
			// 代表一个新的请求，重置状态
			flagClear = true
		}

		if flagClear {
			// 清空buffer
			buffer.Reset()
			flagClear = false
			sourceResponse = ""
			replayResponse = ""
			writeFlag = false
			responseFlag = false
			replayResFlag = false
			//fmt.Println("<<<<<<<<<<<<<<>>>>>>>>>>>>>", flagClear,"<<>>>", buffer.Len())
		}
		buffer.WriteString(content)
		buffer.WriteString("\r\n")

		// 增加对原始响应包的判断，只获取响应结果（个别POST请求接口会传json参数，要过滤掉）
		if strings.HasPrefix(content, "2 ") {
			responseFlag = true
		} else if strings.HasPrefix(content, "3 ") {
			replayResFlag = true
		}
		if responseFlag && !replayResFlag {
			if strings.HasPrefix(content, "HTTP/1.1") && strings.Contains(content, HTTP_STATUS_304) {
				sourceResponse = HTTP_STATUS_304
			} else if strings.HasPrefix(content, "{") || strings.HasPrefix(content, "<") {
				sourceResponse = content
			}
		}
		if replayResFlag {
			if strings.HasPrefix(content, "HTTP/1.1") && strings.Contains(content, HTTP_STATUS_304) {
				replayResponse = HTTP_STATUS_304
			} else if strings.HasPrefix(content, "{") || strings.HasPrefix(content, "<") {
				replayResponse = content
			}
		}
		if sourceResponse != "" && replayResponse != "" {
			//fmt.Println("sourceR===>",sourceResponse)
			//fmt.Println("replayR===>",replayResponse)
			sourceRequestCount++
			// 非json结构直接对比结果字符串
			if strings.HasPrefix(sourceResponse, "<") && strings.HasPrefix(replayResponse, "<") {
				fmt.Println("not json...")
				if !strings.EqualFold(sourceResponse, replayResponse) {
					writeFlag = true
					diffResponseCount++
				}
			} else if strings.EqualFold(sourceResponse, HTTP_STATUS_304) && strings.EqualFold(replayResponse, HTTP_STATUS_304) {
				// 都是304，代表结果一致，进入下一次循环
				fmt.Println("all 304..")
				//fmt.Println("sourceResponse==>", sourceResponse)
				//fmt.Println("replayResponse==>", replayResponse)
				continue
			} else if strings.EqualFold(sourceResponse, HTTP_STATUS_304) {
				fmt.Println("source 304...")
				//fmt.Println("sourceResponse==>", sourceResponse)
				//fmt.Println("replayResponse==>", replayResponse)
				writeFlag = true
				diffResponseCount++
				buffer.WriteString("\r\n")
				buffer.WriteString("<<< 差异结果： >>>\r\n")
				buffer.WriteString("-HTTP" + HTTP_STATUS_304 + "\r\n")
				buffer.WriteString("+" + replayResponse)
				buffer.WriteString("\r\n")
			} else if strings.EqualFold(replayResponse, HTTP_STATUS_304) {
				fmt.Println("replay 304...")
				//fmt.Println("sourceResponse==>", sourceResponse)
				//fmt.Println("replayResponse==>", replayResponse)
				writeFlag = true
				diffResponseCount++
				buffer.WriteString("\r\n")
				buffer.WriteString("<<< 差异结果： >>>\r\n")
				buffer.WriteString("-" + sourceResponse + "\r\n")
				buffer.WriteString("+HTTP" + HTTP_STATUS_304)
				buffer.WriteString("\r\n")
			} else {
				fmt.Println("all not 304...")
				//fmt.Println("sourceResponse==>", sourceResponse)
				//fmt.Println("replayResponse==>", replayResponse)
				// 解析json结果，逐个属性比对
				var (
					sourceJson    map[string]interface{}
					replayJson    map[string]interface{}
					compareResult = &JsonDiff{HasDiff: false, Result: ""}
				)
				sourceErr := json.Unmarshal([]byte(sourceResponse), &sourceJson)
				if sourceErr != nil {
					fmt.Println("sourceResponse to json error.", sourceErr)
					continue
				}
				replayErr := json.Unmarshal([]byte(replayResponse), &replayJson)
				if replayErr != nil {
					fmt.Println("replayResponse to json error.", replayErr)
					continue
				}
				jsonDiffDict(sourceJson, replayJson, 1, compareResult)
				if compareResult.HasDiff {
					writeFlag = true
					diffResponseCount++

					buffer.WriteString("\r\n")
					buffer.WriteString("<<< 结果中差异的字段是： >>>")
					buffer.WriteString(compareResult.Result)
					buffer.WriteString("\r\n")
					//fmt.Println("<<<compareResult>>>", compareResult.Result)
				}
			}
			flagClear = true
		}
		//fmt.Println("writeFlag===>", writeFlag)
		if writeFlag {
			Writefile("soure file lineNo:"+strconv.Itoa(lineNo)+"；代表原文件中该行上面第一组请求 \r\n", w)
			Writefile(buffer.String(), w)
			Writefile("\r\n", w)
			Writefile("<<<<<<<<<<<<<<<<<<++++Separator++++>>>>>>>>>>>>>>>>>>>>>>>> \r\n", w)
			Writefile("\r\n", w)
		}
	}
	if r.Err() != nil {
		fmt.Printf("error: %s\n", r.Err())
	}
	fmt.Printf("sourceRequestCount: %d <> diffResponseCount: %d ==>time cost: %v\n", sourceRequestCount, diffResponseCount, time.Since(startTime))
	fmt.Println("")
	fmt.Println("lineNo--->", lineNo)
	return
}

/**
 * 打开文件
 */
func OpenFile() (*os.File, error) {
	var f *os.File
	var err error
	if CheckFileExist(Settings.diffResult) { // 文件存在
		//fmt.Println("file is exist.")
		result := os.Truncate(Settings.diffResult, 0)
		if result != nil {
			fmt.Println("clear file fail..", err)
			return nil, result
		}
		//fmt.Println("clear file success..")
		f, err = os.OpenFile(Settings.diffResult, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("file open fail..", err)
			return nil, err
		}
	} else { // 文件不存在
		f, err = os.Create(Settings.diffResult)
		if err != nil {
			fmt.Println("file create fail..", err)
			return nil, err
		}
	}
	return f, nil
}

/**
 * 写文件
 */
func Writefile(line string, w *bufio.Writer) {
	// 写文件
	_, errW := w.WriteString(line)
	if errW != nil {
		fmt.Println("write file fail..", errW)
	}
	//fmt.Printf("写入 %d 个字节", n)
	w.Flush()
}

func CheckFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
func main() {
	readfile()
	return
}
