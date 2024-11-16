# FreePDM Documentation
The specs for a file server are minimal (a RPI4/5 should work but is not tested yet). At this moment, I am using a desktop and it works for me. If there are multiple users, then consider using a proper server.

## How to install FreePDM

FreePDM is a PDM that works external to FreeCAD. It can interoperate with FreeCAD files. The server and client is written in Go (https://golang.org), because of simplicity and speed.

## Install (for developers on *nix)

#### Install Go
  - Install Go https://golang.org
  - Make sure you have a directory `bin` inside your home directory and inside your `$PATH`
  - Play a bit with GOPATH in your `.bashrc` file (if you have bash). Make sure to have the line `export GOPATH=/home/user` where user is the user name. To reload the `.bashrc` file you can run `source .bashrc`.

#### Install Samba
  1. Install Samba (a FOSS version of the SMB file server protocol). There are plenty of examples in how to install Samba on your computer.
  2. Adding users (one at the time):
  `sudo useradd -M -g sambashare username` (replace username)
  `sudo smbpasswd -a username`
  `sudo smbpasswd -e username`
  `sudo pdbedit -L`

  3. Edit smb.conf with your editor, such as `sudo nano /etc/samba/smb.conf`
  and in the end of this file, add this:
  ```
  [vaults]
   comment = PDM Vaults
   path = /samba/vaults
   browseable = yes
   read only = no
   writable = yes
   guest ok = yes
   valid users = @sambashare user user1
   ```

#TODO: I don't know whether `guest ok = yes` is the right approach because of security.

- `path = xyz` You can put any path here. This is just how I do this.
- `valid users = xyz` This should note @sambashare and the users.
4. Verify your smb.conf with `testparm`
5. Run `sudo systemctl restart smbd`


#### Install FreePDM
  - `git clone github.com/grd/FreePDM.git`
  - Go to directory FreePDM and type `make all`. This should compile and test everything.

#### Modify FreePDM.toml
Read `cat /etc/group` 
Write down the number after `sambashare` (in my case it is 125) and the the `users`.

Next, add the following text to the end of your `FreePDM.toml`. The file is stored (in Ubuntu) in the `.config/FreePDM` directory:

```
[Users]
vault = 125
user1 = 1001
user2 = 1002
```


### Show icons of FreeCAD files

This option provides the possibility to display a screenshot of the model within the FreeCAD file as the file thumbnail.  

Within FreeCAD:  

`Edit` ➡️ `Preferences` ➡️ `General` ➡️ `Document` ➡️ `Save thumbnail into project when saving document`  

Note: leave the 128 number



[<< Previous Chapter](README.md) | [Content Table](README.md) | [Next Chapter >>](Docker-Compose.md)
