#!/bin/bash

# 设置项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# 创建构建目录
BUILD_DIR="${PROJECT_ROOT}/build"
mkdir -p "${BUILD_DIR}"

# 版本信息
VERSION="1.0.0"

# 定义目标平台
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

# 编译函数
build() {
    local os=$1
    local arch=$2
    local output_name="chat_log_server"
    
    # Windows平台添加.exe后缀
    if [ "$os" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "正在构建 $os/$arch..."
    
    # 设置输出目录
    local output_dir="${BUILD_DIR}/${os}_${arch}"
    mkdir -p "$output_dir"
    
    # 构建可执行文件
    GOOS=$os GOARCH=$arch go build -o "${output_dir}/${output_name}" "${PROJECT_ROOT}/main.go"
    
    # 复制配置文件和README
    cp "${PROJECT_ROOT}/README.md" "$output_dir/"
    
    # 创建版本信息文件
    echo "Chat Log Server v${VERSION}" > "${output_dir}/version.txt"
    echo "Build for $os/$arch" >> "${output_dir}/version.txt"
    
    # 创建压缩包
    local archive_name="chat_log_server_${VERSION}_${os}_${arch}"
    if [ "$os" = "windows" ]; then
        (cd "$BUILD_DIR" && zip -r "${archive_name}.zip" "${os}_${arch}")
    else
        (cd "$BUILD_DIR" && tar czf "${archive_name}.tar.gz" "${os}_${arch}")
    fi
}

# 清理旧的构建文件
echo "清理旧的构建文件..."
rm -rf "$BUILD_DIR"

# 为每个平台构建
for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r os arch <<< "$platform"
    build "$os" "$arch"
done

echo "
构建完成！

使用说明：
1. 根据您的操作系统和架构选择对应的压缩包：
   - MacOS (Intel): chat_log_server_${VERSION}_darwin_amd64.tar.gz
   - MacOS (M1/M2): chat_log_server_${VERSION}_darwin_arm64.tar.gz
   - Linux (x64): chat_log_server_${VERSION}_linux_amd64.tar.gz
   - Linux (ARM64): chat_log_server_${VERSION}_linux_arm64.tar.gz
   - Windows: chat_log_server_${VERSION}_windows_amd64.zip

2. 解压下载的压缩包
   - Linux/MacOS: tar -xzf chat_log_server_*.tar.gz
   - Windows: 使用解压软件解压.zip文件

3. 运行服务器
   - Linux/MacOS: ./chat_log_server
   - Windows: 双击chat_log_server.exe

服务器将在 http://localhost:8080 上启动
"