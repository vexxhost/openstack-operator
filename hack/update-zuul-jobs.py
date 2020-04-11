#!/usr/bin/env python

import copy
import os
import sys

from ruamel import yaml


yaml = yaml.YAML()
images = os.listdir("images")

build_jobs = []
upload_jobs = []

for image in images:
    files = []
    if image != 'openstack-operator':
        files = ['^images/%s/.*' % image]
    deps = ['opendev-buildset-registry']

    job_vars = {
        'docker_images': [
            {
                'context': 'images/%s' % image,
                'repository': 'vexxhost/%s' % image,
            }
        ]
    }
    if image == 'openstack-operator':
        job_vars['docker_images'][0]['context'] = '.'
        job_vars['docker_images'][0]['dockerfile'] = 'images/openstack-operator/Dockerfile'

    build_job = {
        'job': {
            'name': 'openstack-operator:images:build:%s' % image,
            'parent': 'vexxhost-build-docker-image',
            'provides': 'openstack-operator:image:%s' % image,
            'dependencies': deps,
            'vars': job_vars,
        }
    }

    upload_job = {
        'job': {
            'name': 'openstack-operator:images:upload:%s' % image,
            'parent': 'vexxhost-upload-docker-image',
            'provides': 'openstack-operator:image:%s' % image,
            'dependencies': deps,
            'vars': job_vars,
        }
    }

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
