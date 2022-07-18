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
			read -p "SQLite don't have a default connection with LDAP. Are you sure you want to continue?[y/n(n is exit)]"$'\n' sqlitecontinue

			if [[ $sqlitecontinue == "y" ]]; then
				sqlserver="SQLite"
				sqltestcommand="sqlite3"  # can also be sqlite3 --version
				sqlpackages="sqlite3 uwsgi-plugin-sqlite3 SQLitebrowser"
				defaultportnumber=""
				# sqldiff is in unstable # https://manpages.debian.org/unstable/sqlite3/sqldiff.1.en.html
				# sqlite3_analyzer is in unsable # https://manpages.debian.org/unstable/sqlite3-tools/sqlite3_analyzer.1.en.html
			else
				exit
			fi
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

	# Default user name after installation: MySQL == root, SQLite==None, PostgreSQL == postgres, MariaDB == root
	# What to do? Database Admin is the default user? Or Add extra admin user?
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

			read -p "What is your IP address? "$'\n' webip

			echo "webip = \"$webip\"" >> server.conf

			read -p "What is your (web )server_domain name? (default something like web.somename.com)"$'\n' webhostname

			echo "webhostname = \"$webhostname\"" >> server.conf

			# maybe something about admin + password, ports etc
			break

		elif [[ $installwebserver == "n" ]]; then

			echo "installwebserver = \"$installwebserver\"" >> server.conf

			break

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
					printf "When installing OpenLDAP also an Admin Password it requested.\n"
					;;
				2)
					printf "It is recomended to install Apache Directory Studio too.
					Install Eclipse IDE and then Install Apache directory Studio via the Marketplace.\n"
					ldapserver="Apache DS"
					ldaptestcommand=""
					ldappackages="apacheds"
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
					ldappackages="termcap-compat apache2-mpm-worker 389-ds cockpit-389-ds python3-lib389"
					ldapportnumber=389
					;;
			esac

			if [[ $ldapserverc ne 3 ]]; then
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
				# maybe remove this question later on. Inside OpenLDAP a password is requested from menu.

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
			else

				# OpenDJ give his own configuration setup so no configuration from here...
				printf "OpenDJ give it's own configuration setup.\n"

			fi

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
		sudo apt install curl wget apt-transport-https

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

printf "configure $sqlserver.\n"

case $sqlserverc in
	1)
		# MySQL
		# https://www.linuxcapable.com/how-to-install-the-latest-mysql-8-on-debian-11/

		printf "Check policiy\n"

		apt policy mysql-community-server

		printf "Enable $sqlserver at startup\n"

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
		GRANT ALL PRIVILEGES ON $sqldatabasename.* TO $sqldatabaseuser@$sqlhostname;  # Will this work first create user priviledges without existing db...

		EXIT

		# Login as user and create database
		# https://linuxwebdevelopment.com/how-to-create-a-new-database-in-mysql/

		printf "Login as user and create database\n"

		mysql -u $sqldatabaseuser -p${sqldatabaseupassword}

		CREATE DATABASE $sqldatabasename CHARACTER SET utf8;

		EXIT
		;;
	2)
		# SQLite
		:
		;;
	3)
		# PostgreSQL
		# from: https://www.postgresql.org/download/linux/debian/
		# Create the file repository configuration:
		printf "Check policiy\n"

		apt policy postgresql

		printf "Enable $sqlserver at startup\n"

		sudo systemctl enable postgresql

		if [[ $(systemctl is-active postgresql) == "inactive" ]]; then
			sudo systemctl start postgresql
		else
			sudo systemctl restart postgresql
		fi

		# https://stackoverflow.com/questions/33470753/create-mysql-database-and-user-in-bash-script

		# Add sqlhostname to hosts file
		printf "Add sql host name to host file.\n"

		sudo bash -c '127.0.0.1 	$sqlhostname" >> /etc/hosts'

		printf "Create new user.\n"

		# login as admin
		# https://dba.stackexchange.com/questions/198125/sudo-u-postgres-psql-postgres-does-not-ask-for-password
		sudo -u postgres psql postgres

		# create new user
		# https://stackoverflow.com/questions/37239970/connect-to-mysql-server-without-sudo
		# https://www.postgresql.org/docs/current/app-createuser.html
		# https://www.postgresql.org/docs/current/sql-createuser.html
		CREATE USER $sqldatabaseadmin WITH PASSWORD $sqldatabaseapassword;
		ALTER ROLE $sqldatabaseadmin ADMIN;

		# Give database privilidges
		GRANT ALL PRIVILEGES ON $sqldatabasename.* TO $sqldatabaseuser;  # @$sqlhostname
		GRANT ALL PRIVILEGES ON $sqldatabasename.* TO $sqldatabaseadmin;  # @$sqlhostname

		FLUSH PRIVILEGES;

		# Login as user and create database

		printf "Login as user and create database\n"

		CREATE USER $sqldatabaseuser WITH PASSWORD $sqldatabaseupassword;
		# ALTER ROLE $sqldatabaseadmin ADMIN;

		# Encodings: https://www.postgresql.org/docs/14/multibyte.html
		CREATE DATABASE $sqldatabasename WITH ENCODING UTF8;

		\q
		;;
	4)
		# MariaDB
		# https://mariadb.com/docs/deploy/deployment-methods/repo/

		printf "Check policiy\n"

		apt policy mariadb

		printf "Enable $sqlserver at startup\n"

		sudo systemctl enable mariadb

		if [[ $(systemctl is-active mariadb) == "inactive" ]]; then
			sudo systemctl start mariadb
		else
			sudo systemctl restart mariadb
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
		sudo mariadb -u $sqldatabaseadmin -p${sqldatabaseapassword}

		# create new user
		# https://stackoverflow.com/questions/37239970/connect-to-mysql-server-without-sudo
		CREATE USER $sqldatabaseuser@$sqlhostname IDENTIFIED BY $sqldatabaseupassword;

		# Give database privilidges
		GRANT ALL PRIVILEGES ON $sqldatabasename.* TO $sqldatabaseuser@$sqlhostname;  # Will this work first create user priviledges without existing db...

		EXIT

		# Login as user and create database
		# https://linuxwebdevelopment.com/how-to-create-a-new-database-in-mysql/

		printf "Login as user and create database\n"

		mysql -u $sqldatabaseuser -p${sqldatabaseupassword}

		CREATE DATABASE $sqldatabasename CHARACTER SET utf8;

		EXIT

		;;
	5)
		printf "No SQL server is installed.\n"
		;;
esac

# Open port number firewall

printf "Open firewall port: $sqlportnumber for sql.\n"

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

	printf "The following Web server shall be installed: $webserver."
	sleep 1

	testcommand=$webtestcommand
	packages=$webpackages

	if ! [[ $(command -v $testcommand) ]]; then
	  printf "$testcommand could not be found.\n$packages shall be installed. \n"
		sudo apt install -y $packages
		exit
	else
		printf "$packages already installed \n"
	fi

	printf "Create basic site.\n"

	# Set space as the delimiter
	IFS='.'

	# Read the split words into an array based on dot delimiter
	read -ra newarr <<< "$webhostname"

	# concanete each value of the array by using the loop

	for element in "${newarr[@]}";
	do
		if [[ $element == "www" ]]; then
			:
		else
			basicsite="$element"
		fi
	done

	indexfile="<html>
<head>
  <title> Test FreePDM </title>
</head>
<body>
  <p> Welcome to FreePDM!
</body>
</html>"

	docroot="/var/www/$basicsite/"
	sudo mkdir docroot

	cd /var/www/$basicsite/

	echo > "index.html"

	echo $indexfile >> "index.html"

	case $webserverc in
		1)
			# Apache httpd
			# https://ubuntu.com/tutorials/install-and-configure-apache

			printf "Create virtual host conf file.\n"

			vhostfile="<VirtualHost *:80>
				# The ServerName directive sets the request scheme, hostname and port that
				# the server uses to identify itself. This is used when creating
				# redirection URLs. In the context of virtual hosts, the ServerName
				# specifies what hostname must appear in the request's Host: header to
				# match this virtual host. For the default virtual host (this file) this
				# value is not decisive as it is used as a last resort host regardless.
				# However, you must set it for any further virtual host explicitly.
				ServerName $webhostname

				ServerAdmin webmaster@localhost
				DocumentRoot $docroot

				# Available loglevels: trace8, ..., trace1, debug, info, notice, warn,
				# error, crit, alert, emerg.
				# It is also possible to configure the loglevel for particular
				# modules, e.g.
				#LogLevel info ssl:warn

				ErrorLog ${APACHE_LOG_DIR}/error.log
				CustomLog ${APACHE_LOG_DIR}/access.log combined

				# For most configuration files from conf-available/, which are
				# enabled or disabled at a global level, it is possible to
				# include a line for only one particular virtual host. For example the
				# following line enables the CGI configuration for this host only
				# after it has been globally disabled with \"a2disconf\".
				#Include conf-available/serve-cgi-bin.conf
			</VirtualHost>

			# vim: syntax=apache ts=4 sw=4 sts=4 sr noet"

			sudo cd /etc/apache2/sites-available/

			echo > "$basicsite.conf"

			echo $vhostfile >> "$basicsite.conf"

			printf "Enable Apache httpd (Apache2).\n"

			systemctl enable apache2

			if [[ $(systemctl is-active apache2) == "inactive" ]]; then
				sudo systemctl start apache2
			else
				sudo systemctl restart apache2
			fi

			;;
		2)
			# Nginx
			# https://ubuntu.com/tutorials/install-and-configure-nginx

			printf "Create virtual host conf file.\n"

			vhostfile="server {
       listen 80;
       listen [::]:80;

       server_name $webhostname;

       root /var/www/$basicsite;
       index index.html;

       location / {
               try_files $uri $uri/ =404;
       }
}"

			sudo cd /etc/nginx/sites-enabled/

			echo > "$basicsite.conf"

			echo $vhostfile >> "$basicsite.conf"  # use nginx a .conf filetype? can't find back

			;;
		3)
			:
			;;
	esac

	# If statement for IPaddres has always same length and set of dots on same place

	printf "The following Python Web server shall be installed: $webserverpython."
	sleep 1

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

fi

# install LDAP server

# printf command below case
if [[ $installldapserver == "y" ]]; then
	# https://likegeeks.com/linux-ldap-server/
	# phpldapadmin is not default available on debian
	# https://kifarunix.com/install-phpldapadmin-on-debian-10-debian-11/

	printf "The following LDAP server shall be installed: $ldapserver."
	sleep 1

	# Add backports to main and update
	printf "Add Debian backports to main.\n"
	sudo touch /etc/apt/sources.list.d/openldap.list
	echo "deb http://deb.debian.org/debian bullseye-backports main" | sudo tee -a /etc/apt/sources.list.d/openldap.list

	printf "Update repositories.\n"
	apt update

	testcommand=$ldaptestcommand
	packages=$ldappackages

	if ! [[ $(command -v $testcommand) ]]; then
		printf "$testcommand could not be found.\n$packages shall be installed. \n"
		sudo apt install -y $packages
		exit
	else
		printf "$packages already installed \n"
	fi

	case $ldapserverc in
		# Basically all are Java implementations except 389 directory service(c++)
		1)
			# openLDAP
			# https://www.linux.com/topic/desktop/how-install-openldap-ubuntu-server-1804/

			# configure ldap
			printf "Show initial Configuration.\n"

			sudo slapcat

			# More Configuration changes follows later

			printf "Enable openldap.\n"

			systemctl enable slapd

			if [[ $(systemctl is-active slapd) == "inactive" ]]; then
				sudo systemctl start slapd
			else
				sudo systemctl restart slapd
			fi

			;;
		2)
			# ApacheDS
			# configure ldap

			# More Configuration changes follows later

			printf "Enable apacheds.\n"

			systemctl enable apacheds

			if [[ $(systemctl is-active sudo /etc/init.d/apacheds-${version}-default) == "inactive" ]]; then
				sudo /etc/init.d/apacheds-${version}-default start
			else
				sudo /etc/init.d/apacheds-${version}-default restart
			fi

			;;
		3)
			# openDJ
			# https://askubuntu.com/questions/772235/how-to-find-path-to-java

			printf "Create Java links.\n"

			OPENDJ_JAVA_HOME=$(dirname $(which java))
			OPENDJ_JAVA_BIN=$(which java)

			printf "Download and install openDJ.\n"

			# https://github.com/OpenIdentityPlatform/OpenDJ/wiki/Installation-Guide#to-install-from-the-debian-package
			wget --content-disposition https://github.com/OpenIdentityPlatform/OpenDJ/releases/4.5.0/opendj_4.5.0-1_all.deb

			# optional downloads:
			# wget --content-disposition https://github.com/OpenIdentityPlatform/OpenDJ/releases/4.5.0/org.openidentityplatform.opendj.opendj-dsml-servlet.war
			# wget --content-disposition https://github.com/OpenIdentityPlatform/OpenDJ/releases/4.5.0/org.openidentityplatform.opendj.opendj-ldap-toolkit.zip
			# wget --content-disposition https://github.com/OpenIdentityPlatform/OpenDJ/releases/4.5.0/org.openidentityplatform.opendj.opendj-rest2ldap-servlet.war

			sudo dpkg -i opendj_4.5-1_all.deb

			printf "Default location openDJ is /opt/opendj/.
			Config openDJ.\n"

			sudo /opt/opendj/setup

			;;
		4)
			# 389 Directory Server
			# configure ldap

			# More Configuration changes follows later
			# https://directory.fedoraproject.org/docs/389ds/howto/quickstart.html

			printf "Enable 389-ds.\n"

			systemctl enable 389-ds

			if [[ $(systemctl is-active 389-ds) == "inactive" ]]; then
				sudo 389-ds start
			else
				sudo 389-ds restart
			fi

			;;
	esac

fi

# Install other dependecies
