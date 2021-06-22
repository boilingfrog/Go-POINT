<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [gitlab构建CI/CD](#gitlab%E6%9E%84%E5%BB%BAcicd)
  - [准备](#%E5%87%86%E5%A4%87)
  - [docker部署gitlab](#docker%E9%83%A8%E7%BD%B2gitlab)
  - [使用二进制部署gitlab-runner](#%E4%BD%BF%E7%94%A8%E4%BA%8C%E8%BF%9B%E5%88%B6%E9%83%A8%E7%BD%B2gitlab-runner)
  - [gitlab-runner注册](#gitlab-runner%E6%B3%A8%E5%86%8C)
  - [配置Variables](#%E9%85%8D%E7%BD%AEvariables)
  - [简单先来个测试](#%E7%AE%80%E5%8D%95%E5%85%88%E6%9D%A5%E4%B8%AA%E6%B5%8B%E8%AF%95)
  - [开始构建](#%E5%BC%80%E5%A7%8B%E6%9E%84%E5%BB%BA)
  - [遇到的报错](#%E9%81%87%E5%88%B0%E7%9A%84%E6%8A%A5%E9%94%99)

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

安装完成之后使用注册命令注册  

```
$ gitlab-runner register
```

然后会提示输入gitlab的地址以及token信息  

地址信息和token我们在下面可以看到

<img src="/img/gitlab-runner_1.jpg" alt="gitlab-runner" align=center />

根据提示输入信息,需要注意的是里面的`tags`就是我们编写`.gitlab-ci.yml`对应填写的tag

<img src="/img/gitlab-runner_2.jpg" alt="gitlab-runner" align=center />

之后在`gitlab`的`runner`中就可以看到我们注册的`gitlab-runner`了  

<img src="/img/gitlab-runner_3.jpg" alt="gitlab-runner" align=center />

### 配置Variables

在gitlab中可以配置我们gitlab-runner需要的变量，比如我们的docker-hub的密码，gitlab的账号密码等信息  

<img src="/img/gitlab-runner_4.jpg" alt="gitlab-runner" align=center />

### 简单先来个测试

先来个简单的`gitlab-ci.yml`测试下

```
stages:
  - test
  - build

variables:
  GOPROXY: https://goproxy.cn

lint:
  stage: test
  script:
    - echo "hello world lint"
  only:
    - branches
  tags:
    - golang-runner

test:
  stage: test
  script:
    - echo "hello world test"
  only:
    - branches
  cache:
    key: "bazel"
    paths:
      - .cache
  tags:
    - golang-runner
```

<img src="/img/gitlab-runner_5.jpg" alt="gitlab-runner" align=center />

### 开始构建

通过`helm`和`bazel`实现在`gitlab-runner`中k8s应用的自动编译，发布。  

镜像推送到`docker-hub`中,`gitlab-runner`中的`helm`需要配置好，这里我是用了`helm`默认初始化的`charts`结构来发布应用  

`gitlab-ci.yml`

```
stages:
  - test
  - build
  - deploy

variables:
  GOPROXY: https://goproxy.cn

lint:
  stage: test
  script:
    - export GO_PROJECT_PATH="/home/gitlab-runner/goWork/src"
    - mkdir -p $GO_PROJECT_PATH
    - ln -s $(pwd) $GO_PROJECT_PATH/test
    - cd $GO_PROJECT_PATH/test
    - bash build/lint.sh
  only:
    - branches
  tags:
    - golang-runner

test:
  stage: test
  script:
    - go mod vendor
    - bash build/bazel-test.sh
  only:
    - branches
  cache:
    key: "bazel"
    paths:
      - .cache
  tags:
    - golang-runner

build:
  stage: build
  before_script:
    - url_host=`git remote get-url origin | sed -e "s/http:\/\/gitlab-ci-token:.*@//g"`
    - git remote set-url origin "http://$GIT_ACCESS_USER:$GIT_ACCESS_PASSWORD@${url_host}"
    - git config user.name $GIT_ACCESS_USER
    - git config user.email $GIT_ACCESS_EMAIL
    - git fetch --tags --force
  script:
    - docker login -u $DOCKER_ACCESS_USER -p $DOCKER_ACCESS_PASSWORD
    - go mod vendor
    - bash build/bazel-build.sh
  only:
    - master
  cache:
    key: "bazel"
    paths:
      - .cache
  tags:
    - golang-runner



include: '/.gitlab/deploy.yaml'
```

镜像打包

```
test::util:build_docker_images() {
  local docker_registry=$1
  local docker_tag=$2
  local base_image="alpine:3.7"

  query=$(test::util::find_changes)

  if [ "$query" == "" ]; then
    test::util::log "no change and exit..."
    exit 0
  fi


  for b in ${query}; do
    b=${b//\/\/src/"/src"}

    if [[ $b == *test* ]]
        then
        continue
    fi

    local binary_file_path=$(test::util::find_binary "$b")
    local binary_name=$(test::util::get_binary_name "$b")
    local docker_build_path="dockerbuild/${binary_name}"
    local docker_file_path="${docker_build_path}/Dockerfile"
    local docker_image_tag="${docker_registry}/${binary_name}:${docker_tag}"


    test::util::log "Starting docker build for image: ${binary_name}"
    (
      rm -rf "${docker_build_path}"
      mkdir -p "${docker_build_path}"
      cp "${binary_file_path}" "${docker_build_path}/${binary_name}"
      cat <<EOF >"${docker_file_path}"
FROM ${base_image}
COPY ${binary_name} /usr/local/bin/${binary_name}
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
  && apk update --no-cache \
ENTRYPOINT ["/usr/local/bin/${binary_name}"]
EOF
      docker build -q -t "${docker_image_tag}" "${docker_build_path}"
      docker push ${docker_image_tag}
    )


  cat <<EOF >>".gitlab/deploy.yaml"
${binary_name}:
  stage: deploy
  script:
    - bash build/deploy.sh ${docker_registry} ${binary_name} ${docker_tag}
  only:
    - tags
  when: manual
  environment:
    name: test
  tags:
    - golang-runner
EOF

done

  test::util::log "Docker builds done"
}
```

整体的处理思路是

1、通过`bazel`构建go项目。  

2、构建的时候找到有改动的项目，编译，打包镜像，生成`deploy`脚本。    

3、打上tag，推到`gitlab`中。

4、最后通过手动触发项目的`deploy`，通过`helm`发布对应的项目到`k8s`中。

项目的地址[gitlab-runner构建go项目](https://github.com/boilingfrog/gitlab-cicd-test)

### 遇到的报错

`go: writing go.mod cache: mkdir /home/goWork: permission denied
`

给用户`gitlab-runner`添加最用目录的执行权限

```
sudo chown -R $(whoami):gitlab-runner /Users/zhushuyan/go/pkg && sudo chmod -R g+rwx /Users/zhushuyan/go/pkg
```








