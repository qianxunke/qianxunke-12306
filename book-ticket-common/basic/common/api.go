package common

import "strings"

var WhileList *UrlWhileList

//api名单
type UrlWhileList struct {
	Url_whitelist []string `json:"url_whitelist"`
}

//判断是否在白名单内
func (whileList *UrlWhileList) IsInWileList(inUrl string) bool {
	if len(inUrl) == 0 {
		return false //应该没有这种情况
	}
	for _, whileUrl := range whileList.Url_whitelist {
		if strings.Compare(inUrl, whileUrl) == 0 {
			return true
		}
	}
	return false
}
