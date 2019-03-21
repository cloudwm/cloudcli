#!/usr/bin/env bash

( [ -z "${TEST_API_SERVER}" ] || [ -z "${TEST_API_CLIENTID}" ] || [ -z "${TEST_API_SECRET}" ] ) && echo missing required env vars && exit 1

DEBUG_OUTPUT_FILE=${DEBUG_OUTPUT_FILE:-/dev/null}

echo "
##### cloudcli server list #####
"

echo "### running without arguments should fail"
TEMPFILE=`mktemp`
OUT="$(cloudcli server list --no-config)"
[ "$?" == "0" ] && echo "${OUT}" && echo FAILED: exit code should not equal to 0 && exit 1
echo "${OUT}" | grep "ERROR: --api-server flag is required
ERROR: --api-clientid flag is required
ERROR: --api-secret flag is required" >> $DEBUG_OUTPUT_FILE
[ "$?" != "0" ] && echo "${OUT}" && echo FAILED: output should contain error message && exit 1
echo "## OK"

echo "### server list should return human readable summary with at least 1 server"
OUT="$(cloudcli server list --no-config --api-clientid "${TEST_API_CLIENTID}" --api-secret "${TEST_API_SECRET}" --api-server "${TEST_API_SERVER}")"
[ "$?" != "0" ] && echo "${OUT}" && echo FAILED: exit code is not 0 && exit 1
( [ "$(echo "${OUT}" | wc -l)" == "0" ] || [ "$(echo "${OUT}" | wc -l)" == "1" ] ) && echo "${OUT}" && echo FAILED: should return at least 1 server && exit 1
echo "${OUT}" | head -1 | grep "ID " >> $DEBUG_OUTPUT_FILE &&\
echo "${OUT}" | head -1 | grep " NAME " >> $DEBUG_OUTPUT_FILE &&\
echo "${OUT}" | head -1 | grep " DATACENTER " >> $DEBUG_OUTPUT_FILE &&\
echo "${OUT}" | head -1 | grep " POWER" >> $DEBUG_OUTPUT_FILE
[ "$?" != "0" ] && echo "${OUT}" && echo FAILED: invalid header row && exit 1
echo "## OK"

echo "### server list json should return json"
OUT="$(cloudcli server list --format json --no-config --api-clientid "${TEST_API_CLIENTID}" --api-secret "${TEST_API_SECRET}" --api-server "${TEST_API_SERVER}")"
[ "$?" != "0" ] && echo "${OUT}" && echo FAILED: exit code should be 0 && exit 1
echo "${OUT}" | head -1 | grep -- '\[' >> $DEBUG_OUTPUT_FILE &&\
echo "${OUT}" | tail -1 | grep -- '\]' >> $DEBUG_OUTPUT_FILE &&\
echo "${OUT}" | jq . >> $DEBUG_OUTPUT_FILE
[ "$?" != "0" ] && echo "${OUT}" && echo FAILED: invalid json response && exit 1
echo "## OK"

echo "### server list yaml should return yaml"
OUT="$(cloudcli server list --format yaml --no-config --api-clientid "${TEST_API_CLIENTID}" --api-secret "${TEST_API_SECRET}" --api-server "${TEST_API_SERVER}")"
[ "$?" != "0" ] && echo "${OUT}" && echo FAILED: exit code should be 0 && exit 1
echo "${OUT}" | head -1 | grep -- '- ' >> $DEBUG_OUTPUT_FILE &&\
echo "${OUT}" | python3 -c 'import yaml, sys; print(yaml.load(sys.stdin))' >> $DEBUG_OUTPUT_FILE
[ "$?" != "0" ] && echo "${OUT}" && echo FAILED: invalid yaml response && exit 1
echo "## OK"

exit 0
