---
layout: "docs"
page_title: "OpenStack Post-Processor"
description: |-
  The Packer OpenStack post-processor takes an artifact from the OpenStack builder and uploads it to an OpenStack Image Service endpoint (such as a Glance instance).
---

# OpenStack Post-Processor

Type: `vsphere`

The Packer OpenStack post-processor takes an artifact from the
OpenStack builder and uploads it to an OpenStack Image Service
endpoint (such as a Glance instance).

## Configuration

Configuration options are required in order to identify the Image
Service instance and authenticate against it.

### Required:

* `identity_endpoint` (string) - A URL to a service satisfying [the
  Identity API v3](http://developer.openstack.org/api-ref-identity-v3.html).

* `username` (string) - The username used to connect to the OpenStack service.
  If not specified, Packer will use the environment variable
  `OS_USERNAME`, if set.

* `password` (string) - The password used to connect to the OpenStack service.
  If not specified, Packer will use the environment variables
  `OS_PASSWORD`, if set.

* `image_name` (string) - A human-readable name to identify the image
  in the image service.

### Optional:

* `api_key` (string) - The API key used to access OpenStack. Some OpenStack
  installations require this.

* `tenant_id` or `tenant_name` (string). If not specified, Packer will
  use the environment variable `OS_TENANT_NAME`, if set.

* `region` (string) - The name of the region, such as "DFW", in which
  to launch the server to create the AMI.
  If not specified, Packer will use the environment variable
  `OS_REGION_NAME`, if set.
