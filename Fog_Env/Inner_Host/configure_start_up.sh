#!/bin/bash

sudo mv ovs-startup.sh /etc/init.d/
sudo update-rc.d ovs-startup.sh default
