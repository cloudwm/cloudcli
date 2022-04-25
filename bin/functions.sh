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

wait_for() {
  local condition="${1}"
  local wait_for_message="${2}"
  local sleep_seconds="${3}"
  local max_iterations="${4}"
  local num_iterations=0
  local res=0
  while ! eval "${condition}"; do
    num_iterations=$(expr $num_iterations + 1)
    if [ "${num_iterations}" == $(expr $max_iterations + 1) ]; then
      res=1
      break
    fi
    echo "${wait_for_message}"
    sleep $sleep_seconds
  done
  return $res
}

sign_mac_binaries() {
  export amd64_tar_gz="${1}"
  # pulled Apr 18, 2022
  export aws_cli_image="amazon/aws-cli@sha256:579f6355a1f153946f73fec93955573700a2eb0b63f9ae853000830cf6bf351a"
  export aws="docker run -e AWS_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY $aws_cli_image"
  docker pull $aws_cli_image &&\
  export dedicated_host_id=$($aws ec2 describe-hosts | python3 -c "
import sys, json, os
for host in json.load(sys.stdin)['Hosts']:
  if len(host.get('Instances', [])) == 1 and host['Instances'][0].get('InstanceId') == os.environ.get('AWS_MAC_INSTANCE_ID'):
    print(host['HostId'])
    break
  ") &&\
  if [ "${dedicated_host_id}" == "" ]; then
    echo allocating new dedicated host &&\
    export allocate_hosts_res="$($aws ec2 allocate-hosts --availability-zone "${AWS_MAC_INSTANCE_AVAILABILITY_ZONE}" \
      --instance-type mac1.metal --quantity 1)" &&\
    export dedicated_host_id=$(echo $allocate_hosts_res | jq -r '.HostIds[0]') &&\
    $aws ec2 modify-instance-placement --host-id $dedicated_host_id --instance-id $AWS_MAC_INSTANCE_ID &&\
    echo allocated new dedicated host $dedicated_host_id
  else
    echo got existing dedicated host $dedicated_host_id
  fi &&\
  $aws ec2 start-instances --instance-ids $AWS_MAC_INSTANCE_ID &&\
  wait_for '[ "running" == "$($aws ec2 describe-instances --instance-ids $AWS_MAC_INSTANCE_ID | jq -r ".Reservations[0].Instances[0].State.Name" | tee /dev/stderr)" ]' \
    "waiting for instance to be running..." 5 50 &&\
  export ip="$($aws ec2 describe-instances --instance-ids $AWS_MAC_INSTANCE_ID | jq -r '.Reservations[0].Instances[0].NetworkInterfaces[0].Association.PublicIp')" &&\
  echo ip=$ip &&\
  wait_for 'scp -i $AWS_MAC_PEM_KEY_PATH -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $amd64_tar_gz ec2-user@$ip:cloudcli-amd64.tar.gz' \
    "waiting for ssh access to instance..." 5 50 &&\
  ssh -i $AWS_MAC_PEM_KEY_PATH -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ec2-user@$ip "tar -xzvf cloudcli-amd64.tar.gz && ls -lah cloudcli gon-config.json"
  ssh -i $AWS_MAC_PEM_KEY_PATH -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ec2-user@$ip '
    echo "cd /Users/ec2-user && /usr/local/bin/gon -log-level debug gon-config.json" | sudo su -l
  ' &&\
  ssh -i $AWS_MAC_PEM_KEY_PATH -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ec2-user@$ip "ls -lah cloudcli.zip" &&\
  scp -i $AWS_MAC_PEM_KEY_PATH -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ec2-user@$ip:cloudcli.zip ./cloudcli-darwin-amd64.zip
}
