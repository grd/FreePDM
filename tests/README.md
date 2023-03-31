# Howto create a file server
The specs for a file server are minimal. Right now I am using a Raspberry Pi 4 and it works for me, but when you need to have lots of guys then you can use a dedicated server. This are the things you need to do to make it work.

This program uses Ubuntu Server. They support both the Intel and ARM platforms.
First you need to install Ubuntu Server and after that you need to write down the ip address of the server. When you need to have access outside of the LAN you are gonna need port forwarding and also a certificate (certbot for instance). This is the ip address of my server: 10.0.0.11 

servername: 10.0.0.11

path: "/vault/"

### example usernames and passwords:

login name 1: user1, passwd: passwd1

login name 2: user2, passwd: passwd2


The following commands are all in root (use `sudo su`):

```
groupadd vault
adduser user1
adduser user2
adduser myadmin
usermod -aG vault user1
usermod -aG vault user2
usermod -aG vault myadmin

mkdir /vault
chown vault:vault /vault
chmod g+w vault
```

While still in `ssh` have a look at `/etc/group` and write down the number behind the users `vault`, `user1` and `user2`, but keep in mind to use different names. After you completed the following you need to write down in the `FreePDM.conf` file that is stored (in Ubuntu) in the `.cache` or `.config` directory the following, at the bottom:

```
[user]
vault = 1000
user1 = 1001
user2 = 1002
```

### Testing things out in Ubnutu Linux: 

1. Install pip modules: 
```pip install PySide2
pip install appdirs
pip install defusedxml
```
1. System update: `sudo apt update`
1. Install sshfs: `sudo apt install sshfs`
1. Make a directory : `sudo mkdir /mnt/test`
1. Mount your drive: `sudo sshfs -o allow_other user1@10.0.0.11:/vault /mnt/test`
1. Download FreePDM with this command: `git pull https://github.com/grd/FreePDM` 

### Installing a new Vault

In the FreePDM directory go to the directory `tests` and run `python3 create_new_vault.py`. 

### Testing a vault

In the FreePDM directory go to the directory `tests` and run `python3 fileserver.py` and if it runs then it is done.