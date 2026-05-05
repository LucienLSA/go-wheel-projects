package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	// 统计每个key调用「回调函数」的次数，验证缓存是否生效
	loadCounts := make(map[string]int, len(db))
	// 核心：创建缓存组，使用「函数类型适配接口」技巧
	// GetterFunc 是函数类型，将匿名函数强转为 Getter 接口
	gee := NewGroup("scores", 2<<10, GetterFunc(
		// 匿名函数 = 缓存未命中时的「回调函数」(慢速数据源)
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			// 从模拟db中查找数据
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			// key不存在，返回错误
			return nil, fmt.Errorf("%s not exist", key)
		}))

	// 遍历模拟数据源，测试缓存
	for k, v := range db {
		// 1. 第一次调用 Get：缓存未命中，执行回调函数加载数据
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		// 2. 第二次调用 Get：缓存命中，❗❗绝不执行回调函数
		// loadCounts[k] > 1 说明重复调用了回调，缓存失效，测试失败
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}
	// 测试：访问不存在的key，必须返回错误
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
