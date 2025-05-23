# FreePDM Documentation
Docker-Compose is a very interesting platform because of scaleability and security. Docker is light-weight compared to Virtual Machines and requires a lot less maintenance. You can install Docker on the majority of platforms. Docker is some sort of recipe with components that have the right version, like a script and it compiles pretty fast to an immutable container. Another benefit is that you can put the option to automatically reload a crashed container. The trade-off is that when you are interested in Docker you need to know a lot of Docker. That is why we put docker files inside the FreePDM root directory.

## Install

This is the list of steps that you should take. There are two ways to deal with the file sharing: With Samba (SMB) or with sshfs (SSH). When you run Docker on a VPN, have a look at [External Server](ExternalServer.md). Then you need to have SSH because SMB only works on LAN. SSH is also way more secure than SMB (except with the latest versions). When you use SMB be sure to secure the file sharing. In terms of performance, SMB is marginally faster than SSH. SMB is *the* file share protocol used by Windows, so if you use Windows in a LAN we would recommend using SMB. When you have a small organization that works remote, you need to have SSH. The following is only describing how you can set up a SMB share with Docker.

### Install docker and docker-compose
Follow the steps that are required to install [docker](https://docs.docker.com/engine/install/) and [docker-compose](https://docs.docker.com/compose/install/).

### Preparations
Make sure you have a working samba setup, see [Install](Install.md), and also have a network drive that you already mounted.

Type the following to stop samba:
```
sudo systemctl stop smbd
sudo systemctl disable smbd
```
Go to the `FreePDM` directory. Such as `cd FreePDM`

In the root you can find two files that are important for Docker: `Dockerfile` and `docker-compose.yml`. You only need to modify the `docker-compose.yml` file.

### Note

You should check *all the lines* inside all `environment:` sections!

#### The samba group

The `TZ` field you can modify to your own TimeZone. 

You can change the line `path = /home/user/vaults` to anywhere but you need to think about also need to modify this line:
`- /home/user/vaults:/home/user/vaults`
and look at the install page because you also need to modify the samba file.

### Compile docker-compose

`make docker` This updates the docker container to the latest version.

Delete the docker-compose file:
`make docker_stop` or `make docker_rm`

Test whether the network drive is still working and log in. If that is done then docker works.

### Okay. It compiles. Now what?
Check wether the share works:
```
smbclient -L //localhost/vaults -U yourusername
```
This should result in a prompt, not an error.

Open the file browser. Type
```
smb://localhost/vaults
```
That should work after you log in.

Now you can mount the drive. Look at the install page and then #6.

## Monitoring the output of Docker
With docker you can monitor the output. You can see the logs, which means that you can inspect what a possible crash is or who did what. To use the logs you can type:

`docker logs app` when you want to inspect what is happening inside your your app

`docker logs db` for the db output (this is empty now because db is not working yet)

`docker logs samba` for the samba output




[<< Previous Chapter](Install.md) | [Content Table](README.md) | [Next Chapter >>](SetupVirtualServer.md)
