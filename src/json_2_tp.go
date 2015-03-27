// author: zhangliangzhi
// qq:521401
// email: 521401@qq.com
// function: json to tp, cocostudio v1.6 can use TexturePacker in lua game
// 功能：cocostudio导出的ui文件，可以转换为TexturePacker打包的图集
// update time: 2015-03-27 17:35:10
// version: V3.0
// my git home: https://github.com/cocgo

package main

import (
	"bufio"
	"container/list"
	//	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var outputFileName string = "filesName.txt"

func CheckErr(err error) {
	if nil != err {
		panic(err)
	}
}

func GetFullPath(path string) string {
	absolutePath, _ := filepath.Abs(path)
	return absolutePath
}

func PrintFilesName(path string) {
	fullPath := GetFullPath(path)

	listStr := list.New()

	filepath.Walk(fullPath, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		name := fi.Name()
		if outputFileName != name {
			listStr.PushBack(name)
		}

		return nil
	})

	OutputFilesName(listStr)
}

func ConvertToSlice(listStr *list.List) []string {
	sli := []string{}
	for el := listStr.Front(); nil != el; el = el.Next() {
		sli = append(sli, el.Value.(string))
	}

	return sli
}

func OutputFilesName(listStr *list.List) {
	files := ConvertToSlice(listStr)
	//sort.StringSlice(files).Sort()// sort

	// f, err := os.Create(outputFileName)
	// CheckErr(err)
	// defer f.Close()

	// f.WriteString("\xEF\xBB\xBF")
	// writer := csv.NewWriter(f)

	length := len(files)
	for i := 0; i < length; i++ {
		JsonChange(files[i])
	}

	//	writer.Flush()
}

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path, _ = os.Getwd()
	}
	PrintFilesName(path)

	fmt.Println("change json file!")
	fmt.Println("zhangliangzhi，2015-03-27 17:40:26  v3.0 email: 521401@qq.com")
}

func JsonChange(fileName string) {
	if !strings.Contains(fileName, ".json") {
		return
	}
	if strings.Contains(fileName, ".json.tp") {
		return
	}
	fmt.Println(fileName)

	// 读取文件
	fRead, err := os.Open(fileName)
	if err != nil {
		return
	}

	buff := bufio.NewReader(fRead)
	defer fRead.Close()

	// 写文件
	allstring := " "
	writefile := fileName + ".tp"
	fWrite, _ := os.OpenFile(writefile, os.O_CREATE|os.O_TRUNC, 0777)
	defer fWrite.Close()
	for {
		line, err := buff.ReadString('\n')

		// 退出
		if err != nil {
			allstring += line
			break
		}
		if io.EOF == err {

			break
		}
		// 字符串处理
		newLine := LineString(line)
		//fWrite.WriteString(newLine)
		allstring += newLine
	}

	// 去除空格和换行
	allstring = strings.Replace(allstring, " ", "", -1)
	allstring = strings.Replace(allstring, "\n", "", -1)

	// 换字体 "微软雅黑" -> "fonts/d.ttf"
	allstring = strings.Replace(allstring, "微软雅黑", "fonts/d.ttf", -1)

	fWrite.WriteString(allstring)
}

// 处理文字
func LineString(line string) string {
	if len(line) <= 6 {
		return line
	}
	// 1. 行中如果没有 .png 的时候 直接返回
	if !strings.Contains(line, ".png\"") {
		return line
	}

	newline := line
	yhcount := strings.Count(line, "\"")
	ds2Pos := 0
	// 找到倒数第二个 “
	for i := 1; i < yhcount; i++ {
		ygPos := strings.Index(newline, "\"")
		newline = newline[ygPos+1:]
		ds2Pos += ygPos + 1
	}

	// 2. 行中没有 / 的目录的时候
	if !strings.Contains(line, "/") {
		newline = line[:ds2Pos] + "#" + line[ds2Pos:]
		//fmt.Println(newline)
	} else {

		// 3.有相对路径 / 的时候, 要去除
		newline = line
		xgcount := strings.Count(line, "/")
		xg1Pos := 0
		for i := 1; i <= xgcount; i++ {
			xgPos := strings.Index(newline, "/")
			newline = newline[xgPos+1:]
			xg1Pos += xgPos + 1
		}

		newline = line[:ds2Pos] + "#" + line[xg1Pos:]
	}

	return newline
}

func read1(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	chunks := make([]byte, 1024, 1024)
	buf := make([]byte, 1024)
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		chunks = append(chunks, buf[:n]...)
		fmt.Println(string(buf[:n]))
	}
	return string(chunks)
}
