package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

// servers manager

// represent a redis server
type RedisServer struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Auth string `json:"auth"`
}

func (rs *RedisServer) info() string {
	data, err := json.MarshalIndent(*rs, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

// local storage file's path
// const localStoragePath = "./"
const localStorageFile = "cli.cache"

var ServersCache = make([]RedisServer, 0)

// load all servers from local storage
// and as cache
func Init() {
	if len(ServersCache) > 0 {
		return
	}
	filePath := getStoragePath() + localStorageFile
	if isFileExist(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal("open file:", filePath, " error:", err)
		}
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				arr := strings.Split(line, ":")
				if len(arr) == 4 {
					name := arr[0]
					host := arr[1]
					port, _ := strconv.Atoi(arr[2])
					auth := arr[3]
					rs := RedisServer{
						Name: name,
						Host: host,
						Port: port,
						Auth: auth,
					}
					ServersCache = append(ServersCache, rs)
				}
			}
			if err == io.EOF {
				break
			}
		}
		defer file.Close()
	}
}

func ListServers() {
	Init()
	fmt.Println("Servers(Name:Host):")
	for _, s := range ServersCache {
		fmt.Println(s.info())
	}
}

// Add a server to list
func AddServer(name string, host string, port int, auth string) {
	if host == "" {
		fmt.Printf("host is required, you cound use -h to specify it ")
		return
	}
	if name == "" {
		name = host
	}
	Init()
	for i := range ServersCache {
		if (ServersCache[i].Host == host && ServersCache[i].Port == port) || ServersCache[i].Name == name {
			fmt.Printf("Name:%s or Server:%s has exists\n", name, host)
			return
		}
	}

	server := RedisServer{
		Name: name,
		Host: host,
		Port: port,
		Auth: Encrypt(auth)}
	ServersCache = append(ServersCache, server)
	persistent2Local()
}

func persistent2Local() {
	file, _ := os.Create(getStoragePath() + localStorageFile)
	writer := bufio.NewWriter(file)
	defer func() {
		writer.Flush()
		file.Close()
		fmt.Println("save successfully")
	}()
	l := len(ServersCache)
	for i, server := range ServersCache {
		var info = []string{
			server.Name,
			server.Host,
			strconv.Itoa(server.Port),
			server.Auth}
		line := strings.Join(info, ":")
		writer.WriteString(line)
		if i < (l - 1) {
			writer.WriteRune('\n')
		}
	}
}

func UpdateServer(name string, host string, port int, auth string) {
	Init()
	for i := range ServersCache {
		if (ServersCache[i].Host == host && ServersCache[i].Port == port) || ServersCache[i].Name == name {
			if name != "" {
				ServersCache[i].Name = name
			}
			if host != "" {
				ServersCache[i].Host = host
			}
			if port != ServersCache[i].Port {
				ServersCache[i].Port = port
			}
			if auth != "" {
				ServersCache[i].Auth = Encrypt(auth)
			}
			persistent2Local()
			return
		}
	}
	fmt.Printf("Name:%s or Host:%s not found\n", name, host)
}

func RemoveServer(name string, host string) {
	Init()
	var tempServers = make([]RedisServer, 0)
	if name != "" && host != "" {
		for _, s := range ServersCache {
			if s.Name != name || s.Host != host {
				tempServers = append(tempServers, s)
			}
		}
	} else {
		if name != "" {
			for _, s := range ServersCache {
				if s.Name != name {
					tempServers = append(tempServers, s)
				}
			}
		} else if host != "" {
			for _, s := range ServersCache {
				if s.Host != host {
					tempServers = append(tempServers, s)
				}
			}
		}
	}
	ServersCache = tempServers
	persistent2Local()
}

func ConnectServer(name string) {
	Init()
	for _, s := range ServersCache {
		if s.Name == name {
			cmdParser(&s)
			break
		}
	}
}

func ensureSavePathExist(dir string) {
	if !isDirExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func isDirExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return info.IsDir()
}

func isFileExist(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return os.IsExist(err)
	}
	return !info.IsDir()
}

func getStoragePath() string {
	ur, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	osName := runtime.GOOS
	var dir string
	if strings.Contains(osName, "windows") {
		dir = strings.Replace(ur.HomeDir, "\\", "/", -1)
		dir = dir + "/godis-cli/"
	} else {
		dir = ur.HomeDir + "/godis-cli/"
	}
	ensureSavePathExist(dir)
	return dir
}
