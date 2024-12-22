# FreePDM Documentation
The specs for a file server are minimal (a RPI4/5 should work but is not tested yet). At this moment, I am using a desktop and it works for me. If there are multiple users, then consider using a proper server.

## How to install FreePDM

FreePDM is a PDM that works external to FreeCAD. It can interoperate with FreeCAD files. The server and client is written in Go (https://golang.org), because of simplicity and speed.

## Install (for developers on *nix)

#### Install Go
  - Install Go https://golang.org
  - Make sure you have a directory `bin` inside your home directory and in your `$PATH`
  - Play a bit with GOPATH in your `.bashrc` file (if you have bash). Make sure to have the line `export GOPATH=/home/$user` where $user is the user name. Also place the line `export FREEPDM_DIR=/home/$user/FreePDM`. To reload the `.bashrc` file you can run `source ~/.bashrc`.

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

  - **`path`**: You can put any path here. This is just how I do this.
  - **`valid users`**: This should note @sambashare and the users. Each time you add someone you should rerun the program. See #5.

4. Verify your smb.conf with `testparm`
5. Run `sudo systemctl restart smbd`


#### Install FreePDM
  - `git clone github.com/grd/FreePDM.git`
  - Go to directory FreePDM and type `make fstest`. This should test and also put some sample data into the vault.

#### Modify FreePDM.toml
Read `cat /etc/group` 
Write down the number after `sambashare` (in my case it is 125) and the the `users`.

Next, add the following text to the end of your `FreePDM.toml`. The file is stored in the `FreePDM/data` directory:

```
[Users]
vault = 125
user1 = 1001
user2 = 1002
```

#### Show the directory from the client side
The directory shows the files after you mount the directory, so you can open a file. It is versioned, which means that each file has a version. The benefits of this is that it works roughly the same as ordinary commercial PDM's. The downside is that the paths of the files inside an assembly for instance need to be replaced but that is a normal operation inside a PDM.

How you can mount a directory in your PC from a shared directory, depends on your OS.




[<< Previous Chapter](README.md) | [Content Table](README.md) | [Next Chapter >>](Docker-Compose.md)
