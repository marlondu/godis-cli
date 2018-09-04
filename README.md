# godis-cli
a redis command client tool writen in golang

## Dependency

is dependent on redigo mobule, you can install it with:

```sh
go get github.com/gomodule/redigo/redis
```

## Help

```sh
this is a command tool for redis

Usage:
   godis-cli <command> [arguments]

The commands are:
   add    add a new redis server
   list   list all servers
   update update a server
   remove remove a server
   conn   connect to server
   info   check the info of a server

Use "godis-cli help <command>" for more information about a command
```

