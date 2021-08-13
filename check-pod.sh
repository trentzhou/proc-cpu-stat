#!/usr/bin/env bash

# This script is used to check CPU usage details for a pod in kubernetes
# Usage:
#    check-pod.sh NAMESPACE POD
namespace=$1
pod=$2
source_binary=/tmp/proc-cpu-stat
if [ ! -f $source_binary ]; then
    echo "$source_binary does not exist"
    exit 1
fi

kubectl cp -n $namespace $source_binary $pod:$source_binary
kubectl exec -t -n $namespace $pod -- $source_binary
