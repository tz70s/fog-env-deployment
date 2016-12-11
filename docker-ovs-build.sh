#!/bin/bash

# Install ovs-docker utility

cd /usr/bin
sudo wget https://raw.githubusercontent.com/openvswitch/ovs/master/utilities/ovs-docker
chmod a+rwx ovs-docker

# Create OVS bridge

sudo ovs-vsctl add-br ovs-br0

# Set ovs network and dhcp client

sudo ovs-vsctl add-port ens3
sudo dhclient ovs-br0
sudo ifconfig ens3 0


