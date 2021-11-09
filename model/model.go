package model

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type BaseMysqlStruct struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT;comment:'主键'"`
	CreateTime LocalTime `json:"create_time" gorm:"NOT NULL;comment:'创建时间'"`
	UpdateTime LocalTime `json:"update_time" gorm:"NOT NULL;comment:'更新时间'"`
	IsDelete   bool      `json:"is_delete" gorm:"DEFAULT:'0';comment:'是否删除 0 未删除 1删除'"`
}

type LocalTime struct {
	time.Time
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	//格式化秒
	seconds := t.Unix()
	return []byte(strconv.FormatInt(seconds, 10)), nil
}
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t LocalTime) Now() LocalTime {
	return LocalTime{Time: time.Now()}
}
