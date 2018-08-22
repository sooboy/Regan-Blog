# Gin 分析
这里的Gin是指Golang语言里Web框架。[Gin](https://github.com/gin-gonic/gin)的使用简单方便，与node的express非常像。


## 如何开始?
![how to start it ?](../assets/BB920F6D-7F4A-4E58-8C03-9D8D3BF3C18F.png)

最初的最初 调用`gin.Default()`函数生成默认引擎，也就是*Engine结构体。签名如下：![Engine签名](../assets/how_to_start_it.png)

简单概括里面功能：
- 根据当前环境选择打印warnning信息
- 生成 `*Engine`
- 默认使用 `log`、 `recovery` 中间件
- 返回 `*Engine`

往后使用`Engine`注册一个`/ping`路径 并绑定一个`HandlerFunc`.

最后 `Run(:addr)`启动整个引擎

## 那故事应该从`Engine`说起了！
![Engine](../assets/engine_struct.png)

属性字段中去掉标记状态的量，以下几个比较重要：
- `RouterGroup` 管理路由组，实现`IRoutes`以及`IRouter`接口
- `HTMLRender` 模版引擎，默认使用golang下`template`作为模版引擎
- `trees`      存储`路径`以及对应`HandlerChain`详细信息
- `noRoute`、`noMethod` 默认404，405处理handler 可以通过`engine.NoRoute(...HandlerChain)``engine.Method(...HandlerChain)`设置

## 先说 `RouterGroup`
`RouterGroup`在`Engine` 里是内嵌结构体，它是这个样子的：
![RouterGroup](../assets/RouterGroup.png)
里面比较重要是`HandlerChain`:
![HandlerChain](../assets/HandlerChain.png)
`RouterGroup`实现了一下接口：
![IRoute](../assets/IRoute.png)

因为`RouterGroup`是内嵌在`Engine`里的所以它具备了`Use`,`Get`等方法。这些方法本质都是在组合`HandlerChain`、修改`basePath`.
看下面例子：

### Example
```golang
    r := gin.New()
	r.Use(HandlerFn1, HandlerFn2)
```


```golang
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}
```
最初`engine.Handlers`是空的,添加了`HandlerFn1`,`HandlerFn2`后=>[HandlerFn1,HandlerFn2].`basePath`为‘/’
继续添加一下代码
```golang
   r.Get("/",HandlerFn3,HandlerFn4,HandlerIndex)
   admin := r.Group("/admin",HandlerAdminLimit1,HandlerAdminLimit2)
        .Get("/money",HandlerMoneyLimit,HandlerMoney)
        .Get("/vote",HanlderVoteLimit,HandlerVote)
        .Get("/email",HanlderEmailLimit,HandlerEmail)
```
当执行了`r.Get("/",HandlerFn3,HandlerFn4,HandlerIndex)`
```golang
// GET is a shortcut for router.Handle("GET", path, handle).
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle("GET", relativePath, handlers)
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers)
	return group.returnObj()
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

```
将`RouterGroup`里面的`basePath`拼接当前路径
将`RouterGroup`里面的`Handlers`以及当前传入的`HandlerChain`copy到新的`HandlerChain`下


```javascript
{
    path:"/",
    handlers:[HandlerFn1,HandlerFn2,HandlerFn3,HandlerFn4,HandlerIndex]
}
```
将这样的信息（还有其他附加信息）添加到`engine.trees`里。

当执行`   admin := r.Group("/admin",HandlerAdminLimit1,HandlerAdminLimit2)`
则`admin` 这个新的`RouterGroup`的信息是这样
```javascript
{
    basePath:"/admin",
    Handlers:[HandlerFn1,HandlerFn2,HandlerAdminLimit1,HandlerAdminLimit2]
}
```

执行了`.Get("/money",HandlerMoneyLimit,HandlerMoney)`


将`admin`里面的`basePath`拼接当前路径
将`admin`里面的`Handlers`以及当前传入的`HandlerChain`copy到新的`HandlerChain`下


```javascript
{
    path:"/admin/money",
    handlers:[HandlerFn1,HandlerFn2,HandlerAdminLimit1,HandlerAdminLimit2,HandlerMoneyLimit,HandlerMoney]
}
```
将这样的信息（还有其他附加信息）添加到`engine.trees`里,下面以此类推。

### 这些信息如何使用呢？

 当客户端请求“/admin/money”,`Engine`在`tree`里搜索返回`[HandlerFn1,HandlerFn2,HandlerAdminLimit1,HandlerAdminLimit2,HandlerMoneyLimit,HandlerMoney]`所有这些`handler`,使用一个`*Context`上下文
 ```golang
 // ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()

	engine.handleHTTPRequest(c)

	engine.pool.Put(c)
}

func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
 ```

 `Context`可以中断`HandlerChain`的执行。比如，用户权限未认证。

 做到极致，我们可以这样：[gin 代码](../src/gin/demo01/main.go)

 ## 再说下 `tree`

 ```golang
type methodTree struct {
	method string
	root   *node
}

type methodTrees []methodTree

type node struct {
	path      string
	wildChild bool
	nType     nodeType
	maxParams uint8
	indices   string
	children  []*node
	handlers  HandlersChain
	priority  uint32
}
 ```
  














