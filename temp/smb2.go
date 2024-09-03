package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/hirochachacha/go-smb2"
)

var (
	smb_host   = os.Getenv("SMB_HOST")
	smb_mount  = os.Getenv("SMB_MOUNT")
	smb_user   = os.Getenv("SMB_USER")
	smb_passwd = os.Getenv("SMB_PASSWD")
)

func main() {
	if smb_host == "" {
		fmt.Printf("Local environment(s) not set.\n")
		os.Exit(1)
	}
	conn, err := net.Dial("tcp", smb_host+":445")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     smb_user,
			Password: smb_passwd,
			Domain:   "MicrosoftAccount",
		},
	}

	c, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer c.Logoff()

	fs, err := c.Mount(smb_mount)
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	f, err := fs.Create("hello.txt")
	if err != nil {
		panic(err)
	}
	defer fs.Remove("hello.txt")
	defer f.Close()

	_, err = f.Write([]byte("Hello world!"))
	if err != nil {
		panic(err)
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}

	bs, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))

	// Hello world!
}
