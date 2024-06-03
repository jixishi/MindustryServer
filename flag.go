package main

import "flag"

type ptrVal struct {
	*string
	*int
	*bool
	*float64
	*float32
}
type Val struct {
	string
	int
	bool
	float64
	float32
}
type Flag struct {
	v    ptrVal
	sStr string
	lStr string
	dv   Val
	help string
}

func flagInit(f *Flag) {
	if f.v.string != nil {
		flag.StringVar(f.v.string, f.sStr, f.dv.string, "")
		flag.StringVar(f.v.string, f.lStr, f.dv.string, f.help)
	}
	if f.v.bool != nil {
		flag.BoolVar(f.v.bool, f.sStr, f.dv.bool, "")
		flag.BoolVar(f.v.bool, f.lStr, f.dv.bool, f.help)
	}
	if f.v.int != nil {
		flag.IntVar(f.v.int, f.sStr, f.dv.int, "")
		flag.IntVar(f.v.int, f.lStr, f.dv.int, f.help)
	}
}

var newhost string
var chost string

func init() {
	flagInit(&Flag{ptrVal{int: &Cfg.Port}, "p", "port", Val{int: 8080}, "Api接口地址"})
	flagInit(&Flag{ptrVal{string: &newhost}, "a", "add", Val{string: ""}, "添加服务器地址"})
	flagInit(&Flag{ptrVal{string: &chost}, "h", "host", Val{string: ""}, "服务器地址\tip:port"})
}
