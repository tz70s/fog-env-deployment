# OpenVswitch Env Building Log

## Environment
```
Ubuntu 16.04 LTS
Openvswitch 2.5.0
```

## Installation
```Bash
# We can use apt-get in the Ubuntu 16.04 LTS version (2016 Realse)
sudo apt-get install openvswitch-switch

# Ovs Data Base setting
sudo touch /usr/local/etc/ovs-vswitchd.conf
sudo mkdir -p /usr/local/etc/openvswitch
sudo ovsdb-tool create /usr/local/etc/openvswitch/conf.db vswitchd/vswitch.ovsschema

# Start DB
ovsdb-server /usr/local/etc/openvswitch/conf.db \
--remote=punix:/usr/local/var/run/openvswitch/db.sock \
--remote=db:Open_vSwitch,manager_options \
--private-key=db:SSL,private_key \
--certificate=db:SSL,certificate \
--bootstrap-ca-cert=db:SSL,ca_cert --pidfile --detach --log-file

# OVS
sudo ovs-vsctl --no-wait init
sudo ovs-vswitchd --pidfile --detach
sudo ovs-vsctl show
```

## OVS Bridge/Port Setting && Linux Birdge
```Bash
# Add bridge to ovs
sudo ovs-vsctl add-br br0

# Add port to ovs
sudo ovs-vsctl add-port br0 eth0

# Must configure ovs port's ip as 0
sudo ifconfig eth0 0

# Let br0 configured as original ip address
# You can specify ip address or use dhcp

# Specify ip address
sudo ifconfig br0 192.168.8.232 netmask 255.255.255.0
sudo route add default gw 192.168.8.1 br0

# Dhcp
sudo dhclient br0

# Checkout result
ifconfig
route -n
```

## Lab1 - Bidge with virtualbox (vm)
```Bash
# Add virtual port
ip tuntap add mode tap vnet0
ip link set vnet0 up
ovs-vsctl add-port br0 vnet0

# Show Existing virtual port 
ip tuntap

# Delete virtual port
ip link del

# Then configure virtualbox network setting, bridge to vnet0
```

### Reference
* [ovs settings](http://neokentblog.blogspot.tw/2013/10/linuxopenflow-switchopenvswitch.html)
* [ovs virtualbox](http://networkstatic.net/open-vswitch-on-virtualbox/)

