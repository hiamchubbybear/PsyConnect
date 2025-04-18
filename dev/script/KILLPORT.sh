#!/bin/bash

# Define the ports to kill
PORTS=(8080 8081 8082 8083 8888)

echo "Attempting to kill processes on ports: ${PORTS[*]}"

for port in "${PORTS[@]}"; do
    # Find all processes using the port
    process_ids=$(lsof -ti :"$port")
    
    if [[ -n "$process_ids" ]]; then
        echo "Found processes on port $port:"
        ps -p "$process_ids" -o pid,command
        
        # Kill all processes
        kill -9 $process_ids
        echo "Successfully killed processes on port $port"
    else
        echo "No processes found on port $port"
    fi
done

echo "Port cleanup completed"
