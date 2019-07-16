package ordermap

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
)

/*
  有序map
  注意: 如果需要调用json.Marshal方法将对象转为json字节，则只支持下标为 number string
*/

// OrderMap 有序map
type OrderMap struct {
	less OrderMapKeyLess
	data *sync.Map
}

// OrderMapKeyLess 为排序key时的比较函数 - 和系统sort包Less功能类似，这里i和j为map下标
type OrderMapKeyLess func(i, j interface{}) bool

var (
	// DefaultOrderMapKeySort 默认key排序比较 - 转为字符串比较 - 字符串正序排序（从小到大）
	DefaultOrderMapKeySort = func(i, j interface{}) bool {
		return fmt.Sprint(i) < fmt.Sprint(j)
	}
)

// NewOrderMap 创建一个有序map对象
func NewOrderMap(less OrderMapKeyLess) *OrderMap {
	if less == nil {
		less = DefaultOrderMapKeySort
	}
	return &OrderMap{
		data: new(sync.Map),
		less: less,
	}
}

// LoadStringMap 加载一个下标为字符串的map
func (om *OrderMap) LoadStringMap(m map[string]interface{}) (err error) {
	if m == nil {
		err = errors.New("加载的map不能为nil")
		return
	}
	om.data = new(sync.Map)
	for key, val := range m {
		om.data.Store(key, val)
	}
	return
}

// LoadInt64Map 加载一个下标为int64的map - 其它数值类型，可自行转换后再调用此函数
func (om *OrderMap) LoadInt64Map(m map[int64]interface{}) (err error) {
	if m == nil {
		err = errors.New("加载的map不能为nil")
		return
	}
	om.data = new(sync.Map)
	for key, val := range m {
		om.data.Store(key, val)
	}
	return
}

// Set 设置map值 - 可设置多个值
func (om *OrderMap) Set(kv ...interface{}) (err error) {
	l := len(kv)
	if l%2 != 0 {
		return errors.New("kv必须成对出现")
	}
	for i := 0; i < l; i += 2 {
		om.data.Store(kv[i], kv[i+1])
	}
	return
}

// Get 获取一个key的值
func (om *OrderMap) Get(k interface{}) (v interface{}, err error) {
	var ok bool
	v, ok = om.data.Load(k)
	if ok == false {
		return v, errors.New("不存在对应key")
	}
	return
}

// Del 删除一个key
func (om *OrderMap) Del(k interface{}) {
	om.data.Delete(k)
}

// Range 循环遍历 - 有序的 - 和sync.Map的Range相同
func (om *OrderMap) Range(f func(key, value interface{}) bool) {
	// 取出所有key
	keys := om.Keys()
	for _, key := range keys {
		val, _ := om.data.Load(key)
		isStop := f(key, val)
		// 返回false终止迭代
		if isStop == false {
			return
		}
	}
}

// Len 获取map长度
func (om *OrderMap) Len() int {
	count := 0
	om.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// Keys 获取所有key - 有序的
func (om *OrderMap) Keys() (keys []interface{}) {
	// 取出所有key
	keys = make([]interface{}, 0)
	om.data.Range(func(key, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	// key 排序
	if om.less == nil {
		om.less = DefaultOrderMapKeySort
	}
	sort.Slice(keys, func(i, j int) bool {
		return om.less(keys[i], keys[j])
	})
	return
}

// MarshalJSON 将对象转为json字节
func (om *OrderMap) MarshalJSON() ([]byte, error) {
	if om == nil {
		return []byte("null"), nil
	}
	// 取出所有key
	keys := om.Keys()
	// 拼接json字符串
	body := "{"
	for _, key := range keys {
		val, ok := om.data.Load(key)
		if ok == false {
			val = nil
		}
		valByte, err := json.Marshal(val)
		if err != nil {
			return []byte("null"), err
		}
		body += fmt.Sprintf(`"%v":%s,`, key, string(valByte))
	}
	body = strings.TrimRight(body, ",")
	body += "}"

	return []byte(body), nil
}

// UnmarshalJSON 将json字节转为结构体
func (om *OrderMap) UnmarshalJSON(data []byte) (err error) {
	// data设置为新对象
	om.data = new(sync.Map)
	// 尝试转换为字符串下标map
	ms := make(map[string]interface{}, 0)
	err = json.Unmarshal(data, &ms)
	if err != nil {
		goto NUMBER
	}
	// 加载值
	err = om.LoadStringMap(ms)
	if err != nil {
		return
	}
	return
NUMBER:
	mi := make(map[int64]interface{}, 0)
	err = json.Unmarshal(data, &mi)
	if err != nil {
		err = errors.New("尝试将json字符串转为map失败，此json字符串既不是string下标也不是int64下标")
		return
	}
	// 加载值
	err = om.LoadInt64Map(mi)
	if err != nil {
		return
	}
	return
}
