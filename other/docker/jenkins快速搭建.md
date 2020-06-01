## jenkins快速搭建

### 安装

#### 初始化git服务

选择一台虚拟机，首先安装`git`

````
yum install git -y
````

然后初始化仓库

````
# su git
# mkdir repos
# cd repos
# mkdir app.git
# cd app.git
# git --bare init
````