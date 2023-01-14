#!/usr/bin/env bash

set -e

echo "Running $0 $*..."

PARSED=$(getopt -o "" --longoptions balanserip: --name "$0" -- "$@")

eval set -- "$PARSED"
while true; do
  case "$1" in
    --balanserip)
      balanserip="$2"
      shift 2
      ;;
    *) break ;;
  esac
done

if [ ! "${balanserip}" ]; then
  echo "Error: Balanser IP must be specified."
  exit 1
fi

dns_zone_name='stargate'
dns_zone='littlebit.space'

zones=($(yc dns zone list --format json | jq -r '.[].name'))
if [[ " ${zones[*]} " =~ " ${dns_zone_name} " ]]; then
  echo "The DNS zone ${dns_zone_name} already exists. Skipped."
else
  yc dns zone create \
    --name "${dns_zone_name}" \
    --description "Main DNS zone to expose all cluster resoures" \
    --zone "${dns_zone}." \
    --public-visibility
  echo "The DNS zone ${dns_zone_name} created. OK."
fi

a_records=($(yc dns zone list-records --name ${dns_zone_name} --format json | jq -r '.[] | select(.type == "A").name'))
if [[ " ${a_records[*]} " =~ " ${dns_zone}. " ]]; then
  echo "The DNS record A for ${dns_zone} already exists. Skipped."
else
  yc dns zone add-records \
    --name "${dns_zone_name}" \
    --record "*.${dns_zone}. 600 A ${balanserip}"
  echo "The DNS record A for ${dns_zone_name} created. OK."
fi
