#!/bin/bash 

#Stop Service
systemctl stop dispatcher.service

rm -Rf /home/ec2-user/apps/ws-message-dispatcher
mkdir /home/ec2-user/apps/ws-message-dispatcher