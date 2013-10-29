package someutils

// thanks to this for inspiration ... https://gist.github.com/jedy/3357393
 
import (
	"bufio"
	"code.google.com/p/go.crypto/ssh"
//	"crypto"
//	"crypto/rsa"
	"errors"
	"flag"
	"fmt"
	"github.com/laher/uggo"
	"github.com/howeyc/gopass"
	"io"
	"os"
	"strconv"
	"strings"
)

type ScpOptions struct {
	Port *int
	IsRecursive *bool
	IsRemoteTo *bool
	IsRemoteFrom *bool
}

type clientPassword string

func (p clientPassword) Password(user string) (string, error) {
    return string(p), nil
}

func init() {
	Register(Util{
		"scp",
		Scp})
}


//TODO: error for multiple ats or multiple colons
func parseTarget(target string) (string, string, string, error) {
	if strings.Contains(target, ":") {
		//remote
		parts := strings.Split(target, ":")
		userHost := parts[0]
		file := parts[1]
		user := ""
		var host string
		if strings.Contains(userHost, "@") {
			uhParts := strings.Split(userHost, "@")
			user = uhParts[0]
			host = uhParts[1]
		} else {
			host = userHost
		}
		return file, host, user, nil
	} else {
		//local
		return target, "", "", nil
	}
}

func Scp(call []string) error {
	options := ScpOptions{}
	flagSet := flag.NewFlagSet("scp", flag.ContinueOnError)
	options.IsRecursive = flagSet.Bool("r", false, "Recursive copy")
	options.Port = flagSet.Int("P", 22, "Port number")
	options.IsRemoteTo = flagSet.Bool("t", false, "")
	options.IsRemoteFrom = flagSet.Bool("f", false, "")
	helpFlag := flagSet.Bool("help", false, "Show this help")
	err := flagSet.Parse(uggo.Gnuify(call[1:]))
	if err != nil {
		println("Error parsing flags")
		return err
	}
	if *options.IsRecursive {
		return errors.New("This scp does NOT implement 'recursive scp'. Yet.")
	}
	if *options.IsRemoteTo || *options.IsRemoteFrom {
		return errors.New("This scp does NOT implement 'remote scp'. Yet.")
	}
	args := flagSet.Args()
	if *helpFlag || len(args) != 2 {
		println("`scp` [options] [[user@]host1:]file1 [[user@]host2:]file2")
		flagSet.PrintDefaults()
		return nil
	}

	srcFile, srcHost, srcUser, err := parseTarget(args[0])
	if err != nil {
		println("Error parsing source")
		return err
	}
	dstFile, dstHost, dstUser, err := parseTarget(args[1])
	if err != nil {
		println("Error parsing destination")
		return err
	}
	if srcHost != "" {
		if dstHost != "" {
			return errors.New("remote->remote NOT implemented (yet)!")
		}
		//from-scp
		session, err := connect(srcUser, srcHost, *options.Port)
		if err != nil {
			return err
		}
		defer session.Close()
		ce := make (chan error)
		go func() {
			cw, err := session.StdinPipe()
			if err != nil {
				println(err.Error())
			        ce <- err
				return
			}
			defer cw.Close()
			_, err = cw.Write([]byte{0})
			if err != nil {
				println("Write error: "+err.Error())
			        ce <- err
				return
			}
			r, err := session.StdoutPipe()
			if err != nil {
				println("session stdout err: " + err.Error())
			        ce <- err
				return
			}
			//defer r.Close()
			fw, err := os.Create(dstFile)
			if err != nil {
			        ce <- err
				println("File creation error: "+err.Error())
				return
			}
			defer fw.Close()
			scanner := bufio.NewScanner(r)
			scanner.Scan()
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			        ce <- err
				return
			}
			cmdFull := scanner.Text()
			cmd := string(cmdFull[0])
			parts := strings.Split(cmdFull[1:], " ")
			fmt.Printf("Command: %s. Details: %v\n", cmd, parts)
			if cmd == "C" {
				mode, err := strconv.ParseInt(parts[0], 8, 32)
				if err != nil {
					println("Format error: "+err.Error())
					ce <- err
					return
				}
				size, err := strconv.Atoi(parts[1])
				if err != nil {
					println("Format error: "+err.Error())
					ce <- err
					return
				}
				base := parts[2]	
				fmt.Printf("Mode: %d, size: %d, file: %s\n", mode, size, base)
				_, err = cw.Write([]byte{0})
				if err != nil {
					println("Write error: "+err.Error())
					ce <- err
					return
				}
				//todo - buffer
				b := make([]byte, 1)
				for {
					_, err := r.Read(b)
					if err != nil {
						println("Read error: " + err.Error())
						ce <- err
						return
					}
					//fmt.Printf("Read byte: %v\n", b)
					if b[0] != 0 {
						_, err := fw.Write(b)
						if err != nil {
							println("Write error: " + err.Error())
							ce <- err
							return
						}
					} else {
						//TODO: chmod here
						break
					}
				}
				err = fw.Close()
				if err != nil {
					println(err.Error())
					ce <- err
					return
				}
				_, err = cw.Write([]byte{0})
				if err != nil {
					println("Write error: "+err.Error())
					ce <- err
					return
				}
				err = cw.Close()
				if err != nil {
					println(err.Error())
					ce <- err
					return
				}
			} else {
				fmt.Printf("Command '%s' NOT implemented\n", cmd)
				return
			}
		}()
		err = session.Run("/usr/bin/scp -f "+srcFile)
		if err != nil {
			println("Failed to run remote scp: " + err.Error())
		}
		return err
	}

	if dstHost != "" {
		//to-scp
		srcFileInfo, err := os.Stat(srcFile)
		if err != nil {
			return err
		}
		session, err := connect(dstUser, dstHost, *options.Port)
		if err != nil {
			return err
		}
		defer session.Close()
		ce := make (chan error)
		go func() {
			w, err := session.StdinPipe()
			if err != nil {
				println(err.Error())
			        ce <- err
				return
			}
			defer w.Close()
			rdr, err := os.Open(srcFile)
			if err != nil {
			        ce <- err
				println(err.Error())
				return
			}
			defer rdr.Close()
			fmt.Fprintln(w, "C0644", srcFileInfo.Size(), dstFile)
			io.Copy(w, rdr)
			fmt.Fprint(w, "\x00") // terminate with null byte
			err = w.Close()
			if err != nil {
				println(err.Error())
			        ce <- err
				return
			}
			err = rdr.Close()
			if err != nil {
				println(err.Error())
			        ce <- err
				return
			}
		}()
		if err := session.Run("/usr/bin/scp -qrt ./"); err != nil {
			println("Failed to run: " + err.Error())
			return err
		}
		
	}
	return nil
}

func connect(user, host string, port int) (*ssh.Session, error) {
	fmt.Printf("Password for user '%s':\n", user)
 	pass := gopass.GetPasswd()
	password := clientPassword(pass)
	clientConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthPassword(password),
		},
	}
	target := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", target, clientConfig)
	if err != nil {
		println("Failed to dial: " + err.Error())
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		println("Failed to create session: " + err.Error())
	}
	return session, err

}
