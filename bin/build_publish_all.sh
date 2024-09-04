#!/usr/bin/env bash

# this script runs from our Jenkins to build and publish all binaries as latest
# to run it locally (it won't publish because it uses local directory on the Jenkins server to publish):
#   docker build -t cloudcli-build -f Dockerfile.build .
#   export BUILD_ENV_DOCKER_IMAGE_BASE_NAME=cloudcli-build
#   export BUILD_ENV_DOCKER_IMAGE_TAG=latest
#   export PUBLISH_BINARIES_VERSION=v0.0.0
#   export CLOUDCLI_BUILD_ENVIRONMENT_SKIP_DOCKER_PUSH=true

# To sign mac binaries follow the guide at /MAC_SIGN.md
# you need the following env vars are needed to sign the mac binary via the mac VM:
#   export AWS_ACCESS_KEY_ID=
#   export AWS_SECRET_ACCESS_KEY=
#   export AWS_REGION=eu-central-1
#   export AWS_MAC_INSTANCE_AVAILABILITY_ZONE=eu-central-1c
#   export AWS_MAC_INSTANCE_ID=
#   export AWS_MAC_PEM_KEY_PATH=/path/to/mac.pem
#  apple username
#   export AC_USERNAME=
#  apple app-specific password
#   export AC_PASSWORD=
# run the flow:
#   bin/build_publish_all.sh

source bin/functions.sh
build_all_binary_archives "${BUILD_ENV_DOCKER_IMAGE_BASE_NAME}" "${BUILD_ENV_DOCKER_IMAGE_TAG}" &&\
(true || sign_mac_binaries cloudcli-darwin-amd64.tar.gz) &&\
if [ "${PUBLISH_BINARIES_VERSION}" == "" ]; then
  echo Skipping publishing binaries
else
  echo Publishing binaries to version ${PUBLISH_BINARIES_VERSION} &&\
  mkdir -p /var/cloudcli/binaries/${PUBLISH_BINARIES_VERSION} &&\
  cp *.tar.gz /var/cloudcli/binaries/${PUBLISH_BINARIES_VERSION}/ &&\
  cp *.zip /var/cloudcli/binaries/${PUBLISH_BINARIES_VERSION}/ &&\
  echo "${PUBLISH_BINARIES_VERSION}" > /var/cloudcli/binaries/LATEST_VERSION.txt &&\
  rm -f /var/cloudcli/binaries/latest &&\
  ln -s "/var/cloudcli/binaries/${PUBLISH_BINARIES_VERSION}" /var/cloudcli/binaries/latest &&\
  chown -R www-data:www-data /var/cloudcli/binaries
fi
