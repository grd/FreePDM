# Howto create a file server
The specs for a file server are minimal. Right now I am using a Raspberry Pi 4 and it works for me, but when you need to have lots of guys then you can use a dedicated server. This are the things you need to do to make it work.

This program uses Ubuntu Server. They support both the Intel and ARM platforms.
First you need to install Ubuntu Server and after that you need to write down the ip address of the server. When you need to have access outside of the LAN you are gonna need port forwarding and also a certificate (certbot for instance). This is the ip address of my server: 10.0.0.11 

servername: 10.0.0.11

path: "/vault/"

usernames:
inlogname1: user1, passwd: passwd1 (use your own loginname and password)
inlogname2: user2, passwd: passwd2 (use your own loginname and password)

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

Testing things out: Download FreePDM with this command: `git pull https://github.com/grd/FreePDM` then go to the directory `tests` and run `python3 fileserver.py` and if it runs then it is done.