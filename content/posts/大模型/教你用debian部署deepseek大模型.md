---
title: "教你用debian部署deepseek大模型"
date: 2024-08-15
draft: true
tags : [                    # 文章所属标签
    "AI",
    "LLM",
    "deepseek"
]
---
参考：https://github.com/deepseek-ai/DeepSeek-V3?tab=readme-ov-file#63-inference-with-lmdeploy-recommended

1. 安装python基础环境

```
lex@debian:~/deepseek_project$ python3 -m venv ds_venv
lex@debian:~/deepseek_project$ source ds_venv/bin/activate
(ds_venv) lex@debian:~/deepseek_project$
(ds_venv) lex@debian:~/deepseek_project$ pip install lmdeploy

````



2. 安装torch等
conda create -n lmdeploy python=3.8 -y
conda activate lmdeploy
pip install lmdeploy