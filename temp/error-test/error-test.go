package main

import (
	"context"
	"fmt"
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
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer s.Logoff()

	fs, err := s.Mount(smb_mount)
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	_, err = fs.Open("notExist.txt")

	fmt.Println(os.IsNotExist(err)) // true
	fmt.Println(os.IsExist(err))    // false

	fs.WriteFile("hello2.txt", []byte("test"), 0444)
	err = fs.WriteFile("hello2.txt", []byte("test2"), 0444)
	fmt.Println(os.IsPermission(err)) // true

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	_, err = fs.WithContext(ctx).Open("hello.txt")

	fmt.Println(os.IsTimeout(err)) // true
}
