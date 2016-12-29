#!/bin/bash

sudo /usr/share/openvswitch/scripts/ovs-ctl --system-id=vm01 start
sudo ifconfig enp0s25 0
sudo dhclient ovs-br0
