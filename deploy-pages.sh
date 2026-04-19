#!/bin/bash
#
# AnPin AI 静态页面部署脚本
# 将页面文件复制到 Sub2API 的 data/public/ 目录
#
# 用法:
#   chmod +x deploy-pages.sh
#   ./deploy-pages.sh [Sub2API目录]
#
# 例如:
#   ./deploy-pages.sh /opt/sub2api
#   ./deploy-pages.sh              # 默认 /opt/sub2api
#

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

# Sub2API 安装目录（可通过参数指定）
SUB2API_DIR="${1:-/opt/sub2api}"
PUBLIC_DIR="${SUB2API_DIR}/data/public"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 页面文件映射（中文文件名 -> URL 路径）
declare -A PAGE_MAP=(
    ["首页.html"]="home"
    ["模型广场.html"]="pricing"
    ["教程.html"]="tutorial"
    ["文档.html"]="docs"
    ["服务条款.html"]="terms"
    ["隐私政策.html"]="privacy"
    ["使用政策.html"]="aup"
)

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}   AnPin AI 页面部署${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# 检查 Sub2API 目录
if [ ! -d "$SUB2API_DIR" ]; then
    echo -e "${YELLOW}[WARN]${NC} 目录 ${SUB2API_DIR} 不存在"
    read -p "请输入 Sub2API 安装目录: " SUB2API_DIR
    PUBLIC_DIR="${SUB2API_DIR}/data/public"
    if [ ! -d "$SUB2API_DIR" ]; then
        echo -e "${RED}[ERROR]${NC} 目录不存在: ${SUB2API_DIR}"
        exit 1
    fi
fi

# 创建 public 目录
mkdir -p "$PUBLIC_DIR"

echo -e "${CYAN}[INFO]${NC} 部署目录: ${PUBLIC_DIR}"
echo ""

# 复制文件
count=0
for cn_name in "${!PAGE_MAP[@]}"; do
    url_path="${PAGE_MAP[$cn_name]}"
    src="$SCRIPT_DIR/$cn_name"

    if [ -f "$src" ]; then
        cp "$src" "$PUBLIC_DIR/${url_path}.html"
        echo -e "${GREEN}[OK]${NC}   ${cn_name} -> /${url_path}"
        ((count++))
    else
        echo -e "${YELLOW}[SKIP]${NC} ${cn_name} 不存在"
    fi
done

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   部署完成！共 ${count} 个页面${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "  可访问路径:"
echo "    /home      首页"
echo "    /pricing   模型广场"
echo "    /tutorial  教程"
echo "    /docs      文档"
echo "    /terms     服务条款"
echo "    /privacy   隐私政策"
echo "    /aup       使用政策"
echo ""
echo "  文件位置: ${PUBLIC_DIR}"
echo ""
echo -e "${YELLOW}提示: 如果页面未生效，请重启 Sub2API 服务${NC}"
echo ""
