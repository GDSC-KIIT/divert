# parent image
FROM golang:1.15.6-alpine3.12

# workspace directory
WORKDIR /divert

# copy `go.mod` and `go.sum`
ADD ./go.mod ./go.sum ./

# install dependencies
RUN go mod download

RUN apk add build-base

# copy source code
COPY . .

# build executable
RUN go build

# expose ports
EXPOSE 3000

# set entrypoint
ENTRYPOINT [ "./divert" ]