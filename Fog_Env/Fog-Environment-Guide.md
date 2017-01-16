# Fog Environment Guide

## Outline
* System architecture
* Pre-requisition and tools' overview
* Latency emulation
* Testing scenario and tools
* Future work and trouble shooting

### System architecture
Current network partitioned into two subnets; fog-fog-cloud and fog-devices subnets.
* Fog-fog-cloud network ip addresses
```Bash
IP : 192.168.100.0/24
Gateway : 192.168.100.1/24
```
* Fog-devices network ip addresses
```Bash
IP : 192.168.200.0/24
# No gateway
```
* Special addresses and Fogs addresses
```Bash
Special addresses : 192.168.100.2 ~ 192.168.100.30
Fog1 : 192.168.100.31 ~ 192.168.100.60 ( 192.168.200.31 ~ 192.168.200.60 )
Fog2 : 192.168.100.61 ~ 192.168.100.90 ( 192.168.200.61 ~ 192.168.200.90 )
Fog3 : 192.168.100.91 ~ 192.168.100.120 ( 192.168.200.91 ~ 192.168.200.120 )
Fog4 : 192.168.100.121 ~ 192.168.100.150 ( 192.168.200.121 ~ 192.168.200.150 )
```
![image](https://github.com/tz70s/KVM-Docker-OVS-Deployment/blob/master/Fog_Env/arch.jpg)

* Example of setting up a fog node with upstream and downstream interfaces
```Bash
# Upstream interface
auto eth-up
iface eth-up static
  address 192.168.100.x
  netmask 255.255.255.0
  network 192.168.100.0
  gateway 192.168.100.1
  broadcast 192.168.100.255
# Downstream interface
auto eth-down
iface eth-down static
  address 192.168.200.x
  netmask 255.255.255.0
  network 192.168.200.0
  broadcast 192.168.200.255
```
* Checkout
```Bash
ifconfig
route -n
# Destination   Gateway        Netmask        Dev
# 0.0.0.0         192.168.100.1  0.0.0.0        eth-up
# 192.168.100.0   0.0.0.0        255.255.255.0  eth-up
# 192.168.200.0   0.0.0.0        255.255.255.0  eth-down
```
## Pre-requision and tools' overview
* KVM
* OpenvSwitch
* Kubernetes
* Latency emulation

### KVM
KVM can be installed through installing Ubuntu server distribution, remember to choose the option: Virtual Machine Host. Or [KVM Install](https://github.com/tz70s/KVM-Docker-OVS-Deployment/blob/master/kvm-install.sh). Also, following configuration tools are needed.
```Bash
sudo apt-get install virtinst libguestfs
```
* virt-install

First should make sure the target patching ovs bridge is up and the installation image is located.
```Bash
virt-install \
		--virt-type=kvm \
		--name "vm0$1" \
		--ram=2048 \
		--vcpus=1 \
		--os-variant=auto \
		--location /home/$USER/ubuntu-16.04.1-server-amd64.iso \
		--network bridge=ovs-br0,virtualport_type='openvswitch' \
    --network bridge=driver,virtualport_type='openvswitch' \
		--graphics none \
		--disk path=/var/lib/libvirt/images/"vm0$1".img,size=30,bus=virtio \
		--console pty,target_type=serial \
		--extra-args 'console=ttyS0,115200n8 serial'
```
Second step, modify grub to use virsh console
```Bash
sudo vi /etc/default/grub
# Add "console=ttyS0" to following line
GRUB_CMDLINE_LINUX_DEFAULT="quiet splash console=ttyS0"
```
And then set up network configuration and reboot
* virt-sysprep

More easily, use virt-clone and virt-sysprep to create new vm
```Bash
virt-clone \
  --original vm01 \
  --name $1 \
  --file /var/lib/libvirt/images/$1.img
```
```Bash
sudo virt-sysprep -d vm01 --hostname new_hostname
```
* virsh edit
```Bash
# Use virsh edit to modify existed vm settings such as patched bridge
virsh edit target_vm_domain_name
```
* Modify network settings inner vm
* Re-install openssh

openssh will be conflicted inner vm 
```Bash
sudo apt-get install --reinstall openssh-server openssh-client
```
* Virsh command
```Bash
virsh list --all
virsh start target_vm_domain_name
virsh shutdown target_vm_domain_name
virsh suspend target_vm_domain_name
virsh resume target_vm_domain_name
```

### OpenvSwitch
OpenvSwitch is used to simulate network switches/routers.

* Bridges' and Ports' configuration
```Bash
# Show settings
sudo ovs-vsctl show
# create/add/delete bridge/port
sudo ovs-vsctl add-br <target_br_name>
sudo ovs-vsctl add-port <target_br_name> <target_port_name>
sudo ovs-vsctl del-br <target_br_name>
sudo ovs-vsctl del-port <target_port_name>
```

* Flow settings

Flows will be flushed after shutdown 
```Bash
# Dump flows
sudo ovs-ofctl dump-flows <target_bridge>
# Add/delete flow
sudo ovs-ofctl add-flow <target_bridge> priority=100,ip,nw_dst=192.168.100.32,actions=output:4
sudo ovs-ofctl del-flow <target_bridge> ip,nw_dst=192.168.100.32
```

* Reference tutorials, documents

[Tutorial](http://www.rendoumi.com/open-vswitchzhong-ovs-ofctlde-xiang-xi-yong-fa/)

[Docs](http://openvswitch.org/support/dist-docs/ovs-ofctl.8.txt)
