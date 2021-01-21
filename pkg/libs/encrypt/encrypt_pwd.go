package encrypt

import (
	"crypto/md5"
	"fmt"
	"io"
)

//指定两个 salt： salt1 = *~^+%)=-!@&$%   salt2 = "TroilaaliorT"
const (
	salt1 = "*~^+%)=-!@&$%"
	salt2 = "TroilaaliorT"
)

func EncryptPwd(username, pwd string) string {
	// 使用MD5加密密码
	h := md5.New()
	io.WriteString(h, pwd)

	//pwmd5等于e10adc3949ba59abbe56e057f20f883e
	pwmd5 := fmt.Sprintf("%x", h.Sum(nil))

	//salt1+用户名+salt2+MD5拼接
	io.WriteString(h, salt1)
	io.WriteString(h, username)
	io.WriteString(h, salt2)
	io.WriteString(h, pwmd5)

	// 最终加密生成的密码
	lastpwd := fmt.Sprintf("%x", h.Sum(nil))

	return lastpwd
}
