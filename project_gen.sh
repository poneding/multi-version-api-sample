#!/usr/bin/env bash

set -o errexit
set -o pipefail

# Turn colors in this script off by setting the NO_COLOR variable in your
# environment to any value:
#
# $ NO_COLOR=1 test.sh
NO_COLOR=${NO_COLOR:-""}
if [ -z "$NO_COLOR" ]; then
  header=$'\e[1;33m'
  reset=$'\e[0m'
else
  header=''
  reset=''
fi

function header_text {
  echo "$header$*$reset"
}

function gen_project {
  header_text "clean project ..."
  rm -rf multi-version-api-sample
  header_text "starting to generate the multi-version-api-sample project ..."
  mkdir multi-version-api-sample
  cd multi-version-api-sample
  header_text "generate base ..."
  kubebuilder init --domain=poneding.com --repo=github.com/poneding/multi-version-api-sample
  kububuilder edit --multigroup
  kubebuilder create api --group sampleapis --version v1 --kind User --resource --controller=false --make=false
  kubebuilder create webhook --group sampleapis --version v1 --kind User --defaulting --programmatic-validation
  kubebuilder create webhook --group sampleapis --version v1 --kind User --conversion --force
  kubebuilder create api --group sampleapis --version v2 --kind User --resource --controller
  kubebuilder create webhook --group sampleapis --version v2 --kind User --defaulting --programmatic-validation
  kubebuilder create webhook --group sampleapis --version v2 --kind User --conversion --force
  go mod tidy
  make
}


gen_project