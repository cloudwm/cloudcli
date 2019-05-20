#!/usr/bin/env bash

( [ -z "${TEST_API_SERVER}" ] || [ -z "${TEST_API_CLIENTID}" ] || [ -z "${TEST_API_SECRET}" ] ) && echo missing required env vars && exit 1

! which cloudcli && echo missing cloudcli binary in PATH && exit 1

if [ "${DEBUG_OUTPUT_FILE}" == "" ]; then
    export DEBUG_OUTPUT_FILE=`mktemp`
fi

echo '########################' >> $DEBUG_OUTPUT_FILE
echo Start: `date` >> $DEBUG_OUTPUT_FILE
echo '########################' >> $DEBUG_OUTPUT_FILE

echo '# Running all tests' | tee -a $DEBUG_OUTPUT_FILE
echo '# Writing to debug file: '$DEBUG_OUTPUT_FILE

! which sshpass && sudo apt-get install -y sshpass
! sudo pip install ruamel.yaml && echo failed to verify dependencies && exit 1

echo "-----" &&\
tests/test_init.sh &&\
echo "-----" &&\
tests/test_server_list.sh &&\
echo "-----" &&\
tests/test_server_options.sh &&\
echo "-----" &&\
python3 tests/test_server_create.py
RES="$?"
echo '########################' >> $DEBUG_OUTPUT_FILE
echo End: `date` >> $DEBUG_OUTPUT_FILE
echo '########################' >> $DEBUG_OUTPUT_FILE
echo | tee -a $DEBUG_OUTPUT_FILE

if [ "$RES" != "0" ]; then
    echo tests failed | tee -a $DEBUG_OUTPUT_FILE
    echo Check the debug output at $DEBUG_OUTPUT_FILE
    exit 1
else
    echo Great Success! | tee -a $DEBUG_OUTPUT_FILE
    echo | tee -a $DEBUG_OUTPUT_FILE
    echo Check the debug output at $DEBUG_OUTPUT_FILE
    exit 0
fi
