sudo systemctl stop docker
sudo rm -rf /var/lib/docker
sudo rm -rf /var/lib/containerd
sudo systemctl start docker
