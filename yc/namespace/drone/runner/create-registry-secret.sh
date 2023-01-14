#!/usr/bin/env bash

set -e

tmp_file="$(mktemp /tmp/service-accountkey.XXXXXX)"

yc iam key create \
  --description "Auth credentials to access the docker registry from the Drone CI" \
  --service-account-name drone \
  -o "${tmp_file}"

valueKey=`cat ${tmp_file}`

valueSeviceCredentials=$(echo "json_key:$valueKey" | base64 -w0)
echo "Base64 encoded docker credentials:"
echo "{
  \"auths\": {
    \"cr.yandex\": {
      \"auth\": \"${valueSeviceCredentials}\"
    }
  }
}" | base64 -w0

rm "${tmp_file}"
