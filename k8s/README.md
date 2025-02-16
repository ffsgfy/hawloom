Example Minikube setup:
```sh
# CD to the base directory
cd ..

# Start the database
docker compose up -d db

# Start the Minikube cluster in the same network as the database
minikube start --driver=docker --network=hawloom_hawloom-net --static-ip=172.28.0.200

# Create secret and config
k8s/create_secret.sh
k8s/create_config.sh

# Create the deployment
kubectl apply -f k8s/hawloom.yaml

# Optionally, configure Fluent Bit
kubectl apply -f k8s/fluent-bit.yaml

# Wait a bit, then view the logs
kubectl logs -f fluent-bit-sink
```
