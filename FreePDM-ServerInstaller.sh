!#/bin/bash

printf "Welcome to the server Installer from FreePDM:\n"

# TODO: Create install Conf file
# TODO: Install ssl-server
# TODO: Install Web-server - Optional
# TODO: Install SQL-server
# TODO: Install LDAP-server - Optional
# TODO: Install Other dependencies

printf "There are a set of (optional )dependencies that are configured now. This dependecies are:
- A SSL server
- A web server (Optional)
- A SQL server
- A LDAP server (Optional)\n"

sleep 1

printf "Do you want to install a Webserver? (y / n)\n"

read installwebserver

if installwebserver == "y" ]]; then
	printf "What backend do you want to install? (1 - 3)\n
1 - Apache Httpd
2 - Nginx
3 - option 3\n"

	read webservertype

	printf "What python web backend do you want to install? (1 - 4)\n
1 - Django
2 - Pyramid
3 - Falcon
4 - WebPy\n"

	read webserverpython
fi

printf "What SQL server do you want to install? (1 - 3)\n
1 - MySQL
2 - SQLite
3 - PostgreSQL\n"

read installsqlserver

printf "Do you want to install a LDAP server? (y / n)\n"

read installldapserver

if [[ $installldapserver == "y" ]]; then
	:
fi

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


