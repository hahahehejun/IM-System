package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/hahahehejun/IM-System/common"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前客户端模式
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net dial err:", err)
		return nil
	}

	client.conn = conn
	return client
}

func (client *Client) selectUser() {
	command := &common.Command{
		Type:      0,
		Parameter: make(map[string]string),
	}
	msg := common.Build(command)
	fmt.Println(msg)
	_, err := client.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("conn write err:", err)
		return
	}
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>请输入聊天内容 exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) > 0 {
			command := &common.Command{
				Type:      1,
				Parameter: make(map[string]string),
			}

			command.Parameter["chatMsg"] = chatMsg + "\n"
			msg := common.Build(command)

			_, err := client.conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("conn write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容 exit退出")
		fmt.Scanln(&chatMsg)
	}

}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.selectUser()
	fmt.Println(">>>>请输入聊天对象 exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>请输入聊天内容 exit退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) > 0 {
				command := &common.Command{
					Type:      2,
					Parameter: make(map[string]string),
				}
				command.Parameter["toUser"] = remoteName
				command.Parameter["chatMsg"] = chatMsg + "\n\n"
				msg := common.Build(command)
				_, err := client.conn.Write([]byte(msg + "\n"))
				if err != nil {
					fmt.Println("conn write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>请输入聊天内容 exit退出")
			fmt.Scanln(&chatMsg)
		}
		client.selectUser()
		fmt.Println(">>>>请输入聊天对象 exit退出")
		fmt.Scanln(&remoteName)
	}

}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>请输入用户名：")
	fmt.Scanln(&client.Name)

	command := &common.Command{
		Type:      3,
		Parameter: make(map[string]string),
	}

	command.Parameter["newName"] = client.Name
	msg := common.Build(command)

	_, err := client.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("conn write err:", err)
		return false
	}
	return true
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("4. 退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("命令错误")
		return false
	}
}

//处理服务器响应消息
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {
		}
		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}

	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器ip(默认127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "服务器端口(默认8888)")
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("链接服务器失败")
		return
	}

	go client.DealResponse()
	fmt.Println("链接服务器成功")

	client.Run()
}
