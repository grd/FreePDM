# FreePDM Documentation
Docker-Compose is a very interesting platform because of immutability, scaleability and security.

## Install Docker-Compose

This is the list of steps that you should take.

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

You can change the lines:
```
- GO_DIRECTORY=/home/user/FreePDM
- GO_WORKDIR=/home/user/FreePDM/cmd/serverapp
```
But you should check *all the lines* inside all `environment:` sections!

#### The samba group

The `TZ` field you can modify to your own TimeZone. 

You can change the line `path = /home/user/vaults` to anywhere but you need to think about also need to modify this line:
`- /home/user/vaults:/home/user/vaults`
and look at the install page because you also need to modify the samba file.

### Compile docker-compose
```
type `make docker`
```

Delete the docker-compose file:
```
type `make docker_stop` or `make docker_rm`
```

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

[<< Previous Chapter](Install.md) | [Content Table](README.md) | [Next Chapter >>](SetupVirtualServer.md)
