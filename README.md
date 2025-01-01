<center>
<img src="./lbx.png" height="100px">  
<h1>LBX: Loadbalancer in GO</h1>
<h3> We Balance Loads, and Occasionally Your Sanity </h3>
<br>
</center>

This project implements a simple load balancer in Go that distributes incoming requests across multiple backend servers using a **Round Robin** algorithm. The load balancer reads configuration from a YAML file and forwards requests to the backend servers in a rotating manner.

## Example Configuration

The load balancer configuration is specified in a YAML file (`config.yml`). Below is an example configuration:

```yaml
entry_point: 3000
servers:
  - name: server1
    url: http://server1:3000
  - name: server2
    url: http://server2:3000
  - name: server3
    url: http://server3:3000
```

- `entry_point`: The port where the load balancer will listen for incoming traffic.
- `servers`: A list of backend servers where the load balancer will forward requests.
  - `name`: The name of the server.
  - `url`: The URL of the backend server that will handle the requests.

## Build and Run the Load Balancer

### 1. Docker Setup

You can run the load balancer as a Docker container by using the following command:

```bash
docker run --name lbx -v config.yml:/app/config.yml -p 3000:3000 subrotokumar/lbx
```

### Parameters:
- `--name lbx`: Assigns the name `lbx` to the container.
- `-v config.yml:/app/config.yml`: Mounts the local `config.yml` file into the container.
- `-p 3000:3000`: Exposes port `3000` on your local machine and maps it to port `3000` in the container.
- `subrotokumar/lbx`: The Docker image for the load balancer.

Once the container is running, the load balancer will start listening on port `3000`. Requests to this port will be distributed across the backend servers defined in the configuration file using the round-robin algorithm.

## Load Balancer Logic

The load balancer will:
1. Receive incoming requests on the configured entry point (e.g., port `3000`).
2. Distribute these requests to the backend servers in a **round-robin** fashion:
   - After sending a request to `server1`, the next request will go to `server2`, then `server3`, and so on.
   - Once all servers have received a request, the cycle repeats starting from `server1`.

## Example Use Case

1. You have three backend servers: `server1`, `server2`, and `server3`, each running on port `3000`.
2. The load balancer receives incoming HTTP requests on port `3000` and forwards them to each of these servers in a round-robin order.
3. If `server1` handles a request, the next request will go to `server2`, and the next to `server3`. After that, it starts over at `server1`.

## Docker Compose Example (Optional)

You can also run the load balancer and the backend servers using Docker Compose. Here's an example `docker-compose.yml` file:

```yaml
version: "3.7"

services:
  lbx:
    image: subrotokumar/lbx
    ports:
      - "3000:3000"
    volumes:
      - ./config.yml:/app/config.yml
    depends_on:
      - server1
      - server2
      - server3

  server1:
    image: some-http-server-image
    environment:
      - SERVER_NAME=server1
    expose:
      - "3000"

  server2:
    image: some-http-server-image
    environment:
      - SERVER_NAME=server2
    expose:
      - "3000"

  server3:
    image: some-http-server-image
    environment:
      - SERVER_NAME=server3
    expose:
      - "3000"
```

To start the services, run:

```bash
docker-compose up
```

This will start the load balancer and the three backend servers, with the load balancer forwarding traffic to them in round-robin order.

## License

This project is open-source and available under the MIT License.

---

This README provides an overview of how to configure and run the load balancer, including Docker setup instructions, configuration examples, and the load balancing logic.