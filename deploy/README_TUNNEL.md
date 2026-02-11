# Distributed Deployment via SSH Tunnels

This method solves the "different subnet" problem where the Server (`.54`) cannot initiate connections to your Laptops (`.179`, `.249`), but your Laptops can SSH into the Server.

## Architecture
- **Server C (192.168.50.54):** Acts as the central hub. It reaches Laptops via **SSH Reverse Tunnels** (localhost ports).
- **Laptops A & B:** Connect to Server C normally.

## Step 1: Prepare Config Files
We use **Split-Horizon DNS** (different configs for different machines).

### On Laptops A & B
Use `deploy/cluster_laptops.json`. Rename it to `cluster.json`.
It points to public IPs for everyone.

### On Server C
Use `deploy/cluster_server.json`. Name it `cluster.json`.
It points to `127.0.0.1` for nodes 1-10 (tunnel endpoints) and its own IP for 11-15.

## Step 2: Establish Tunnels
Before running any nodes, set up the network.

### On Laptop A (Nodes 1-5)
Run:
```bash
bash deploy/start_tunnel_A.sh
```
Keep this terminal open!

### On Laptop B (Nodes 6-10)
Run:
```bash
bash deploy/start_tunnel_B.sh
```
Keep this terminal open!

## Step 3: Run the Nodes
Now connections are possible.

### Server C (Nodes 11-15)
SSH into `.54`. Ensure `cluster.json` is the **Server** version.
```bash
./shard-server -id 11 -config cluster.json &
# ... repeat for 12-15
```

### Laptop A (Nodes 1-5)
Ensure `cluster.json` is the **Laptop** version.
```bash
./shard-server -id 1 -config cluster.json &
# ... repeat for 2-5
```

### Laptop B (Nodes 6-10)
Ensure `cluster.json` is the **Laptop** version.
```bash
./shard-server -id 6 -config cluster.json &
# ... repeat for 7-10
```

## Step 4: Verification
- **Server Logs:** Should show successful connections to `127.0.0.1:7xxx`.
- **Laptop Logs:** Should show connections to `192.168.50.54:7xxx`.
