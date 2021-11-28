# Go Web基础



## web 概念

### Request

include post , get , cookie,  url







### Response





### Conn

用户的每次请求链接





### Handler

处理请求和生成返回信息的处理逻辑



## Http 包执行流程

![image-20210903113614403](https://cdn.jsdelivr.net/gh/flyingchase/Private-Img@master/uPic/image-20210903113614403.png)

- 创建 Listen Socket 监听窗口
- Listen Socket 接收客户端请求——>client Socket 
- Client socket 读取 HTTP请求的协议头 交给对象的 handler 处理 再通过 client socket 写给客户端



#### 监听端口

- `ListenAndServe`
  初始化 serve 对象，调用`net.Listen("tcp",addr)` 底层使用 tcp 协议搭建服务来监听端口
- 







#### 接收客户端请求

- 调用`srv.Serve(net.Listener)`函数 

- 每次请求创建 Conn





#### 分配 Handler

- conn 解析 request `c.readRequest()` 获取 handler 即为调用 ListenAndServe 的第二个参数 
  - 为空 nil 则默认`handler=DefaultServeMax `  为路由器 匹配 url 跳转到对应的 handle 函数
  - 路由请求规则`“/”`： 
    - 跳转到 `http.HandleFunc("/",selfFunc)` 中的自定义函数 
    - DefaultServeMux 调用 ServeHTTP 方法 内部及调用 selfFunc







































































































# 实战

## 创建 web server

使用 http.ListenAndServer()

- 第一个参数网络地址	

  - “” 则为所有网络接口的 80 端口

- 第二个参数 handler

  - nil 则为 DefaultServeMux(multipledxer) 看做路由器

  

使用http.Server

- struct
  - Addr 字段表示网络地址
  - Handler 字段
    - nil defaulSerMux
  - ListenAndServe() 方法

``` go
// way One 
http.ListenAndServe("localhost:9090",nil)
// way Two
server:=http.Server{Addr: "localhost:8080",Handler: nil}
server.ListenAndServe()
```

以上只能执行 http 而非 https

需要分别加上：

- http.ListenAndServeTLS()
- Server.ListenAndServeTLS()



#### Handler

是一个 interface 接口

- 定义了一个方法ServeHTTP()

  - HTTPResponseWriter
  - 指向 Request（struct）的指针

  ```go
  type Handler interface {
     ServeHTTP(ResponseWriter, *Request)
  }
  ```

  







#### DefaultServeMux

Multiplexer 多路复用器（可被视为路由器） 是 ServerMux 的一个指针变量

- 也是一个 Handler
- 转发调用其他 handler
- 调用 http.Handle 函数实际上调用的是 DefaultServeMux 上的 Hanler 方法









不指定 server struct 中的 handler 字段值

- 可以使用 http.Handle 将某个 Handler 附加到 DefaultServeMux
  - http 包有一个 Handle 函数
  - ServerMux struct 也有一个 Handle 方法













#### http.Handle

func Handle(pattern string, handler Handler)

-  第二个参数是 handler（注意是*指针*）

  ``` go
  type Handler interface {
   ServeHTTP(ResponseWriter, *Request)
  // 实现 ServerHTTP 方法的类型均可视为 handler
   }
  ```

  

```go
server := http.Server{
   Addr: "localhost:1090",
   Handler: nil,  // use DefaultServeMux
}
http.Handle("/wo",&mh)
server.ListenAndServe()
```



#### http.HandleFunc

Handler函数行为与hanlder 类似 将 f 适配为 handler 使得handler 具有方法 f  类似*类型转换*

作用即为： Handler 函数转化为 Handler 内部还是调用 http.Handle 函数 

- Handler 函数的签名与 ServeHTTP 方法的签名一样，接收：
  - 一个 http.ResponseWriter
  - 一个 指向 http.Request 的指针

```go
// 第二个参数是 func 但是不要带()  带（）就直接执行了
http.HandleFunc("/home",welcome)

func welcome(w http.ResponseWriter, r *http.Request) {
   w.Write([]byte("Home!"))
}
```

HandleFunc 源码：

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    // 调用 DefaultServeMux 的 HandleFunc
   DefaultServeMux.HandleFunc(pattern, handler) 
}

// DefaultServeMux 的 HandleFunc 第二个参数是 Handler函数（不同于 http.Handle 第二个参数是 Handler）
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
    // 内部还是调用 http.Handle 函数
	mux.Handle(pattern, HandlerFunc(handler))
}
```





## 内置 Handlers

### NotFoundHandler

```go
// NotFound replies to the request with an HTTP 404 not found error.
func NotFound(w ResponseWriter, r *Request) { Error(w, "404 page not found", StatusNotFound) }

// NotFoundHandler returns a simple request handler
// that replies to each request with a ``404 page not found'' reply.
func NotFoundHandler() Handler { return HandlerFunc(NotFound) }
```

给每个请求的响应均为404

### RedirectHandler

```go
// Redirect to a fixed URL
type redirectHandler struct {
   url  string
   code int
}

// The provided code should be in the 3xx range and is usually
// StatusMovedPermanently, StatusFound or StatusSeeOther.
func RedirectHandler(url string, code int) Handler {
	return &redirectHandler{url, code}
}
```

将每个请求使用给定的状态码code——>指定的 url 跳转到提供的第一个参数



### StripPrefix

```go
func StripPrefix(prefix string, h Handler) Handler
```

从请求的 URL去掉指定的前缀prefix 再调用第二个参数 handler h

- 请求的 URL 与前缀 prefix 不符合 404
- h handler 将会在 请求 url 被去除 prefix 后调用 用于接收请求



### TimeoutHandler

```go
func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler
```

- time.Duration 表示一段时间 alias int64 的别名
- 返回 handler 在指定时间 dt duration 内运行传入的 h
  - h 将要被修饰的 handler
  - msg 超时则返回 msg 信息表示响应时间过长
  - dt 处理 h 的允许时间



### FileServer

```go
func FileServer(root FileSystem) Handler
```

返回 handler 基于 root 文件系统来响应请求  root 是字符串作为根目录

```go
// filestystem 可以自定义  通常委托给 
type FileSystem interface {
   Open(name string) (File, error)
}
```



eg 使用内置 filehandler 实现 handleFunc 

```go
http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
      http.ServeFile(w,r,"wwwroot"+r.URL.Path)
   })
   
http.ListenAndServe(":9090",nil)
http.ListenAndServe("8080",http.FileServer(http.Dir("wwwwroot")))
```













## 请求 Request

### HTTP 请求











### Request







### URL

请求信息的第一行里面的信息

指向 url>URL类型的指针 url>URL 是一个 struct

scheme://[userinfo@]host/path[?query]\[#fragment]通用格式 



*Query*

- 查询字符串



### Header

`map[string]/[]string`类型

设置 key 时候创建空的[]string 作为 value 第一个元素就是新的 header 值

key 添加元素执行 append 操作





### Body



io.ReadCloser 接口

- reader 接口

  - []byte 返回 byte 的数量 可选的错误

- closer 接口

  - 返回可选的错误

  





## UpLoad









### Form 表单



#### 表单发送请求

- html 表单里面的数据以 name-value 对的形式 通过 method 规定post/get请求发送出去

- 数据内容存储在 POST 请求的 Body里面

  - name-value 对的格式 通过表单的`Content Type`指定 `enctype`属性

  - entry 属性默认值 `application/x-www-form-urlcoded` 

  - enrty 属性设置为`multipart/form-data`   (大量数据、上传文件)

    - 每个 name-value 对转换为 MIME消息 每部分各自有 Content Type 和 Content Disposition

    

    

    

- method 属性设置 POST 和 GET
  - GET 请求没有 Body     数据通过 URL 编码的 name-value 对发送



#### Form 字段

- Request 上的函数允许我们从 URL 或/和 Body 中提取数据，通过这些字段：
  - Form  是 url.Values 类型——>type Values map[string]\[]string 类型
  - PostForm
  - MultipartForm

- Form 里面的数据是 key-value 对

  - 每个 key 对应一个切片 可以有多个值

- 通常的做法是：

  - 调用 ParseForm 或 ParseMultipartForm 来解析 Request
  - 相应的访问 Form、PostForm 或 MultipartForm 字段

- ```go
  func main() {
  	server:=http.Server {
  		Addr : "localhost:8080",
  	}
  	http.HandleFunc("/process", func(w http.ResponseWriter,r *http.Request){
  		r.ParseForm() // 解析 request
  
  		fmt.Fprintln(w,r.Form)
  	})
  	server.ListenAndServe()
  }
  
  // index.html 输出 是一个 map[string][]string
  // map[first_name:[wo] last_name:[456] uploaded:[Go Web.md]]
  ```

  



#### PostForm 字段

只读取表单的 key-value 对 不需读取 url 的 kv 对  使用 PostForm 字段

当 url 和 form 中均有 key对应的 Value 时候 Form 字段显示所有的 values 表单在前 url 在后

```html
map[first_name:[D] firtst_name:[Nick] last_name:[as]]  // D 为表单 Nick 为 url
```

- 只支持`"application/x-www-form-urlencoded"` 

```go
fmt.Fprintln(w,r.PostForm)  // 使用 PostForm字段
// map[first_name:[qw] last_name:[q]]   只显示表单输入的 key对应的 values
```



#### MultipartForm 字段

- 首先调用` ParseMultipartForm ` 方法
  - 该方法会在必要时调用 `ParseForm `方法
    - 参数是需要读取数据的长度 字节数
    - MultipartForm 只包含*表单*的 key-value 对
    - 返回类型是一个 struct 而不是 map。这个 struct 里有两个 map：
      - key 是 string，value 是 []string
      - 空的（key 是 string，value 是文件）

```go
func main() {
   server:=http.Server {
      Addr : "localhost:8080",
   }
   http.HandleFunc("/process", func(w http.ResponseWriter,r *http.Request){
      r.ParseMultipartForm(1024) // 解析 request  使用 ParseMultipartForm 解析 需要传入长度（字节数）

      fmt.Fprintln(w,r.MultipartForm) 
   })
   server.ListenAndServe()
}
// index 输出  struct 两个 map 第一个有数据第二个是空的
&{map[first_name:[12] last_name:[12]] map[]}
```

#### FormValue&PostFormValue 字段

FormValue 方法会返回 Form 字段中指定 key 对应的*第一个 value*

- 无需调用 ParseForm 或 ParseMultipartForm

PostFormValue 方法只能读取 PostForm

- FormValue 和 PostFormValue 都会调用 ParseMultipartForm 方法
- 表单的 `enctype` 设为 multipart/form-data，无法通过 FormValue 获得想要的值。



#### 文件 Files



- 调用 `ParseMultipartForm` 方法

- 从 `File`字段获得 FileHeader 调用 Open 方法获得文件

- 使用`ioutil.ReadAll`函数将文件内容读取到 []byte中

  ```go
  func process(w http.ResponseWriter, r *http.Request) {
     r.ParseMultipartForm(1024) // 最大传递字节
     fileHead := r.MultipartForm.File["uploaded"][0] // 读取指定 body 内容
     file, err := fileHead.Open()
     if err == nil {
        data, err := ioutil.ReadAll(file)
        if err == nil {
           fmt.Fprintln(w, string(data))
        }
     }
  }
  
  func main() {
     server := http.Server{
        Addr: "localhost:8080",
     }
     http.HandleFunc("/process", process)
     server.ListenAndServe()
  }
  ```





*FormFile*

- 返回对应 Key的第一个文件
  - 返回指定 Key 的第一个 Value



#### POST Json





#### MultipartReader

```go
func (r *Request) MultipartReader() (*multipart.Reader, error)
```









```go
type ResponseWriter interface {
   Header() Header

   Write([]byte) (int, error)

   WriteHeader(statusCode int)
}

// 是一个接口 response 实现了其内部的所有函数 所以 ResponseWriter 可视为 response 的指针
```

### ResponseWrite





### 内置 Response

- NotFound 函数，包装一个 404 状态码和一个额外的信息
- ServeFile 函数，从文件系统提供文件，返回给请求者
- ServeContent 函数，它可以把实现了 io.ReadSeeker 接口的任何东西里面的内容返回给请求者
- 还可以处理 Range 请求（范围请求），如果只请求了资源的一部分内容，那么 ServeContent 就可以如此响应。而 ServeFile 或 io.Copy 则不行。
- Redirect 函数，告诉客户端重定向到另一个 URL



```go
func process(w http.ResponseWriter, r *http.Request) {
   r.ParseMultipartForm(1024)                      // 最大传递字节
   fileHead := r.MultipartForm.File["uploaded"][0] // 读取指定 body 内容
   file, err := fileHead.Open()
   if err == nil {
      data, err := ioutil.ReadAll(file)
      if err == nil {
         fmt.Fprintln(w, string(data))
      }
   }
}
func writeExample(w http.ResponseWriter, r *http.Request) {
   str := `<html>
<head><title>Go Web</title></head>
<body><h1>Hello World</h1></body>
</html>`
   w.Write([]byte(str))
}
func writeHeaderExampl(w http.ResponseWriter, r *http.Request) {
   w.WriteHeader(501)
   fmt.Fprintln(w, "No such service, try next door")
}
func main() {
   server := http.Server{
      Addr: "localhost:8080",
   }
   http.HandleFunc("/write", writeHeaderExampl)
   http.HandleFunc("/redirect", headerEXample)
   http.HandleFunc("/json", jsonExample)
   server.ListenAndServe()
}
func headerEXample(w http.ResponseWriter, r *http.Request) {
   w.Header().Set("Location", "http://google.com")
   w.WriteHeader(302)

}

type Post struct {
   User    string
   Threads []string
}

func jsonExample(w http.ResponseWriter, r *http.Request) {
   w.Header().Set("Content-Type", "application/json")
   post := &Post{
      User : "wlzhou",
      Threads: []string {"first","second","third"},
   }
   json,_:=json.Marshal(post)
   w.Write(json)
}
```





## 模板

Web 模板即为 HTML 页面（预先设置好的）

`text/template` `html/template`模板库

### 模板引擎

合并模板和上下文数据产生 HTML

- 生成 HTML 写入`ResponseWriter` 再加入 HTTP响应返回给客户端
- ![IG3Raz](https://cdn.jsdelivr.net/gh/flyingchase/Private-Img@master/uPic/IG3Raz.png)





*ParseFiles*

- 解析模板文件 创建解析好的模板 struct 
- 是 template struct 上的 ParseFiles 上的方法调用



创建新的模板 名称为文件名  

``` go
 t, _ := template.ParseFiles("tmpl.html")

```





*ParseGlob*

- 模式匹配  根目录下匹配  

  ``` go
   t, _ := template.ParseGlob("*.html")
  ```

  







*Parse*

- 上述两个函数均会调用



### Action

模板中嵌入的命令 两组花括号之间{{}}

条件 迭代 设置 包含 定义

### 参数、管道、变量

*参数：*

- 模板中的值 
  - bool 整数 string struct key 变量 方法 

*管道：*

Unix 管道类似

- 把参数输出发送到下一个参数
- |隔开



### 函数

内置函数有：

​	define template block html js urlquery 

index print len with



自定义函数：

``` go
template.Funcs(funcMap FuncMap) *Template

type FuncMap map[string]interface{}

```



### 模板组合









## 路由

Controller

- Main() 设置类工作
- Controller：
  - 静态资源
  - 不同的请求发送给不同的 controller 处理







### 路由参数

- 静态路由：
  - 一个路径对应一个页面
    - /home /index







- 带参路由：
  - 依据路由参数 创建出一族不同的页面
    - /companies/123
    - /companies/homeAbout































































































