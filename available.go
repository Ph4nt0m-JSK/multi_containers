package xdd

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/cdle/sillyGirl/core"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type UserInfoResult struct {
	Data struct {
		JdVvipCocoonInfo struct {
			JdVvipCocoon struct {
				DisplayType   int    `json:"displayType"`
				HitTypeList   []int  `json:"hitTypeList"`
				Link          string `json:"link"`
				Price         string `json:"price"`
				Qualification int    `json:"qualification"`
				SellingPoints string `json:"sellingPoints"`
			} `json:"JdVvipCocoon"`
			JdVvipCocoonStatus string `json:"JdVvipCocoonStatus"`
		} `json:"JdVvipCocoonInfo"`
		JdVvipInfo struct {
			JdVvipStatus string `json:"jdVvipStatus"`
		} `json:"JdVvipInfo"`
		AssetInfo struct {
			AccountBalance string `json:"accountBalance"`
			BaitiaoInfo    struct {
				AvailableLimit     string `json:"availableLimit"`
				BaiTiaoStatus      string `json:"baiTiaoStatus"`
				Bill               string `json:"bill"`
				BillOverStatus     string `json:"billOverStatus"`
				Outstanding7Amount string `json:"outstanding7Amount"`
				OverDueAmount      string `json:"overDueAmount"`
				OverDueCount       string `json:"overDueCount"`
				UnpaidForAll       string `json:"unpaidForAll"`
				UnpaidForMonth     string `json:"unpaidForMonth"`
			} `json:"baitiaoInfo"`
			BeanNum    string `json:"beanNum"`
			CouponNum  string `json:"couponNum"`
			CouponRed  string `json:"couponRed"`
			RedBalance string `json:"redBalance"`
		} `json:"assetInfo"`
		FavInfo struct {
			FavDpNum    string `json:"favDpNum"`
			FavGoodsNum string `json:"favGoodsNum"`
			FavShopNum  string `json:"favShopNum"`
			FootNum     string `json:"footNum"`
			IsGoodsRed  string `json:"isGoodsRed"`
			IsShopRed   string `json:"isShopRed"`
		} `json:"favInfo"`
		GrowHelperCoupon struct {
			AddDays     int     `json:"addDays"`
			BatchID     int     `json:"batchId"`
			CouponKind  int     `json:"couponKind"`
			CouponModel int     `json:"couponModel"`
			CouponStyle int     `json:"couponStyle"`
			CouponType  int     `json:"couponType"`
			Discount    float64 `json:"discount"`
			LimitType   int     `json:"limitType"`
			MsgType     int     `json:"msgType"`
			Quota       float64 `json:"quota"`
			RoleID      int     `json:"roleId"`
			State       int     `json:"state"`
			Status      int     `json:"status"`
		} `json:"growHelperCoupon"`
		KplInfo struct {
			KplInfoStatus string `json:"kplInfoStatus"`
			Mopenbp17     string `json:"mopenbp17"`
			Mopenbp22     string `json:"mopenbp22"`
		} `json:"kplInfo"`
		OrderInfo struct {
			CommentCount     string        `json:"commentCount"`
			Logistics        []interface{} `json:"logistics"`
			OrderCountStatus string        `json:"orderCountStatus"`
			ReceiveCount     string        `json:"receiveCount"`
			WaitPayCount     string        `json:"waitPayCount"`
		} `json:"orderInfo"`
		PlusPromotion struct {
			Status int `json:"status"`
		} `json:"plusPromotion"`
		UserInfo struct {
			BaseInfo struct {
				AccountType    string `json:"accountType"`
				BaseInfoStatus string `json:"baseInfoStatus"`
				CurPin         string `json:"curPin"`
				DefinePin      string `json:"definePin"`
				HeadImageURL   string `json:"headImageUrl"`
				LevelName      string `json:"levelName"`
				Nickname       string `json:"nickname"`
				Pinlist        string `json:"pinlist"`
				UserLevel      string `json:"userLevel"`
			} `json:"baseInfo"`
			IsHideNavi     string `json:"isHideNavi"`
			IsHomeWhite    string `json:"isHomeWhite"`
			IsJTH          string `json:"isJTH"`
			IsKaiPu        string `json:"isKaiPu"`
			IsPlusVip      string `json:"isPlusVip"`
			IsQQFans       string `json:"isQQFans"`
			IsRealNameAuth string `json:"isRealNameAuth"`
			IsWxFans       string `json:"isWxFans"`
			Jvalue         string `json:"jvalue"`
			OrderFlag      string `json:"orderFlag"`
			PlusInfo       struct {
			} `json:"plusInfo"`
			XbScore string `json:"xbScore"`
		} `json:"userInfo"`
		UserLifeCycle struct {
			IdentityID      string `json:"identityId"`
			LifeCycleStatus string `json:"lifeCycleStatus"`
			TrackID         string `json:"trackId"`
		} `json:"userLifeCycle"`
	} `json:"data"`
	Msg       string `json:"msg"`
	Retcode   string `json:"retcode"`
	Timestamp int64  `json:"timestamp"`
}

var ua = `Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5 UCBrowser/13.4.2.1122`

func initCookie() {
	cks := GetJdCookies()
	l := len(cks)
	for i := 0; i < l-1; i++ {
		if cks[i].Available == True && !CookieOK(&cks[i]) {
			if pt_key, err := cks[i].OutPool(); err == nil && pt_key != "" {
				i--
			}
		}
	}
	go func() {
		Save <- &JdCookie{}
	}()

}

func CookieOK(ck *JdCookie) bool {
	// fmt.Println(ck.PtPin)
	cookie := "pt_key=" + ck.PtKey + ";pt_pin=" + ck.PtPin + ";"
	// fmt.Println(cookie)
	// jdzz(cookie, make(chan int64))
	if ck == nil {
		return true
	}
	req := httplib.Get("https://me-api.jd.com/user_new/info/GetJDUserInfoUnion")
	req.Header("Cookie", cookie)
	req.Header("Accept", "*/*")
	req.Header("Accept-Language", "zh-cn,")
	req.Header("Connection", "keep-alive,")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://home.m.jd.com/myJd/newhome.action?sceneval=2&ufc=&")
	req.Header("Host", "me-api.jd.com")
	req.Header("User-Agent", "jdapp;iPhone;9.4.4;14.3;network/4g;Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1")
	data, err := req.Bytes()
	if err != nil {
		return false
	}
	ui := &UserInfoResult{}
	if nil != json.Unmarshal(data, ui) {
		return false
	}
	switch ui.Retcode {
	case "1001": //ck.BeanNum
		if ui.Msg == "not login" {
			if len(ck.Wskey) > 0 {
				var pinky = ck.Wskey
				msg, err := GetWsKey(pinky)
				if err != nil {
					logs.Error(err)
				}
				//JdCookie{}.Push(fmt.Sprintf("????????????wskey---%s", msg))
				//??????????????????
				if strings.Contains(msg, "??????") {
					//ck.Push(fmt.Sprintf("Wskey???????????????%s", ck.PtPin))
					//(&JdCookie{}).Push(fmt.Sprintf("Wskey?????????%s", ck.PtPin))
				} else {
					ptKey := FetchJdCookieValue("pt_key", msg)
					ptPin := FetchJdCookieValue("pt_pin", msg)
					logs.Info(ptPin)
					ck := JdCookie{
						PtKey: ptKey,
						PtPin: ptPin,
					}
					if nck, err := GetJdCookie(ptPin); err == nil {
						nck.InPool(ck.PtKey)
						msg := fmt.Sprintf("???????????????%s", ck.PtPin)
						//(&JdCookie{}).Push(msg)
						logs.Info(msg)
					} else {
						//nck.Update(Available, False)
						//(&JdCookie{}).Push(fmt.Sprintf("???????????????%s", nck.PtPin))
					}
				}

			} else {
				//ck.Push(fmt.Sprintf("???????????????%s", ck.PtPin))
				//JdCookie{}.Push(fmt.Sprintf("???????????????%s", ck.Nickname))
			}
			return false
		}
	case "0":
		if url.QueryEscape(ui.Data.UserInfo.BaseInfo.CurPin) != ck.PtPin {
			return av2(cookie)
		}
		if ui.Data.UserInfo.BaseInfo.Nickname != ck.Nickname || ui.Data.AssetInfo.BeanNum != ck.BeanNum || ui.Data.UserInfo.BaseInfo.UserLevel != ck.UserLevel || ui.Data.UserInfo.BaseInfo.LevelName != ck.LevelName {
			ck.Updates(JdCookie{
				Nickname:  ui.Data.UserInfo.BaseInfo.Nickname,
				BeanNum:   ui.Data.AssetInfo.BeanNum,
				Available: True,
				UserLevel: ui.Data.UserInfo.BaseInfo.UserLevel,
				LevelName: ui.Data.UserInfo.BaseInfo.LevelName,
			})
			ck.UserLevel = ui.Data.UserInfo.BaseInfo.UserLevel
			ck.LevelName = ui.Data.UserInfo.BaseInfo.LevelName
			ck.Nickname = ui.Data.UserInfo.BaseInfo.Nickname
			ck.BeanNum = ui.Data.AssetInfo.BeanNum
		}
		return true
	}
	return av2(cookie)
}

func WsKeyOK2(ck *JdCookie) (bool, string) {
	rsp, err := GetWsKey(ck.Wskey)
	if err != nil {
		logs.Error(err)
		return false, rsp
	}
	if strings.Contains(rsp, "fake") {
		logs.Error(err)
		return false, rsp
	}
	return true, rsp
}

func av2(cookie string) bool {
	req := httplib.Get(`https://m.jingxi.com/user/info/GetJDUserBaseInfo?_=1629334995401&sceneval=2&g_login_type=1&g_ty=ls`)
	req.Header("User-Agent", ua)
	req.Header("Host", "m.jingxi.com")
	req.Header("Accept", "*/*")
	req.Header("Connection", "keep-alive")
	req.Header("Accept-Language", "zh-cn")
	req.Header("Accept-Encoding", "gzip, deflate, br")
	req.Header("Referer", "https://st.jingxi.com/my/userinfo.html?&ptag=7205.12.4")
	req.Header("Cookie", cookie)
	data, err := req.String()
	if err != nil {
		return true
	}
	return !strings.Contains(data, "login")
}
func updateCookie(s core.Sender) {
	cks := GetJdCookies()
	l := len(cks)
	logs.Info(l)
	xx := 0
	yy := 0
	s.Reply("????????????????????????Wskey", core.E)
	for i := range cks {
		if len(cks[i].Wskey) > 0 {
			time.Sleep(10 * time.Second)
			ck := cks[i]
			rsp, err := GetWsKey(ck.Wskey)
			if err != nil {
				logs.Error(err)
			}
			if strings.Contains(rsp, "fake") {
				s.Reply(fmt.Sprintf("Wskey?????????%s", ck.PtPin), core.E)
				//ck.Push(fmt.Sprintf("Wskey???????????????%s", ck.PtPin))
			} else {
				ptKey := FetchJdCookieValue("pt_key", rsp)
				ptPin := FetchJdCookieValue("pt_pin", rsp)
				ck := JdCookie{
					PtKey: ptKey,
					PtPin: ptPin,
				}
				if ptPin != "" || ptKey != "" {
					if nck, err := GetJdCookie(ck.PtPin); err == nil {
						xx++
						nck.InPool(ck.PtKey)
						nck.Update(Available, True)
					} else {
						yy++
						ck.Update(Available, False)
					}
				}
				go func() {
					Save <- &JdCookie{}
				}()
			}
		}
	}
	s.Reply(fmt.Sprintf("??????CK??????????????????%d???,?????????????????????%d???", xx, yy), core.E)
}
func FetchJdCookieValue(key string, cookies string) string {
	match := regexp.MustCompile(key + `=([^;]*);{0,1}`).FindStringSubmatch(cookies)
	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}
