#!/bin/bash

# $1 is the $1 and the $2 is the IP
for IP in $(seq 1 254)
do
	printf "%-15s ==> " "${1}.${IP}"
	ping -c 1 "${1}.${IP}" > /dev/null 2>&1 

	if [ $? -eq 0 ]
	then
		printf "${RED}up${WHITE}\n"
	else
		printf "down\n"
	fi
done