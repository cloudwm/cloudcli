#!/usr/bin/env python3

# this script runs from our Jenkins to terminate old mac dedicated hosts (older then 25 hours)
# to run it locally following env vars are needed:
#   export AWS_ACCESS_KEY_ID=
#   export AWS_SECRET_ACCESS_KEY=
#   export AWS_REGION=eu-central-1
# run the flow:
#   python3 bin/terminate_old_mac_hosts.py

import subprocess, json, datetime, time


def stop_instance(aws, instance):
    print("Stopping instance: {}".format(instance))
    subprocess.check_call([*aws, 'ec2', 'stop-instances', '--instance-ids', instance['InstanceId']])
    state = None
    for i in range(100):
        res = json.loads(subprocess.check_output([*aws, 'ec2', 'describe-instances', '--instance-ids', instance['InstanceId']]).decode())
        state = res['Reservations'][0]['Instances'][0]['State']['Name']
        if state == 'stopped':
            break
        else:
            print('Waiting for instance to stop ({})...'.format(state))
            time.sleep(30)
    assert state == 'stopped'


def release_host(aws, host):
    print("Releasing host: {}".format(host['HostId']))
    for instance in host['Instances']:
        stop_instance(aws, instance)
    res = json.loads(subprocess.check_output([*aws, 'ec2', 'release-hosts', '--host-ids', host['HostId']]).decode())
    assert len(res['Successful']) == 1, res


aws_cli_image = "amazon/aws-cli@sha256:579f6355a1f153946f73fec93955573700a2eb0b63f9ae853000830cf6bf351a"
aws = ["docker", "run", "-e", "AWS_REGION", "-e", "AWS_ACCESS_KEY_ID", "-e", "AWS_SECRET_ACCESS_KEY", aws_cli_image]
subprocess.check_call(['docker', 'pull', aws_cli_image])
num_unreleasable_hosts = 0
num_released_hosts = 0
for host in json.loads(subprocess.check_output([*aws, 'ec2', 'describe-hosts']).decode())['Hosts']:
    allocation_time = datetime.datetime.strptime(host['AllocationTime'].split('+')[0], '%Y-%m-%dT%H:%M:%S')
    if allocation_time + datetime.timedelta(hours=25) <= datetime.datetime.now():
        release_host(aws, host)
        num_released_hosts += 1
    else:
        num_unreleasable_hosts += 1
print('done. num_unreleasable_hosts={} num_released_hosts={}'.format(num_unreleasable_hosts, num_released_hosts))
