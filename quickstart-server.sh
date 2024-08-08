# Server should be running fedora 
dnf install git fish -y

dnf -y install dnf-plugins-core
dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
dnf install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
systemctl start docker
systemctl enable docker
git clone https://github.com/mycorrhiza-inc/kessler
cd kessler

