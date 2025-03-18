#!/bin/bash
set -eo pipefail

if [ -z "$1" ]; then
    echo "No build is set" 1>&2
    exit 1
fi

BINARY_NAME="$1"
BINARY_NAME_UBI7="dist/${BINARY_NAME}-ubi7"
IMAGE="${BINARY_NAME}:ubi7"
CONTAINER_NAME="${BINARY_NAME}-container"


# build the linux binary
echo "Build linux binary of pc"
go build -v

# build the windows binary
echo "Build windows binary of $BINARY_NAME"
GOOS=windows GOARCH=amd64 go build -v

# check binary and if it works, then upload to artifactory
if ./$BINARY_NAME -h >/dev/null; then
    echo "Push to artifactory"

    artifactory-upload.sh -lf=$BINARY_NAME       -tr=scptools-bin-dev-local     -tf="tools/$BINARY_NAME"
    artifactory-upload.sh -lf=$BINARY_NAME       -tr=scptools-bin-dev-local     -tf="ocp-stable-4.16/clients/$BINARY_NAME"

    artifactory-upload.sh -lf=$BINARY_NAME.exe   -tr=scptools-bin-dev-local     -tf="tools/$BINARY_NAME"
    artifactory-upload.sh -lf=$BINARY_NAME.exe   -tr=scptools-bin-dev-local     -tf="ocp-stable-4.16/clients/$BINARY_NAME"

    echo "Copy it to share folder PEWI4124://Daten"
    cp $BINARY_NAME $BINARY_NAME.exe  /gast-drive-d/Daten/
fi

echo
echo
echo "#  B U I L D   I M A G E   T O O L   F O R   U B I 7"
echo "bin in $PWD"

# build ubi7 binary in image
/usr/bin/podman build -t "$IMAGE" -f Dockerfile

echo "##########  copy binary from container to local  ##########"
if podman ps -a | rg "$CONTAINER_NAME" >/dev/null; then
    podman rm "$CONTAINER_NAME"
fi
podman create --name "$CONTAINER_NAME" "localhost/$IMAGE"
podman cp "$CONTAINER_NAME":/app/dist/$BINARY_NAME "$BINARY_NAME_UBI7"
scp "$BINARY_NAME_UBI7" cid-scp0-tls-v01-mgmt:
podman rm "$CONTAINER_NAME"

artifactory-upload.sh -lf="$BINARY_NAME_UBI7"   -tr=scptools-bin-dev-local   -tf="tools/$BINARY_NAME"
artifactory-upload.sh -lf="$BINARY_NAME_UBI7"   -tr=scptools-bin-dev-local   -tf="ocp-stable-4.16/clients/$BINARY_NAME"

