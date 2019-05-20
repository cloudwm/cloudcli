#!/usr/bin/env bash

( [ -z "${TEST_API_SERVER}" ] || [ -z "${TEST_API_CLIENTID}" ] || [ -z "${TEST_API_SECRET}" ] ) && echo missing required env vars && exit 1

DEBUG_OUTPUT_FILE=${DEBUG_OUTPUT_FILE:-/dev/null}

echo "
##### cloudcli init #####
"

rm -f ~/.cloudcli.schema.json

echo "### init with test server should succeed"
OUT="$(cloudcli init --api-clientid "${TEST_API_CLIENTID}" --api-secret "${TEST_API_SECRET}" --api-server "${TEST_API_SERVER}")"
[ "$?" != "0" ] && echo "${OUT}" && echo FAILED: exit code should equal to 0 && exit 1
echo "## OK"

exit 0
