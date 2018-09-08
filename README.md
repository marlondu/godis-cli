# godis-cli
A redis command client tool written in golang

## Dependency

Is dependent on redigo module, you can install it with:

```sh
go get github.com/gomodule/redigo/redis
```

## Download & Install
- You can download binary file from [here](https://github.com/marlondu/godis-cli/releases)
- Or if you have installed golang environment, you can install it with source code:

```sh
$ go get github.com/marlondu/godis-cli
$ cd ${GOPATH}/github.com/marlondu/godis-cli
$ go build .
# for macOS
$ sudo sh install.sh
```

## Help

```sh
this is a command tool for redis

Usage:
   godis-cli [arguments] <command>

The commands are:
   add    add a new redis server
   list   list all servers
   update update a server
   remove remove a server
   conn   connect to server
   info   check the info of a server

Use "godis-cli help <command>" for more information about a command
```

