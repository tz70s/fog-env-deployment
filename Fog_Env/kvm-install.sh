#!/bin/bash

# Check if cpus can support hardware virtualization
# If not, you may open it in BIOS

HV=`egrep -c '(vmx|svm)' /proc/cpuinfo`

if
	[ "$HV" -ne 0 ]
then 
	echo "Hardware can be virtualized"

else
	echo "Hardware cannot be virtualized"
	exit 1
fi

# Start kvm and dependencies installations
echo
echo "Start kvm and dependencies installations"

sudo apt-get install qemu-kvm libvirt-bin ubuntu-vm-builder bridge-utils virt-manager

echo "Automatically reboot? y/n"
read ANS

if [ "$ANS" = 'y' ]; then
	sudo reboot
fi

# Check the installation finished and libvirtd successfully added in groups

LIST=`groups`
X='libvirtd'
echo ""

if [[ $LIST =~ (^|[[:space:]])$X($|[[:space:]]) ]]; then
   echo "Verified user in group : success"
   echo
else
   echo "libvirtd is not in groups, please reboot"
   echo
fi

# Verify Installation

echo "Verify installation"
echo
virsh list --all


