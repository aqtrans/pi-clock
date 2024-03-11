#!/bin/bash

export PATH=$PATH:/usr/local/go/bin

set -e

# Build sdl-clock in a balenalib 'cross-build' container
# Create an accessible container from this, then copy the built binary out
# Remove both the container and image

#podman build -t sdl-clock .
#podman create --name built-clock sdl-clock
#podman cp built-clock:/go/src/app/sdl-clock ./

#podman rm built-clock
#podman rmi sdl-clock

## Do this all via Bash
#export GOLANG_VERSION="go1.22.1.linux-armv6l"

## Install required packages
#apt-get -y libsdl2-2.0-0 libsdl2-dev libsdl2-image-2.0-0:armhf libsdl2-image-dev:armhf libsdl2-ttf-2.0-0:armhf libsdl2-ttf-dev:armhf
####aptitude -y -q install libsdl2-2.0-0 libsdl2-dev libsdl2-image-2.0-0 libsdl2-image-dev libsdl2-ttf-2.0-0 libsdl2-ttf-dev &&

## Install Golang - COMMENTED OUT TO KEEP THE TOOLCHAIN STABLE ON DEBIAN/BUSTER
##rm -f go1.22.1.linux-armv6l.tar.gz &&
##wget --quiet https://go.dev/dl/go1.22.1.linux-armv6l.tar.gz &&
##rm -rf /usr/local/go &&
##tar -C /usr/local -xzf go1.22.1.linux-armv6l.tar.gz &&

cd /home/deploy/pi-clock/

## Compile
go get -v
go build -v -o sdl-clock

go version

## Put into place and restart
systemctl stop sdl-clock
cp ./sdl-clock /home/sdl-clock/clock/sdl-clock
chown sdl-clock: /home/sdl-clock/clock/sdl-clock
systemctl restart sdl-clock
systemctl status sdl-clock
