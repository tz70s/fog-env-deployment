#!/bin/bash

# Prerequisites

KERNAL_VERSION=`uname -r`
VIR_NUM=${KERNAL_VERSION:0:3}

echo "The linux kernel verison : $KERNAL_VERSION"

if [ $VIR_NUM '>' 3.1  ]; then
	echo "Linux kernel version correct"
else
	echo "Please update the linux kernel"
	exit 1
fi

# Update apt sources

sudo apt-get update

# Install apt dependencies

sudo apt-get install apt-transport-https ca-certificates

# Downloads GPG key

sudo apt-key adv \
	--keyserver hkp://ha.pool.sks-keyservers.net:80 \
	--recv-keys 58118E89F3A912897C070ADBF76221572C52609D

# APT Repository

echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" | sudo tee /etc/apt/sources.list.d/docker.list
sudo apt-get update

# Check if dockker engine install

apt-cache policy docker-engine

# Install dependencies for auf storage driver

sudo apt-get update
sudo apt-get install linux-image-extra-$(uname -r) linux-image-extra-virtual

# Install Docker

sudo apt-get update
sudo apt-get install -y docker-engine
sudo service docker start

# Executing the Docker command without sudo

sudo groupadd docker
sudo usermod -aG docker $USER

# Need to Reboot

echo "You need to reboot to have the usermode change"
echo "Reboot ? y/n"

read ANS

if [ $ANS = 'y' ]; then
	sudo reboot
else
	echo "Remember to reboot later"
fi


