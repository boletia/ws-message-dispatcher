#!/bin/bash 
source /home/ec2-user/.bashrc

sudo chown -Rv ec2-user:ec2-user /home/ec2-user/apps/ws-message-dispatcher
sudo chmod -Rv 775 /home/ec2-user/apps/ws-message-dispatcher

## Go to the deployment directory
cd /home/ec2-user/apps/ws-message-dispatcher

# Run asset precompilation
make build
mv build/ws-message-dispatcher .