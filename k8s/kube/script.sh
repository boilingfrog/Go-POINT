#!/bin/bash

# 设置 KUBECONFIG 环境变量
export KUBECONFIG=./kube/config-test:./kube/config-production

# 列出所有上下文
kubectl config get-contexts

# 提示用户选择上下文
echo "请输入要使用的 k8s 集群的 NAME："
read context_name

# 切换到选择的上下文
kubectl config use-context $context_name

# 显示当前上下文
kubectl config current-context