package xdd

import (
	"github.com/cdle/sillyGirl/core"
)

var Parallel = "parallel"
var Config Yaml
var GhProxy = "https://ghproxy.com/"
var Xdd = core.NewBucket("xdd")
var Web = core.NewBucket("web")

func init() {
	core.Tail = ""
}
