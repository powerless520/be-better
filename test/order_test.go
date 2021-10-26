package test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestDemo(t *testing.T) {
	fmt.Println(GetTableNameFromOrderId("10160062026159d37fb4d218687063431633276799"))
	fmt.Println("")
	fmt.Println(GetTableNameFromOrderId(""))
}

// GetTableNameFromOrderId 根据订单号计算order表名
func GetTableNameFromOrderId(order_id string) string {
	unixtime_str_int64 := time.Now().Unix()
	var err error
	if err = VerifyValidOrderId(order_id); err == nil {
		unixtime_str := order_id[len(order_id)-10 : len(order_id)]
		unixtime_str_int64, err = strconv.ParseInt(unixtime_str, 10, 64)
		if err != nil {
			unixtime_str_int64 = time.Now().Unix()
		}
	}
	unixtime := time.Unix(unixtime_str_int64, 0)
	year, week := unixtime.ISOWeek()
	year_week := strconv.Itoa(year*100 + week)
	return "t_order_" + year_week
}

func VerifyValidOrderId(order_id string) (err error) {
	if order_id == "" {
		return errors.New("订单号为空")
	}
	if len(order_id) < 10 {
		return errors.New("订单号错误")
	}
	return nil
}
