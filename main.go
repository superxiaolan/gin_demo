package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"log"
	"net/http"
)

//自定义中间件，相当于Java里的拦截器
func MyMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		//通过自定义中间件，设置的值，在后续的处理函数中可以获取这里的参数
		context.Set("usersession","userid-1")
		//前置处理
		context.Next() //放行
		//后置处理
	}
}


func main()  {

	//创建一个服务
	ginServer :=gin.Default()
	ginServer.Use(favicon.New("./img.ico"))

	//加载静态页面
	ginServer.LoadHTMLGlob("templates/*")
	//加载静态资源文件
	ginServer.Static("/static","./static")

	//连接数据库的代码


	//访问地址，处理我们的请求  Request Response
	ginServer.GET("/index",func(context *gin.Context){
		/*context.JSON(200,gin.H{"msg":"hello,world!"})*/
		context.HTML(http.StatusOK,"index.html",gin.H{
			"msg":"hello,world!-- from Go background data",
		})
	})

	//接收前端传递过来的参数

	//方式一：通过url传递参数
	//url?userid=xxx&username=xxxx
	//eg: http://localhost:8082/user/info?name=张三&age=18
	/*ginServer.GET("/user/info",func(context *gin.Context){
		//用Query获取前端传递过来的参数
		userid := context.Query("userid")
		username := context.Query("username")

		context.JSON(200,gin.H{
			"userid":userid,
			"username":username,
		})
	})*/

	//方式二：通过表单传递参数
	//url/user/info/xxx/xxxx
	//eg: http://localhost:8082/user/info/张三/18
	ginServer.GET("/user/info/:username/:age", MyMiddleWare(), func(context *gin.Context){

		//取出中间件设置的值
		usersession := context.MustGet("usersession").(string)
		log.Println("======================>",usersession)

		//用Param获取前端传递过来的参数
		username := context.Param("username")
		age := context.Param("age")

		context.JSON(http.StatusOK,gin.H{
			"username":username,
			"age":age,
		})
	})


	//前端给后端传数据 json格式
	ginServer.POST("/json",func(context *gin.Context){
		//request body
		//json格式的数据
		//{"name":"张三","age":18}

		data, _ := context.GetRawData()

		var m map[string]interface{}

		//包装为json格式
		_ = json.Unmarshal(data,&m)

		context.JSON(http.StatusOK,m)
	})

	//支持函数式编程
	ginServer.POST("/user/add",func(context *gin.Context){
		username := context.PostForm("username")
		password := context.PostForm("password")
		mobile_number := context.PostForm("mobile_number")
		context.JSON(http.StatusOK,gin.H{
			"username":username,
			"password":password,
			"mobile_number":mobile_number,
		})
	})


	//路由
	ginServer.GET("/test",func(context *gin.Context){
		//重定向
		context.Redirect(http.StatusMovedPermanently,"http://www.baidu.com")
	})

	ginServer.NoRoute(func(context *gin.Context){
		context.HTML(http.StatusNotFound,"404.html",nil)
	})

	//路由组 /user/add  /user/login  /user/logout
	userGroup := ginServer.Group("/user")
	{
		userGroup.GET("/add")
		userGroup.GET("/login")
		userGroup.GET("/logout")
	}

	//路由组之订单 /order/add  /order/delete  /order/update
	orderGroup := ginServer.Group("/order")
	{
		orderGroup.GET("/add")
		orderGroup.GET("/delete")
		orderGroup.GET("/update")
	}

	//服务器端口
	ginServer.Run(":8082")
}
