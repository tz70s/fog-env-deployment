#!/bin/bash

sudo /usr/share/openvswitch/scripts/ovs-ctl --system-id=fog-node start
sudo ifconfig enp0s25 0
sudo dhclient ovs-br0
