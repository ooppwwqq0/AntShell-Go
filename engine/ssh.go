package engine

//
//import (
//	"AntShell-Go/models"
//	"AntShell-Go/ssh"
//	"AntShell-Go/ssh/terminal"
//	"fmt"
//	"github.com/mitchellh/go-homedir"
//	"io/ioutil"
//	"log"
//	"os"
//	"time"
//)
//
//func SSH(host models.Hosts) {
//	sshHost := host.Ip
//	sshUser := host.User
//	sshPassword := host.Passwd
//	sshType := "key"              //password 或者 key
//	sshKeyPath := "~/.ssh/id_rsa" //ssh id_rsa.id 路径"
//	sshPort := host.Porgt
//
//	//创建sshp登陆配置
//	config := &ssh.ClientConfig{
//		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
//		User:            sshUser,
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
//		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
//	}
//	if sshType == "password" {
//		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
//	} else {
//		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
//	}
//
//	//dial 获取ssh client
//	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
//	client, err := ssh.Dial("tcp", addr, config)
//	if err != nil {
//		log.Fatal("创建ssh client 失败", err)
//	}
//	//defer client.Close()
//
//	//创建ssh-session
//	session, err := client.NewSession()
//	if err != nil {
//		log.Fatal("创建ssh session 失败", err)
//	}
//	defer session.Close()
//
//	//执行远程命令
//	combo, err := session.CombinedOutput("uptime")
//	if err != nil {
//		log.Fatal("远程执行cmd 失败", err)
//	}
//	log.Println("命令输出:", string(combo))
//	session.Close()
//	//创建ssh-session
//	session, err = client.NewSession()
//	if err != nil {
//		log.Fatal("创建ssh session 失败", err)
//	}
//	defer session.Close()
//	fd := int(os.Stdin.Fd())
//
//	state, err := terminal.MakeRaw(fd)
//	if err != nil {
//		panic(err)
//	}
//
//	defer terminal.Restore(fd, state)
//
//	session.Stdout = os.Stdout
//	session.Stderr = os.Stderr
//	session.Stdin = os.Stdin
//
//	modes := ssh.TerminalModes{
//		ssh.ECHO:          0,     // disable echoing
//		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
//		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
//	}
//
//	// Request pseudo terminal
//	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
//		log.Fatal("request for pseudo terminal failed: ", err)
//	}
//
//	//if err := session.Shell(); err != nil {
//	//	log.Fatal("failed to start shell: ", err)
//	//}
//
//	session.Wait()
//
//}
//
//func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
//	keyPath, err := homedir.Expand(kPath)
//	if err != nil {
//		log.Fatal("find key's home dir failed", err)
//	}
//	key, err := ioutil.ReadFile(keyPath)
//	if err != nil {
//		log.Fatal("ssh key file read failed", err)
//	}
//	// Create the Signer for this private key.
//	signer, err := ssh.ParsePrivateKey(key)
//	if err != nil {
//		log.Fatal("ssh key signer failed", err)
//	}
//	return ssh.PublicKeys(signer)
//}
//
////
////
////func (client *ClientSSH) ConnectionBack() {
////	client.GetSSHInfo()
////	session, err := client.GetSession()
////
////	modes := ssh.TerminalModes{
////		ssh.ECHO:          1,
////		ssh.TTY_OP_ISPEED: 14400,
////		ssh.TTY_OP_OSPEED: 14400,
////	}
////	h, w := utils.GetSttySize()
////
////
////	var modeList []byte
////	for k, v := range modes {
////		kv := struct {
////			Key byte
////			Val uint32
////		}{k, v}
////		modeList = append(modeList, ssh.Marshal(&kv)...)
////	}
////	modeList = append(modeList, 0)
////
////	type ptyRequestMsg struct {
////		Term     string
////		Columns  uint32
////		Rows     uint32
////		Width    uint32
////		Height   uint32
////		Modelist string
////	}
////
////	req := ptyRequestMsg{
////		Term:     "xterm",
////		Columns:  uint32(w),
////		Rows:     uint32(h),
////		Width:    uint32(w * 8),
////		Height:   uint32(h * 8),
////		Modelist: string(modeList),
////	}
////	ok, err := session.Chan.SendRequest("pty-req", true, ssh.Marshal(&req))
////	if !ok || err != nil {
////		log.Println(err)
////	}
////
////
////
////	ok, err := session.Chan.SendRequest("shell", true, nil)
////	if !ok || err != nil {
////		log.Println(err)
////	}
////
////	flag := make(chan bool)
////
////	br := bufio.NewReader(session.Chan)
////	//buf := []byte{}
////	t := time.NewTimer(time.Microsecond * 1000)
////	defer t.Stop()
////	// 构建一个信道, 一端将数据远程主机的数据写入, 一段读取数据写入ws
////	r := make(chan string)
////	defer session.Chan.Close()
////	defer session.Close()
////
////	go func() {
////		session.Chan.Write([]byte("uptime\n"))
////		session.Chan.Write([]byte("sudo -iu root\n"))
////	}()
////
////	go func() {
////
////		var data = make([]byte, 1024)
////		for {
////			n, _ := br.Read(data)
////			r <- string(data[:n])
////		}
////	}()
////
////	go func() {
////
////		for {
////			t := time.NewTimer(time.Second)
////			select {
////			case d := <-r:
////				if strings.HasPrefix(d, "uptime") {
////					continue
////				}
////				fmt.Print(d)
////			case <-t.C:
////				t.Stop()
////				flag <- true
////				return
////			}
////		}
////	}()
////
////	defer func() {
////		if err := recover(); err != nil {
////			log.Println(err)
////		}
////	}()
////	<-flag
////
////	session.Stdout = os.Stdout
////	session.Stderr = os.Stderr
////	session.Stdin = os.Stdin
////
////	fd := int(os.Stdin.Fd())
////
////	state, err := terminal.MakeRaw(fd)
////	if err != nil {
////		panic(err)
////	}
////
////	defer terminal.Restore(fd, state)
////
////	sigs := make(chan os.Signal, 1)
////	signal.Notify(sigs, syscall.SIGWINCH)
////	go func() {
////		for {
////			<-sigs
////			w, h, _ := terminal.GetSize(fd)
////			session.WindowChange(h, w)
////		}
////	}()
////
////	session.OutStart()
////	session.Wait()
////}
