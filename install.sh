echo "OS version : Ubuntu 16.04 LTS"
echo "Start installing openvswitch"

echo "Install openvswitch"
sudo apt-get install openvswitch-switch
echo "Open-Vswitch Installed"

echo "Start setting database"
sudo touch /usr/local/etc/ovs-vswitchd.conf
sudo mkdir -p /usr/local/etc/openvswitch
sudo ovsdb-tool create /usr/local/etc/openvswitch/conf.db vswitchd/vwitch.ovsschema
echo "Finish database setting"

echo "Start Database" 
ovsdb-server /usr/local/etc/openvswitch/conf.db \                                                                   
--remote=punix:/usr/local/var/run/openvswitch/db.sock \
--remote=db:Open_vSwitch,manager_options \
--private-key=db:SSL,private_key \
--certificate=db:SSL,certificate \
--bootstrap-ca-cert=db:SSL,ca_cert --pidfile --detach --log-file 

echo "Start switch"
sudo ovs-vsctl --no-wait init
sudo ovs-vswitchd --pidfile --detach
sudo ovs-vsctl show

echo "Installation finished"

