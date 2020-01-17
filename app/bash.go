package app



import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
)

type Bash struct {
	User        string
	Host        string
	Port        int
	Debug       bool
	client *ssh.Client
	PrefixDebug string
}

func NewBash(user string, host string, port int, debug bool, prefixDebug string) Bash {
	return Bash{
		User: user,
		Host: host,
		Port: port,
		Debug: debug,
		PrefixDebug: prefixDebug,
	}
}
func (b *Bash) runOutput(command string) (*string, error) {
	raw, err := b.run(command)

	if err == nil {
		output := string(raw)

		outputLength := len(output)
		if outputLength > 0 {
			output = output[:outputLength-1]
		}

		return &output, err
	}

	return nil, err
}
func (b *Bash) close()error {
	var err error

	if b.client != nil {
		err = b.client.Close()
	}

	return err
}

func (b *Bash) multiRun(commands []string) error {
	var err error
	for _, command := range commands {
		_, err = b.run(command)

		if err != nil {
			break
		}
	}

	return err
}
func (b *Bash) run(command string) ([]byte, error) {
	if b.Debug == true {
		fmt.Println(b.PrefixDebug + command)
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

	out, err := session.CombinedOutput(command)

	return out, err
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


	return client, err
}
func (b *Bash) keyPath() string {
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return path.Join(home, ".ssh/id_rsa")
	}
	return ""
}

