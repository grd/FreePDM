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

sleep 1

# Add line about check for existing server
printf "Do you want to install a Webserver? (y / n)\n"

read installwebserver

if [[ $installwebserver == "y" ]]; then
	printf "What backend do you want to install? (1 - 3)\n
1 - Apache Httpd
2 - Nginx
3 - option 3\n"

	read webserverc  # webserver case

	printf "What python web backend do you want to install? (1 - 4)\n
1 - Django
2 - Pyramid
3 - Falcon
4 - WebPy\n"

	read webserverpythonc  # webserver python case

	printf "What is your server name?\n"

	read webservername

	printf "What is your (web )server_domain OR IP address? (default something like web.somename.com)\n"

	read webhostname

	# maybe something about admin + password, ports etc
elif [[ $installwebserver == "n" ]]; then
	:
else
	printf "$installwebserver is not 'y' OR 'n'.\n"
fi

# Add line about check for existing server
printf "What SQL server do you want to install? (1 - 3)\n
1 - MySQL
2 - SQLite
3 - PostgreSQL\n"

read sqlserverc  # sqlserver case

read -p "Enter SQL Username:" sqlservername

printf "What is your (sql )server_domain OR IP address? (default something like sql.somename.com)\n"

read sqlhostname

# maybe something about admin + password, ports etc

# Add line about check for existing server
printf "Do you want to install a LDAP server? (y / n)\n"

read installldapserver

if [[ $installldapserver == "y" ]]; then
	printf "What LDAP server do you want to install? (1 - 4)\n
1 - open LDAP
2 - Apache DS
3 - openDJ
4 - 389 Directory server\n"

	read ldapserverc

	read -p "Enter LDAP Username:" ldapusername

	# read -sp "Enter LDAP Password:" ldappw1  # Silent
	read -n "Enter LDAP Password:" ldappw1  # With asterix

	read -n "Re-enter LDAP Password:" ldappw2  # With asterix

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
if ! (( command -V $testcommand )); then  #
  printf "$testcommand could not be found.\n$packages shall be installed. \n"
	sudo apt install -y $packages  # Somehow apt install fials
	exit
fi

# Install of a webserver

if [[ $installwebserver == "y" ]]; then

	case $webserverc in
		1)
			webserver="Apache httpd"
			testcommand=""
			packages=""

			;;
		2)
			webserver="Nginx"
			testcommand=""
			packages=""
			;;
		3)
			webserver="option3"
			testcommand=""
			packages=""
			;;
	esac

	printf "The following Web server shall be installed: $webserver."
	sleep 1

	# If statement for IPaddres has always same length and set of dots on same place

	case $webserverpythonc in
		1)
			webserver="Django"
			testcommand=""
			packages=""
			;;
		2)
			webserver="Pyramid"
			testcommand=""
			packages=""
			;;
		3)
			webserver="Falcon"
			testcommand=""
			packages=""
			;;
		4)
			webserver="WebPy"
			testcommand=""
			packages=""
			;;
		esac

	printf "The following Python Web server shall be installed: $webserverpython."
	sleep 1

fi


# install of SQL server

# Check if SQL server already exist. if yes  add database to existing server?
# work only with selected sql server

case $sqlserverc in
	1)
		sqlserver="MySQL"
		testcommand=""
		packages="mysql-server"
		;;
	2)
		sqlserver="SQLite"
		testcommand=""
		packages="sqlite3"
		;;
	3)
		sqlserver="postgreSQL"
		testcommand=""
		packages="postgresql postgresql-contrib"
		# https://www.geeksforgeeks.org/install-postgresql-on-linux/
		;;
esac

printf "The following SQL server shall be installed: $sqlserver.\n"

# install LDAP server
# https://www.howtoforge.com/how-to-install-openldap-on-debian-11/

if [[ $installldapserver == "y" ]]; then

	case $ldapserverc in
		1)
			ldapserver="open LDAP"
			testcommand=""
			packages="slapd ldap-utils"
			;;
		2)
			ldapserver="Apache DS"
			testcommand=""
			packages=""
			;;
		3)
			ldapserver="OpenDJ"
			testcommand=""
			packages=""
			;;
		4)
			ldapserver="389 Directory server"
			testcommand=""
			packages=""
			;;
	esac

	printf "The following LDAP server shall be installed: $ldapserver."
	sleep 1
fi

# Install other dependecies
