#!/bin/sh

# Build sdl-clock in a balenalib 'cross-build' container
# Create an accessible container from this, then copy the built binary out
# Remove both the container and image

podman build -t sdl-clock .
podman create --name built-clock sdl-clock
podman cp built-clock:/go/src/app/sdl-clock ./

podman rm built-clock
podman rmi sdl-clock