# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang
MAINTAINER PierreZ

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/pierrez/ipx800

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/influxdb/influxdb/client
RUN go install github.com/pierrez/ipx800

ENV INFLUX_USER mon_super_user
ENV INFLUX_PWD mon_super_pwd

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/ipx800