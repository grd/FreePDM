# FreePDM Documentation
Docker-Compose is a very interesting platform because of scaleability and security. In this case I am afraid it doesn't really suit very well because it is complicated and also right now it doesn't have the database options yet.

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
Go to the `FreePDM/docker` directory.

Open the file `docker-compose.yml`. The TZ field you can modify to your own TimeZone. Right now the last line inside the volumes section is fixed and hard coded. This should of course in the future be modifyable.

Type `make start`

Test whether the network drive is still working and log in. If that is done then docker works.

[<< Previous Chapter](Install.md) | [Content Table](README.md) | [Next Chapter >>](SetupVirtualServer.md)
