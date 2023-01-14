#!/bin/bash
for file in yc/namespace/buildbarn/configuration/*.yaml; do
  kubectl apply -f $file
done
