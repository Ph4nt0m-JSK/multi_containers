package xdd

//
//import (
//	"io/ioutil"
//
//	"github.com/beego/beego/v2/core/logs"
//	"gopkg.in/yaml.v2"
//)
//
//type Yaml struct {
//	Containers []Container
//	// Tasks              []Task
//	Qrcode              string
//	Master              string
//	Mode                string
//	Static              string
//	Database            string
//	QywxKey             string `yaml:"qywx_key"`
//	Resident            string
//	UserAgent           string `yaml:"user_agent"`
//	Theme               string
//	TelegramBotToken    string `yaml:"telegram_bot_token"`
//	TelegramUserID      int    `yaml:"telegram_user_id"`
//	QQID                int64  `yaml:"qquid"`
//	QQGroupID           int64  `yaml:"qqgid"`
//	SMSAddress          string `yaml:"SMSAddress"`
//	ApiToken            string `yaml:"ApiToken"`
//	DefaultPriority     int    `yaml:"default_priority"`
//	NoGhproxy           bool   `yaml:"no_ghproxy"`
//	QbotPublicMode      bool   `yaml:"qbot_public_mode"`
//	DailyAssetPushCron  string `yaml:"daily_asset_push_cron"`
//	Version             string `yaml:"version"`
//	Node                string
//	Npm                 string
//	Python              string
//	Pip                 string
//	NoAdmin             bool   `yaml:"no_admin"`
//	QbotConfigFile      string `yaml:"qbot_config_file"`
//	HttpProxyServerPort int    `yaml:"http_proxy_server_port"`
//	Wskey               bool   `yaml:"Wskey"`
//	CTime               string `yaml:"AtTime"`
//}
//
//var Balance = "balance"
//var Parallel = "parallel"
//var GhProxy = "https://ghproxy.com/"
//var Cdle = false
//
//var Config Yaml
//
//func init() {
//	content, err := ioutil.ReadFile("E:\\work\\JD\\sillyGirl\\develop\\multi_containers\\cogradient.yaml")
//	if err != nil {
//		logs.Warn("解析config.yaml读取错误: %v", err)
//	}
//	if yaml.Unmarshal(content, &Config) != nil {
//		logs.Warn("解析config.yaml出错: %v", err)
//	}
//
//}
