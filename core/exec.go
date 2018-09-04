package core

import (
	"bufio"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// created: 2018/9/4
// execute redis commands and return result

func cmdParser(server *RedisServer) {
	reader := bufio.NewReader(os.Stdin)
	var conn redis.Conn
	var err error
	// create connection with redis
	address := server.Host + ":" + strconv.Itoa(server.Port)
	if server.Auth != "" {
		option := redis.DialPassword(server.Auth)
		conn, err = redis.Dial("tcp", address, option, redis.DialConnectTimeout(3 * time.Second))
	}else {
		conn, err = redis.Dial("tcp", address, redis.DialConnectTimeout(3 * time.Second))
	}
	// release connection when function exit
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	if err != nil {
		log.Fatal(err)
	}
	for   {
		line := readLine(reader, server)
		expr := "[\\s]+"
		regex, err := regexp.Compile(expr)
		if err != nil {
			fmt.Println("regex error: ", err)
			break
		}
		commands := regex.Split(line, -1)
		l := len(commands)
		if l > 1 {
			cmd := commands[0]
			params := make([]interface{},l - 1)
			for i, c := range commands[1:] {
				params[i] = c
			}
			reply, err := conn.Do(cmd, params...)
			if err != nil {
				fmt.Printf("execute command:%s error:%v\n", cmd, err)
				continue
			}
			printReply(reply)
		}else if l == 1 {
			if commands[0] == "exit" {
				fmt.Println("exit ...")
				break
			}else {
				fmt.Println("invalid command or parameters are required")
			}
		}
	}

}

func printReply(reply interface{}) {
	switch reply.(type) {
	case []byte:
		bytes := reply.([]byte)
		fmt.Println(string(bytes))
		break
	case string, int, int64, float64, bool:
		fmt.Println(reply)
		break
	case nil:
		fmt.Println()
		break
	case []interface{}:
		values := reply.([]interface{})
		for _, v := range values {
			printReply(v)
		}
		break
	default:
		break
	}
}

func readLine(reader *bufio.Reader, server *RedisServer) string {
	fmt.Printf("%s>", server.Name)
	data,_, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	line := string(data)
	return strings.TrimSpace(line)
}