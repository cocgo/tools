// author: zhangliangzhi
// 功能：生成文件的md5值，并保存到flist
// 生成MD5 521401@qq.com
// v1 生成md5
// v2 去除.svn
// v3 改目录 2015-03-27 17:22:04
// my git home: https://github.com/cocgo

package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var needMD5 = true

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
	fmd5    string
	curName string
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

	var md5fromfile = ""
	if needMD5 {
		fbuff, _ := ioutil.ReadFile(path)
		md5fromfile = hex.EncodeToString(byte2string(md5.Sum(fbuff)))
	}

	// h := md5.New()
	// h.Write(fbuff)                                     // 需要加密的字符串为 123456
	// fmt.Printf("%s\n", hex.EncodeToString(h.Sum(nil))) // 输出加密结果
	var cname = path[len(nowpath)+1:]
	cname = strings.Replace(cname, "\\", "/", -1)

	inoFile := &sysFile{
		fName:   path,
		fType:   tp,
		fPerm:   f.Mode(),
		fMtime:  f.ModTime(),
		fSize:   f.Size(),
		fmd5:    md5fromfile,
		curName: cname,
	}
	self.files = append(self.files, inoFile)
	return nil
}

func main() {
	fmt.Println("请稍后, MakeMD5 v3.0版本 正在生成md5...")
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

	var fcount = 0
	var dircount = 0
	var allsize int64 = 0
	var dirPaths = ""
	var fileInfoList = ""
	for _, v := range self.files {
		//		fmt.Println("{name = \"", v.curName, "\", code = \"", v.fmd5, "\", size = ", v.fSize, "},")

		// 去除.svn 目录下的所有文件, svn目录下不放进来
		if strings.Contains(v.curName, ".svn") {
			continue
		}

		if v.fType == 0 {
			dirPaths += fmt.Sprint("\t\t{name = \"", v.curName, "\"},\r")
			dircount += 1
		} else {
			if v.curName != "MakeMD5.exe" && v.curName != "flist" {
				fileInfoList += fmt.Sprint("\t\t{name = \"", v.curName, "\", code = \"", v.fmd5, "\", size = ", v.fSize, "},\r")
				fcount += 1
				allsize += v.fSize
			}
		}

	}

	// 写flist.txt
	wfile := nowpath + "\\flist"
	fwrite, _ := os.OpenFile(wfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer fwrite.Close()

	fwrite.WriteString("local flist = {\r\tappVersion = 1,\r\tversion = \"1.0.2\",\r\tdirPaths = {\r")
	fwrite.WriteString(dirPaths)
	fwrite.WriteString("\t},\r\tfileInfoList = {\r")
	fwrite.WriteString(fileInfoList)
	fwrite.WriteString("\t},\r}\r\rreturn flist")

	fmt.Println("outPath:", wfile, "\nMakeMD5 v3, \nmd5值生成到lua格式的flist文件中 2015-03-27 17:23:22   zhangliangzhi, email: 521401@qq.com")
	fmt.Println("dir count:", dircount, " file count:", fcount, "allSize:", allsize)
}
