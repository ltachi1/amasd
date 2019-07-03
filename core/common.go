// 通用工具类
package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
	"crypto/md5"
)

//验证给定的参数是否是url格式
func IsUrl(url string) bool {
	urlReg := regexp.MustCompile(`^https?://\w.+$`)
	return urlReg.MatchString(url)
}

//验证邮箱
func IsEmail(email string) bool {
	emailReg := regexp.MustCompile(`^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`)
	return emailReg.MatchString(email)
}

//判断元素是否在数字数组中
func InIntArray(element int, target []int) bool {
	for _, e := range target {
		if element == e {
			return true
		}
	}
	return false
}

//判断元素是否在字符串数组中
func InStringArray(element string, array []string) bool {
	for _, e := range array {
		if element == e {
			return true
		}
	}
	return false
}

//判断请求是否为ajax
func IsAjax(c *gin.Context) bool {
	if strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest" {
		return true
	}
	return false
}

//错误跳转
func Error(c *gin.Context, message gin.H) {
	c.HTML(http.StatusOK, "index/error", gin.H{"message": message["msg"]})
	c.Abort()
}

//字符串转换字节数组
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//字节数组转字符串
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//格式化时间字符串
func Time2String(timestamp Timestamp, format string) string {
	return time.Unix(int64(timestamp), 0).Format(format)
}

//格式化字符串类型的时间戳
func FormatDateByString(timestamp string, format string) string {
	r, _ := strconv.ParseInt(timestamp, 10, 64)
	return time.Unix(r, 0).Format(format)
}

//日期格式化成时间戳
func DateToTimestamp(date string) int {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", date, loc)
	return int(theTime.Unix())
}

//字符串数组转换成int类型数组
func StringArrayToInt(str []string) []int {
	intArray := make([]int, len(str))
	for i := 0; i < len(str); i++ {
		intArray[i], _ = strconv.Atoi(str[i])
	}
	return intArray
}

//计算页数
func CalculationPages(totalCount int, pageSize int) int {
	return int(math.Ceil(float64(totalCount) / float64(pageSize)))
}

//拼接批量更新sql
func JoinBatchUpdateSql(table string, fields []B, whereField string) string {
	sql := ""
	ids := make([]string, 0)
	final := map[string][]string{}
	for _, field := range fields {
		ids = append(ids, field[whereField])
		for k, v := range field {
			if k != whereField {
				if _, exists := final[k]; exists {
					final[k] = append(final[k], fmt.Sprintf("WHEN %s THEN \"%s\"", field[whereField], v))
				} else {
					final[k] = []string{
						fmt.Sprintf("WHEN %s THEN \"%s\"", field[whereField], v),
					}
				}
			}
		}
	}
	for k, v := range final {
		sql += fmt.Sprintf("%s = CASE %s\n%s\nELSE %s END ,", k, whereField, strings.Join(v, "\n"), k)
	}
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s in (%s)", table, sql[0:len(sql)-1], whereField, strings.Join(ids, ","))
}

//计算两个时间戳的时间差
func TimeDifference(startTimestamp int, endTimestamp int) string {
	difference := endTimestamp - startTimestamp
	day := difference / (24 * 3600)
	hours, seconds := 0, 0
	if difference%(24*3600) > 0 {
		hours = (difference % (24 * 3600)) / 3600
		if difference%(24*3600)%3600 > 0 {
			seconds = difference % (24 * 3600) % 3600 / 60
		}
	}
	if day == 0 {
		if hours == 0 {
			return fmt.Sprintf("%d分钟", seconds)
		}
		return fmt.Sprintf("%d小时%d分钟", hours, seconds)
	}
	return fmt.Sprintf("%d天%d小时%d分钟", day, hours, seconds)

}

func PageResponse(items interface{}, page int, pageSize int, totalCount int) A {
	return A{
		"data": items,
		"meta": gin.H{
			"page":    page,
			"total":   totalCount,
			"pages":   CalculationPages(totalCount, pageSize),
			"perpage": pageSize,
		},
	}
}

func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
