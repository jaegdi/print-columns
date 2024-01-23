#!/bin/bash
set -e

echo "Build linux binary of pc (print columns)"
go build
./pc -man > pc-ReadMe.md 2>/dev/null

echo "Build windows binary of pc"
GOOS=windows GOARCH=amd64 go build

echo "generate ReadMe.md"
./pc -h > ReadMe.md

echo "Push to artifactory"
artifactory-upload.sh  -lf=pc             -tr=scptools-bin-develop   -tf=tools/pc
artifactory-upload.sh  -lf=pc.exe         -tr=scptools-bin-develop   -tf=tools/pc
artifactory-upload.sh  -lf=pc-ReadMe.md   -tr=scptools-bin-develop   -tf=tools/pc

echo "Copy it to share folder PEWI4124://Daten"
cp pc pc.exe  /gast-drive-d/Daten
cp pc-ReadMe.md  /gast-drive-d/Daten
