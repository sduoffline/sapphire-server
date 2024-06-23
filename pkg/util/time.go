package util

import (
	"log/slog"
	"strconv"
	"time"
)

func FormatTimeStr(time time.Time) string {
	sub := time.Sub(time.UTC())
	slog.Info("sub: ", sub)

	// 超过一周
	// 返回"x月x日"
	if sub.Hours() > 24*7 {
		return time.Format("01-02")
	}

	// 超过一天
	// 返回"几天前"
	if sub.Hours() > 24 {
		days := int(sub.Hours() / 24)
		return strconv.Itoa(days) + "天前"
	}

	// 超过一小时
	// 返回"几小时前"
	if sub.Hours() > 1 {
		hours := int(sub.Hours())
		return strconv.Itoa(hours) + "小时前"
	}

	// 超过一分钟
	// 返回"几分钟前"
	if sub.Minutes() > 1 {
		minutes := int(sub.Minutes())
		return strconv.Itoa(minutes) + "分钟前"
	}

	// 返回"刚刚"
	return "刚刚"
}
