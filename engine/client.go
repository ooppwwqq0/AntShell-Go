package engine

import (
	"AntShell-Go/config"
	"AntShell-Go/models"
	"AntShell-Go/ssh"
	"AntShell-Go/ssh/terminal"
	"AntShell-Go/utils"
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	BastionOn  = 1
	BastionOff = 0
	TypeKey    = "authKey"
	TypePasswd = "password"
	TokenCache = "~/.antshell/cache"
)

type ClientSSH struct {
	host       models.Hosts
	config     *ssh.ClientConfig
	authMethod ssh.AuthMethod
	session    *ssh.Session
	c          config.Config

	sshHost     string
	sshUser     string
	sshPassword string
	sshType     string
	sshKeyPath  string
	sshPort     int
}

func (client *ClientSSH) Init(host models.Hosts, c config.Config) {
	client.host = host
	client.c = c
	//创建ssh登陆配置
	client.config = &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            client.host.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
}

func (client *ClientSSH) AuthKey() {
	keyPath, _ := homedir.Expand(client.c.Default.Key_Path)
	if keyPath != "" && utils.IsFile(keyPath) {
		key, err := ioutil.ReadFile(keyPath)
		if err != nil {
			logs.Error("ssh key file read failed", err)
		}
		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			logs.Error("ssh key signer failed", err)
		}
		client.authMethod = ssh.PublicKeys(signer)
		client.sshType = TypeKey
	}
}

func (client *ClientSSH) GetSSHInfo() {
	if client.host.Bastion == BastionOn {
		bastion := client.c.Bastion
		client.sshHost = bastion.Bastion_Host
		client.sshPort, _ = strconv.Atoi(bastion.Bastion_Port)
		client.sshUser = bastion.Bastion_User
		client.sshPassword = GetBastionPasswd(client.c.Bastion)
		client.config.User = bastion.Bastion_User
		client.sshType = TypePasswd
	} else {
		client.sshHost = client.host.Ip
		client.sshPort = client.host.Port
		client.sshUser = client.host.User
		client.sshPassword = client.host.Passwd
		client.AuthKey()
	}
	if client.sshType == TypePasswd {
		client.config.Auth = []ssh.AuthMethod{ssh.Password(client.sshPassword)}
	} else {
		client.config.Auth = []ssh.AuthMethod{client.authMethod}
	}
}

func (client *ClientSSH) GetSession() (session *ssh.Session, err error) {
	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", client.sshHost, client.sshPort)
	c, err := ssh.Dial("tcp", addr, client.config)
	if err != nil {
		logs.Error("创建ssh client 失败", err)
		os.Exit(1)
	}
	session, err = c.NewSession()
	if err != nil {
		logs.Error("创建ssh session 失败", err)
	}
	return
}

func (client *ClientSSH) Connection(sudo string, path string) {
	client.GetSSHInfo()
	session, err := client.GetSession()
	client.session = session

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	h, w := utils.GetSttySize()

	// Request pseudo terminal
	if err := session.RequestPty("xterm", h, w, modes); err != nil {
		logs.Error("request for pseudo terminal failed: ", err)
	}

	ok, err := session.Chan.SendRequest("shell", true, nil)
	if !ok || err != nil {
		logs.Error(err)
	}

	defer session.Chan.Close()
	defer session.Close()

	// 命令发送协程
	go func() {
		if client.host.Bastion == BastionOn {
			session.Chan.Write([]byte(client.host.Ip + " " + strconv.Itoa(client.host.Port) + "\r"))
		}
		session.Chan.Write([]byte("export LANG=en_US.UTF-8;export LC_ALL=en_US.UTF-8;export LC_CTYPE=en_US.UTF-8\r"))
		if client.host.Sudo != "" {
			session.Chan.Write([]byte("sudo -iu " + utils.IF(sudo != "", sudo, client.host.Sudo).(string) + "\r"))
		}
		if path != "" {
			session.Chan.Write([]byte("cd " + path + "\r"))
		} else if client.host.Path != "" {
			session.Chan.Write([]byte("cd " + client.host.Path + "\r"))
		}
	}()

	// 释放终端之前读取终端返回信息发送管道
	r := make(chan string)
	go func() {
		br := bufio.NewReader(session.Chan)
		var data = make([]byte, 1024)
		for {
			n, _ := br.Read(data)
			r <- string(data[:n])
		}
	}()

	// 读取管道信息输出
	flag := make(chan bool)
	go func() {
		for {
			t := time.NewTimer(time.Millisecond * 100)
			select {
			case d := <-r:
				//if strings.HasPrefix(d, "cd") || strings.HasPrefix(d, "export") || strings.HasPrefix(d, "sudo") {
				//	continue
				//}
				fmt.Print(d)
				t = time.NewTimer(time.Millisecond * 100)
			case <-t.C:
				t.Stop()
				flag <- true
				return
			}
		}
	}()

	// 错误回收
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
		}
	}()

	// 命令执行完成将终端释放给用户
	<-flag

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	fd := int(os.Stdin.Fd())

	state, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

	defer terminal.Restore(fd, state)

	// 捕获信号，感知终端大小变化
	go client.changeWindowsBySignal(fd)

	session.OutStart()
	session.Wait()
}

func (client *ClientSSH) changeWindowsBySignal(fd int) {
	// 捕获信号，感知终端大小变化
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGWINCH)
	for {
		<-sigs
		w, h, _ := terminal.GetSize(fd)
		client.session.WindowChange(h, w)
	}
}

func (client *ClientSSH) conn() {
	client.GetSSHInfo()

	session, err := client.GetSession()
	defer session.Close()
	fd := int(os.Stdin.Fd())

	state, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

	defer terminal.Restore(fd, state)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	w, h, err := terminal.GetSize(fd)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm", h, w, modes); err != nil {
		logs.Error("request for pseudo terminal failed: ", err)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGWINCH)
	go func() {
		for {
			<-sigs
			w, h, _ := terminal.GetSize(fd)
			session.WindowChange(h, w)
		}
	}()

	if err := session.Shell(); err != nil {
		logs.Error("failed to start shell: ", err)
	}

	session.Wait()
}

func GetBastionPasswd(bastion config.BastionSection) (bastionPasswd string) {

	if bastion.Bastion_Totp != "" {
		// token 缓存
		cacheToken, _ := utils.ReadByFile(TokenCache)
		var token string
		if cacheToken != "" {
			var n int
			for {
				token = utils.GetPasswdByTotp(bastion.Bastion_Totp)
				if cacheToken != token {
					break
				}
				if n%5 == 0 {
					logs.Info("等待动态token更新...")
				}
				n++
				time.Sleep(time.Second * 1)
			}
			if n != 0 {
				logs.Info("动态token已更新...，已等待：", n)
			}
		} else {
			token = utils.GetPasswdByTotp(bastion.Bastion_Totp)
		}

		go utils.CreateAndWrite(TokenCache, token)
		if bastion.Bastion_Passwd_Prefix != "" {
			bastionPasswd = bastion.Bastion_Passwd_Prefix + token
		} else {
			var bastionPasswdPrefix string
			for {
				fmt.Print(fmt.Sprintf(
					"%s@%s:%s's PIN:******",
					bastion.Bastion_User, bastion.Bastion_Host, bastion.Bastion_Port,
				))
				fmt.Scanln(&bastionPasswdPrefix)
				if bastionPasswdPrefix != "" {
					break
				}
			}
			bastionPasswd = bastionPasswdPrefix + token
		}
	} else {
		msg := fmt.Sprintf(
			"%s@%s:%s's password: PIN:****** + Token:",
			bastion.Bastion_User, bastion.Bastion_Host, bastion.Bastion_Port,
		)
		for {
			fmt.Print(msg)
			fmt.Scanln(&bastionPasswd)
			if bastionPasswd != "" {
				break
			}
		}
	}
	return
}
