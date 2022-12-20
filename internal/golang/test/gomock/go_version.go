package gomock

import "study/internal/golang/test/gomock/spider"

// GetGoVersion 获取 Go 的最新版本
func GetGoVersion(s spider.Spider) string {
	body := s.GetBody()
	return body
}
