## jenkins快速搭建

### 安装

#### 初始化git服务

选择一台虚拟机，首先安装`git`

````
yum install git -y
````

然后初始化仓库

````
# useradd git
# passwd git
# su - git
# mkdir repos
# cd repos
# mkdir app.git
# cd app.git
# git --bare init
````

使用其他的机器查看配置的`git`服务器是否可以成功的拉取

````
# mkdir test
# cd test
# git clone git@192.168.56.203:/home/git/repos/app.git
````

#### 安装jenkins

````
sudo wget -O /etc/yum.repos.d/jenkins.repo https://pkg.jenkins.io/redhat-stable/jenkins.repo
sudo rpm --import https://pkg.jenkins.io/redhat-stable/jenkins.io.key
yum install jenkins
````

启动的时候需要依赖`java`环境，先去安装下。下载直接到官网下载即可`http://www.oracle.com/technetwork/java/javase/downloads/jdk8-downloads-2133151.html`

1、创建安装目录

````
mkdir /usr/local/java/
````

2、解压至安装目录

````
tar -zxvf jdk-8u171-linux-x64.tar.gz -C /usr/local/java/

````

3、设置环境变量

````
vim /etc/profile
````

在末尾添加

````
export JAVA_HOME=/usr/local/java/jdk1.8.0_171
export JRE_HOME=${JAVA_HOME}/jre
export CLASSPATH=.:${JAVA_HOME}/lib:${JRE_HOME}/lib
export PATH=${JAVA_HOME}/bin:$PATH
````

使环境变量生效

````
source /etc/profile
````

添加软链接

````
ln -s /usr/local/java/jdk1.8.0_171/bin/java /usr/bin/java
````

查看版本

````
java -version
````


#### 启动jenkins

````
# systemctl start jenkins
````