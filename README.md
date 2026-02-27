## blog功能需求
Done:
- 用户认证系统
- 文章管理(增删改查)
- JWT身份验证
- 集成Markdown编辑器
- 文件上传(图片等)
- 分页处理
- 文章分类
- 文章是否公开(不公开为仅自己可见)
- 管理员后台
    - 分类管理
    - 用户管理
- 评论系统
- 全文搜索功能(使用数据库全文索引)

feat：实现用户管理模块
Todo:
- 阅读数，点赞数，评论数统计
- 基于 Elasticsearch 提供全文搜索接口
- API文档生成


## 介绍
基于Golang的Gin框架和Gorm框架实现
Bootstrap前端样式
MySQL数据存储与Reids缓存
集成protobuf与gRPC构建微服务，用于内部调用 


## Linux Installation
1.  更新系统包: 
sudo apt update && sudo apt upgrade -y
2.  安装Go
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
3.  安装MySQL
sudo apt install mysql-server -y
4.  安装Git
sudo apt install git -y
5.  安装Redis
sudo apt install redis-server -y
sudo systemctl start redis-server
sudo systemctl enable redis-server
6.  安装必要的构建工具
sudo apt install build-essential -y
7.  安全安装MySQL
sudo mysql_secure_installation
8.  登录MySQL
sudo mysql -u root -p
9.  在MySQL中创建数据库和用户
CREATE DATABASE blog;
CREATE USER 'blog_user'@'%' IDENTIFIED BY '123456';
GRANT ALL PRIVILEGES ON blog.* TO 'blog_user'@'%';
FLUSH PRIVILEGES;
EXIT;
10.  克隆项目代码
git clone https://github.com/wrr5/blog.git
cd blog
11. 安装Go依赖
go mod tidy
12. 构建应用（生产环境建议禁用调试信息和优化二进制大小）
go build -ldflags="-w -s" -o main
13. 测试运行
./main
14. 创建systemd服务文件，让应用在后台运行并自动重启
sudo nano /etc/systemd/system/blog.service
文件内容:
[Unit]
Description=Gin GORM Application
After=mysql.service

[Service]
Type=simple
User=your_username
WorkingDirectory=/home/wrr/blog
ExecStart=/home/wrr/blog/main
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
15. 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable blog.service
sudo systemctl start blog.service
sudo systemctl status blog.service  # 检查状态
sudo systemctl stop blog.service # 停止服务
17. 查看虚拟机ip
ip addr show

## 更新项目
1.  进入项目目录
cd /home/wrr/blog
2.  从Git拉取最新代码
git pull origin main
3.  下载依赖
go mod tidy
4.  构建项目
go build -ldflags="-w -s" -o main
5.  重启服务
sudo systemctl restart gin-app.service
6.  检查服务状态
sudo systemctl status gin-app.service


## Windows交叉编译构建Linux应用
$env:GOOS="linux"
$env:GOARCH="amd64" 
$env:CGO_ENABLED="0"
go build -buildvcs=false -ldflags="-s -w" -o order-service main.go

