# blog
功能需求
Done:
- 用户认证系统
- 文章管理(增删改查)
- JWT身份验证
- 集成Markdown编辑器
- 文件上传(图片等)
- 管理员后台
- 分页处理

Todo:
- 文章分类和标签
    文章创建/编辑表单：需要支持选择分类（下拉选择或单选）和输入标签（输入框，支持多选，可以是自由输入或从已有标签中选择）
    文章列表页：可以显示每篇文章的分类和标签，并且可以点击分类或标签来筛选文章。
    文章详情页：显示文章的分类和标签，同样可以点击以筛选同类文章。
    侧边栏或导航：可以展示分类列表和标签云，方便用户按分类或标签浏览。
- 可读权限管理
- 评论系统
- 全文搜索功能(使用 Elasticsearch 或数据库全文索引)
- API文档生成


#### 介绍
基于Golang的Gin框架和Gorm框架实现
Bootstrap前端样式
MySQL数据存储


#### 安装教程

# 更新系统包
sudo apt update && sudo apt upgrade -y
# 安装Go
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
# 安装MySQL
sudo apt install mysql-server -y
# 安装Git
sudo apt install git -y
# 安装必要的构建工具
sudo apt install build-essential -y
# 安全安装MySQL
sudo mysql_secure_installation
# 登录MySQL
sudo mysql -u root -p
# 在MySQL中创建数据库和用户
CREATE DATABASE blog;
CREATE USER 'blog_user'@'%' IDENTIFIED BY '123456';
GRANT ALL PRIVILEGES ON blog.* TO 'blog_user'@'%';
FLUSH PRIVILEGES;
EXIT;
# 克隆项目代码
git clone https://github.com/wrr5/blog.git
cd blog
# 安装Go依赖
go mod tidy
# 构建应用（生产环境建议禁用调试信息和优化二进制大小）
go build -ldflags="-w -s" -o main
# 测试运行
./main
# 创建systemd服务文件，让应用在后台运行并自动重启
sudo nano /etc/systemd/system/blog.service
# 文件内容
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
# 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable blog.service
sudo systemctl start blog.service
sudo systemctl status blog.service  # 检查状态
sudo systemctl stop blog.service # 停止服务
# 查看虚拟机ip
ip addr show

# 更新项目
# 进入项目目录
cd /home/wrr/blog
# 从Git拉取最新代码
git pull origin main
# 下载依赖
go mod tidy
# 构建项目
go build -ldflags="-w -s" -o main
# 重启服务
sudo systemctl restart gin-app.service
# 检查服务状态
sudo systemctl status gin-app.service


#### 使用说明

1.  xxxx
2.  xxxx
3.  xxxx
