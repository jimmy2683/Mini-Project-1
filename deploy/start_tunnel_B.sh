#!/bin/bash
# Run this on Laptop B (192.168.31.249) running Nodes 6-10

SERVER_USER="cs23btech11048"
SERVER_IP="192.168.50.54"

echo "Starting Reverse SSH Tunnel for Nodes 6-10..."
echo "This allows Server C ($SERVER_IP) to reach local ports 7006-7010."

ssh -N -R 7006:localhost:7006 \
       -R 7007:localhost:7007 \
       -R 7008:localhost:7008 \
       -R 7009:localhost:7009 \
       -R 7010:localhost:7010 \
       $SERVER_USER@$SERVER_IP

echo "Tunnel stopped."
