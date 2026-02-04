# Kubernetes & CI/CD Setup Guide

## 1. Server Setup (One-time)

Run these commands on your physical server:

### Install K3s (Lightweight Kubernetes)

```bash
curl -sfL https://get.k3s.io | sh -
# Get kubeconfig
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $(id -u):$(id -g) ~/.kube/config
```

### Install Helm

```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

### Setup Infrastructure

Access your server and clone this repo (or copy the k8s folder).
Then follow the commands in `k8s/infrastructure/README.md` to install MySQL, Kafka, Redis.

### Application Deployment (Initial)

```bash
# Apply Base Configs
kubectl apply -f k8s/base/namespace.yaml
kubectl apply -f k8s/base/configmaps.yaml
# EDIT secrets.yaml first!
kubectl apply -f k8s/base/secrets.yaml

# Apply Services (Run for all folders)
kubectl apply -f k8s/services/identity/
kubectl apply -f k8s/services/profile/
# ... etc
kubectl apply -f k8s/services/webapp/
kubectl apply -f k8s/ingress.yaml
```

## 2. GitHub Self-Hosted Runner Setup

Vì bạn đang dùng Cloudflare Tunnel và Server nằm sau NAT, giải pháp tốt nhất là cài đặt **GitHub Action Runner** ngay trên Server.

### Tại sao?

- **Không cần SSH Key**: Runner sẽ chủ động kết nối lên GitHub để nhận job.
- **Không cần mở port 22**: An toàn hơn.

### Cài đặt Runner

1. Vào GitHub Repo -> **Settings** -> **Actions** -> **Runners**.
2. Nhấn **New self-hosted runner**.
3. Chọn **Linux**.
4. Chạy các lệnh được cung cấp trên Server của bạn. Ví dụ:

   ```bash
   # Create a folder
   mkdir actions-runner && cd actions-runner
   # Download (Link version might change, check GitHub UI)
   curl -o actions-runner-linux-x64-2.311.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-x64-2.311.0.tar.gz
   # Extract
   tar xzf ./actions-runner-linux-x64-2.311.0.tar.gz

   # Configure (Replace TOKEN with the one in GitHub UI)
   ./config.sh --url https://github.com/hiamchubbybear/psyconnect --token <YOUR_TOKEN>

   # Install as Service (Important! So it runs in background)
   sudo ./svc.sh install
   sudo ./svc.sh start
   ```

## 3. Deployment Flow

1. Push code to `main`.
2. Runner trên server nhận lệnh.
3. Runner tự build Docker image -> Push lên Registry.
4. Runner tự chạy `kubectl` để restart Pods.
