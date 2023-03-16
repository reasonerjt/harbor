#!/bin/bash
set -xe

export POSTGRESQL_HOST=$1
export REGISTRY_URL=http://$1:5000
export CHROME_BIN=chromium-browser
export GO111MODULE=auto
#export DISPLAY=:99.0
#sh -e /etc/init.d/xvfb start

sudo docker-compose -f ./make/docker-compose.test.yml up -d
sleep 10
./tests/pushimage.sh
docker ps

DIR="$(cd "$(dirname "$0")" && pwd)"
echo "current dir: `pwd`"
echo "GOPATH: $GOPATH, GOROOT: $GOROOT"
find . -maxdepth 2
go test -race -i ./src/core ./src/jobservice
sudo -E env "PATH=$PATH" "POSTGRES_MIGRATION_SCRIPTS_PATH=$DIR/../../make/migrations/postgresql/" ./tests/coverage4gotest.sh
#goveralls -coverprofile=profile.cov -service=github || true
