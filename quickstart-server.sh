# Server should be running fedora 
dnf install git fish tmux micro neovim -y
chsh -s /usr/bin/fish
dnf -y install dnf-plugins-core
dnf-3 config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
dnf install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
systemctl start docker
systemctl enable docker

systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target
# Package specific stuff
cd / 
mkdir mycorrhiza
cd mycorrhiza
git clone https://github.com/mycorrhiza-inc/kessler
cd kessler
