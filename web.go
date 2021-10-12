package xdd

import (
	"encoding/json"
	"fmt"
	_ "github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strings"
)

type CkLogin struct {
	Ck    string `json:"ck"`
	WsKey string `json:"wsKey"`
	QQ    int    `json:"qq"`
	Token string `json:"token"`
}

var theme = ""

func init() {
	server := core.Server
	server.GET(Web.Get("path", "/web"), func(c *gin.Context) {
		c.String(200, "11111111111111111111111111")
		if Config.Theme == "" {
			Config.Theme = GhProxy + "https://raw.githubusercontent.com/xiaeroc/xdd/master/myTheme/kuduan.html"
		}
		if theme != "" {
			c.Writer.WriteString(theme)
			return
		}
		if strings.Contains(Config.Theme, "http") {
			s, _ := httplib.Get(Config.Theme).String()
			if s != "" {
				theme = s
				c.Writer.WriteString(s)
				return
			}

		}
		f, err := os.Open(Config.Theme)
		if err == nil {
			d, _ := ioutil.ReadAll(f)
			theme = string(d)
			c.Writer.WriteString(string(d))
			return
		}
	})

	server.GET(Web.Get("ckLogin", "/ckLogin"), func(c *gin.Context) {
		data, _ := c.GetRawData()
		ckLogin := CkLogin{}
		err := json.Unmarshal(data, &ckLogin)
		if err != nil {
			c.String(200, "登录失败")
			return
		}
		if ckLogin.Ck != "" {
			ptKey := FetchJdCookieValue("pt_key", ckLogin.Ck)
			ptPin := FetchJdCookieValue("pt_pin", ckLogin.Ck)
			if ptKey != "" && ptPin != "" {
				ck := JdCookie{
					PtKey: ptKey,
					PtPin: ptPin,
					QQ:    ckLogin.QQ,
					Hack:  False,
				}
				if CookieOK(&ck) {
					if HasKey(ck.PtKey) {
						c.String(200, "重复添加")
						return
					} else {
						if nck, err := GetJdCookie(ck.PtPin); err == nil {
							err := nck.InPool(ck.PtKey)
							if err != nil {
								return
							}
							c.String(200, "更新账号")
							core.NotifyMasters(fmt.Sprintf("网页：更新账号，%s", ck.PtPin))
						} else {
							err := NewJdCookie(&ck)
							if err != nil {
								return
							}
							core.NotifyMasters(fmt.Sprintf(fmt.Sprintf("网页：添加账号，%s", ck.PtPin)))
							c.String(200, "添加账号")
						}
						return
					}
				} else {
					c.String(200, "无效账号")
					return
				}
			}
		}
	})

	server.GET(Web.Get("SMSLogin", "/SMSLogin"), func(c *gin.Context) {
		data, _ := c.GetRawData()
		ckLogin := CkLogin{}
		err := json.Unmarshal(data, &ckLogin)
		if err != nil {
			c.String(200, "登录失败")
			return
		}
		ptKey := FetchJdCookieValue("pt_key", ckLogin.Ck)
		ptPin := FetchJdCookieValue("pt_pin", ckLogin.Ck)
		ck := &JdCookie{
			PtKey: ptKey,
			PtPin: ptPin,
			Hack:  False,
			QQ:    ckLogin.QQ,
		}

		if ptKey != "" && ptPin != "" {
			if CookieOK(ck) {
				if !HasPin(ptPin) {
					NewJdCookie(ck)
					msg := fmt.Sprintf("来自短信的添加,账号：%s,QQ: %d", ck.PtPin, ck.QQ)
					core.NotifyMasters(msg)
				} else if !HasKey(ptKey) {
					ck, _ := GetJdCookie(ptPin)
					ck.InPool(ptKey)
					if ckLogin.QQ != 0 {
						ck.Update(QQ, ckLogin.QQ)
					}
					msg := fmt.Sprintf("来自短信的更新,账号：%s", ck.PtPin)
					core.NotifyMasters(msg)
				}

			} else {

			}
		}
	})
}
