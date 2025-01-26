# FreePDM Documentation
The specs for a file server are minimal (a RPI4/5 should work but is not tested yet). At this moment, I am using a desktop and it works for me. If there are multiple users, then consider using a [proper server](https://github.com/grd/FreePDM/blob/main/ConceptOfDesign/FreePDM_Install/Docker-Compose.md).

## How to install FreePDM on an external VPN

FreePDM is a PDM that works external to FreeCAD. It can interoperate with FreeCAD files. The server and client is written in Go (https://golang.org), because of simplicity and speed.

This is about installing FreePDM in an external VPN. It works perfect with for instance NordVPN.
The consequences of using an external VPN is that you can't use Samba (SMB), because SMB only works in a LAN, so you have to use SSH with SSHFS. The benefit is that the install becomes a lot easier and it works more secure and over the internet instead of your LAN.

## Install (for developers on *nix)

#### Setting up some environment variables
- Inside the file `.bashrc` (if you have bash), place two lines at the end:
- `export GOPATH=/home/$user` where $user is the user name.
- `export FREEPDM_DIR=/home/$user/FreePDM` (nto necessary when you write in the /data/FreePDM.toml the line vaultsData = "[data dir]/vaults")
- To reload the `.bashrc` file you can run `source ~/.bashrc`.

#### Install Go
  - Install Go https://golang.org
  - Make sure you have a directory `bin` inside your home directory and in your `$PATH`

#### Open ports
The program ufw helps you to open the ports. First you need to know whether the program is active.
`sudo ufw status`
After that you need to open the ports.
```
  sudo ufw allow 22
  sudo ufw allow 8080/tcp
```
Now enable UFW with the following command.
`sudo ufw enable`

#### Install SSH
  1. SSH should be installed when you use Debian of Ubuntu Server. If not then there are plenty of examples in how to install SSH on your VPN.
  2. Adding the vault group
  `sudo groupadd vault`
  3. Adding users (one at the time):
```
  sudo useradd -M -g vault username  (replace username)
  sudo passwd username
  id username (note the number after uid)
```
  4. Install sshfs
  `sudo apt install sshfs`


#### Install FreePDM
  - `git clone github.com/grd/FreePDM.git`

#### Modify FreePDM.toml
Read `cat /etc/group` 
Write down the number after `vault` (in my case it is 2000) and the the `users` the names and UID's of the users.

Next, modify `FreePDM.toml` that is stored in the `FreePDM/data` directory:

```
VaultsDirectory = "/home/user/vaults" (this directory is the vault directory. Change it to you  like, such as "/vaults")
LogFile = ""
LogLevel = ""

[Users]
vault = 2000
user1 = 1001
user2 = 1002
```

The fields that start with Log are ignored ATM.

### Create the vault directory
When you are in the root directory and you want to create the vaults inside that directory you need to write:
`sudo mkdir /vaults` 
If you are in your own directory you can write:
`mkdir /home/user/vaults` where `user` stands for the directory that you own.

These two lines change who has access to the directory and the directory mode.
`sudo chown user:vault /vaults`
`sudo chmod 0775 /vaults`
When you own the directory you can write:
`chown user:vault /home/user/vaults`
`chmod 0775 /home/user/vaults`


#### build and test
Go to directory FreePDM and type:
  - `make createvault` Create a vault inside the vaults directory. You should first try with the vault `testpdm`
  -  `make vaultstest`. This should test and also put some sample data into the vault `testpdm`.

#### Start running the server
- `make pdmserver` runs the server.

If it runs then it is okay. You can stop using the server with pressing Ctrl-C.

### Use docker instead of running from bare metal
First install docker. There are plenty of examples on how to install docker and docker-compose.

Go to the FreePDM directory and edit (with nano for instance) the file docker-compose.yml and remove all the lines below the samba group because now we are using SSH.
In the app group go to "volumes" and add the following line there:
`- /vaults:/vaults`  This must match with the vaults directory that you wrote before. Or in case of you having `/vaults` inside your directory:
`- /home/user/vaults:/home/user/vaults`

To start docker you can type `make docker`


#### Show the directory from the client side
The directory shows the files after you mount the directory, so you can open a file. It is versioned, which means that each file has a version. The benefits of this is that it works roughly the same as ordinary commercial PDM's. The downside is that the paths of the files inside an assembly for instance need to be replaced but that is a normal operation inside a PDM.

`sshfs username@servername:/path/to/vaults /mnt` After this you need to write the password and if it all works well you should not have any error.



#### Start running the terminal application
- `make pdmterm` should run the termimal program. When you are running `pdmterm` and you run the command `help` you see all the commands that should run.

#### Client application
In the future there will be a client application. Right now I am still investigating how to do it right. See https://github.com/grd/FreePDM/discussions/77. It should run with the command `make pmdclient` but right now there is nothing.



[<< Previous Chapter](README.md) | [Content Table](README.md) | [Next Chapter >>](Docker-Compose.md)
