# FreePDM
***Documentation***

## Setup VirtualBox

### General

For setting up a VM in VirtualBox there are multiple tutorials out in the wild.
For example use the one from [it's foss](https://itsfoss.com/install-linux-mint-in-virtualbox/).
This one installs mind, but the principles are the same for Debian or RedHat.

### Network

Before you can have acces to your Virtal Machine / Server you have to setup your VM network.
Below a screenshot of network preference page.
You have to set it as a bridged adapter and set it's name to the network adapter that

- your system is using
- a second adapter that has also network acces.

![vbnetworksettings](./figures/VBNetworksettings.png)

In this case the standard network adapter is used.

You can check with [_nmap_](https://nmap.org) if you can see that the VM has acces to your network.  
If you compare both _IP addresses_ (one from _nmap_ and one from [setup static ip address](#setup-static-ip-address)) that both _IP addresses_ are the same.  
Now you can acces your VM via [ssh](#ssh)

![zenmap](./figures/zenmap.png)

### Setup static IP address

Setting up a _static_ _IP address_ is recommended, but not required.
If preferred to leave the _IP address_ as is, see below how to check your _IP address_.
For changing the _IP address_ to _static_ one see for example [this blog](https://www.rosehosting.com/blog/how-to-configure-static-ip-address-on-ubuntu-20-04/) or [here](https://linuxconfig.org/how-to-setup-a-static-ip-address-on-debian-linux).

How to check your _IP address_:  
The command ```ip -a | grep inet``` gives your current (local )_IP address_.
![getipaddress](./figures/get_IPadress.png)
You can acces this your server from your _terminal_ or via _Putty_. 

### ssh

_Ssh_ is a  Secure Shell protocol.
As long as you have installed it on both your server and on you system you should be able to remote login.

Login works as follows:
```ssh USERNAME@###.###.###.###``` then typing in your password, accepting the handshake (if not done before). 
Then you are ready to go.

![ssh_login](./figures/ssh_login.png)

[<< Previous Chapter](commands.md) | [Content Table](README.md) | [Next Chapter >>]()
