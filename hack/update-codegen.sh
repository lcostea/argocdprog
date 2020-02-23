#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
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

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." >/dev/null 2>&1 && pwd )"
CODEGEN_VERSION=$(go list -m k8s.io/code-generator | awk '{print $2}' | head -1)
CODEGEN_PKG=$(echo `go env GOPATH`"/pkg/mod/k8s.io/code-generator@${CODEGEN_VERSION}")


TEMP_DIR=$(mktemp -d)
cleanup() {
    rm -rf ${TEMP_DIR}
}
trap "cleanup" EXIT SIGINT


chmod +x ${CODEGEN_PKG}/generate-groups.sh


bash "${CODEGEN_PKG}"/generate-groups.sh "deepcopy,client,informer,lister" \
  lcostea.io/knapps/pkg/generated lcostea.io/knapps/pkg/apis \
  knapps:v1alpha1 \
  --output-base "${TEMP_DIR}" \
  --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt

cp -r "${TEMP_DIR}/lcostea.io/knapps/." "${SCRIPT_ROOT}/"