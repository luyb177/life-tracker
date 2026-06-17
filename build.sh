#!/bin/bash
set -e

REGISTRY=crpi-u5azhs6neq326bz0.cn-hangzhou.personal.cr.aliyuncs.com
NAMESPACE=yub_lu
IMAGE=life-tracker
TAG=0.0.1

echo "=== Building linux/amd64 image ==="
docker build --platform linux/amd64 \
  -t ${REGISTRY}/${NAMESPACE}/${IMAGE}:${TAG} \
  -f backend/Dockerfile backend/

echo "=== Pushing to ${REGISTRY} ==="
docker push ${REGISTRY}/${NAMESPACE}/${IMAGE}:${TAG}

echo "=== Done: ${REGISTRY}/${NAMESPACE}/${IMAGE}:${TAG} ==="
