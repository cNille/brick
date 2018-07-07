#!/bin/bash

# ================
# INSTALL NODE (with nvm)
# ================

sudo apt-get install -y git-core curl

curl https://raw.githubusercontent.com/creationix/nvm/v0.23.3/install.sh | bash

echo "source /home/vagrant/.nvm/nvm.sh" >> /home/vagrant/.profile
source /home/vagrant/.profile

nvm install stable
nvm alias default stable

# ================
# INSTALL GO
# ================

#!/bin/bash
set -e
VERSION="1.7.1"
DFILE="go$VERSION.linux-amd64.tar.gz"
if [ -d "$HOME/.go" ] || [ -d "$HOME/go" ]; then
    echo "Installation directories already exist. Exiting."
    exit 1
fi
echo "Downloading $DFILE ..."
wget https://storage.googleapis.com/golang/$DFILE -O /tmp/go.tar.gz
if [ $? -ne 0 ]; then
    echo "Download failed! Exiting."
    exit 1
fi
echo "Extracting ..."
tar -C "$HOME" -xzf /tmp/go.tar.gz
mv "$HOME/go" "$HOME/.go"
touch "$HOME/.bashrc"
{
    echo '# GoLang'
    echo 'export GOROOT=$HOME/.go'
    echo 'export PATH=$PATH:$GOROOT/bin'
    echo 'export GOPATH=/vagrant'
    echo 'export PATH=$PATH:$GOPATH/bin'
} >> "$HOME/.bashrc"
export GOROOT=$HOME/.go
export PATH=$PATH:$GOROOT/bin
export GOPATH=/vagrant
export PATH=$PATH:$GOPATH/bin
rm -f /tmp/go.tar.gz
