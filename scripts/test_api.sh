#!/bin/bash

# 设置API基础URL
BASE_URL="http://localhost:8080/api/v1"

# 测试获取日志接口 - 无参数
echo "测试获取日志接口 - 无参数"
curl -X GET "${BASE_URL}/logs"
echo -e "\n"

# 测试获取日志接口 - 带日期范围
echo "测试获取日志接口 - 带日期范围"
curl -X GET "${BASE_URL}/logs?start_date=2025-03-01&end_date=2025-03-03"
echo -e "\n"

# 测试创建日志接口 - 基本测试
echo "测试创建日志接口 - 基本测试"
curl -X POST "${BASE_URL}/logs" \
  -H "Content-Type: application/json" \
  -d '{"text":"这是一条测试日志","validation_result":true}'
echo -e "\n"

# 测试创建日志接口 - 空文本
echo "测试创建日志接口 - 空文本"
curl -X POST "${BASE_URL}/logs" \
  -H "Content-Type: application/json" \
  -d '{"text":"","validation_result":true}'
echo -e "\n"

# 测试创建日志接口 - 长文本
echo "测试创建日志接口 - 长文本"
curl -X POST "${BASE_URL}/logs" \
  -H "Content-Type: application/json" \
  -d '{"text":"这是一条非常长的测试日志内容，用于测试API处理长文本的能力。这是一条非常长的测试日志内容，用于测试API处理长文本的能力。这是一条非常长的测试日志内容，用于测试API处理长文本的能力。","validation_result":true}'
echo -e "\n"

# 测试创建日志接口 - 特殊字符
echo "测试创建日志接口 - 特殊字符"
curl -X POST "${BASE_URL}/logs" \
  -H "Content-Type: application/json" \
  -d '{"text":"特殊字符测试：!@#$%^&*()_+-=[]{}|;:,.<>?","validation_result":true}'
echo -e "\n"

# 测试创建日志接口 - 无验证结果
echo "测试创建日志接口 - 无验证结果"
curl -X POST "${BASE_URL}/logs" \
  -H "Content-Type: application/json" \
  -d '{"text":"这是一条没有验证结果的测试日志"}'
echo -e "\n"

# 测试创建日志接口 - 错误的JSON格式
echo "测试创建日志接口 - 错误的JSON格式"
curl -X POST "${BASE_URL}/logs" \
  -H "Content-Type: application/json" \
  -d '{text:"格式错误的JSON"'
echo -e "\n"

# 为脚本添加执行权限
chmod +x "$0"