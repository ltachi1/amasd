// 类型别名,所有都用单字母标识，没有特殊含义，只是为了方便编写
package core

import "time"

type A map[string]interface{}
type B map[string]string
type C map[string]B
type Timestamp int

func (t Timestamp) MarshalJSON() ([]byte, error) {
	//当返回时间为空时，需特殊处理
	if t == 0 {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Unix(int64(t), 0).Format("2006-01-02 15:04:05") + `"`), nil
}
