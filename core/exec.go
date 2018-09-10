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
	conn, err := obtainConnection(server)
	// release connection when function exit
	defer func() {
		conn.Flush()
		conn.Close()
	}()
	if err != nil {
		log.Fatal(err)
	}
	for {
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
			params := make([]interface{}, l-1)
			for i, c := range commands[1:] {
				params[i] = c
			}
			err = ensureConnected(&conn, server)
			if err != nil {
				fmt.Printf("obtain a new connection error:%v\n", err)
				continue
			}
			reply, err := conn.Do(cmd, params...)
			if err != nil {
				fmt.Printf("execute command:%s error:%v\n", cmd, err)
				continue
			}
			printReply(reply)
		} else if l == 1 {
			if commands[0] == "exit" {
				fmt.Println("exit ...")
				break
			} else {
				fmt.Println("invalid command or parameters are required")
			}
		}
	}

}

func obtainConnection(server *RedisServer) (conn redis.Conn, err error) {
	// create connection with redis
	address := server.Host + ":" + strconv.Itoa(server.Port)
	options := []redis.DialOption{
		redis.DialKeepAlive(15 * time.Minute),
		redis.DialConnectTimeout(3 * time.Second),
	}
	if server.Auth != "" {
		options = append(options, redis.DialPassword(Decrypt(server.Auth)))
	}
	conn, err = redis.Dial("tcp", address, options...)
	return
}

func ensureConnected(conn *redis.Conn, server *RedisServer) error {
	err := (*conn).Err()
	if err != nil {
		// connection is not usable
		// reload connection
		connection, err := obtainConnection(server)
		if err != nil {
			return err
		} else {
			*conn = connection
		}
	}
	return nil
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
	data, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	line := string(data)
	return strings.TrimSpace(line)
}
