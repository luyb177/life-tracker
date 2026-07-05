#!/bin/bash
set -e

# ============================================================
# Life Tracker — 构建 & 推送到阿里云容器镜像服务
# ============================================================
# 用法:
#   ./build.sh              # 默认 latest
#   ./build.sh 0.0.2         # 指定版本号
#   ./build.sh 0.0.2 nopush  # 只构建不推送
# ============================================================

VERSION="${1:-latest}"
NOPUSH="${2:-}"

REGISTRY="crpi-u5azhs6neq326bz0.cn-hangzhou.personal.cr.aliyuncs.com"
NAMESPACE="yub_lu"
IMAGE_NAME="life-tracker"
FULL_IMAGE="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${VERSION}"

DOCKERFILE="backend/Dockerfile"
CONTEXT="backend"

echo "=========================================="
echo "  Life Tracker Build & Push"
echo "  Image: ${FULL_IMAGE}"
echo "=========================================="

# 1. 构建 Linux AMD64 镜像
echo ""
echo "[1/3] Building linux/amd64 image ..."
docker build \
    --platform linux/amd64 \
    -t "${IMAGE_NAME}:${VERSION}" \
    -f "${DOCKERFILE}" \
    "${CONTEXT}"

# 2. 打 Tag
echo ""
echo "[2/3] Tagging ..."
docker tag "${IMAGE_NAME}:${VERSION}" "${FULL_IMAGE}"

# 3. 推送（可选）
if [ "${NOPUSH}" = "nopush" ]; then
    echo ""
    echo "[3/3] Skip push (nopush mode)"
    echo ""
    echo "=========================================="
    echo "  Done — local image: ${IMAGE_NAME}:${VERSION}"
    echo "  To push manually:"
    echo "    docker push ${FULL_IMAGE}"
    echo "=========================================="
else
    echo ""
    echo "[3/3] Pushing to registry ..."
    docker push "${FULL_IMAGE}"

    echo ""
    echo "=========================================="
    echo "  Done — pushed: ${FULL_IMAGE}"
    echo "  On server: docker pull ${FULL_IMAGE}"
    echo "=========================================="
fi
