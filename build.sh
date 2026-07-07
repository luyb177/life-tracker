#!/bin/bash
set -e

# ============================================================
# Life Tracker — 构建后端 + 前端镜像并推送到阿里云容器镜像服务
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
BACKEND_IMAGE_NAME="life-tracker"
FRONTEND_IMAGE_NAME="life-tracker-web"
BACKEND_IMAGE="${REGISTRY}/${NAMESPACE}/${BACKEND_IMAGE_NAME}:${VERSION}"
FRONTEND_IMAGE="${REGISTRY}/${NAMESPACE}/${FRONTEND_IMAGE_NAME}:${VERSION}"

BACKEND_DOCKERFILE="backend/Dockerfile"
BACKEND_CONTEXT="backend"
FRONTEND_DOCKERFILE="frontend/Dockerfile"
FRONTEND_CONTEXT="frontend"

echo "=========================================="
echo "  Life Tracker Build & Push"
echo "  Backend:  ${BACKEND_IMAGE}"
echo "  Frontend: ${FRONTEND_IMAGE}"
echo "=========================================="

# 1. 构建 Linux AMD64 镜像
echo ""
echo "[1/3] Building linux/amd64 images ..."
docker build \
    --platform linux/amd64 \
    -t "${BACKEND_IMAGE_NAME}:${VERSION}" \
    -f "${BACKEND_DOCKERFILE}" \
    "${BACKEND_CONTEXT}"

docker build \
    --platform linux/amd64 \
    -t "${FRONTEND_IMAGE_NAME}:${VERSION}" \
    -f "${FRONTEND_DOCKERFILE}" \
    "${FRONTEND_CONTEXT}"

# 2. 打 Tag
echo ""
echo "[2/3] Tagging ..."
docker tag "${BACKEND_IMAGE_NAME}:${VERSION}" "${BACKEND_IMAGE}"
docker tag "${FRONTEND_IMAGE_NAME}:${VERSION}" "${FRONTEND_IMAGE}"

# 3. 推送（可选）
if [ "${NOPUSH}" = "nopush" ]; then
    echo ""
    echo "[3/3] Skip push (nopush mode)"
    echo ""
    echo "=========================================="
    echo "  Done — local images:"
    echo "    ${BACKEND_IMAGE_NAME}:${VERSION}"
    echo "    ${FRONTEND_IMAGE_NAME}:${VERSION}"
    echo "  To push manually:"
    echo "    docker push ${BACKEND_IMAGE}"
    echo "    docker push ${FRONTEND_IMAGE}"
    echo "=========================================="
else
    echo ""
    echo "[3/3] Pushing to registry ..."
    docker push "${BACKEND_IMAGE}"
    docker push "${FRONTEND_IMAGE}"

    echo ""
    echo "=========================================="
    echo "  Done — pushed:"
    echo "    ${BACKEND_IMAGE}"
    echo "    ${FRONTEND_IMAGE}"
    echo "  On server:"
    echo "    docker compose pull && docker compose up -d"
    echo "=========================================="
fi
