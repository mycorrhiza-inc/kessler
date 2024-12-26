set -euo pipefail
# Install utility packages to make debugging on the server easier.
dnf install git fish tmux micro neovim btop wget -y
chsh -s /usr/bin/fish
git clone https://github.om/LazyVim/starter ~/.config/nvim
rm -rf ~/.config/nvim/.git

mkdir /root/install-artifacts
cd /root/install-artifacts
# Install k8s and helm
dnf install -y kubernetes1.32 kubernetes1.32-kubeadm kubernetes1.30-client 
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh

systemctl enable containerd
systemctl start containerd
# Enable IP forwarding
echo "1" > /proc/sys/net/ipv4/ip_forward
# Make IP forwarding persistent across reboots
echo "net.ipv4.ip_forward = 1" > /etc/sysctl.d/99-kubernetes-cri.conf
sysctl --system

# this works
kubeadm init

# Then I try this and get an error
helm list
# Error: Kubernetes cluster unreachable: Get "http://localhost:8080/version": dial tcp [::1]:8080: connect: connection refused

# # Apparently docker makes offical changes to the Open Container Initiative (OCI) standard. So we are not using docker here
# wget https://github.com/containerd/containerd/releases/download/v2.0.1/containerd-2.0.1-linux-amd64.tar.gz
# tar Cxzvf /usr/local containerd-*
# systemctl daemon-reload
# systemctl enable --now containerd
#
# wget https://github.com/opencontainers/runc/releases/download/v1.2.3/runc.amd64 e
# install -m 755 runc.amd64 /usr/local/sbin/runc
#
# wget https://github.com/containernetworking/plugins/releases/download/v1.6.1/cni-plugins-linux-amd64-v1.6.1.tgz
# mkdir -p /opt/cni/bin
# tar Cxzvf /opt/cni/bin cni-plugins-linux-amd64-v1.6.1.tgz
# # Although apparently you can also install it with dnf using the containerd distributed by docker like so
# dnf -y install dnf-plugins-core
# dnf-3 config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
# dnf install containerd.io  -y
# systemctl daemon-reload
# systemctl enable --now containerd

cd / 
mkdir mycorrhiza
cd mycorrhiza
git clone https://github.com/mycorrhiza-inc/kessler
cd kessler
git fetch --all




# manually do some magic to copy k8s/secret.yml
# helm install kessler ./k8s -f k8s/values-prod.yaml
