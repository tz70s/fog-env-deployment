# DOCKER, KVM, OpenvSwitch Deployment

* A build-up scripts and documents for multi-tenant fog-sdn environment
* Deep view references and tutorial in documents folder

### OS VERSION
```bash
# Both Host OS & VM
Ubuntu Server 16.04.1 LTS
```
### WORK FLOW

* [Install KVM](https://github.com/tz70s/KVM-Docker-OVS-Deployment/blob/master/Documents/kvm-install.md)
```bash
./kvm-install.sh
```

* [Install OpenvSwitch](https://github.com/tz70s/KVM-Docker-OVS-Deployment/blob/master/Documents/ovs-install.md)
```bash
./ovs-install.sh
```

* [Create KVM VM and Bridge with OVS](http://blog.codybunch.com/2016/10/14/KVM-and-OVS-on-Ubuntu-1604/)
```bash
# First Configure ovs bridge and ovs network setting
# e.g.

sudo ovs-vsctl add-br ovs-br0 
# The bridge name have to match with kvm-build.sh (virt-install)

# Then configure network from osv bridge to physical nic
# e.g. ( I use dhcp )

# add original ethernet into ovs-br
sudo ovs-vsctl add-port ovs-br0 eth0
sudo dhclient ovs-br0 
sudo ifconfig eth0 0
```

```bash
# Then build kvm vm with ovs
./kvm-build.sh

# KVM will automatically add ovs-port to ovs-bridge and vnic to vm
# The ip of this configuration is dhcpclient
# You should check out if there is a dhcp server can be connected
```

* [Install Docker on KVM ubuntu-server](https://github.com/tz70s/KVM-Docker-OVS-Deployment/blob/master/Documents/docker-install.md)
```bash
# Execute in VM node
./docker-install.sh
```

* [Config Docker Network and Bridge with OVS](https://github.com/tz70s/KVM-Docker-OVS-Deployment/blob/master/Documents/docker-ovs.md)
```bash
# First Read the documents, and pre-install ovs-docker, docker images to avoid the network connection can't link from outside world

./docker-ovs-build.sh
./docker-run.sh

```

### TROUBLE SHOOT
* See [Config Docker Network and Bridge with OVS](https://github.com/tz70s/KVM-Docker-OVS-Deployment/blob/master/Documents/docker-ovs.md)
* [Ubuntu kvm hvm trouble](http://qkxue.net/info/189364/ubuntu-kvm-ubuntu-quot-Couldn-find-hvm-kernel-for-Ubuntu-tree-quot-04-64-39)
* After in-vm ovs configuration, it can't access outside the subnet
