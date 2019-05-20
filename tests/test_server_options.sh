#!/usr/bin/env bash

( [ -z "${TEST_API_SERVER}" ] || [ -z "${TEST_API_CLIENTID}" ] || [ -z "${TEST_API_SECRET}" ] ) && echo missing required env vars && exit 1

DEBUG_OUTPUT_FILE=${DEBUG_OUTPUT_FILE:-/dev/null}

echo "
##### cloudcli server options #####
"

echo "### running without arguments should fail"
TEMPFILE=`mktemp`
OUT="$(cloudcli server options --no-config)"
[ "$?" == "0" ] && echo "${OUT}" && echo FAILED: exit code should not equal to 0 && exit 1
echo "${OUT}" | grep "ERROR: --api-server flag is required
ERROR: --api-clientid flag is required
ERROR: --api-secret flag is required" >> $DEBUG_OUTPUT_FILE
[ "$?" != "0" ] && echo "${OUT}" && echo FAILED: output should contain error message && exit 1
echo "## OK"

# TODO: add ram to list of options, it currently fails due to change in RAM response
SERVER_OPTIONS="billing cpu datacenter disk image network traffic"

echo "### server options should return human responses for all arg options"

for A in $SERVER_OPTIONS; do
    cloudcli server options --no-config --api-clientid "${TEST_API_CLIENTID}" --api-secret "${TEST_API_SECRET}" --api-server "${TEST_API_SERVER}" \
        --${A} --cache >> ${DEBUG_OUTPUT_FILE}
    [ "$?" != "0" ] && echo FAILED: arg --${A} returned failure exit code && exit 1
done

echo "### server options should return successfully for all arg options"

for F in yaml json; do
    for A in billing cpu datacenter disk image network traffic; do
        cloudcli server options --no-config --api-clientid "${TEST_API_CLIENTID}" --api-secret "${TEST_API_SECRET}" --api-server "${TEST_API_SERVER}" \
            --${A} --cache --format "${F}" >> ${DEBUG_OUTPUT_FILE}
        [ "$?" != "0" ] && echo FAILED: arg --${A} returned failure exit code && exit 1
    done
done
echo "## OK"

exit 0
