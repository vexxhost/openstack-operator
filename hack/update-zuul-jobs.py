#!/usr/bin/env python

# Copyright 2020 VEXXHOST, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import copy
import os
import sys

import dockerfile
from ruamel import yaml


yaml = yaml.YAML()
images = os.listdir("images")

build_jobs = []
upload_jobs = []

for image in images:
    files = []
    if image != 'openstack-operator':
        files = ['^images/%s/.*' % image]
    build_deps = ['openstack-operator:images:build:openstack-operator']
    upload_deps = ['openstack-operator:images:upload:openstack-operator']

    job_vars = {
        'docker_images': [
            {
                'context': 'images/%s' % image,
                'repository': 'vexxhost/%s' % image,
            }
        ]
    }

    # Parse the Docker file to see if we have multiple targets
    targets = []
    for line in dockerfile.parse_file('images/%s/Dockerfile' % image):
        if line.cmd != 'from':
            continue
        if len(line.value) >= 3 and line.value[1].lower() != 'as':
            continue
        if line.value[0] != image:
            continue
        targets.append(line.value[2])

    # Update images if we have more than 1 target
    if targets:
        job_vars['docker_images'] = [
            {
                'context': 'images/%s' % image,
                'repository': 'vexxhost/%s' % target,
                'target': target,
            } for target in targets
        ]

    if image == 'openstack-operator':
        job_vars['docker_images'][0]['context'] = '.'
        job_vars['docker_images'][0]['dockerfile'] = 'images/openstack-operator/Dockerfile'

    build_job = {
        'job': {
            'name': 'openstack-operator:images:build:%s' % image,
            'parent': 'vexxhost-build-docker-image',
            'provides': 'openstack-operator:image:%s' % image,
            'vars': job_vars,
        }
    }

    upload_job = {
        'job': {
            'name': 'openstack-operator:images:upload:%s' % image,
            'parent': 'vexxhost-upload-docker-image',
            'provides': 'openstack-operator:image:%s' % image,
            'vars': job_vars,
        }
    }

    if image != 'openstack-operator':
        build_job['job']['dependencies'] = build_deps
        upload_job['job']['dependencies'] = upload_deps

    promote_job = {
        'job': {
            'name': 'openstack-operator:images:promote:%s' % image,
            'parent': 'vexxhost-promote-docker-image',
            'vars': job_vars,
        }
    }

    if files:
        build_job['job']['files'] = files
        upload_job['job']['files'] = files
        promote_job['job']['files'] = files

    project_config = {
        'project': {
            'check': {'jobs': [build_job['job']['name']]},
            'gate': {'jobs': [upload_job['job']['name']]},
            'promote': {'jobs': [promote_job['job']['name']]},
        }
    }

    config = [
        build_job,
        upload_job,
        promote_job,
        project_config
    ]

    if image == 'openstack-operator':
        build_jobs.append(build_job['job']['name'])
        upload_jobs.append(upload_job['job']['name'])
    else:
        build_jobs.append({'name': build_job['job']['name'], 'soft': True})
        upload_jobs.append({'name': upload_job['job']['name'], 'soft': True})

    with open("zuul.d/%s-jobs.yaml" % image, "w+") as fd:
        yaml.dump(config, fd)


with open("zuul.d/functional-jobs.yaml") as fd:
    data = yaml.load(fd)

for obj in data:
    if 'project' in obj:
        for job in obj['project']['check']['jobs']:
            if 'openstack-operator:functional' in job:
                job['openstack-operator:functional']['dependencies'] = \
                    build_jobs
        for job in obj['project']['gate']['jobs']:
            if 'openstack-operator:functional' in job:
                job['openstack-operator:functional']['dependencies'] = \
                    upload_jobs

with open("zuul.d/functional-jobs.yaml", "w+") as fd:
    yaml.dump(data, fd)
