<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [mac中virtualBox添加主机网络报错](#mac%E4%B8%ADvirtualbox%E6%B7%BB%E5%8A%A0%E4%B8%BB%E6%9C%BA%E7%BD%91%E7%BB%9C%E6%8A%A5%E9%94%99)
  - [现场复原](#%E7%8E%B0%E5%9C%BA%E5%A4%8D%E5%8E%9F)
  - [解决方法](#%E8%A7%A3%E5%86%B3%E6%96%B9%E6%B3%95)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## mac中virtualBox添加主机网络报错


### 现场复原

`virtual box`添加主机网络报错

```go
VBoxNetAdpCtl: Error while adding new interface: failed to open /dev/vboxnetctl: No such file or directory
```

### 解决方法

```go
sudo "/Library/Application Support/VirtualBox/LaunchDaemons/VirtualBoxStartup.sh" restart
```

### 参考
【VBoxNetAdpCtl: Error while adding new interface: failed to open /dev/vboxnetctl: No such file or directory】https://github.com/gasolin/foxbox/issues/32