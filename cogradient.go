package xdd

import (
	"fmt"
	"github.com/cdle/sillyGirl/core"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"strconv"
	"strings"
)


var xdd = core.NewBucket("xdd")
var spec = "0 " + strconv.Itoa(rand.Intn(59)) + " " + Config.CTime + "/12 * * ?" //cron表达式，每秒一次
func init() {
	core.AddCommand("xdd", []core.Function{
		{
			Rules: []string{`raw ^同步`},
			Cron:  xdd.Get("update", "59 * * * *"),
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				CogradientContainers(s)
				return nil
			},
		},
		{
			Rules: []string{`raw ^状态`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				return count()
			},
		},
		{

			Rules: []string{`raw ^转换`, `raw ^zh`},
			Cron:  xdd.Get("update", spec),
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				fmt.Println("开始wskey转换")
				updateCookie(s)
				return nil
			},
		},
		{

			Rules: []string{`raw ^优先级`},
			Admin: true,
			Handle: func(s core.Sender) interface{} {
				fmt.Println("开始wskey转换")
				updateCookie(s)
				return nil
			},
		},
	})
}

func CogradientContainers(s core.Sender) {
	cks := GetJdCookies(func(sb *gorm.DB) *gorm.DB {
		return sb.Where(fmt.Sprintf("%s >= ? and %s != ?", Priority, Hack), 0, True)
	})
	tmp := []JdCookie{}
	for _, ck := range cks {
		if ck.Priority >= 0 && ck.Hack != True {
			tmp = append(tmp, ck)
		}
	}
	cks = tmp

	if Config.Mode == Parallel {
		for i := range Config.Containers {
			(&Config.Containers[i]).read()
		}
		for i := range Config.Containers {
			(&Config.Containers[i]).Write(cks)
		}
	} else {
		resident := []JdCookie{}
		if Config.Resident != "" {
			tmp := cks
			cks = []JdCookie{}
			for _, ck := range tmp {
				if strings.Contains(Config.Resident, ck.PtPin) {
					resident = append(resident, ck)
				} else {
					cks = append(cks, ck)
				}
			}
		}
		type balance struct {
			Container Container
			Weigth    float64
			Ready     []JdCookie
			Should    int
		}
		availables := []Container{}
		parallels := []Container{}
		bs := []balance{}
		for i := range Config.Containers {
			(&Config.Containers[i]).read()
			if Config.Containers[i].Available {
				if Config.Containers[i].Mode == Parallel {
					parallels = append(parallels, Config.Containers[i])
				} else {
					availables = append(availables, Config.Containers[i])
					bs = append(bs, balance{
						Container: Config.Containers[i],
						Weigth:    float64(Config.Containers[i].Weigth),
					})
				}
			}
		}
		bat := cks
		for {
			left := []JdCookie{}
			l := len(cks)
			total := 0.0
			for i := range bs {
				total += float64(bs[i].Weigth)
			}
			for i := range bs {
				if bs[i].Weigth == 0 {
					bs[i].Should = 0
				} else {
					bs[i].Should = int(math.Ceil(bs[i].Weigth / total * float64(l)))
				}

			}
			a := 0
			for i := range bs {
				j := bs[i].Should
				if j == 0 {
					continue
				}
				s := 0
				if bs[i].Container.Limit > 0 && j > bs[i].Container.Limit {
					s = a + bs[i].Container.Limit
					left = append(left, cks[s:a+j]...)
					bs[i].Weigth = 0
				} else {
					s = a + j
				}
				if s > l {
					s = l
				}
				bs[i].Ready = append(bs[i].Ready, cks[a:s]...)
				a += j
				if a >= l-1 {
					break
				}
			}
			if len(left) != 0 {
				cks = left
				continue
			}
			break
		}
		for i := range bs {
			bs[i].Container.Write(append(resident, bs[i].Ready...))
		}
		for i := range parallels {
			parallels[i].Write(append(resident, bat...))
		}
	}
	s.Reply("账号同步容器完成", core.E)
}

func count() interface{} {
	zs := 0
	yx := 0
	wx := 0
	tl := 0
	ts := 0
	tc := 0
	dt := Date()
	cks := GetJdCookies()
	for _, ck := range cks {
		zs++
		if ck.Available == True {
			yx++
		} else {
			wx++
		}
		if ck.CreateAt == dt {
			tc++
		}
	}
	jps := []JdCookiePool{}
	db.Find(&jps)
	for _, jp := range jps {
		if jp.CreateAt == dt {
			ts++
		}
		if jp.LoseAt == dt {
			tl++
		}
	}
	return fmt.Sprintf("总数%d,有效%d,无效%d,今日失效%d,今日扫码%d,今日新增%d", zs, yx, wx, tl, ts, tc)
}
