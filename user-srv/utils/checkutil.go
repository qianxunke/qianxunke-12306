package checkutil

import (
	"regexp"
)

const (
	regular = "^(13[0-9]|14[57]|15[0-35-9]|18[07-9])\\\\d{8}$"
)

func ValiTephone(mobileNum string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}
