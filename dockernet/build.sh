#!/bin/bash

set -eu 
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/config.sh

BUILDDIR="$2"
mkdir -p $BUILDDIR


build_local_and_docker() {
   module="$1"
   title=$(printf "$module" | awk '{ print toupper($0) }')

   echo "Building $title Docker...  "
   if [[ "$module" == "elys" ]]; then
      image=Dockerfile
   else
      image=dockernet/dockerfiles/Dockerfile.$module
   fi

   DOCKER_BUILDKIT=1 docker build --tag elyszone:$module -f $image . | true
   docker_build_succeeded=${PIPESTATUS[0]}

   if [[ "$docker_build_succeeded" == "0" ]]; then
      echo "Done" 
   else
      echo "Failed"
   fi
   return $docker_build_succeeded
}

# build docker images and local binaries
while getopts ebrhf flag; do
   case "${flag}" in
      e) build_local_and_docker elys ;;
      f) build_local_and_docker feeder ;;
      b) build_local_and_docker band ;;
      r) build_local_and_docker relayer ;; 
      h) echo "Building Hermes Docker... ";
         docker build --tag elyszone:hermes -f dockernet/dockerfiles/Dockerfile.hermes . ;
         echo "Done" ;;
   esac
done