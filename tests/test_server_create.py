#!/usr/bin/env python3
import sys
import subprocess
import os
import csv
import binascii
import time
import datetime
from ruamel import yaml


def echo_title(title):
    print("\n##### {} #####\n".format(title), file=sys.stderr)


def echo_subtitle(subtitle):
    print("### {}".format(subtitle), file=sys.stderr)


def echo_ok(msg, start_time=None):
    print("## OK: {}{}".format(msg, get_elapsed_time_message(start_time)), file=sys.stderr)


def echo_info(info, start_time=None):
    print(info + get_elapsed_time_message(start_time), file=sys.stderr)


def echo_failed(msg):
    print("## Failed: {}".format(msg), file=sys.stderr)


def get_api_args():
    return "--no-config --api-clientid \"{}\" --api-secret \"{}\" --api-server \"{}\"".format(
        os.environ.get("TEST_API_CLIENTID", ""),
        os.environ.get("TEST_API_SECRET", ""),
        os.environ.get("TEST_API_SERVER", ""),
    )


def get_elapsed_time_message(start_time=None):
    return " (Elapsed time: {} seconds)".format((datetime.datetime.now() - start_time).total_seconds()) if start_time else ""


def create_server(csv_report_writer, test):
    echo_info("Creating server {} in datacenter {} with args: {}".format(test['server_name'], test['datacenter'], test['args']))
    csv_report_writer.writerow(["start_create", test['server_name'], test['datacenter'], "args=", test['args']])
    exitcode, output = subprocess.getstatusoutput(
        "cloudcli server create {} --format yaml --datacenter \"{}\" --name \"{}\" --password \"{}\" {}".format(
            get_api_args(), test['datacenter'], test['server_name'], test['password'], test['args']
        )
    )
    command_id = None
    if exitcode == 0:
        data = yaml.safe_load(output)
        if list(data.keys()) == ['command_ids'] and len(data['command_ids'])==1:
            command_id = data['command_ids'][0]
    return command_id


def assert_running_without_cloudcli_credentials_should_fail(context):
    exitcode, output = subprocess.getstatusoutput("cloudcli server create --no-config")
    assert exitcode != 0
    for line in [
        'ERROR: --api-server flag is required',
        'ERROR: --api-clientid flag is required',
        'ERROR: --api-secret flag is required',
    ]:
        assert line in output


def assert_running_without_required_flags_should_fail(context):
    exitcode, output = subprocess.getstatusoutput("cloudcli server create {}".format(get_api_args()))
    assert exitcode != 0
    assert 'Error: required flag(s) "datacenter", "image", "name", "password" not set' in output


def assert_after_server_powered_on(tests):
    all_have_expected_output_lines = True
    for test in tests:
        test_create_log = test['create_log'] = None
        if test.get('create_command_id'):
            exitcode, output = subprocess.getstatusoutput("cloudcli queue detail {} --id {} --log".format(
                get_api_args(), test['create_command_id']
            ))
            if exitcode == 0:
                test_create_log = test['create_log'] = output
        else:
            test_create_log = test['create_log'] = None
        if test.get('expected_output_lines'):
            if test_create_log:
                for line in test['expected_output_lines'].splitlines():
                    if line.strip() not in test_create_log:
                        echo_failed("server {} does not have expected output line in create log: \"{}\"".format(
                            test['server_name'], line.strip()
                        ))
                        echo_info(test_create_log)
                        all_have_expected_output_lines = False
            else:
                echo_failed("server {} does not have expected output lines in create log".format(test['server_name']))
                all_have_expected_output_lines = False
        ssh_user = test.get('ssh', {}).get('user')
        ssh_hostname_network = test.get('ssh', {}).get('hostname', {}).get('from-network')
        ssh_hostname = None
        if ssh_hostname_network and test_create_log:
            for line in test_create_log.splitlines():
                if line.strip().startswith(ssh_hostname_network):
                    ssh_hostname = line.strip().split(' - ')[1].split(' @ ')[0]
                    break
        if ssh_hostname and ssh_user and test.get('expected_server_commands'):
            for command in test['expected_server_commands']:
                cmd = command.get('cmd')
                if cmd:
                    echo_info("Running command on server {}: {}".format(test['server_name'], cmd))
                    exitcode, output = subprocess.getstatusoutput(
                        "sshpass -p {} ssh -o StrictHostKeyChecking=no root@{} \"{}\"".format(test['password'], ssh_hostname, cmd)
                    )
                    if exitcode == 0:
                        echo_ok("Server command completed successfully")
                        expected_output_lines = command.get('expected_output_lines')
                        if expected_output_lines:
                            for line in expected_output_lines.splitlines():
                                if line.strip() not in output:
                                    echo_failed("Line not found in output: {}".format(line.strip()))
                                    all_have_expected_output_lines = False
                    else:
                        echo_info(output)
                        echo_failed("Failed to run server command")
                        all_have_expected_output_lines = False
    return all_have_expected_output_lines


def assert_running_with_various_flags_should_create_servers(context):
    test_ids = context.get('test-ids')
    timestamp = subprocess.check_output("date +%Y-%m-%d-%H-%m-%s", shell=True).decode().strip()
    csv_report_filename="tests/output/test_server_create_flags-{}.csv".format(timestamp)
    start_time = datetime.datetime.now()
    tests = []
    with open("tests/test_server_create.yaml") as f:
        for test in yaml.safe_load(f):
            test_id = test.get('test-id')
            if test_ids and test_id not in test_ids: continue
            tests.append(test)
    echo_info("Running test IDs:")
    for test in tests:
        echo_info("  {}".format(test.get('test-id')))
    with open(csv_report_filename, "w") as csv_report_file:
        csv_report_writer = csv.writer(csv_report_file)
        csv_report_writer.writerow(["event_name", "server_name", "datacenter", "args"])
        echo_info("Writing test report to {}".format(csv_report_filename))
        for test_number, test in enumerate(tests, 1):
            test['server_name'] = "test-{}-{}".format(timestamp, test_number)
            test['password'] = 'Aa{}'.format(binascii.hexlify(os.urandom(6)).decode())
            test['create_command_id'] = create_server(csv_report_writer, test)
            if test['create_command_id']:
                csv_report_writer.writerow(["create_success", test['server_name'], test['datacenter'], "create_command_id=", test['create_command_id']])
                echo_ok("command id: {}".format(test['create_command_id']))
            else:
                csv_report_writer.writerow(["create_failed", test['server_name'], test['datacenter']])
                echo_failed("Failed to create server")
        echo_info("waiting for servers to be powered on", start_time=start_time)
        all_powered_on = False
        for i in range(1, 50):
            for test in tests:
                exitcode, output = subprocess.getstatusoutput(
                    "cloudcli server list {} | grep \" {} \" | grep \" {} \" | grep ' on$'".format(
                        get_api_args(), test['server_name'], test['datacenter']
                    )
                )
                was_powered_on = test.get('powered on')
                is_powered_on = test['powered on'] = exitcode == 0
                if not was_powered_on and is_powered_on:
                    echo_ok("Server powered on successfully: {} in datacenter {}".format(test['server_name'], test['datacenter']), start_time=start_time)
            all_powered_on = all([test['powered on'] for test in tests if test.get('create_command_id')])
            if all_powered_on:
                break
            echo_info("Sleeping 10 seconds ({}/50)".format(i), start_time=start_time)
            time.sleep(10)
        if all_powered_on:
            csv_report_writer.writerow(["wait_powered_on_success", "all servers powered on"])
            echo_ok("All servers powered on", start_time=start_time)
            all_have_expected_output_lines = assert_after_server_powered_on(tests)
        else:
            csv_report_writer.writerow(["wait_powered_on_failed", "some servers failed to be powered on"])
            echo_failed("Some servers failed to be powered on in the allocated time")
            all_have_expected_output_lines = False
        echo_info("Deleting servers", start_time=start_time)
        all_terminated = True
        for test in tests:
            csv_report_writer.writerow(["start_delete", test['server_name']])
            exitcode, output = subprocess.getstatusoutput(
                "cloudcli server terminate {} --force --name \"{}\"".format(
                    get_api_args(), test['server_name']
                )
            )
            if exitcode == 0:
                csv_report_writer.writerow(["delete_success", test['server_name']])
                echo_ok("server {} was terminated".format(test['server_name']))
            else:
                csv_report_writer.writerow(["delete_failed", test['server_name']])
                echo_failed("server {} was not terminated".format(test['server_name']))
                all_terminated = False
        assert all_powered_on and all_terminated and all_have_expected_output_lines
        echo_ok("All servers were powered on, terminated and have the expected output lines", start_time=start_time)


def main(context):
    echo_title("cloudcli server create")
    for assertion in [
        "assert_running_without_cloudcli_credentials_should_fail",
        "assert_running_without_required_flags_should_fail",
        "assert_running_with_various_flags_should_create_servers"
    ]:
        echo_subtitle(assertion)
        globals()[assertion](context)
        echo_ok("Success")


if __name__ == "__main__":
    context = {
        'test-ids': sys.argv[1:]
    }
    main(context)
