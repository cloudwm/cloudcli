get_build_environment_docker_image() {
  local go_os="${1}"
  local go_arch="${2}"
  local image_base_name="${3}"
  local image_tag="${4}"
  if [ "${image_tag}" == "master" ]; then
    local image_tag="latest"
  fi
  echo "${image_base_name}${go_os}-${go_arch}:${image_tag}"
}

make_build_environment() {
  local go_os="${1}"
  local go_arch="${2}"
  local docker_image="${3}"
  if docker pull "${docker_image}"; then
    docker build --build-arg GOOS="${go_os}" --build-arg GOARCH="${go_arch}" \
                 -t "${docker_image}" -f ./Dockerfile.build \
                 --cache-from "${docker_image}" \
                 .
  else
    docker build --build-arg GOOS="${go_os}" --build-arg GOARCH="${go_arch}" \
                 -t "${docker_image}"  -f ./Dockerfile.build \
                 .
  fi &&\
  if [ "${CLOUDCLI_BUILD_ENVIRONMENT_SKIP_DOCKER_PUSH}" != "true" ]; then
    docker push "${docker_image}"
  fi &&\
  docker run --rm -v "`pwd`:/go/src/github.com/cloudwm/cli" "${docker_image}" dep ensure
}

get_binary_ext() {
  local go_os="${1}"
  if [ "${go_os}" == "windows" ]; then
    echo ".exe"
  else
    echo ""
  fi
}

build_binary() {
  local go_os="${1}"
  local go_arch="${2}"
  local image_base_name="${3}"
  local image_tag="${4}"
  local ext="$(get_binary_ext "${go_os}")"
  local docker_image="$(get_build_environment_docker_image "${go_os}" "${go_arch}" "${image_base_name}" "${image_tag}")"
  make_build_environment "${go_os}" "${go_arch}" "${docker_image}" &&\
  docker run --rm -v "`pwd`:/go/src/github.com/cloudwm/cli" "${docker_image}" go build -o "cloudcli${ext}" main.go
}

build_binary_archive() {
  local go_os="${1}"
  local go_arch="${2}"
  local image_base_name="${3}"
  local image_tag="${4}"
  local ext="$(get_binary_ext "${go_os}")"
  build_binary "${go_os}" "${go_arch}" "${image_base_name}" "${image_tag}" &&\
  if [ "${ext}" == ".exe" ]; then
    zip cloudcli-${go_os}-${go_arch}.zip cloudcli.exe &&\
    echo Great Success! && echo Created cloudcli binary archive: cloudcli-${go_os}-${go_arch}.zip
  else
    [ "${ext}" != "" ] && echo invalid extension && return 1
    tar -czvf cloudcli-${go_os}-${go_arch}.tar.gz cloudcli &&\
    echo Great Success! && echo Created cloudcli binary archive: cloudcli-${go_os}-${go_arch}.tar.gz
  fi
}

build_all_binary_archives() {
  local image_base_name="${1}"
  local image_tag="${2}"
  # build_binary_archive darwin 386 "${image_base_name}" "${image_tag}" &&\
  build_binary_archive darwin amd64 "${image_base_name}" "${image_tag}" &&\
  build_binary_archive linux 386 "${image_base_name}" "${image_tag}" &&\
  build_binary_archive linux amd64 "${image_base_name}" "${image_tag}" &&\
  build_binary_archive windows 386 "${image_base_name}" "${image_tag}" &&\
  build_binary_archive windows amd64 "${image_base_name}" "${image_tag}" &&\
  echo && echo && echo Great Success! All binaries compiled and archived
}

run_tests() {
  local image_base_name="${1}"
  local image_tag="${2}"
  local server_git_repo_url="${3}"
  local server_git_branch="${4}"
  local cloudcli_api_server="${5}"
  local cloudcli_server_port="${6}"
  if [ "${server_git_repo_url}" != "" ] && [ "${server_git_branch}" != "" ] && [ "${cloudcli_api_server}" != "" ] && [ "${cloudcli_server_port}" != "" ]; then
    echo Starting a local cloudcli server for testing
    local serverdir=/etc/kamatera/cloudcli-server
    rm -rf $serverdir
    mkdir -p $serverdir
    local workdir="$(pwd)"
    export TEST_API_SERVER=http://localhost:$cloudcli_server_port
    cd $serverdir &&\
    git clone -b $server_git_branch --depth 1 $server_git_repo_url . &&\
    docker build -t cloudcli-server . &&\
    ( docker rm -f cloudcli-server || true ) &&\
    docker run --rm --name cloudcli-server -d --sig-proxy=false -e CLOUDCLI_PROVIDER=proxy -e CLOUDCLI_API_SERVER=$cloudcli_api_server -p $cloudcli_server_port:80 cloudcli-server &&\
    cd "$workdir"
  else
    echo Testing with default server: $TEST_API_SERVER
  fi &&\
  build_binary linux amd64 "${image_base_name}" "${image_tag}" &&\
  bin/test.sh all
}

sign_mac_binaries() {
  amd64_tar_gz="${1}"
  # pulled Apr 18, 2022
  AWS_CLI_IMAGE="amazon/aws-cli@sha256:579f6355a1f153946f73fec93955573700a2eb0b63f9ae853000830cf6bf351a"
  alias aws="docker run -e AWS_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY $AWS_CLI_IMAGE"
  docker pull $AWS_CLI_IMAGE &&\
  ALLOCATE_HOSTS_RES="$(aws ec2 allocate-hosts --availability-zone "${AWS_MAC_INSTANCE_AVAILABILITY_ZONE}" \
    --instance-type mac1.metal --quantity 1)" &&\
  DEDICATED_HOST_ID=$(echo $ALLOCATE_HOSTS_RES | jq -r '.HostIds[0]') &&\
  echo DEDICATED_HOST_ID=$DEDICATED_HOST_ID &&\
  aws ec2 modify-instance-placement --host-id $DEDICATED_HOST_ID --instance-id $AWS_MAC_INSTANCE_ID &&\
  aws ec2 start-instances --instance-ids $AWS_MAC_INSTANCE_ID &&\
  while [ "running" != "$(aws ec2 describe-instances --instance-ids $AWS_MAC_INSTANCE_ID | jq -r '.Reservations[0].Instances[0].State.Name' | tee /dev/stderr)" ]; do
    echo Waiting for instance to be running...
    sleep 5
  done &&\
  IP="$(aws ec2 describe-instances --instance-ids $AWS_MAC_INSTANCE_ID | jq -r '.Reservations[0].Instances[0].NetworkInterfaces[0].Association.PublicIp')" &&\
  echo IP=$IP &&\
  scp -i $AWS_MAC_PEM_KEY_PATH $amd64_tar_gz ec2-user@$IP:/cloudcli-amd64.tar.gz
  ssh -i $AWS_MAC_PEM_KEY_PATH ec2-user@$IP "
    tar -xzf cloudcli-amd64.tar.gz &&\
  "
}

stop_mac_instance_release_host() {
  # pulled Apr 18, 2022
  AWS_CLI_IMAGE="amazon/aws-cli@sha256:579f6355a1f153946f73fec93955573700a2eb0b63f9ae853000830cf6bf351a"
  alias aws="docker run -e AWS_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY $AWS_CLI_IMAGE"
  docker pull $AWS_CLI_IMAGE &&\
  aws ec2 stop-instances --instance-ids $AWS_MAC_INSTANCE_ID &&\
  while [ "stopped" != "$(aws ec2 describe-instances --instance-ids $AWS_MAC_INSTANCE_ID | jq -r '.Reservations[0].Instances[0].State.Name' | tee /dev/stderr)" ]; do
    echo Waiting for instance to be stopped...
    sleep 5
  done &&\
  DEDICATED_HOST_ID="$(aws ec2 describe-instances --instance-ids $AWS_MAC_INSTANCE_ID | jq -r '.Reservations[0].Instances[0].Placement.HostId')" &&\
  echo DEDICATED_HOST_ID=$DEDICATED_HOST_ID &&\
  aws ec2 release-hosts --host-ids $DEDICATED_HOST_ID
}
