package console

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
)

//type console struct {
//	printCommand *func(command string)
//}

type bash struct {
	user        string
	host        string
	port        int
	client *ssh.Client
	printCommand *func(command string)
}

type result struct {
	raw *[]byte
}

func (r result) ToString() *string  {
	output := string(*r.raw)

	outputLength := len(output)
	if outputLength > 0 {
		output = output[:outputLength-1]
	}

	return &output
}
func (r result) Raw() *[]byte  {
	return r.raw
}

type resulter interface {
	ToString() *string
	Raw() *[]byte
}

type BashConsoler interface {
	Run(command string) (resulter, error)
	MultiRun(commands []string) error



	getClient() (*ssh.Client, error)
	keyPath() string
	close()error
}

func NewBash(user string, host string, port int, printCommand *func(command string)) *BashConsoler {
	var b BashConsoler
	b = &bash{
		user: user,
		host: host,
		port: port,
		printCommand: printCommand,
	}

	return &b
}


func (b *bash) Run(command string) (resulter, error) {
	if b.printCommand != nil {
		(*b.printCommand)(command)
	}

	if b.client == nil {
		var err error

		b.client, err = b.getClient()

		if err != nil {
			return nil, err
		}
	}

	session, err := b.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	raw, err := session.CombinedOutput(command)

	var output resulter
	output = result{&raw}

	return output, err
}

func (b *bash) MultiRun(commands []string) error {
	var err error
	for _, command := range commands {
		_, err = b.Run(command)

		if err != nil {
			break
		}
	}

	return err
}














func (b *bash) getClient() (*ssh.Client, error) {
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
		User:            b.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	addr := fmt.Sprintf("%s:%d", b.host, b.port)

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}


	return client, err
}
func (b *bash) keyPath() string {
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return path.Join(home, ".ssh/id_rsa")
	}
	return ""
}

func (b *bash) close()error {
	var err error

	if b.client != nil {
		err = b.client.Close()
	}

	return err
}

