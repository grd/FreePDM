// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Tests for the SMB share functionality.
package main

import (
	"fmt"
	iofs "io/fs"
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

	fmt.Println("here")
	matches, err := iofs.Glob(fs.DirFS("."), "*")
	if err != nil {
		panic(err)
	}
	for _, match := range matches {
		fmt.Println(match)
	}

	err = RunTheProgram(fs)
	if err != nil {
		panic(err)
	}
}
