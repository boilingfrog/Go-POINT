## ansible

### 前言

用到了就总结下吧

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
 
 
#### 常用到的指令

查看ip是否可用
```go
ansible all -m ping 
```
执行
```go
ansible-playbook xxxx.yml  
``` 
查看这个 playbook 的执行会影响到哪些 hosts  
 ```go
ansible-playbook playbook.yml --list-hosts
```




