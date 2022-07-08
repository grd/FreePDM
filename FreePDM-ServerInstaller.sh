!# /bin/bash
# :copyright: Copyright 2022 by the FreePDM team
# :license: MIT License.

printf "Welcome to the server Installer from FreePDM:\n"

# TODO: Create install Conf file
# TODO: Install SQL-server
# TODO: Install Web-server - Optional
# TODO: Install LDAP-server - Optional
# TODO: Install Other dependencies
# TODO: Test Test Test...

printf "There are a set of (optional )dependencies that are configured now. This dependecies are:
- A SSH server
- A SQL server
- A web server (Optional)
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
			# https://www.geeksforgeeks.org/bash-scripting-how-to-read-a-file-line-by-line/
			setconf="read"

			break
		elif [[ $reuseconf == "n" ]]; then
			oldconf="server$(date +%y%m%d-T%H%M%S).conf"

			$(mv $conffile $oldconf)

			printf "$conffile renamed to $oldconf\n"

			echo > $conffile
			echo "!# /bin/bash" >> server.conf

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
	echo "!# /bin/bash" >> server.conf

	echo "# Config file FreePDM Installer" >> server.conf

	setconf="write"
fi

sleep 1

# Add line about check for existing server
if [[ $setconf == "read" ]]; then
	# This is most complex because it can switch to read at some point in the process
	:
elif [[ $setconf == "write" ]]; then

	# -- SQL Server --
	echo "# -- SQL server --" >> server.conf

	# Add line about check for existing server
	printf "What SQL server do you want to install? (1 - 5)\n
1 - MySQL Community Edition
2 - SQLite
3 - PostgreSQL(default)
4 - MariaDB
5 - Other\n"

	read sqlserverc  # sqlserver case

	if [[ $sqlserverc == "" || $sqlserverc > 4 || $sqlserverc < 1 ]]; then
		sqlserverc=5
		printf "SQL server was emtpty OR outside range and is replaced by it's default\n"
	fi

	case $sqlserverc in
		1)
			sqlserver="MySQL"
			sqltestcommand="mysql"
			# packages="mysql-server mysql-client"
			sqlpackages="mysql-community-server"  # as alternative see link installation
			defaultportnumber=3306  # https://kinsta.com/knowledgebase/mysql-port/
			# mysql is replaced by mariadb see: https://wiki.debian.org/MySql
			# https://www.digitalocean.com/community/tutorials/how-to-install-the-latest-mysql-on-debian-10
			;;
		2)
			sqlserver="SQLite"
			sqltestcommand="sqlite3"  # can also be sqlite3 --version
			sqlpackages="sqlite3 uwsgi-plugin-sqlite3 SQLitebrowser"
			defaultportnumber=""
			# sqldiff is in unstable # https://manpages.debian.org/unstable/sqlite3/sqldiff.1.en.html
			# sqlite3_analyzer is in unsable # https://manpages.debian.org/unstable/sqlite3-tools/sqlite3_analyzer.1.en.html
			;;
		3)
			# https://www.geeksforgeeks.org/install-postgresql-on-linux/
			# https://sqlserverguides.com/postgresql-installation-on-linux/
			sqlserver="postgreSQL"
			sqltestcommand="postgres"
			sqlpackages="postgresql postgresql-contrib"
			defaultportnumber=5432
			;;
		4)
			printf "What MariaDB version do you want to install? (1 - 2)\n
1 - MariaDB Enterprice Edition(default)
2 - MariaDB Community Edition\n"

			read mariadbedition

			if [[ $mariadbedition == "" || $mariadbedition > 2 || $mariadbedition < 1 ]]; then
				sqlserverc=1
				printf "MariaDB Edition was emtpty OR outside range and is replaced by it's default\n"
			fi

			sqlserver="mariadb"
			sqltestcommand="mariadb"  # can also be sqlite3 --version
			sqlpackages=""
			defaultportnumber="3306"
			;;
		5)
			printf "For custom SQL servers You have to install them selfes.\n"  # Add extra info later on
			;;
	esac

	echo "sqlserverc = $sqlserverc" >> server.conf

	read -p "Enter SQL Servername:"$'\n' sqlservername

	echo "sqlservername = \"$sqlservername\"" >> server.conf

	read -p "What is your (sql )server_domain OR IP address? (default something like sql.somename.com)"$'\n' sqlhostname
	# Can i check if something is an IP addres?, Is there a need for?
	# https://stackoverflow.com/questions/23675400/validating-an-ip-address-using-bash-script

	echo "sqlhostname = \"$sqlhostname\"" >> server.conf

	read -p "The default portnumber is $defaultportnumber. Do you want to change it?(for current port leave empty)"$'\n' sqlportnumber

	if [[ $sqlportnumber == "" ]]; then
		portnuber=$defaultportnumber
	fi

	echo "sqlportnumber = \"$sqlportnumber\"" >> server.conf

	read -p "What is the Database root directory?"$'\n' -r sqlrootdirectory  # Adding default location? second disc?

	echo "sqlrootdirectory = \"$sqlrootdirectory\"" >> server.conf

	read -p "What is your Database name?"$'\n' sqldatabasename

	echo "sqldatabasename = \"$sqldatabasename\"" >> server.conf

	read -p "What is your Database admin(root acces) name?(for current user leave empty)"$'\n' sqldatabaseadmin

	if [[ $sqldatabaseadmin == "" ]]; then
		sqldatabaseadmin=$(whoami)
	fi

	# Don't save admin name to file
	# echo "sqldatabaseadmin = \"$sqldatabaseadmin\"" >> server.conf

	# printf "What is your Database user password?\n"
	read -p "What is your Database admin password?"$'\n' -s -r sqldatabaseapassword

	# Don't save admin password to file
	# echo "sqldatabaseapassword = \"$sqldatabaseapassword\"" >> server.conf

	read -p "What is your Database user name?"$'\n' sqldatabaseuser

	# Don't save user name to file
	# echo "sqldatabaseuser = \"$sqldatabaseuser\"" >> server.conf

	# printf "What is your Database user password?\n"
	read -p "What is your Database user password?"$'\n' -s -r sqldatabaseupassword

	# Don't save user password to file
	# echo "sqldatabaseupassword = \"$sqldatabaseupassword\"" >> server.conf

	# -- Web Server --
	echo "# -- Web server --" >> server.conf

	while :
	do

		read -p "Do you want to install a Webserver? (y / n)"$'\n' installwebserver

		if [[ $installwebserver == "y" ]]; then

			echo "installwebserver = \"$installwebserver\"" >> server.conf

			printf "A Web Interface is not implemented yet. Choose this only if you want to create one by yourself.
If so Feel free to create a Pull Request.\n"

			printf "What backend do you want to install? (1 - 2)\n
1 - Apache Httpd (default)
2 - Nginx\n"

			read webserverc  # webserver case

			if [[ $webserverc == "" || $webserverc > 2 || $webserverc < 1 ]]; then
				webserverc=1
				printf "http server was emtpty OR outside range and is replaced by it's default. \n"
			fi

			echo "webserverc = $webserverc" >> server.conf

			case $webserverc in
				1)
					# https://ubuntu.com/tutorials/install-and-configure-apache#1-overview
					webserver="Apache httpd"
					webtestcommand="apache2"  # should also work with apachectl -v
					webpackages="apache2"
					;;
				2)
					# https://ubuntu.com/tutorials/install-and-configure-nginx#1-overview
					webserver="Nginx"
					webtestcommand="nginx"
					webpackages="nginx"
					;;
				3)
					# https://www.hostinger.com/tutorials/how-to-install-tomcat-on-ubuntu/
					webserver="Appache tomcat"  # Java
					webtestcommand=""
					webpackages=""
						;;
			esac

			# python webserver default has to be choosen!
			printf "What python web backend do you want to install? (1 - 5)\n
1 - Django
2 - Pyramid (Default)
3 - Falcon
4 - WebPy
5 - None\n"

			read webserverpythonc  # webserver python case

			if [[ $webserverpythonc == "" || $webserverpythonc > 5 || $webserverpythonc < 1 ]]; then
				webserverpythonc=2
				printf "python web server was emtpty OR outside range and is replaced by it's default. \n"
			fi

			echo "webserverpythonc = $webserverpythonc" >> server.conf

			case $webserverpythonc in
				# Make use of package manager OR Virtual Environment...?
				1)
					# https://www.digitalocean.com/community/tutorials/how-to-install-the-django-web-framework-on-ubuntu-20-04
					webserver="Django"
					webptestcommand="django-admin --version"
					webppackages="django"
					;;
				2)
					# https://www.digitalocean.com/community/tutorials/how-to-use-the-pyramid-framework-to-build-your-python-web-app-on-ubuntu
					webserver="Pyramid"
					webptestcommand=""  # added later
					webppackages="pyramid"
					;;
				3)
					# https://www.digitalocean.com/community/tutorials/how-to-deploy-falcon-web-applications-with-gunicorn-and-nginx-on-ubuntu-16-04
					webserver="Falcon"
					webptestcommand=""
					webppackages="cython falcon gunicorn"
					;;
				4)
					webserver="WebPy"
					webptestcommand=""
					webppackages="web.py"  # "python-webpy"
					;;
				5)
					# No Python webserver
					webserver=""
					webptestcommand=""
					webppackages=""
					;;
			esac

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

	# -- LDAP Server --
	echo "# -- LDAP server --" >> server.conf

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

			case $ldapserverc in
				# Basically all are Java implementations except 389 directory service
				1)
					# https://www.howtoforge.com/how-to-install-openldap-on-debian-11/
					ldapserver="open LDAP"
					ldaptestcommand="slapd"  # https://serverfault.com/questions/839948/how-to-check-the-version-of-openldap-installed-in-command-line
					ldappackages="slapd ldap-utils"
					ldapportnumber=389
					;;
				2)
					ldapserver="Apache DS"
					ldaptestcommand=""
					ldappackages="apacheds"
					ldapportnumber=389
					;;
				3)
					# https://backstage.forgerock.com/docs/opendj/2.6/install-guide/
					ldapserver="OpenDJ"
					ldaptestcommand=""
					ldappackages=""
					ldapportnumber=389
					;;
				4)
					# https://directory.fedoraproject.org/docs/389ds/howto/howto-debianubuntu.html
					ldapserver="389 Directory server"
					ldaptestcommand=""
					ldappackages="termcap-compat apache2-mpm-worker"
					ldapportnumber=389
					;;
			esac

			read -p "ldap organisation name, o=" o

			echo "o = \"$o\"" >> server.conf

			read -p "ldap domain component, dc=" dcfulldomain

			# https://linuxize.com/post/bash-concatenate-strings/
			# https://www.geeksforgeeks.org/bash-scripting-split-string/
			# Set space as the delimiter
			IFS='.'

			# Read the split words into an array based on dot delimiter
			read -ra newarr <<< "$dcfulldomain"

			# concanete each value of the array by using the loop
			# lastelement=${newarr[(( $length - 1 ))]}
			printf "$lastelement\n"

			var=""
			for element in "${newarr[@]}";
			do
			  if [[ $element == ${newarr[(( $length - 1 ))]} ]]; then
			    var+="dc=$element"
			  else
			    var+="dc=$element,"
			  fi
			done
			printf "$dcfulldomain is replaced by $dc.\n"

			echo "dc = \"$dc\"" >> server.conf

			read -p "ldap common name, cn=" cn

			echo "cn = \"$cn\"" >> server.conf

			# Admin name and passwords are not stored
			read -p "Enter LDAP admin name:"$'\n' ldapadminname

			# Don't save admin name to file
			# echo "ldapadminname = \"$ldapadminname\"" >> server.conf

			# read -sp "Enter LDAP Password:" ldappw1  # Silent
			read -p "Enter LDAP admin password:"$'\n' -s -r ldapadminpassword1  # With asterix

			if [[ $sqldatabaseadmin == "" ]]; then
				sqldatabaseadmin=$(whoami)
			fi

			# Don't save user password to file
			# echo "ldapadminpassword1 = \"ldapadminpassword1\"" >> server.conf

			read -p "Enter LDAP username:"$'\n' ldapusername

			# read -sp "Enter LDAP Password:" ldappw1  # Silent
			read -p "Enter LDAP user password:"$'\n' -s -r ldapuserpassword1  # With asterix

			# Don't save user password to file
			# echo "ldapuserpassword1 = \"ldapuserpassword1\"" >> server.conf

			read -p "Re-enter LDAP user password:"$'\n' -s -r ldapuserpassword2  # With asterix

			# Don't save user password to file
			# echo "ldapuserpassword2 = \"ldapuserpassword2\"" >> server.conf

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


# Install of a SSH server if not available
printf "Install SSH-server\n"

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


# Install of python3 and pip3 if not available
printf "Install Python and pip\n"

testcommand="python3"
packages="python3"
# if ! [[ $(command -v $the_command) &> /dev/null ]]; then
if ! [[ $(command -v $testcommand) ]]; then
  printf "$testcommand could not be found.\n$packages shall be installed. \n"
	sudo apt install -y $packages
	exit
else
	printf "$packages already installed\n"
fi

testcommand="pip3"
packages="python3-pip"
# if ! [[ $(command -v $the_command) &> /dev/null ]]; then
if ! [[ $(command -v $testcommand) ]]; then
  printf "$testcommand could not be found.\n$packages shall be installed. \n"
	sudo apt install -y $packages
	pip install --no-warn-script-location pip_search
	exit
else
	printf "$packages already installed\n"
fi

# install of SQL server

# Check if SQL server already exist. if yes  add database to existing server?
# work only with selected sql server

printf "The following SQL server shall be installed: $sqlserver.\n"

case $sqlserverc in
	1)
		# MySQL
		# https://www.linuxcapable.com/how-to-install-the-latest-mysql-8-on-debian-11/
		printf "Update packages AND upgrade packages.\n"
		sudo apt update && sudo apt upgrade

		printf "Install required packages for MySQL.\n"
		sudo apt install software-properties-common apt-transport-https wget ca-certificates gnupg2 -y

		# Import the repository signing key:
		printf "Add MySQL repository signing key.\n"
		sudo wget -O- http://repo.mysql.com/RPM-GPG-KEY-mysql-2022 | gpg --dearmor | sudo tee /usr/share/keyrings/mysql.gpg

		printf "Add MySQL Repository.\n"
		echo 'deb [signed-by=/usr/share/keyrings/mysql.gpg] http://repo.mysql.com/apt/debian bullseye mysql-8.0' | sudo tee /etc/apt/sources.list.d/mysql.list

		printf "Update packages.\n"
		sudo apt update

		printf "Install MySQL.\n"
		;;
	2)
		# SQLite
		printf "Update packages.\n"
		sudo apt update
		;;
	3)
		# PostgreSQL
		# from: https://www.postgresql.org/download/linux/debian/
		# Create the file repository configuration:
		printf "Add postgreSQL repository.\n"
		sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'

		# Import the repository signing key:
		printf "Add Repository signing key "
		wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -

		printf "Update packages.\n"
		sudo apt update

		printf "Install PostgreSQL.\n"
		;;
	4)
		# MariaDB
		# https://mariadb.com/docs/deploy/deployment-methods/repo/
		printf "Update packages AND upgrade packages.\n"
		sudo apt update && sudo apt upgrade

		printf "Install required packages for MariaDB.\n"
		sudo apt install wget apt-transport-https

		case $mariadbedition in
			1 )
				# Enterprice Edition
				wget https://dlm.mariadb.com/enterprise-release-helpers/mariadb_es_repo_setup

				echo "cfcd35671125d657a212d92b93be7b1f4ad2fda58dfa8b5ab4b601bf3afa4eae  mariadb_es_repo_setup" | sha256sum -c -

				chmod +x mariadb_es_repo_setup

				sudo ./mariadb_es_repo_setup --token="CUSTOMER_DOWNLOAD_TOKEN" --apply
				;;
			2)
				# Community Edition
				wget https://downloads.mariadb.com/MariaDB/mariadb_repo_setup

				echo "d4e4635eeb79b0e96483bd70703209c63da55a236eadd7397f769ee434d92ca8  mariadb_repo_setup" | sha256sum -c -

				chmod +x mariadb_repo_setup

				sudo ./mariadb_repo_setup
				;;
		esac
		;;
	5)
		printf "assumed a SQL server is already installed"
		;;
esac

testcommand=$sqltestcommand
packages=$sqlpackages

if ! [[ $(command -v $testcommand) ]]; then
	printf "$testcommand could not be found.\n$packages shall be installed. \n"
	sudo apt install -y $packages
	exit
else
	printf "$packages already installed \n"
fi

printf "configure $packages.\n"

case $sqlserverc in
	1)
		# MySQL
		# https://www.linuxcapable.com/how-to-install-the-latest-mysql-8-on-debian-11/

		printf "Check policiy\n"

		apt policy mysql-community-server

		printf "Enable My-SQL at startup\n"

		sudo systemctl enable mysql

		if [[ $(systemctl is-active mysql) == "inactive" ]]; then
			sudo systemctl start mysql
		else
			sudo systemctl restart mysql
		fi

		# https://fedingo.com/how-to-automate-mysql_secure_installation-script/
		# sudo mysql_secure_installation  # Make use of mysql_secure_installation?

		sudo bash -c 'echo "bind-address=127.0.0.1" >> /etc/mysql/mysql.conf.d/mysqld.cnf'

		# https://stackoverflow.com/questions/33470753/create-mysql-database-and-user-in-bash-script

		# Add sqlhostname to hosts file
		printf "Add sql host name to host file.\n"

		sudo bash -c '127.0.0.1 	$sqlhostname" >> /etc/hosts'

		printf "Create new user.\n"

		# login as admin
		# https://stackoverflow.com/questions/33470753/create-mysql-database-and-user-in-bash-script
		sudo mysql -u $sqldatabaseadmin -p${sqldatabaseapassword}

		# create new user
		# https://stackoverflow.com/questions/37239970/connect-to-mysql-server-without-sudo
		CREATE USER $sqldatabaseuser@$sqlhostname IDENTIFIED BY $sqldatabaseupassword;

		# Give database privilidges
		GRANT ALL PRIVILEGES ON database_name.* TO $sqldatabaseuser@$sqlhostname;

		EXIT

		# Login as user and create database

		printf "Login as user and create database\n"

		mysql -u $sqldatabaseuser -p${sqldatabaseupassword}

		CREATE DATABASE $sqldatabasename CHARACTER SET utf8;

		exit
		;;
	2)
		# SQLite
		:
		;;
	3)
		# PostgreSQL
		# from: https://www.postgresql.org/download/linux/debian/
		# Create the file repository configuration:
		:
		;;
	4)
		# MariaDB
		# https://mariadb.com/docs/deploy/deployment-methods/repo/
		:
		;;
	5)
		:
		;;
esac

if [[ $(ufw -v) ]]; then
	ufw allow $sqlportnumber
	exit
elif [[ $(iptables -v) ]]; then
	iptables –A INPUT –p tcp –dport $sqlportnumber –j ACCEPT
else
	printf "No firewall is installed \n"
fi

# Install of a webserver

if [[ $installwebserver == "y" ]]; then

	case $webserverc in
		1)
			:
			;;
		2)
			:
			;;
		3)
			:
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
			:
			;;
		2)
			:
			;;
		3)
			:
			;;
		4)
			:
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

# install LDAP server
# printf command below case
if [[ $installldapserver == "y" ]]; then
	# https://likegeeks.com/linux-ldap-server/
	# phpldapadmin is not default available on debian
	# https://kifarunix.com/install-phpldapadmin-on-debian-10-debian-11/

	case $ldapserverc in
		# Basically all are Java implementations except 389 directory service(c++)
		1)
			:
			;;
		2)
			:
			;;
		3)
			:
			;;
		4)
			:
			;;
	esac

	printf "The following LDAP server shall be installed: $ldapserver."
	sleep 1

	if ! [[ $(command -v $testcommand) ]]; then
		printf "$testcommand could not be found.\n$packages shall be installed. \n"
		sudo apt install -y $packages
		exit
	else
		printf "$packages already installed \n"
	fi

	printf "Enable ldap.\n"

	systemctl enable slapd

fi

# Install other dependecies
