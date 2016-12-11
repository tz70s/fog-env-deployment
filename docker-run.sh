#!/bin/bash

# Create container

docker pull ansible/ubuntu14.04-ansible

docker run -itd --privileged --name=container1 --net=none ansible/ubuntu14.04-ansible

# Create ovs docker connection
# Without prec-onfigure ip address

sudo ovs-docker add-port ovs-br0 eth0 container1

# Then do the following actions in container

docker start -i container1

# mv /sbin/dhclient /usr/sbin/dhclient
# dhclient eth0
# ignore the error log
