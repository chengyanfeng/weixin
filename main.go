package main

import (
	_ "weixin/routers"
	"github.com/astaxie/beego"
	"weixin/controllers"
	//"github.com/astaxie/beego/toolbox"


)

func main() {
	//端口
	beego.BConfig.Listen.HTTPPort = 98                     //端口设置
	beego.BConfig.RecoverPanic = true                        //开启异常捕获
	beego.BConfig.EnableErrorsShow = true
	beego.BConfig.CopyRequestBody = true

	//自动注册路由
	beego.AutoRouter(&controllers.MainController{})
	//开启缓存
	controllers.InitCache()
	go func() {
		//开启协程
		//开启定时任务
		/*crontab()*/

	}()
	beego.Run()

}
/*func crontab() {
	toolbox.AddTask("pd", toolbox.NewTask("pd", "0 *//*1 * * * *", func() error { //每10分钟运行以下函数
		controllers.Get_user_message();
		return nil
	}))
	toolbox.StartTask() //开启定时任务
}*/

