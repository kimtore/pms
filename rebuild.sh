#!/bin/bash
# Re-build PMS

./autogen.sh && ./configure && make

if [ $? -ne 0 ]; then
	echo
	echo "Build failed. Bug reports to <kimtjen@gmail.com>"
	echo
	exit 1
fi

echo
echo "Build finished."
echo

echo -en "Perform \`sudo make install'? [y/N] "
read answer

if [ "$answer" == "y" ] || [ "$answer" == "Y" ]; then
	sudo make install
else
	echo
	echo "Not installing for all users."
	echo
fi

./pms -v
