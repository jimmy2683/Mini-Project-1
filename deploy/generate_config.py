#!/usr/bin/env python3
import json
import sys

def main():
    print("Generating cluster.json for 15 nodes on 3 servers.")
    
    # Get server IPs (default to localhost for testing)
    server1 = input("Enter IP for Server 1 (Nodes 1-5) [localhost]: ") or "127.0.0.1"
    server2 = input("Enter IP for Server 2 (Nodes 6-10) [localhost]: ") or "127.0.0.1"
    server3 = input("Enter IP for Server 3 (Nodes 11-15) [localhost]: ") or "127.0.0.1"
    
    base_port = 7000
    peers = {}
    
    # helper to track next available port per IP
    # next_port[ip] = current_max_used_port
    next_port = {
        server1: base_port,
        server2: base_port, 
        server3: base_port
    }

    # Nodes 1-5 on Server 1
    for i in range(1, 6):
        port = next_port[server1] + 1
        peers[i] = f"{server1}:{port}"
        next_port[server1] = port

    # Nodes 6-10 on Server 2
    for i in range(6, 11):
        # If server2 is new (not in next_port dict yet, though initialized above), start at base
        # If it is same as server1, it continues incrementing from where server1 left off
        port = next_port[server2] + 1
        peers[i] = f"{server2}:{port}"
        next_port[server2] = port

    # Nodes 11-15 on Server 3
    for i in range(11, 16):
        port = next_port[server3] + 1
        peers[i] = f"{server3}:{port}"
        next_port[server3] = port

    config = {"peers": peers}
    
    with open("cluster.json", "w") as f:
        json.dump(config, f, indent=4)
        
    print(f"Generated cluster.json with {len(peers)} nodes.")
    print("Server 1 (Nodes 1-5):", server1)
    print("Server 2 (Nodes 6-10):", server2)
    print("Server 3 (Nodes 11-15):", server3)
    print("Copy this file to all servers.")

if __name__ == "__main__":
    main()
