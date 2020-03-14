<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Markdown自动生成目录](#markdown%E8%87%AA%E5%8A%A8%E7%94%9F%E6%88%90%E7%9B%AE%E5%BD%95)
  - [使用npm语法生成](#%E4%BD%BF%E7%94%A8npm%E8%AF%AD%E6%B3%95%E7%94%9F%E6%88%90)
    - [1、安装npm](#1%E5%AE%89%E8%A3%85npm)
    - [2、安装doctoc插件](#2%E5%AE%89%E8%A3%85doctoc%E6%8F%92%E4%BB%B6)
    - [3、执行生成](#3%E6%89%A7%E8%A1%8C%E7%94%9F%E6%88%90)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## Markdown自动生成目录

### 使用npm语法生成

#### 1、安装npm

我的系统是deepin，其他系统的可自行google
````
sudo apt install node
````

#### 2、安装doctoc插件

````
npm i doctoc -g //install 简写 i
````

#### 3、执行生成

切换到md的目录下面，执行下面的命令就能生成
````
$ doctoc markdown自动生成目录.md 

DocToccing single file "markdown自动生成目录.md" for github.com.

==================

"markdown自动生成目录.md" will be updated

Everything is OK.
````

### 参考
【markdown文件生成目录的方式】https://www.jianshu.com/p/b0a18eb32d09  