// 检查png重名 521401@qq.com
// Author: 张良志
// qq:521401
// update time: 2015-03-27 17:53:37
// version: v2.0
// 功能：检查png重名
// function: check png is rename
// my git home: https://github.com/cocgo

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	IsDirectory = iota
	IsRegular
	IsSymlink
)

type sysFile struct {
	fType   int
	fName   string
	fLink   string
	fSize   int64
	fMtime  time.Time
	fPerm   os.FileMode
	curName string
	name    string
}

type F struct {
	files []*sysFile
}

func GetFullPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	line := path

	newline := line
	yhcount := strings.Count(line, "\\")
	ds1Pos := 0
	// 找到倒数第一个 \
	for i := 1; i <= yhcount; i++ {
		ygPos := strings.Index(newline, "\\")
		newline = newline[ygPos+1:]
		ds1Pos += ygPos + 1
	}
	line = line[:ds1Pos-1]
	return line
}

var nowpath = GetFullPath()

func byte2string(in [16]byte) []byte {
	tmp := make([]byte, 16)
	for _, value := range in {
		tmp = append(tmp, value)
	}

	return tmp[16:]
}

func (self *F) visit(path string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if path == nowpath {
		return nil
	}
	var tp int
	if f.IsDir() {
		tp = IsDirectory
	} else if (f.Mode() & os.ModeSymlink) > 0 {
		tp = IsSymlink
	} else {
		tp = IsRegular
	}

	var cname = path[len(nowpath)+1:]
	cname = strings.Replace(cname, "\\", "/", -1)

	// 没有路径，只有文件名
	filename := cname

	// 找到倒数第1个 /
	xgcount := strings.Count(filename, "/")
	for i := 1; i <= xgcount; i++ {
		xgPos := strings.Index(filename, "/")
		filename = filename[xgPos+1:]
	}

	inoFile := &sysFile{
		fName:   path, // 绝对路径
		fType:   tp,
		fPerm:   f.Mode(),
		fMtime:  f.ModTime(),
		fSize:   f.Size(),
		curName: cname,    // 相对路径文件名
		name:    filename, // 文件名
	}
	self.files = append(self.files, inoFile)
	return nil
}

// 检查根目录是否有重名
func checkReName() {

}

func main() {
	fmt.Println("重名检查中...\n")
	root := GetFullPath()
	self := F{
		files: make([]*sysFile, 0),
	}
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		return self.visit(path, f, err)
	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

	//var mapname = map[string]string
	mapname := make(map[string]string)

	var cmcount = 0

	var dircount = 0
	var dirPaths = ""
	//	var fileInfoList = ""
	for _, v := range self.files {
		//		fmt.Println("{name = \"", v.curName, "\", code = \"", v.fmd5, "\", size = ", v.fSize, "},")

		// 去除.svn 目录下的所有文件, svn目录下不放进来
		if strings.Contains(v.curName, ".svn") {
			continue
		}
		if strings.Contains(v.curName, ".git") {
			continue
		}

		//fmt.Println(v.fName)
		if v.fType == 0 {
			dirPaths += fmt.Sprint("\t\t{name = \"", v.curName, "\"},\r")
			dircount += 1
		} else {
			name, ok := mapname[v.name]
			if ok {
				fmt.Println(v.curName, " | ", name)
				cmcount++
			}
			mapname[v.name] = v.curName
		}

	}

	fmt.Println("该目录下一共有", cmcount, "个重名文件")
	fmt.Println("\n\n重名检查工具 v1.0. \nzhangliangzhi  email: 521401@qq.com 2015-03-27 17:47:35 \n\n请按任意键继续. . .")

	pause()
}
func pause() {
	anystring := ""
	_, err1 := fmt.Scanln(&anystring)
	if nil == err1 {
		os.Exit(0)
	}
}
