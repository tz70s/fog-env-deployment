#!/bin/bash

# Start create VM

# Check vm name, ovs brdige name ....

# Check number of vm nodes

if [ -z $1 ]
then
	echo "Please enter at least 1-8 nodes number ... exit"
	exit 0
else
	if [ $1 -gt 8 ]
	then
		echo "Can't build > 8 nodes ... exit"
		exit 0
	else
		if [ $1 -le 0 ]
		then
			echo "Can't build <= 0 nodes ... exit"
			exit 0
		fi
		echo "$1 vms start building"
	fi
fi

# Check image

if [ -s /home/$USER/ubuntu-16.04.1-server-amd64.iso ]
then
	echo "Check image complete"
else
	echo "Image file not exist ... exit"
	exit 0
fi

for ((i = 1; i <= $1; i++))
do
	sudo virt-install \
		--virt-type=kvm \
		--name "vm0$1" \
		--ram=2048 \
		--vcpus=2 \
		--os-variant=auto \
		--location /home/$USER/ubuntu-16.04.1-server-amd64.iso \
		--network bridge=ovs-br0,virtualport_type='openvswitch' \
		--graphics none \
		--disk path=/var/lib/libvirt/images/"vm0$1".img,size=30,bus=virtio \
		--console pty,target_type=serial \
		--extra-args 'console=ttyS0,115200n8 serial'
done
