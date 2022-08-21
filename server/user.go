package server

import (
	"fmt"
	"net"

	"github.com/hahahehejun/IM-System/common"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	//将用户加入onlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播用户上线
	this.server.BroadCast(this, "已上线")

}

func (this *User) Offline() {
	//将用户加入onlineMap
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播用户上线
	this.server.BroadCast(this, "已下线")
}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {
	fmt.Println(msg)
	command := common.Parse(msg)
	if command == nil {
		fmt.Println("json parse err")
		return
	}
	if command.Type == 0 {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Name + "] 在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if command.Type == 3 {
		newName, has := command.Parameter["newName"]
		if !has {
			this.SendMsg("参数异常")
		}
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("该用户名已存在")
		} else {
			this.server.mapLock.Lock()

			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this

			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMsg("用户名更新成功")
		}
	} else if command.Type == 2 {
		remoteName, has := command.Parameter["toUser"]
		if !has {
			this.SendMsg("消息格式错误")
			return
		}
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("用户不存在")
			return
		}
		content, hasMsg := command.Parameter["chatMsg"]
		if !hasMsg || content == "" {
			this.SendMsg("消息为空")
			return
		}
		remoteUser.SendMsg("[" + this.Name + "]: " + content)
	} else if command.Type == 1 {
		content, hasMsg := command.Parameter["chatMsg"]
		if !hasMsg || content == "" {
			this.SendMsg("消息为空")
			return
		}
		this.server.BroadCast(this, content)
	} else {
		this.SendMsg("命令异常")
	}
}
