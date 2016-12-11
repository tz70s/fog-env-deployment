#!/bin/bash

# Start create VM

# Check vm name, ovs brdige name ....

virt-install \
	--virt-type=kvm \
	--name vm01 \
	--ram=1024 \
	--vcpus=1 \
	--os-variant=auto \
	--hvm \
	--cdrom /home/tzuchiao/ubuntu-16.04.1-server-amd64.iso \
	--network bridge=ovs-br0,virtualport_type='openvswitch' \
	--graphics vnc \
	--disk path=/var/lib/libvirt/images/vm02.img,size=20,bus=virtio \
	
