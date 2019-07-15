package ordermap

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

// go test -v -count=1 map_test.go map.go

func TestOrderMap_OrderMap(t *testing.T) {
	// 要加载的map
	data := make(map[string]interface{}, 0)
	data["1"] = "ha1"
	data["2"] = "ha2"
	data["10"] = "h10"
	data["3"] = "ha3"
	// key转数字，从小到大排序
	less := func(i, j interface{}) bool {
		ii, _ := strconv.Atoi(fmt.Sprint(i))
		jj, _ := strconv.Atoi(fmt.Sprint(j))
		return ii < jj
	}
	// 有序map对象
	om := NewOrderMap(less)
	err := om.LoadStringMap(data)
	if err != nil {
		t.Error(err)
		return
	}
	// 设置一个值
	err = om.Set("1", "这是设置后的值")
	t.Log(err)
	err = om.Set("8", "ha8")
	t.Log(err)
	// 删除一个值
	om.Del("2")
	// 获取一个值
	v, err := om.Get("10")
	t.Log(v, err)

	js, err := json.Marshal(om)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(js))
}

func TestOrderMap_OrderMapUnmarshal(t *testing.T) {
	str := `{"1":"这是设置后的值","3":"ha3","8":"ha8","10":"h10"}`
	om := NewOrderMap(nil)
	err := json.Unmarshal([]byte(str), om)
	if err != nil {
		t.Error(err)
		return
	}
	// 获取一个值
	v, err := om.Get("8")
	t.Log(v, err)

	js, err := json.Marshal(om)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(js))
}

func TestOrderMap_Range(t *testing.T) {
	// 要加载的map
	data := make(map[string]interface{}, 0)
	data["1"] = "ha1"
	data["2"] = "ha2"
	data["10"] = "h10"
	data["3"] = "ha3"
	// key转数字，从小到大排序
	less := func(i, j interface{}) bool {
		ii, _ := strconv.Atoi(fmt.Sprint(i))
		jj, _ := strconv.Atoi(fmt.Sprint(j))
		return ii < jj
	}
	// 有序map对象
	om := NewOrderMap(less)
	err := om.LoadStringMap(data)
	if err != nil {
		t.Error(err)
		return
	}
	om.Range(func(k, v interface{}) bool {
		t.Logf("%v:%v \n", k, v)
		return true
	})

	t.Log("map长度 ", om.Len())
}
