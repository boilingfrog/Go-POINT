<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
- [git-fls](#git-fls)
  - [什么是git-fls](#%E4%BB%80%E4%B9%88%E6%98%AFgit-fls)
  - [基本的命令](#%E5%9F%BA%E6%9C%AC%E7%9A%84%E5%91%BD%E4%BB%A4)
<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## git-fls

### 什么是git-fls

`Git LFS`（Large File Storage, 大文件存储）是可以把音乐、图片、视频等指定的任意文件存在 Git 仓库之外，而在 Git 仓库中用一个占用空间 1KB 不到的文本指针来代替的小工具。通过把大文件存储在 Git 仓库之外，可以减小 Git 仓库本身的体积，使克隆 Git 仓库的速度加快，也使得 Git 不会因为仓库中充满大文件而损失性能。  

使用 Git LFS，在默认情况下，只有当前签出的 commit 下的 LFS 对象的当前版本会被下载。此外，我们也可以做配置，只取由 Git LFS 管理的某些特定文件的实际内容，而对于其他由 Git LFS 管理的文件则只保留文件指针，从而节省带宽，加快克隆仓库的速度；也可以配置一次获取大文件的最近版本，从而能方便地检查大文件的近期变动。  

### 基本的命令

下载  

```
git lfs install
```

在仓库中选择要使用lfs管理的文件，通过扩展名，管理这一类的文件  

```
git lfs track "*.psd"
```

添加`.gitattributes`

```
git add .gitattributes
```

然后正常提交文件到仓库中  

```
git add file.psd
git commit -m "Add design file"
git push origin master
```




