<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [ansible](#ansible)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [常用到的指令](#%E5%B8%B8%E7%94%A8%E5%88%B0%E7%9A%84%E6%8C%87%E4%BB%A4)
    - [查看ip是否可用](#%E6%9F%A5%E7%9C%8Bip%E6%98%AF%E5%90%A6%E5%8F%AF%E7%94%A8)
    - [执行](#%E6%89%A7%E8%A1%8C)
    - [执行,查看日志输出](#%E6%89%A7%E8%A1%8C%E6%9F%A5%E7%9C%8B%E6%97%A5%E5%BF%97%E8%BE%93%E5%87%BA)
    - [查看这个 playbook 的执行会影响到哪些 hosts](#%E6%9F%A5%E7%9C%8B%E8%BF%99%E4%B8%AA-playbook-%E7%9A%84%E6%89%A7%E8%A1%8C%E4%BC%9A%E5%BD%B1%E5%93%8D%E5%88%B0%E5%93%AA%E4%BA%9B-hosts)
  - [设置服务器免密登录](#%E8%AE%BE%E7%BD%AE%E6%9C%8D%E5%8A%A1%E5%99%A8%E5%85%8D%E5%AF%86%E7%99%BB%E5%BD%95)
  - [ansible了解](#ansible%E4%BA%86%E8%A7%A3)
    - [变量名的使用](#%E5%8F%98%E9%87%8F%E5%90%8D%E7%9A%84%E4%BD%BF%E7%94%A8)
    - [playbooks了解](#playbooks%E4%BA%86%E8%A7%A3)
      - [Handlers](#handlers)
      - [task](#task)
      - [register使用](#register%E4%BD%BF%E7%94%A8)
      - [set_fact使用](#set_fact%E4%BD%BF%E7%94%A8)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## ansible

### 前言

用到了就总结下吧

### 常用到的指令

#### 查看ip是否可用
```go
ansible all -m ping 
```
#### 执行
```go
ansible-playbook xxxx.yml  
``` 
#### 执行,查看日志输出
```go
ansible-playbook xxxx.yml -vvv 
``` 
#### 查看这个 playbook 的执行会影响到哪些 hosts  
 ```go
ansible-playbook playbook.yml --list-hosts
```

### 设置服务器免密登录

添加本机的pub,公钥到目标服务器`~/.ssh/authorized_keys`中，然后设置权限`chmod 600 /root/.ssh/authorized_keys`  

### ansible了解

`Ansible`是使用`Python`开发的自动化运维工具，如果这么说比较抽象的话，那么可以说`Ansible`可以让服务器管理人员使用文本来管理服务器，编写一段配置文件，在不同的机器上执行。  

#### 变量名的使用

在使用变量之前最好先知道什么是合法的变量名. 变量名可以为字母,数字以及下划线.变量始终应该以字母开头. “foo_port”是个合法的变量名.”foo5”也是. “foo-port”, “foo port”, “foo.port” 和 “12”则不是合法的变量名.  

#### playbooks了解

`Playbooks`可用于声明配置,更强大的地方在于,在 `playbooks` 中可以编排有序的执行过程,甚至于做到在多组机器间,来回有序的执行特别指定的步骤.并且可以同步或异步的发起任务.  

在运行 `playbook` 时（从上到下执行）,如果一个 `host` 执行 `task` 失败,这个 `host` 将会从整个 `playbook` 的 `rotation` 中移除. 如果发生执行失败的情况,请修正 `playbook` 中的错误,然后重新执行即可.   

`modules` 具有”幂等”性.重复多次执行`playbook`是安全的。  

比如对于创建文件夹，如果不存在就创建，存在了就不创建了。  

##### Handlers

在发生改变时执行的操作  

Handlers 也是一些 task 的列表,通过名字来引用,它们和一般的 task 并没有什么区别.Handlers 是由通知者进行 notify, 如果没有被 notify,handlers 不会执行.不管有多少个通知者进行了 notify,等到 play 中的所有 task 执行完成之后,handlers 也只会被执行一次.  

```go
handlers:
    - name: restart memcached
      service:  name=memcached state=restarted
    - name: restart apache
      service: name=apache state=restarted
```

Handlers 最佳的应用场景是用来重启服务,或者触发系统重启操作.除此以外很少用到了.
 
##### task

对于playbook,我们一般使用 include 语句引用 task 文件的方法，将playbook进行拆分。 

##### register使用

`register`的作用一般用于获取命令输出和判断执行是否成功。  

`register`可以存储指定命令的输出结果到一个自定义的变量中，我们可以通过访问这个自定义的变量来获取命令的输出，然后判断是否执行成功。  

````
- name: Check than logfile exists
  stat: path={{ DATA_PATH }}/mongos/log/mongo.log
  register: logfile_start
  when: MONGO_SYSYTEMLOG_DESTIANTION == "file"

- name: Create log if missing
  file:
    state: touch
    dest: "{{ DATA_PATH }}/mongos/log/mongo.log"
    owner: mongod
    group: mongod
    mode: 0644
  when: ( MONGO_SYSYTEMLOG_DESTIANTION == "file"
        and logfile_start is defined
        and not logfile_start.stat.exists )
````

通过判断`logfile_start`来判断目标目录是否存在。  

##### set_fact使用

`set_fact`用来做变量的赋值。  

````
- name: 注册replicaset_host变量
  set_fact:
    replicaset_host: []

- name: 循环处理host
  set_fact:
    replicaset_host: "{{replicaset_host}} + [ '{{ item }}:{{ MONGO_NET_PORT }}' ]"
  with_items: "{{ groups['mongo'] }}"
````

比如上面注册了一个`replicaset_host`数组，下面通过`with_items`循环对`replicaset_host`进行了赋值操作，之后后面的task就可以直接使用这个变量了。

```
- name: 初始化副本集
  mongodb_replicaset:
    login_host: localhost
    login_port: "{{ MONGO_NET_PORT }}"
    login_user: "{{ MONGO_ROOT_USERNAME }}"
    login_password: "{{ MONGO_ROOT_PASSWORD }}"
    replica_set: mongos
    members: "{{ replicaset_host }}"
```
 




