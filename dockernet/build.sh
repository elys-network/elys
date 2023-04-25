#!/bin/bash

set -eu 
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/config.sh

BUILDDIR="$2"
mkdir -p $BUILDDIR


build_local_and_docker() {
   module="$1"
   folder="$2"
   title=$(printf "$module" | awk '{ print toupper($0) }')

   printf '%s' "Building $title Locally...  "
   cwd=$PWD
   cd $folder
   GOBIN=$BUILDDIR go install -mod=readonly -trimpath ./... | grep -v -E "deprecated|keychain" | true
   local_build_succeeded=${PIPESTATUS[0]}
   cd $cwd
   

   if [[ "$local_build_succeeded" == "0" ]]; then
      echo "Done" 
   else
      echo "Failed"
      return $local_build_succeeded
   fi

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
while getopts sgojthrn flag; do
   case "${flag}" in
      g) build_local_and_docker band deps/band ;;
      r) build_local_and_docker relayer deps/relayer ;;  
      h) echo "Building Hermes Docker... ";
         docker build --tag elyszone:hermes -f dockernet/dockerfiles/Dockerfile.hermes . ;

         printf '%s' "Building Hermes Locally... ";
         cd deps/hermes; 
         cargo build --release --target-dir $BUILDDIR/hermes; 
         cd ../..
         echo "Done" ;;
   esac
done