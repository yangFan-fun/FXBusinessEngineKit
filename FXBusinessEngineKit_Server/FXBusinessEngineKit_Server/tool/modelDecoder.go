package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

// JsonDecoder 自动序列化json到指定类型
// jsonStr 接受来自客户端传递的json参数
// decoder 解析到的类型对象
func JsonDecoder[T any](jsonStr string, decoder T) T {
	var obj T
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		fmt.Println("解析器解析出现异常", err)
		return obj
	}
	fmt.Println("解析器解析成功", obj)
	return obj
}

// BindableModel 将请求中的Body参数绑定到模型
func BindableModel[T any](model T, context *gin.Context) (T, error) {
	obj := model
	bindErr := context.ShouldBindJSON(&obj)
	objStr, _ := json.Marshal(obj)
	// 读取数据时，指针会移动至EOF，会导致下一次无法读取
	// 将数据恢复
	context.Request.Body = io.NopCloser(bytes.NewBufferString(string(objStr)))
	if bindErr != nil {
		return obj, bindErr
	}
	return obj, nil
}
