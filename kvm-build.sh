#!/bin/bash

# Create a VM

virt-install \
	--virt-type=kvm \
	--name vm01 \
	--ram=1024 \
	--vcpus=1 \
	--os-variant=auto \
	--hvm \
	--cdrom /home/tzuchiao/ubuntu-16.04.1-server-amd64.iso \
	--network network=default,model=virtio \
	--graphics vnc \
	--disk path=/var/lib/libvirt/images/vm01.img,size=20,bus=virtio \
	
