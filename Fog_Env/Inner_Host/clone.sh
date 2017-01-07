#!/bin/bash

virt-clone \
  --original vm01 \
  --name $1 \
  --file /var/lib/libvirt/images/$1.img
