set -euo pipefail
# Install utility packages to make debugging on the server easier.
dnf install git fish tmux micro neovim btop wget -y
chsh -s /usr/bin/fish
# git clone https://github.com/LazyVim/starter ~/.config/nvim
# rm -rf ~/.config/nvim/.git




# Trying something different, k8s is seeminly hugely overkill, and all the interoperable components 
# I feel could create a complexity nightmare, trying something a bit simpler with a regular k3s install 
curl -sfL https://get.k3s.io | sh -

export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh

helm list


# Install k8s dashboard
helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/
helm install testing-dashboard kubernetes-dashboard/kubernetes-dashboard
# To make acessible publicly
# kubectl -n default port-forward svc/test-dashboard-kong-proxy 8443:443 --address 0.0.0.0

# Create a service account with cluster-admin privileges
kubectl create serviceaccount nicole -n default
kubectl create clusterrolebinding nicole-admin-binding \
    --clusterrole=cluster-admin \
    --serviceaccount=default:nicole

# Create a token for the service account
kubectl create token nicole

kubectl create namespace traefik
helm install traefik traefik/traefik --namespace traefik --values k8s/helm-traefik-values.yaml


# helm repo add jetstack https://charts.jetstack.io
# helm repo update
# helm install \
#  cert-manager jetstack/cert-manager \
#   --namespace cert-manager \
#   --create-namespace \
#   --set installCRDs=true


cd / 
mkdir mycorrhiza
cd mycorrhiza
git clone https://github.com/mycorrhiza-inc/kessler
cd kessler
git fetch --all




# manually do some magic to copy k8s/secret.yml
helm install kessler ./k8s -f k8s/values-prod.yaml





# Now the service account 'nicole' has full cluster administration privileges
# WARNING: This gives complete unrestricted access - only use in test environments!
#
# mkdir /root/install-artifacts
# cd /root/install-artifacts
# # Install k8s and helm
# dnf install -y kubernetes1.32 kubernetes1.32-kubeadm kubernetes1.30-client 
#
# systemctl enable containerd
# systemctl start containerd
# # Enable IP forwarding
# echo "1" > /proc/sys/net/ipv4/ip_forward
# # Make IP forwarding persistent across reboots
# echo "net.ipv4.ip_forward = 1" > /etc/sysctl.d/99-kubernetes-cri.conf
# sysctl --system
#
# # Initialize kubernetes
# kubeadm init


# Now helm should work


