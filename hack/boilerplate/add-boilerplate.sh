#!/bin/bash

# Copyright 2018 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Modifications Copyright 2020 Paulhindemith
#
# The original source code can be referenced from the link below.
# https://github.com/knative/serving/blob/9f64f866d633b5cf7ffc4e50e3bc327fd9a3a924/hack/boilerplate/add-boilerplate.sh
# The change history can be obtained by looking at the differences from the
# following commit that added as the original source code.
# 550faea6bb43f0d7fa6a214dc29b5e9760bfe066

USAGE=$(cat <<EOF
Add boilerplate.<ext>.txt to all .<ext> files missing it in a directory.
Usage: (from repository root)
        ./vendor/github.com/paulhindemith/dev-infra/hack/boilerplate/add-boilerplate.sh <ext> <DIR>
Example: (from repository root)
        ./vendor/github.com/paulhindemith/dev-infra/hack/boilerplate/add-boilerplate.sh go cmd
EOF
)

set -e

if [[ -z $1 || -z $2 ]]; then
  echo "${USAGE}"
  exit 1
fi

grep -r -L -P "Copyright \d+ Paulhindemith" $2  \
  | grep -P "\.$1\$" \
  | xargs -I {} sh -c \
  "cat $(dirname ${BASH_SOURCE[0]})/boilerplate.$1.txt {} > /tmp/boilerplate && mv /tmp/boilerplate {}"
