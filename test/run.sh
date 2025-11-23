#!/bin/bash
IMAGE=mattermost-nats
TAG=latest

log() {
    echo -e $* >&2
}

err() {
    echo "Error!"
    echo -e $* >&2
}

usage() {
    log "$(basename $0): Usage"
    log ""
    log "\tbuild\t\t: Build packages and docker"
    log "\tbuild-docker\t: Bring docker image only"
    log "\tdown\t\t: Down all services"
    log "\tstart\t\t: Start all services"
    log "\tstop\t\t: Stop all services"
    log "\t-h\t\t: Help"
    echo ""
}

# build_docker() {
#     pushd ../mattermost/server/build
#     docker build --platform linux/amd64 -t $IMAGE:$TAG .
#     popd
# }

build_mm() {
    pushd ../mattermost/server
    #go get github.com/samber/lo
    #go get github.com/kelseyhightower/envconfig
    #go get github.com/nats-io/nats.go@v1.34.0

    make build
    make package
    cp dist/mattermost-team-linux-amd64.tar.gz ./build/
    popd
}

build_docker() {
    pushd ./build
    cp -f ../../mattermost/server/dist/mattermost-team-linux-amd64.tar.gz ./
    docker build --platform linux/amd64 -t $IMAGE:$TAG .
    popd
}

execute() {
    local task=${1}
    case ${task} in
        build-mm)
            build_mm
            ;;
        build)
            build_mm
            build_docker
            ;;
        build-docker)
            build_docker
            ;;
        up)
            docker-compose up --build -d
            ;;
        down)
            docker-compose down
            ;;
        start)
            docker-compose start
            ;;
        stop)
            docker-compose stop
            ;;
        *)
            echo "invalid task: ${task}"
            usage
            exit 1
            ;;
    esac
}

main() {
    execute $@
}

main $@
