dnf install -y kubernetes1.32 kubernetes1.32-kubeadm kubernetes1.30-client 
dnf install git fish tmux micro neovim -y
chsh -s /usr/bin/fish
git clone https://github.om/LazyVim/starter ~/.config/nvim
rm -rf ~/.config/nvim/.git

cd / 
mkdir mycorrhiza
cd mycorrhiza
git clone --mirror https://github.com/mycorrhiza-inc/kessler
cd kessler
