# Docker container on OpenvSwitch

##TROUBLE SHOOT

* While in-vm ovs configuration, it can't access outside the subnet
* [docker container dhcp issue](http://unix.stackexchange.com/questions/155990/docker-bridges-dhcp)
* After exiting container, the network configuration will be disappear. It may be configure in /etc/network/interface as booting initialization?

##REFERNCES

* [ovs-docker](https://developer.ibm.com/recipes/tutorials/using-ovs-bridge-for-docker-networking/)

##NOTICE

* Ensure while ovs-docker execution, containers should run simultaneously
