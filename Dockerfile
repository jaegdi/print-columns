# Use Red Hat UBI 7 as the base image
FROM registry.access.redhat.com/ubi7/ubi:latest

USER 0

# Install dependencies
RUN yum install -y \
    wget \
    tar \
    gcc \
    make \
    # libX11-devel \
    glibc-devel \
    && yum clean all

# Install Go manually
ENV GO_VERSION=1.23.7
ENV CGO_CFLAGS="-std=c99 -Wno-implicit-function-declaration"
ENV CGO_ENABLED=0
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOLANGTGZ="go${GO_VERSION}.linux-amd64.tar.gz"

RUN wget -q https://go.dev/dl/${GOLANGTGZ} \
 && tar -C /usr/local -xzf ${GOLANGTGZ} \
 && rm ${GOLANGTGZ}


# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the pc binary with updated C99 and GCC flags
RUN rm -f go.mod go.sum pc pc.exe\
 && go mod init pc \
 && go mod tidy \
 && echo "######## Now building the binary ##########" \
 && if go build -v -o ./dist ; then echo "######## Binary built successfully ##########"; fi

# Set the entrypoint to the built binary
CMD ["sleep", "infinity"]
