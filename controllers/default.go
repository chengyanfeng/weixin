package controllers

import (
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)
type P map[string]interface{}
type p P
var localCache cache.Cache
func InitCache() {
c, err := cache.NewCache("memory", `{"interval":60}`)
//c, err := cache.NewCache("file", `{"CachePath":"./dhcache","FileSuffix":".cache","DirectoryLevel":2,"EmbedExpiry":120}`)
if err != nil {

Error(err)
} else {
localCache = c
fmt.Println("afafafa")
}
}

type MainController struct {
	beego.Controller
}


func (c *MainController) Get() {
	var token=""
	if localCache.Get("ticket")!=nil {
		fmt.Println("ticket 是从缓存里拿的")
		ticket:=localCache.Get("ticket").(string)
		c.Ctx.WriteString(string(ticket))
	}else {
if localCache.Get("token") ==nil{
//获取微信token
	response_token, _ := http.Get("https://api.weixin.qq.com/cgi-bin/token?appid=wxe658b8718e8a0faa&secret=4543fa8cf296881457ac8a7869227a9c&grant_type=client_credential")
	defer response_token.Body.Close()
	token_body, _ := ioutil.ReadAll(response_token.Body)
	p := *JsonDecode([]byte(string(token_body)))
	token= p["access_token"].(string);
	localCache.Put("token", token, 100*time.Minute)
	fmt.Println("这是从重新获取拿的")
		}else {
		token = localCache.Get("token").(string)
	fmt.Println("这是从缓存里拿的")
		}
//从token获取微信ticket
	response_ticket, _ := http.Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token="+token+"&type=jsapi")
	defer response_ticket.Body.Close()
	ticket_body, _ := ioutil.ReadAll(response_ticket.Body)
		p :=*JsonDecode([]byte(string(ticket_body)))
	ticket := p["ticket"].(string)
		localCache.Put("ticket", ticket, 100*time.Minute)
		fmt.Println("ticket 是从重新拿的")
		c.Ctx.WriteString(string(ticket))
}

}

func (c *MainController) Get_user_token() {
	var access_token=""
	var openid=""
	if localCache.Get("access_token") ==nil {
	fmt.Println("进去预先一步了")
	url:= c.GetString("url")
		fmt.Println("url:"+url)
	begin:=strings.Index(url, "code=")
	end:=strings.Index(url, "&")
	code:=url[begin:end]
		fmt.Println("code:"+code)
	//获取微信token
	response_token,_:=http.Get("https://api.weixin.qq.com/sns/oauth2/access_token?appid=wxe658b8718e8a0faa&secret=4543fa8cf296881457ac8a7869227a9c&"+code+"&grant_type=authorization_code")
	defer response_token.Body.Close()
	token_body,_:=ioutil.ReadAll(response_token.Body)

	p:=*JsonDecode([]byte(string(token_body)))

	refresh_token:=p["refresh_token"].(string)
	fmt.Print(refresh_token)
	//获取刷新token
	refresh_token_token,_:=http.Get("https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=wxe658b8718e8a0faa&grant_type=refresh_token&refresh_token="+refresh_token)
	defer refresh_token_token.Body.Close()
	ticket_body,_:=ioutil.ReadAll(refresh_token_token.Body)
	p=*JsonDecode([]byte(string(ticket_body)))
	access_token=p["access_token"].(string)
	openid=p["openid"].(string)
	//放到缓存里
	localCache.Put("access_token", access_token, 100*time.Minute)
	localCache.Put("openid", openid, 100*time.Minute)
	}else {
		access_token=localCache.Get("access_token").(string)
		openid=localCache.Get("openid").(string)
		fmt.Println("或许的是缓存里的")
	}
	fmt.Println("ytjytuytjt")
	fmt.Println(access_token)
	response_usermessage,_:=http.Get("https://api.weixin.qq.com/sns/userinfo?access_token="+access_token+"&openid="+openid+"&lang=zh_CN")
	defer response_usermessage.Body.Close()
	usermessage,_:=ioutil.ReadAll(response_usermessage.Body)
	c.Ctx.WriteString(string(usermessage))
}


func  Get_user_message(){
	access_token:=localCache.Get("access_token")
	openid:=localCache.Get("openid")
	fmt.Println(access_token)
	response_usermessage,_:=http.Get("https://api.weixin.qq.com/sns/userinfo?access_token="+access_token.(string)+"&openid="+openid.(string)+"&lang=zh_CN")
	defer response_usermessage.Body.Close()
	usermessage,_:=ioutil.ReadAll(response_usermessage.Body)
		fmt.Println(usermessage)
}




func JsonDecode(b []byte) (p *P) {
	p = &P{}
	err := json.Unmarshal(b, p)
	if err != nil {
		Error("JsonDecode", string(b), err)
	}
	return
}


// 记录err信息
func Error(v ...interface{}) {
	beego.Error(v)
}



func JsonEncode(v interface{}) (r string) {
	b, err := json.Marshal(v)
	if err != nil {
		Error(err)
	}
	r = string(b)
	return
}
