#!/bin/bash
# Database Backup Script for Gomall
# 使用方法: ./backup.sh [backup_dir]

set -e

# 配置
BACKUP_DIR=${1:-./backups}
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/gomall_${DATE}.sql.gz"
DB_HOST=${GOMALL_DB_HOST:-localhost}
DB_PORT=${GOMALL_DB_PORT:-3306}
DB_USER=${GOMALL_DB_USER:-root}
DB_PASSWORD=${GOMALL_DB_PASSWORD:-password}
DB_NAME=${GOMALL_DB_NAME:-Gomall}

# 保留最近多少天的备份
RETENTION_DAYS=7

# 创建备份目录
mkdir -p "${BACKUP_DIR}"

echo "开始备份数据库: ${DB_NAME}"
echo "备份文件: ${BACKUP_FILE}"

# 执行备份
mysqldump -h "${DB_HOST}" -P "${DB_PORT}" -u "${DB_USER}" -p"${DB_PASSWORD}" \
    --single-transaction --routines --triggers \
    "${DB_NAME}" | gzip > "${BACKUP_FILE}"

# 验证备份
if [ -f "${BACKUP_FILE}" ]; then
    FILE_SIZE=$(stat -f%z "${BACKUP_FILE}" 2>/dev/null || stat -c%s "${BACKUP_FILE}" 2>/dev/null || echo "unknown")
    echo "备份成功! 文件大小: ${FILE_SIZE} bytes"
else
    echo "备份失败!"
    exit 1
fi

# 清理旧备份
echo "清理 ${RETENTION_DAYS} 天前的备份..."
find "${BACKUP_DIR}" -name "gomall_*.sql.gz" -mtime +${RETENTION_DAYS} -delete

# 列出当前备份
echo ""
echo "当前备份文件:"
ls -lh "${BACKUP_DIR}"/gomall_*.sql.gz 2>/dev/null || echo "无备份文件"

echo ""
echo "备份完成!"
