!#/bin/bash

printf "Welcome to the server Installer from FreePDM:\n"

# TODO: Install ssl-server
# TODO: Install Web-server - Optional
# TODO: Install SQL-server
# TODO: Install LDAP-server - Optional
# TODO: Install Other dependencies

printf "What do you want to install?\n"

sleep 1

printf "Do you want to install a Webserver? (y / n)\n"

read installwebserver

printf "What SQL server do you want to install?\n
1 - option 1
2 - option 2
3 - option 3\n"

read installsqlserver

printf "Do you want to install a LDAP server?\n"

read installldapserver

# from here start installing

printf "Installing start within a few seconds\n"

sleep 3


# Install of a webserver

if [[ $installwebserver == "y" ]]; then
	printf "Install of a webserver start now\n"
	sleeps 1
fi


# install of SQL server

case $installsqlserver in
	1)
		sqlserver="type1 sql server"
		;;
	2)
		sqlserver="type2 sqlserver"
		;;
	3)
		sqlserver="type3 sqlserver"
		;;
esac

# install LDAP server

if [[ $installldapserver == "y" ]]; then
	printf "Install LDAP server start now\n"
	sleep 1
fi

# Install other dependecies


