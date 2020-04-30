#!/bin/bash

## Print deployment info
DEPLOYMENT_TIME=$( date -u "+%Y/%m/%d %H:%M:%S" )
echo "Deployment finished at: "$DEPLOYMENT_TIME" UTC" > /home/ec2-user/apps/ws-message-dispatcher/deployment_time.txt

sudo chown -Rv ec2-user:ec2-user /home/ec2-user/apps/ws-message-dispatcher
sudo chmod -Rv 775 /home/ec2-user/apps/ws-message-dispatcher

sudo systemctl start dispatcher.service
sudo service nginx restart