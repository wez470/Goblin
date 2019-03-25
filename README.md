# Goblin

A TCP load balancer written in Go. Currently has static configuration and
uses a round robin algorithm for balancing.

# Building / Running
First, configure conf.yaml to include the backend addresses you would like
to balance between as well the port you would like to run the Goblin on.

Next, run `go run .`