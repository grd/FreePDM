package main

import (
	"fmt"
	"net"
	"os"

	"github.com/hirochachacha/go-smb2"
)

var (
	smb_host   = os.Getenv("SMB_HOST")
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

	s, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer s.Logoff()

	names, err := s.ListSharenames()
	if err != nil {
		panic(err)
	}

	for _, name := range names {
		fmt.Println(name)
	}
}
