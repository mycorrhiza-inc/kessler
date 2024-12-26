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

# Apparently docker makes offical changes to the Open Container Initiative (OCI) standard. So we are not using docker here
wget https://github.com/containerd/containerd/releases/download/v2.0.1/containerd-2.0.1-linux-amd64.tar.gz
tar Cxzvf /usr/local containerd-*
systemctl daemon-reload
systemctl enable --now containerd

wget https://github.com/opencontainers/runc/releases/download/v1.2.3/runc.amd64
install -m 755 runc.amd64 /usr/local/sbin/runc

wget https://github.com/containernetworking/plugins/releases/download/v1.6.1/cni-plugins-linux-amd64-v1.6.1.tgz
mkdir -p /opt/cni/bin
tar Cxzvf /opt/cni/bin cni-plugins-linux-amd64-v1.6.1.tgz
# # Install docker
# dnf -y install dnf-plugins-core
# dnf-3 config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
# dnf install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
# systemctl start docker
# systemctl enable docker
#
#
# # Install a docker shim, apparently it lets k8s control docker
# wget https://github.com/Mirantis/cri-dockerd/releases/download/v0.3.16/cri-dockerd-0.3.16-3.fc36.x86_64.rpm
# dnf install -y ./cri-dockerd-0.3.16-3.fc36.x86_64.rpm

cd / 
mkdir mycorrhiza
cd mycorrhiza
git clone https://github.com/mycorrhiza-inc/kessler
cd kessler
git fetch --all




# manually do some magic to copy k8s/secret.yml
# helm install kessler ./k8s -f k8s/values-prod.yaml
