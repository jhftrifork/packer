# -*- mode: ruby -*-
# vi: set ft=ruby :

$script = <<SCRIPT
SRCROOT="/opt/go"

# Install Go
sudo apt-get update
sudo apt-get install -y build-essential mercurial bzr
sudo hg clone -u release https://code.google.com/p/go ${SRCROOT}
cd ${SRCROOT}/src
sudo ./all.bash

# Setup the GOPATH
sudo mkdir -p /opt/gopath
cat <<EOF >/tmp/gopath.sh
export GOPATH="/opt/gopath"
export PATH="/opt/go/bin:/opt/gopath/bin:\$PATH"
EOF
sudo mv /tmp/gopath.sh /etc/profile.d/gopath.sh
sudo chmod 0755 /etc/profile.d/gopath.sh

# Make sure the gopath is usable by vagrant
sudo chown -R vagrant:vagrant $SRCROOT
sudo chown -R vagrant:vagrant /opt/gopath

# Install some other stuff we need
sudo apt-get install -y curl git-core zip

cd /opt/gopath/src/github.com/mitchellh/packer
make updatedeps
make
go get -u github.com/mitchellh/gox
SCRIPT

Vagrant.configure(2) do |config|
  config.vm.box = "chef/ubuntu-12.04"

  config.vm.provision "shell", inline: $script

  config.vm.synced_folder ".", "/opt/gopath/src/github.com/mitchellh/packer"

  ["vmware_fusion", "vmware_workstation"].each do |p|
    config.vm.provider "p" do |v|
      v.vmx["memsize"] = "2048"
      v.vmx["numvcpus"] = "2"
      v.vmx["cpuid.coresPerSocket"] = "1"
    end
  end
end
