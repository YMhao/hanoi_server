package user_api

import (
	"testing"
)

func TestFormat(t *testing.T) {
	testGroup := []string{
		"xxx@qq.com",
		"xxxsdfkm",
		"@outlook.com",
		"13800138000",
	}

	for _, str := range testGroup {
		t.Log(str, getUserNameType(str))
	}

	md5Group := []string{
		"f75a3761d76e84a775a63737c4a29f39",
		"ac59075b964b0715",
	}

	for _, str := range md5Group {
		t.Log(str, checkPasswdFormat(str))
	}

}
