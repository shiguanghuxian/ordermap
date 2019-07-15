# OrderMap
go 有序map，用于json输出有序key的对象和后端有序循环map取值

## 测试

### 后台循环

```

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

=== RUN   TestOrderMap_Range
--- PASS: TestOrderMap_Range (0.00s)
    map_test.go:92: 1:ha1 
    map_test.go:92: 2:ha2 
    map_test.go:92: 3:ha3 
    map_test.go:92: 10:h10 
    map_test.go:96: map长度  4
PASS
ok      command-line-arguments  0.005s
```

### JSON序列化

```

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


=== RUN   TestOrderMap_OrderMap
--- PASS: TestOrderMap_OrderMap (0.00s)
    map_test.go:34: <nil>
    map_test.go:36: <nil>
    map_test.go:41: h10 <nil>
    map_test.go:48: {"1":"这是设置后的值","3":"ha3","8":"ha8","10":"h10"}
```