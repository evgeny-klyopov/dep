package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
)

type Bash struct {
	User string
	Host string
	Port int
}

func NewBash(user string, host string, port int) Bash {
	return Bash{
		User: user,
		Host: host,
		Port: port,
	}
}


func (b *Bash) run(commands []string) error {
	client, err := b.getClient()

	if err != nil {
		return err
	}

	for name, command := range commands {
		func() {
			fmt.Println(name)

			session, err := client.NewSession()
			if err != nil {
				//exit(err, " connect #" + strconv.Itoa(num) + " level 1")
			}
			defer session.Close()

			b, err := session.CombinedOutput(command)
			if err != nil {
				//exit(err, " connect #" + strconv.Itoa(num) + " level 2")
			}

			fmt.Println("----------")
			fmt.Println("Output")
			fmt.Print(string(b))
			fmt.Println("----------")
			fmt.Println("")
		}()
	}

	return nil
}

func (b *Bash) getClient() (*ssh.Client, error) {
	var pk = b.keyPath()

	key, err := ioutil.ReadFile(pk)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            b.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	addr := fmt.Sprintf("%s:%d", b.Host, b.Port)

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client, err
}
func (b *Bash) keyPath() string {
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return path.Join(home, ".ssh/id_rsa")
	}
	return ""
}

