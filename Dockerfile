FROM balenalib/raspberry-pi2-debian-golang:latest-buster-build

RUN [ "cross-build-start" ]

WORKDIR /go/src/app
COPY . .

RUN apt-get update  
RUN apt-get install libsdl2-image-dev libsdl2-mixer-dev libsdl2-ttf-dev libsdl2-gfx-dev

RUN go get -d -v
RUN go build -v

RUN [ "cross-build-end" ] 
