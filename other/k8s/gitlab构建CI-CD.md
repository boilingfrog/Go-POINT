<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [gitlab构建CI/CD](#gitlab%E6%9E%84%E5%BB%BAcicd)
  - [准备](#%E5%87%86%E5%A4%87)
  - [docker部署gitlab](#docker%E9%83%A8%E7%BD%B2gitlab)
  - [docker部署gitlab-runner](#docker%E9%83%A8%E7%BD%B2gitlab-runner)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## gitlab构建CI/CD

### 准备 

### docker部署gitlab

通过docker-compose启动gitlab

```go
version: '3'
services:
  gitlab:
    image: 'gitlab/gitlab-ce:latest'
    restart: always
    hostname: '1.1.1.1'
    environment:
      TZ: 'Asia/Shanghai'
      GITLAB_OMNIBUS_CONFIG: |
        external_url 'http://1.1.1.1:9001'
        gitlab_rails['gitlab_shell_ssh_port'] = 1022
        unicorn['port'] = 8888
        nginx['listen_port'] = 9001
    ports:
      - '9001:9001'
      - '443:443'
      - '1022:22'
    volumes:
      - ./config:/etc/gitlab
      - ./data:/var/opt/gitlab
      - ./losg:/var/log/gitlab
```
### 使用二进制部署gitlab-runner

可参考官方的安装方式[Install GitLab Runner manually on GNU/Linux](https://docs.gitlab.com/runner/install/linux-manually.html)

### gitlab-runner注册


