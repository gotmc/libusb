# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VER = "2"
Vagrant.configure(VAGRANTFILE_API_VER) do |config|
  config.vm.box = "ubuntu/trusty64"

  # Add Linux related disk images to shared folder
  # config.vm.synced_folder "~/Disk Images/Linux", "/vagrant_disk_images"

  # Enable provisioning with a shell script.
  config.vm.provision "shell", inline: <<-SHELL
    apt-get update
    apt-get install -y git curl vim tree
    apt-get install -y pkg-config libusb-1.0-0 libusb-1.0-0-dev
    echo "### Setting up dotfiles"
    su vagrant -c 'mkdir -p development/go/src'
    su vagrant -l -c 'cd ~/development && git clone https://github.com/matthewrankin/dotfiles.git'
    su vagrant -l -c 'cd ~/development/dotfiles && ./deploy-dot-files.py'
    echo '### Installing Go 1.8'
    cd /tmp && wget -q https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz
    tar -C /usr/local -xzf /tmp/go1.8.1.linux-amd64.tar.gz
    rm /tmp/go1.8.1.linux-amd64.tar.gz
    su vagrant -l -c '/usr/local/go/bin/go get -u github.com/golang/lint/golint'
    su vagrant -l -c '/usr/local/go/bin/go get github.com/gotmc/libusb'
    echo '### Show version of libusb installed'
    su vagrant -l -c 'pkg-config --modversion libusb-1.0'
  SHELL
end
