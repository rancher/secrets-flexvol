#!/bin/bash

mkdir -p /usr/libexec/kubernetes/kubelet-plugins/volume/exec/rancher~secrets
cp /usr/bin/secrets-flexvol /usr/libexec/kubernetes/kubelet-plugins/volume/exec/rancher~secrets/secrets

/bin/cat
