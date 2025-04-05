#!/bin/bash

# 设置 KUBECONFIG 环境变量
export KUBECONFIG=$(ls ./kube/config-*.yaml | tr '\n' ':' | sed 's/:$//')


# 获取所有上下文
contexts=$(kubectl config get-contexts -o name)

# 重命名上下文
for context in $contexts; do
  for file in ./kube/config-*.yaml; do
    if grep -q $context $file; then
      suffix=$(basename $file | sed 's/config-//;s/.yaml//')
      kubectl config rename-context $context $suffix
    fi
  done
done

# 列出所有上下文
kubectl config get-contexts

# 合并 kubeconfig 文件
kubectl config view --merge --flatten > ~/.kube/config

export KUBECONFIG=~/.kube/config