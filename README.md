# BRICK

For isolating the processes some linux kernel features are used. Therefore the 
go-binary should be build with `GOOS=linux` and run from a lightweight linux-vm.

A vm is used through the Vagrantfile.


## Usage

```
# Compile your go
GOOS=linux go build
```

```
# Initate vm
vagrant up

# SSH into it
vagrant ssh
```

```
# Jump to folder which is in sync with your bricket folder.
cd /vagrant

# Initialize brick
# Run twice, some bug somewhere...
sudo ./brick init

# Start bash through the container
sudo ./brick run /bin/bash
```

Voila! You now have a mini container to that is isolated from the filesystem 
and other processes. You can check that by typing `ps`, but there is a bug with the mounting so you have first to run `mount -t proc proc /proc`


## creds

https://www.infoq.com/articles/build-a-container-golang
