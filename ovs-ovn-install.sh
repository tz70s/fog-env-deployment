#!/bin/bash


# Pre-requisite packages

sudo apt-get update
sudo apt-get install -y build-essential fakeroot debhelper \
                    autoconf automake bzip2 libssl-dev \
                    openssl graphviz python-all procps \
                    python-dev python-setuptools \
                    python-twisted-conch libtool git dh-autoreconf \
                    linux-headers-$(uname -r)
sudo easy_install -U pip

git clone https://github.com/openvswitch/ovs.git
cd ovs

# Start Compiling and Installation

sudo ./boot.sh
sudo ./configure --prefix=/usr --localstatedir=/var  --sysconfdir=/etc --enable-ssl --with-linux=/lib/modules/`uname -r`/build
sudo make -j3

sudo make install
sudo make modules_install

# OVS Python libraries

sudo pip install ovs


# Create a depmod.d file to use OVS kernel modules from this repo instead of upstream linux

sudo cat > /etc/depmod.d/openvswitch.conf << EOF
override openvswitch * extra
override vport-* * extra
EOF

# Copy a startup script and start OVS

sudo depmod -a
sudo cp debian/openvswitch-switch.init /etc/init.d/openvswitch-switch
sudo /etc/init.d/openvswitch-switch force-reload-kmod


