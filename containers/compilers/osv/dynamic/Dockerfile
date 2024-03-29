FROM ubuntu:18.04

# Install prerequisites
RUN apt-get update -y
RUN apt-get install -y curl git qemu

# Install Go
RUN curl https://storage.googleapis.com/golang/go1.18.1.linux-amd64.tar.gz | tar xz -C /usr/local && \
    mv /usr/local/go /usr/local/go1.18.1 && \
    ln -s /usr/local/go1.18.1 /usr/local/go
ENV GOPATH=/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:/usr/local/go/bin:$PATH

# Build Capstan from source (use mikelangelo-project fork that supports package management)
RUN go get github.com/mikelangelo-project/capstan && \      
    go install github.com/mikelangelo-project/capstan

# Copy files needed by docker container
COPY docker_files/root /root

# Create mount point directory
RUN mkdir /project_directory

# Compose boot image and copy it to /project_directory (kernctl expects it there as a result)
CMD cd /project_directory && \
    capstan pull mike/osv-loader && \
    capstan package compose kernctl/dynamic-image --pull-missing --size $MAX_IMAGE_SIZE && \	
    cp /root/.capstan/repository/kernctl/dynamic-image/dynamic-image.qemu /project_directory/boot.qcow2

#
# NOTES
#
# Build this container with:
# docker build -t bhojpur/compilers-osv-dynamic:v0.0 . --no-cache
#
# Run this container with:
# docker run -ti --volume="$PWD:/project_directory" --env MAX_IMAGE_SIZE=500MB bhojpur/compilers-osv-dynamic:v0.0
#