package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var filepath = flag.String("filepath","defautl", "文件路径")
var diffresult = flag.String("diffresult", "result.log", "对比结果文件")

func readfile()  {
	flag.Parse()
	fmt.Println("-filepath:", *filepath)
	file, err := os.Open(*filepath)
	if err != nil {
		fmt.Println("read file fail..", err)
		return
	}
	defer file.Close()

	r := bufio.NewReader(file)

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
	for {
		line,_, err := r.ReadLine()
		if err == io.EOF {
			fmt.Println("EOF break..")
			break
		}
		if err != nil {
			fmt.Println("readline fail..", err)
		}

		/**
		1、记录文件内容
		2、找到两次返回结果进行比较
		3、如果两次结果相同，则清空记录的内容，重新记录
		4、如果两次结果不同，则写入文件
		*/
		// 记录文件行号
		lineNo ++
		content := string(line)
		fmt.Println("")
		fmt.Println("line==>>>", content)
		fmt.Println("flagClear==>>>", flagClear)
		if flagClear {
			// 清空buffer
			buffer.Reset()
			// 重置计数器
			flagClear = false
			sourceResponse = ""
			replayResponse = ""
			writeFlag = false
			fmt.Println("<<<<<<<<<<<<<<>>>>>>>>>>>>>", flagClear,"<<>>>", buffer.Len())
		}
		buffer.WriteString(content)
		buffer.WriteString("\r\n")
		if strings.HasPrefix(content, "{") || strings.HasPrefix(content, "<") {
			if sourceResponse == "" {
				sourceResponse = content
			}else if replayResponse == "" {
				replayResponse = content
			 }
		}
		fmt.Println("sourceRes===>", sourceResponse)
		fmt.Println("replayRes==>>", replayResponse)
		if sourceResponse != "" && replayResponse != "" {
			sourceRequestCount ++
			if !strings.EqualFold(sourceResponse, replayResponse) {
				writeFlag = true
				diffResponseCount ++
			}
			flagClear = true
		}
		fmt.Println("writeFlag===>", writeFlag)
		if writeFlag {
			Writefile("soure file lineNo:"+strconv.Itoa(lineNo)+"\r\n",w)
			Writefile(buffer.String(), w)
			Writefile("\r\n",w)
			Writefile("<<<<<<<<<<<<<<<<<<++++Separator++++>>>>>>>>>>>>>>>>>>>>>>>> \r\n",w)
			Writefile("\r\n",w)
		}

		//Writefile(line, w)
		//count ++
		//if count > 103 {
		//	break
		//}
	}
	fmt.Printf("sourceRequestCount: %d <> diffResponseCount: %d ", sourceRequestCount, diffResponseCount)
	fmt.Println("")
	return
}

/**
 * 打开文件
 */
func OpenFile() (*os.File, error) {
	var f *os.File
	var err error
	if CheckFileExist(*diffresult) { // 文件存在
		fmt.Println("file is exist.")
		result := os.Truncate(*diffresult, 0)
		if result != nil {
			fmt.Println("clear file fail..", err)
			return nil, result
		}
		fmt.Println("clear file success..")
		f, err = os.OpenFile(*diffresult, os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("file open fail..", err)
			return nil, err
		}
	} else { // 文件不存在
		f, err = os.Create(*diffresult)
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
	n, errW := w.WriteString(line)
	if errW != nil {
		fmt.Println("write file fail..", errW)
	}
	fmt.Printf("写入 %d 个字节", n)
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
