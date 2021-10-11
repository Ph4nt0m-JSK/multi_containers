package xdd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/buger/jsonparser"
)

type Container struct {
	Type         string
	Name         string
	Default      bool
	Address      string
	Username     string
	Password     string
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Help         string
	Path         string
	Version      string
	QlVersion    string
	Token        string
	Available    bool
	Delete       []string
	Weigth       int
	Mode         string
	Reader       *bufio.Reader
	Config       string
	Limit        int
}
type Yaml struct {
	Containers      []Container
	Database        string
	DefaultPriority int `yaml:"default_priority"`
	Resident        string
	Mode            string
}

func init() {

	content, err := ioutil.ReadFile(Xdd.Get("cogradient", "/etc/sillyGirl/develop/multi_containers/cogradient.yaml"))
	if err != nil {
		logs.Warn("解析config.yaml读取错误: %v", err)
	}
	if yaml.Unmarshal(content, &Config) != nil {
		logs.Warn("解析config.yaml出错: %v", err)
	}

	for i := range Config.Containers {
		if Config.Containers[i].Weigth == 0 {
			Config.Containers[i].Weigth = 1
		}
		Config.Containers[i].Type = ""
		vv := regexp.MustCompile(`^(https?://[\.\w]+:?\d*)`).FindStringSubmatch(Config.Containers[i].Address)
		if len(vv) == 2 {
			Config.Containers[i].Address = vv[1]
		} else {
			logs.Warn("%s地址错误", Config.Containers[i].Type)
		}
		if Config.Containers[i].getToken() == nil {
			logs.Info("青龙登录成功")
		} else {
			logs.Warn("青龙登录失败")
		}

	}
}

func (c *Container) Write(cks []JdCookie) error {
	if len(c.Delete) > 0 {
		c.request("/open/envs", DELETE, fmt.Sprintf(`[%s]`, strings.Join(c.Delete, ",")))
	}
	hh := []string{}
	if len(cks) != 0 {
		for _, ck := range cks {
			if ck.Available == True {
				hh = append(hh, fmt.Sprintf(`{"name":"JD_COOKIE","value":"pt_key=%s;pt_pin=%s;","remarks":"%s"}`, ck.PtKey, ck.PtPin, ck.Nickname))
			}
		}
		sprintf := fmt.Sprintf(`[%s]`, strings.Join(hh, ","))
		c.request("/open/envs", POST, sprintf)
		type AutoGenerated struct {
			Code int `json:"code"`
			Data []struct {
				Value     string  `json:"value"`
				ID        string  `json:"_id"`
				Created   int64   `json:"created"`
				Status    int     `json:"status"`
				Timestamp string  `json:"timestamp"`
				Position  float64 `json:"position"`
				Name      string  `json:"name"`
				Remarks   string  `json:"remarks,omitempty"`
			} `json:"data"`
		}
		help := getQLHelp(len(cks))
		if c.Help != False {
			for k := range help {
				var data, err = c.request("/open/envs?searchValue=" + k)
				a := AutoGenerated{}
				err = json.Unmarshal(data, &a)
				if err != nil {
					continue
				}
				toDelete := []string{}
				for _, env := range a.Data {
					toDelete = append(toDelete, fmt.Sprintf("\"%s\"", env.ID))
				}
				if len(toDelete) > 0 {
					c.request("/open/envs", DELETE, fmt.Sprintf(`[%s]`, strings.Join(toDelete, ",")))
				}
			}
			for k, v := range help {
				if v == "" {
					v = "&"
				}
				r := map[string]string{
					"name":  k,
					"value": v,
				}
				d, _ := json.Marshal(r)
				c.request("/open/envs", POST, fmt.Sprintf(`[%s]`, string(d)))
			}
		}
	}
	return nil
}

func (c *Container) read() error {
	c.Available = true
	type AutoGenerated struct {
		Code int `json:"code"`
		Data []struct {
			Value     string  `json:"value"`
			ID        string  `json:"_id"`
			Created   int64   `json:"created"`
			Status    int     `json:"status"`
			Timestamp string  `json:"timestamp"`
			Position  float64 `json:"position"`
			Name      string  `json:"name"`
			Remarks   string  `json:"remarks,omitempty"`
		} `json:"data"`
	}
	var data, err = c.request("/open/envs?searchValue=JD_COOKIE")
	a := AutoGenerated{}
	err = json.Unmarshal(data, &a)
	if err != nil {
		c.Available = false
		return err
	}
	c.Delete = []string{}
	for _, env := range a.Data {
		c.Delete = append(c.Delete, fmt.Sprintf("\"%s\"", env.ID))
		res := regexp.MustCompile(`pt_key=(\S+);pt_pin=([^\s;]+);?`).FindAllStringSubmatch(env.Value, -1)
		for _, v := range res {
			CheckIn(v[2], v[1])
		}
	}
	return nil

	return nil
}

func (c *Container) getToken() error {
	req := httplib.Get(c.Address + fmt.Sprintf("/open/auth/token?client_id=%s&client_secret=%s", c.ClientId, c.ClientSecret))
	req.Header("Content-Type", "application/json;charset=UTF-8")
	if rsp, err := req.Response(); err == nil {
		data, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		c.Token, _ = jsonparser.GetString(data, "token")
		if c.Token == "" {
			c.Token, _ = jsonparser.GetString(data, "data", "token")
		}
	} else {
		return err
	}
	return nil
}

func (c *Container) request(ss ...string) ([]byte, error) {
	var api, method, body string
	for _, s := range ss {
		if s == GET || s == POST || s == PUT || s == DELETE {
			method = s
		} else if strings.Contains(s, "/open/") {
			api = s
		} else {
			body = s
		}
	}
	var req *httplib.BeegoHTTPRequest
	var i = 0
	for {
		i++
		switch method {
		case POST:
			req = httplib.Post(c.Address + api)
		case PUT:
			req = httplib.Put(c.Address + api)
		case DELETE:
			req = httplib.Delete(c.Address + api)
		default:
			req = httplib.Get(c.Address + api)
		}
		req.Header("Authorization", "Bearer "+c.Token)
		if body != "" {
			req.Header("Content-Type", "application/json;charset=UTF-8")
			req.Body(body)
		}
		if data, err := req.Bytes(); err == nil {
			code, _ := jsonparser.GetInt(data, "code")
			if code == 200 {
				return data, nil
			} else {
				logs.Warn(string(data))
				if i >= 5 {
					return nil, errors.New("异常")
				}
				c.getToken()
			}
		}
	}
	return []byte{}, nil
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)
