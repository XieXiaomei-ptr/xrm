package main

import (
	"log"
)

const (
	DIR     = "dir"
	LAN     = "lan"
	VER     = "v"
	VERSION = "xrm 0.0.1"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	args := MakeArgs()
	args.AddStringArg(DIR, "", "The directory in which your code base is located. Multiple paths are supported.\nWhen multiple paths are specified at the same time, a separate #.\nFor example: -dir=/data/my_project/code1#/data/my_project/code2", "#")
	args.AddStringArg(LAN, "", "Specify the programming language contained in dir, which supports c, c++, go, and js.\nIf you specify multiple languages at the same time, you need to separate with #.\nFor Example: -lan=c#c++#js#go", "#")
	args.AddBoolArg(VER, false, "Version")
	if !args.Parse() {
		return
	}
	ver := args.GetBoolVar(VER)
	if ver {
		log.Println(VERSION)
	}
	lans := args.GetStringVar(LAN)
	dirs := args.GetStringVar(DIR)
	if len(lans) == 0 || len(dirs) == 0 {
		if !ver {
			log.Println("[error] You must specify both -" + DIR + " and -" + LAN + " parameters")
		}
		return
	}
	xrm := MakeXrm()
	for i := 0; i < len(lans); i++ {
		xrm.Config(lans[i])
	}
	xrm.ConfigCodeDir(dirs)
	retCode := xrm.Execute()
	if retCode {
		log.Println("It's done!")
	} else {
		log.Println("[error] execute failed!")
	}
}
