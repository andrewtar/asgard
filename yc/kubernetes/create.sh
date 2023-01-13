#!/usr/bin/env bash

# Based on the instruction https://cloud.yandex.com/en-ru/docs/managed-kubernetes/tutorials/new-kubernetes-project.
# This script is idempotent. 

folder_id=$(yc config get folder-id)
cluster_name='asgard'
cluster_version='1.23' # yc managed-kubernetes list-versions
cluster_nodes='asgard-nodes'
key_name='hofund'
network_name='asgard-network'
subnet_name='asgard-subnet'
account_name='odin'

accounts=($(yc iam service-account list --format json | jq -r '.[].name'))
if [[ " ${accounts[*]} " =~ " ${account_name} " ]]; then
  echo "The account ${account_name} already exists. Skipped."
else
  yc iam service-account create \
    --name ${account_name} \
    --description "Kubernetes resource account"
  echo "The account ${account_name} created. OK."
fi
roles=(
  editor
  container-registry.images.puller
  alb.editor
  vpc.publicAdmin
  certificate-manager.certificates.downloader
  compute.viewer
)
account_id=$(yc iam service-account get --name ${account_name} --format json | jq .id -r)
for role in "${roles[@]}"
do
  yc resource-manager folder add-access-binding \
    --id $folder_id \
    --role "${role}" \
    --subject serviceAccount:$account_id
  echo "Role binding"${role}" for account ${account_name}. OK."
done

keys=($(yc kms symmetric-key list --format json | jq -r '.[].name'))
if [[ " ${keys[*]} " =~ " ${key_name} " ]]; then
  echo "The symmetric-key ${key_name} already exists. Skipped."
else
  yc kms symmetric-key create \
    --name "${key_name}" \
    --description "Encrypts the kubernetes secrets" \
    --default-algorithm aes-256 \
    --rotation-period 24h
  echo "The symmetric-key ${key_name} created. OK."
fi

networks=($(yc vpc network list --format json | jq -r '.[].name'))
if [[ " ${networks[*]} " =~ " ${network_name} " ]]; then
  echo "The network ${network_name} already exists. Skipped."
else
  yc vpc network create \
    --name "${network_name}" \
    --description "Kubernetes network"
  echo "The network ${network_name} created. OK."
fi

subnets=($(yc vpc subnet list --format json | jq -r '.[].name'))
if [[ " ${subnets[*]} " =~ " ${subnet_name} " ]]; then
  echo "The subnet ${subnet_name} already exists. Skipped."
else
  yc vpc subnet create \
    --name "${subnet_name}" \
    --description "Kubernetes subnet" \
    --network-name "${network_name}" \
    --zone ru-central1-a \
    --range 192.168.0.0/24
  echo "The subnet ${subnet_name} created. OK."
fi

clusters=($(yc managed-kubernetes cluster list --format json | jq -r '.[].name'))
if [[ " ${clusters[*]} " =~ " ${cluster_name} " ]]; then
  echo "The kubernetes cluster ${cluster_name} already exists. Skipped."
else
  yc managed-kubernetes cluster create \
    --name "${cluster_name}" \
    --description 'My kubernetes cluster' \
    --network-name "${network_name}" \
    --subnet-name "${subnet_name}" \
    --zone 'ru-central1-a' \
    --cluster-ipv4-range '10.0.0.0/16' \
    --service-ipv4-range '172.16.0.0/16' \
    --release-channel 'regular' \
    --version "${cluster_version}" \
    --service-account-name "${account_name}" \
    --node-service-account-name "${account_name}" \
    --auto-upgrade \
    --kms-key-name "${key_name}" \
    --public-ip
  echo "The kubernetes cluster ${cluster_name} created. OK."
fi

nodes=($(yc managed-kubernetes node-group list --format json | jq -r '.[].name'))
if [[ " ${nodes[*]} " =~ " ${cluster_nodes} " ]]; then
  echo "Nodes for ${cluster_nodes} for cluster already exists. Skipped."
else
  yc managed-kubernetes node-group create \
    --name "${cluster_nodes}" \
    --description 'Kubernetes nodes' \
    --cluster-name "${cluster_name}" \
    --platform 'standard-v3' \
    --memory '16GB' \
    --cores 4 \
    --core-fraction 50 \
    --anytime-maintenance-window \
    --disk-size '50GB' \
    --disk-type 'network-hdd' \
    --fixed-size 1 \
    --location "subnet-name=${subnet_name},zone=ru-central1-a" \
    --version "${cluster_version}"
  echo "Nodes ${cluster_nodes} created. OK."
fi

yc managed-kubernetes cluster get-credentials \
  "${cluster_name}" \
  --external \
  --force

echo "The kubernetes cluster setup is finished."
