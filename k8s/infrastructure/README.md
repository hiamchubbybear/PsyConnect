# Infrastructure Setup

We use Helm to install standard services.

## Prerequisites

- Kubernetes Cluster (k3s)
- Helm installed (`curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash`)

## Installation Steps

1. Add Bitnami repo:

   ```bash
   helm repo add bitnami https://charts.bitnami.com/bitnami
   helm repo update
   ```

2. Install Services:

   ```bash
   # Create namespace first if not exists
   kubectl create namespace psyconnect

   # Kafka
   helm install kafka bitnami/kafka -n psyconnect -f k8s/infrastructure/kafka-values.yaml

   # MySQL
   helm install mysql bitnami/mysql -n psyconnect -f k8s/infrastructure/mysql-values.yaml

   # Redis
   helm install redis bitnami/redis -n psyconnect -f k8s/infrastructure/redis-values.yaml

   # Nginx Ingress (if creating manually, otherwise k3s might have traefik)
   # If using Nginx:
   # helm upgrade --install ingress-nginx ingress-nginx \
   #  --repo https://kubernetes.github.io/ingress-nginx \
   #  --namespace ingress-nginx --create-namespace
   ```

3. Verify:
   ```bash
   kubectl get pods -n psyconnect
   ```
