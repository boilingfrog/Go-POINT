### linux上面修改hostname

### 临时方法
````
hostname xh01
```

### 永久的
修改所有三个主机名：静态、瞬态和灵活主机名：
````
[root@localhost ~]# hostnamectl set-hostname xh00
[root@localhost ~]# hostnamectl --pretty
[root@localhost ~]# hostnamectl --static
xh00
[root@localhost ~]# hostnamectl --transient
````