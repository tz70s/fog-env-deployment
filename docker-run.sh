#!/bin/bash

# Create container

docker pull ansible/ubuntu14.04-ansible

docker run --name=container1 --net=none ansible/ubuntu14.04-ansible

