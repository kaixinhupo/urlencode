# urlencode
## 1.安装
```
go get github.com/kaixinhupo/urlencode
```

## 2.使用
对所有请求参数编码
```go
package main

func test()  {
     params := make(map[string]string)
     params["foo"]="bar"
     params["local"]="中国"
     query := urlencode.UrlEncode(params,"gbk")
     println(query)
     //foo=bar&local=%d6%d0%b9%fa 
}

```
对字符串进行编码

```go
package main
func test()  {
     println(urlencode.Encode("中国","gbk"))
     //%d6%d0%b9%fa
}

```

解码
```go
package main

func test () {
    str:="foo=bar&local=%d6%d0%b9%fa"
    params := urlencode.UrlDecode(str,"gbk")
    fmt.Println(params)
    //map[foo->bar local->中国]

    println(urlencode.Decode("%d6%d0%b9%fa","gbk"))
    // 中国
}
```