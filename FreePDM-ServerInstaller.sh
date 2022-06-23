!# /bin/bash
# :copyright: Copyright 2022 by the FreePDM team
# :license: MIT License.

printf "Welcome to the server Installer from FreePDM:\n"

# TODO: Create install Conf file
# TODO: Install ssh-server
# TODO: Install Web-server - Optional
# TODO: Install SQL-server
# TODO: Install LDAP-server - Optional
# TODO: Install Other dependencies

printf "There are a set of (optional )dependencies that are configured now. This dependecies are:
- A SSH server
- A web server (Optional)
- A SQL server
- A LDAP server (Optional)\n"

# https://stackoverflow.com/questions/4662938/create-text-file-and-fill-it-using-bash
conffile="server.conf"
if [[ -e $conffile ]]; then
	printf "$conffile already exist\n"

	while :
	do
	read -p "Do you want to reuse this configuration file? (y / n)"$'\n' reuseconf

		if [[ $reuseconf == "y" ]]; then
			# Read the file
			setconf="read"

			break
		elif [[ $reuseconf == "n" ]]; then
			oldconf="server$(date +%y%m%d-T%H%M%S).conf"

			$(mv $conffile $oldconf)

			printf "$conffile renamed to $oldconf\n"

			echo > $conffile

			echo "# Config file FreePDM Installer" >> server.conf

			printf "New $conffile created\n"

			setconf="write"
			break
		else
			printf "$reuseconf is not 'y' OR 'n'.\n"
		fi

	done

else
	echo > $conffile

	echo "# Config file FreePDM Installer" >> server.conf

	setconf="write"
fi

sleep 1

# Add line about check for existing server
if [[ $setconf == "read" ]]; then
	# This is most complex because it can switch to read at some point in the process
	:
elif [[ $setconf == "write" ]]; then
	while :
	do
		# -- SQL Server --
		read -p "Do you want to install a Webserver? (y / n)"$'\n' installwebserver

		if [[ $installwebserver == "y" ]]; then

			echo "installwebserver = \"$installwebserver\"" >> server.conf

			printf "What backend do you want to install? (1 - 2)\n
1 - Apache Httpd (default)
2 - Nginx\n"

			read webserverc  # webserver case

			if [[ $webserverc == "" || $webserverc > 2 || $webserverc < 1 ]]; then
				webserverc=1
				printf "http server was emtpty OR outside range and is replaced by it's default. \n"
			fi

			echo "webserverc = $webserverc" >> server.conf

			# python webserver default has to be choosen!
			printf "What python web backend do you want to install? (1 - 4)\n
1 - Django
2 - Pyramid (Default)
3 - Falcon
4 - WebPy\n"

			read webserverpythonc  # webserver python case

			if [[ $webserverpythonc == "" || $webserverpythonc > 4 || $webserverpythonc < 1 ]]; then
				webserverpythonc=2
				printf "python web server was emtpty OR outside range and is replaced by it's default. \n"
			fi

			echo "webserverpythonc = $webserverpythonc" >> server.conf

			read -p "What is your server name?"$'\n' webservername

			echo "webservername = \"$webservername\"" >> server.conf

			read -p "What is your (web )server_domain OR IP address? (default something like web.somename.com)"$'\n' webhostname

			echo "webhostname = \"$webhostname\"" >> server.conf

			# maybe something about admin + password, ports etc
			break

		elif [[ $installwebserver == "n" ]]; then

			echo "installwebserver = \"$installwebserver\"" >> server.conf

			break
			:
		else
			printf "$installwebserver is not 'y' OR 'n'.\n"
		fi
	done

	# -- SQL Server --
	# Add line about check for existing server
	printf "What SQL server do you want to install? (1 - 3)\n
	1 - MySQL
	2 - SQLite
	3 - PostgreSQL(default)\n"

	read sqlserverc  # sqlserver case

	if [[ $sqlserverc == "" || $sqlserverc > 3 || $sqlserverc < 1 ]]; then
		sqlserverc=3
		printf "SQL server was emtpty OR outside range and is replaced by it's default\n"
	fi

	echo "sqlserverc = $sqlserverc" >> server.conf

	read -p "Enter SQL Username:"$'\n' sqlservername

	echo "sqlservername = \"$sqlservername\"" >> server.conf

	read -p "What is your (sql )server_domain OR IP address? (default something like sql.somename.com)"$'\n' sqlhostname

	echo "sqlhostname = \"$sqlhostname\"" >> server.conf

	# maybe something about admin + password, ports etc

	# -- LDAP Server --
	# Add line about check for existing server

	while :
	do

		read -p "Do you want to install a LDAP server? (y / n)"$'\n' installldapserver

		if [[ $installldapserver == "y" ]]; then

			echo "installldapserver = \"$installldapserver\"" >> server.conf

			printf "What LDAP server do you want to install? (1 - 4)\n
1 - open LDAP (Default)
2 - Apache DS
3 - openDJ
4 - 389 Directory server\n"

			read ldapserverc

			if [[ $ldapserverc == "" || $ldapserverc > 4 || $ldapserverc < 1  ]]; then
				ldapserverc=1
				printf "LDAP server was emtpty OR outside range and is replaced by it's default. \n"
			fi

			echo "ldapserverc = $ldapserverc" >> server.conf

			# user name and passwords are not stored
			read -p "Enter LDAP Username:"$'\n' ldapusername

			# read -sp "Enter LDAP Password:" ldappw1  # Silent
			read -p "Enter LDAP Password:"$'\n' -s ldappw1  # With asterix

			read -p "Re-enter LDAP Password:"$'\n' -s ldappw2  # With asterix

			break

		elif [[ $installldapserver == "n" ]]; then
			echo "installldapserver = \"$installldapserver\"" >> server.conf

			break

		else
			printf "$installldapserver is not 'y' OR 'n'.\n"
		fi

	done

else
	:
fi

# Show cofiguration summery


# from here start installing
printf "Installing start within a few seconds\n"

sleep 3

printf "Update repositories.\n"
sudo apt update
# printf "Upgrade repositories.\n"  # upgrade don't work yet
# sudo apt upgrade

# Install of a SSH server
testcommand="ssh"
packages="openssh-server"
# if ! [[ $(command -v $the_command) &> /dev/null ]]; then
if ! [[ $(command -v $testcommand) ]]; then
  printf "$testcommand could not be found.\n$packages shall be installed. \n"
	sudo apt install -y $packages
	exit
else
	printf "$packages already installed\n"
fi

# Install of a webserver

if [[ $installwebserver == "y" ]]; then

	case $webserverc in
		1)
			# https://ubuntu.com/tutorials/install-and-configure-apache#1-overview
			webserver="Apache httpd"
			testcommand="apache2"  # should also work with apachectl -v
			packages="apache2"
			;;
		2)
			# https://ubuntu.com/tutorials/install-and-configure-nginx#1-overview
			webserver="Nginx"
			testcommand="nginx"
			packages="nginx"
			;;
		3)
			# https://www.hostinger.com/tutorials/how-to-install-tomcat-on-ubuntu/
			webserver="Appache tomcat"  # Java
			testcommand=""
			packages=""
				;;
	esac

	printf "The following Web server shall be installed: $webserver."
	sleep 1

	if ! [[ $(command -v $testcommand) ]]; then
	  printf "$testcommand could not be found.\n$packages shall be installed. \n"
		sudo apt install -y $packages
		exit
	else
		printf "$packages already installed \n"
	fi

	# If statement for IPaddres has always same length and set of dots on same place

	case $webserverpythonc in
		# Make use of package manager OR Virtual Environment...?
		1)
			# https://www.digitalocean.com/community/tutorials/how-to-install-the-django-web-framework-on-ubuntu-20-04
			webserver="Django"
			testcommand=""  # "django-admin --version"
			packages=""  # "python3-django"
			;;
		2)
			# https://www.digitalocean.com/community/tutorials/how-to-use-the-pyramid-framework-to-build-your-python-web-app-on-ubuntu
			webserver="Pyramid"
			testcommand=""
			packages=""
			;;
		3)
			# https://www.digitalocean.com/community/tutorials/how-to-deploy-falcon-web-applications-with-gunicorn-and-nginx-on-ubuntu-16-04
			webserver="Falcon"
			testcommand=""
			packages=""
			;;
		4)
			webserver="WebPy"
			testcommand=""
			packages=""  # "python-webpy"
			;;
	esac

	# if ! (( command -V $testcommand )); then  #
	# 	printf "$testcommand could not be found.\n$packages shall be installed. \n"
	# 	sudo apt install -y $packages
	# 	exit
	# else
	# 	printf "$packages already installed \n"
	# fi

	printf "The following Python Web server shall be installed: $webserverpython."
	sleep 1

fi


# install of SQL server

# Check if SQL server already exist. if yes  add database to existing server?
# work only with selected sql server

case $sqlserverc in
	1)
		sqlserver="MySQL"
		testcommand="mysql"
		packages="mysql-server mysql-client"
		# mysql is replaced by mariadb see: https://wiki.debian.org/MySql
		# https://www.digitalocean.com/community/tutorials/how-to-install-the-latest-mysql-on-debian-10
		;;
	2)
		sqlserver="SQLite"
		testcommand="sqlite3"  # can also be sqlite3 --version
		packages="sqlite3 uwsgi-plugin-sqlite3 SQLitebrowser"
		# sqldiff is in unstable # https://manpages.debian.org/unstable/sqlite3/sqldiff.1.en.html
		# sqlite3_analyzer is in unsable # https://manpages.debian.org/unstable/sqlite3-tools/sqlite3_analyzer.1.en.html
		;;
	3)
		# https://www.geeksforgeeks.org/install-postgresql-on-linux/
		# https://sqlserverguides.com/postgresql-installation-on-linux/
		sqlserver="postgreSQL"
		testcommand="postgres"
		packages="postgresql postgresql-contrib"

		# from: https://www.postgresql.org/download/linux/debian/
		# Create the file repository configuration:
		printf "Add postgreSQL repository to list.\n"
		sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'

		# Import the repository signing key:
		wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
		;;
esac

printf "The following SQL server shall be installed: $sqlserver.\n"

printf "Update packages.\n"
sudo apt-get update

if ! [[ $(command -v $testcommand) ]]; then
	printf "$testcommand could not be found.\n$packages shall be installed. \n"
	sudo apt install -y $packages
	exit
else
	printf "$packages already installed \n"
fi

# install LDAP server

if [[ $installldapserver == "y" ]]; then

	case $ldapserverc in
		# Basically all are Java implementations except 389 directory service
		1)
			# https://www.howtoforge.com/how-to-install-openldap-on-debian-11/
			ldapserver="open LDAP"
			testcommand="slapd"  # https://serverfault.com/questions/839948/how-to-check-the-version-of-openldap-installed-in-command-line
			packages="slapd ldap-utils"
			;;
		2)
			ldapserver="Apache DS"
			testcommand=""
			packages="apacheds"
			;;
		3)
			# https://backstage.forgerock.com/docs/opendj/2.6/install-guide/
			ldapserver="OpenDJ"
			testcommand=""
			packages=""
			;;
		4)
			# https://directory.fedoraproject.org/docs/389ds/howto/howto-debianubuntu.html
			ldapserver="389 Directory server"
			testcommand=""
			packages="termcap-compat apache2-mpm-worker"
			;;
	esac

	printf "The following LDAP server shall be installed: $ldapserver."
	sleep 1
fi

if ! [[ $(command -v $testcommand) ]]; then
	printf "$testcommand could not be found.\n$packages shall be installed. \n"
	sudo apt install -y $packages
	exit
else
	printf "$packages already installed \n"
fi

# Install other dependecies
