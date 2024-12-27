# FreePDM Documentation
The specs for a file server are minimal (a RPI4/5 should work but is not tested yet). At this moment, I am using a desktop and it works for me. If there are multiple users, then consider using a proper server.

## How to install FreePDM

FreePDM is a PDM that works external to FreeCAD. It can interoperate with FreeCAD files. The server and client is written in Go (https://golang.org), because of simplicity and speed.

## Install (for developers on *nix)

#### Setting up some environment variables
- Inside the file `.bashrc` (if you have bash), place two lines at the end:
- `export GOPATH=/home/$user` where $user is the user name.
- `export FREEPDM_DIR=/home/$user/FreePDM`
- To reload the `.bashrc` file you can run `source ~/.bashrc`.

#### Install Go
  - Install Go https://golang.org
  - Make sure you have a directory `bin` inside your home directory and in your `$PATH`

#### Open ports
The program ufw helps you to open the ports. First you need to know whether the program is active.
`sudo ufw status`
If the anwer is 'Status: inactive' you need to activate ufw.
`sudo ufw enable`
After that you need to open the ports.
```
  sudo ufw allow 8080/tcp
  sudo ufw allow Samba
```

#### Install Samba
  1. Install Samba (a FOSS version of the SMB file server protocol). There are plenty of examples in how to install Samba on your computer.
  2. Adding users (one at the time):
```
  sudo useradd -M -g sambashare username  (replace username)
  sudo smbpasswd -a username
  sudo smbpasswd -e username
  sudo pdbedit -L
```
  3. Edit smb.conf with your editor, such as `sudo nano /etc/samba/smb.conf`
  and in the end of this file, add this:
```
  [vaults]
   comment = PDM Vaults
   path = /samba/vaults
   browseable = yes
   read only = no
   writable = yes
   guest ok = no
   valid users = @sambashare user user1
```

  - **`path`**: You can put any path here. This is just an example.
  - **`valid users`**: This should note @sambashare and the users. Each time you add someone (#2) you should rerun the program (#5).

4. Verify your smb.conf with `testparm`
5. Run `sudo systemctl restart smbd`
6. Test with `smbclient //localhost/vaults -U $user` (where $user is the user name)
If you see the `vaults` then it is working. After that you can mount the directory with:
`sudo mount -t cifs //localhost/vaults /mnt/vaults -o username=admin,password=<password>,uid=1000,gid=1000,vers=3.0,rw,iocharset=utf8`. Make sure that the directory `/mnt/vaults` (or the directory that you put into) exist and that the user name and password are correct.
When this is working you can try to replace `//localhost/vaults` with your ip address, such as `//192.168.0.15/vaults`.

#### Install FreePDM
  - `git clone github.com/grd/FreePDM.git`

#### Modify FreePDM.toml
Read `cat /etc/group` 
Write down the number after `sambashare` (in my case it is 125) and the the `users`.

Next, modify `FreePDM.toml` that is stored in the `FreePDM/data` directory:

```
VaultsDirectory = "/home/user/vaults"
LogFile = ""
LogLevel = ""

[Users]
vault = 125
user1 = 1001
user2 = 1002
```

The fields that start with Log are ignored ATM.

#### build and test
Go to directory FreePDM and type:
  - `make createvault` Create a vault inside the vaults directory. You should try with the vault `testpdm`
  -  `make vaultstest`. This should test and also put some sample data into the vault `testpdm`.

#### Start running the server
- `make pdmserver` runs the server.

#### Show the director from the client side
The directory shows the files after you mount the directory, so you can open a file. It is versioned, which means that each file has a version. The benefits of this is that it works roughly the same as ordinary commercial PDM's. The downside is that the paths of the files inside an assembly for instance need to be replaced but that is a normal operation inside a PDM.

#### Start running the terminal application
- `make pdmterm` should run the termimal program. When you run the command `help` you see all the commands that should run.

#### Client application
In the future there will be a client application. Right now I am still investigating how to do it right. See https://github.com/grd/FreePDM/discussions/77. It should run with the command `make pmdclient` but right now there is nothing.



[<< Previous Chapter](README.md) | [Content Table](README.md) | [Next Chapter >>](Docker-Compose.md)
