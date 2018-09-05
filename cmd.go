package main

import (
	"flag"
	"fmt"
	"github.com/marlondu/godis-cli/core"
	"strings"
)

// command operator manager
// Entrance of this program

//
//var cmds = []string{
//	"add",
//	"list",
//	"update",
//	"remove",
// 	"conn"
// 	"info"
//	"help",
//}

var (
	h = flag.String("h", "", "the server host of redis, default \"\"")
	p = flag.Int("p", 6379, "redis server's port, default is 6379")
	a = flag.String("a", "", "the auth of redis server,default nil")
	n = flag.String("n", "", "the name of redis server, you should specify when you add")
)

func main() {
	flag.Parse()
	cmds := flag.Args()
	if len(cmds) == 0 {
		fmt.Println("Use command \"godis-cli help\" for more information")
		return
	}
	dealCmds(cmds)
}

func dealCmds(cmds []string) {
	l := len(cmds)
	var mainCmd = cmds[0]
	switch mainCmd {
	case "help":
		if l > 1 {
			helpCmd(cmds[1])
		} else {
			printHelpTips()
		}
		break
	case "add":
		core.AddServer(*n, *h, *p, *a)
		break
	case "list":
		core.ListServers()
		break
	case "update":
		core.UpdateServer(*n, *h, *p, *a)
		break
	case "remove":
		core.RemoveServer(*n, *h)
		break
	case "conn":
		core.ConnectServer(*n)
		break
	case "info":
		core.InfoServer(*n, *h)
		break
	default:
		fmt.Println("Use command \"godis-cli help\" for more information")
		break
	}
}

func helpCmd(cmd string) {
	cmd = strings.ToLower(cmd)
	switch cmd {
	case "help":
		printHelpTips()
		break
	case "list", "add", "update", "remove", "info":
		printCmdHelpTips()
		break
	default:
		break

	}
}

// print help tip
func printHelpTips() {
	var tips = `this is a command tool for redis

Usage:
	godis-cli <command> [arguments]

The commands are:
	add		add a new redis server
	list	list all servers
	update	update a server
	remove 	remove a server
	conn	connect to server
	info 	check the info of a server

Use "godis-cli help <command>" for more information about a command`

	fmt.Println(tips)
}

func printCmdHelpTips() {
	var tips = `the supported commands are list, add, update, remove, info
Usage:
	godis-cli [parameters] <commands>

Params:
	-n server's name
	-h server's host
	-p server's port
	-a server' password to authorize

Examples:
	godis-cli list
	godis-cli -n local-redis conn
	godis-cli -n prod-redis -h 172.43.10.5 -a ****** add
`
	fmt.Println(tips)
}
