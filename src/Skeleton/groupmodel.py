import os

# Get the list of all files and directories
path = "~"
dir_list = os.listdir()
 
print("Files and directories in '", path, "' :")
 
# prints all files
for x in dir_list:
	# if x.endswith("FCStd"):
		# do something
	print(x)


