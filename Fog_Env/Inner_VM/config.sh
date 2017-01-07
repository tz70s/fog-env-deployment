#!/bin/bash

# sudo -s
# sudo passwd root
# sudo passwd -u root
# exit

# Relogin as root

sudo usermod -l $1 vm01
sudo usermod -d /home/$1 -m $1
sudo hostname $1
