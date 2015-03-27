// author: zhangliangzhi
// 功能：异或加密结果
// xor Result Verification
// email: 521401@qq.com
// time: 2015-03-27 17:36:27
// my git home: https://github.com/cocgo

package main

import (
	"fmt"
)

func myxor(str string) string {
	var retstr string
	for i := 0; i < len(str); i++ {
		onechar := str[i] ^ 'z'
		retstr = retstr + string(onechar)
	}

	return retstr
}

func main() {
	var cmdid uint16 = 41

	str := string(cmdid)
	retstr := myxor(str)

	fmt.Println(retstr)
	fmt.Println(myxor(retstr))
}
