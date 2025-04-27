
帮忙写一个shell脚本，完成如下任务 
1. 把root目录下的文件nem-7202.zip解压
2. 把解压过的文件复制到pod里：
    查看pod方法：
    [root@matrix-01-152 ~]# kubectl get pods -n dc
NAME                                      READY   STATUS    RESTARTS   AGE
cert-manager-7c46c49b6c-xprbm             1/1     Running   0          7m47s
cpm-594fc7cb9c-glp4x                      1/1     Running   0          7m43s
dc1-5d9cfb68cd-d7w8p                      1/1     Running   0          5m58s
dc2-975f8fdd9-42jmc                       1/1     Running   0          23h
dc3-7cf4f88fcd-plm2t                      1/1     Running   0          67s
复制到dc1-xxx、dc2-xxx、dc3-xxx的pod里的/opt/sdn/virgo/bin/lib目录下
3. 分别进上面的pod里重启服务`systemctrl restart sdnc.service`
4. 清理工作目录

脚本如下

```bash
#!/bin/bash

# 打印执行的命令
set -x

# 定义变量
ZIP_FILE="/root/nem-7202.zip"
TEMP_DIR="/tmp/nem_extract"
POD_NAMESPACE="dc"
TARGET_DIR="/opt/sdn/virgo/bin/lib"
SERVICE_NAME="sdnc.service"

# 1. 解压文件
echo "解压文件 $ZIP_FILE..."
echo "执行: mkdir -p \"$TEMP_DIR\" && unzip -o \"$ZIP_FILE\" -d \"$TEMP_DIR\""
mkdir -p "$TEMP_DIR"
unzip -o "$ZIP_FILE" -d "$TEMP_DIR"
if [ $? -ne 0 ]; then
    echo "解压失败!"
    exit 1
fi

# 获取解压后的文件列表
echo "执行: ls \"$TEMP_DIR\""
EXTRACTED_FILES=("$TEMP_DIR"/*)
if [ ${#EXTRACTED_FILES[@]} -eq 0 ]; then
    echo "解压目录为空!"
    exit 1
fi

# 2. 获取目标POD列表
echo "获取目标POD..."
echo "执行: kubectl get pods -n \"$POD_NAMESPACE\" --no-headers | grep -E 'dc1-|dc2-|dc3-' | awk '{print \$1}'"
PODS=$(kubectl get pods -n "$POD_NAMESPACE" --no-headers | grep -E 'dc1-|dc2-|dc3-' | awk '{print $1}')
if [ -z "$PODS" ]; then
    echo "没有找到目标POD!"
    exit 1
fi

# 3. 复制文件到每个POD并重启服务
for POD in $PODS; do
    echo "处理POD: $POD"

    # 复制文件
    for FILE in "${EXTRACTED_FILES[@]}"; do
        kubectl cp "$FILE" -n $POD_NAMESPACE "$POD:$TARGET_DIR/"
        if [ $? -ne 0 ]; then
            echo "复制文件 $FILE 到 $POD 失败!"
            continue
        fi
    done

    # 重启服务
    echo "执行: kubectl exec -n \"$POD_NAMESPACE\" \"$POD\" -- systemctl restart \"$SERVICE_NAME\""
    kubectl exec -n "$POD_NAMESPACE" "$POD" -- systemctl restart "$SERVICE_NAME"
    if [ $? -ne 0 ]; then
        echo "重启服务 $SERVICE_NAME 失败!"
    else
        echo "服务 $SERVICE_NAME 重启成功!"
    fi
done

# 4. 清理工作目录
echo "清理临时目录 $TEMP_DIR..."
echo "执行: rm -rf \"$TEMP_DIR\""
rm -rf "$TEMP_DIR"

echo "所有操作完成!"

```

#!/bin/bash

# 定义变量
ZIP_FILE="/root/nem-7202.zip"
TEMP_DIR="/tmp/nem_extract"
POD_NAMESPACE="dc"
TARGET_DIR="/opt/sdn/virgo/bin/lib"
SERVICE_NAME="sdnc.service"

# 1. 解压文件
echo "解压文件 $ZIP_FILE..."
mkdir -p "$TEMP_DIR"
unzip -o "$ZIP_FILE" -d "$TEMP_DIR"
if [ $? -ne 0 ]; then
    echo "解压失败!"
    exit 1
fi

# 获取解压后的文件列表
EXTRACTED_FILES=("$TEMP_DIR"/*)
if [ ${#EXTRACTED_FILES[@]} -eq 0 ]; then
    echo "解压目录为空!"
    exit 1
fi

# 2. 获取目标POD列表
echo "获取目标POD..."
PODS=$(kubectl get pods -n "$POD_NAMESPACE" --no-headers | grep -E 'dc1-|dc2-|dc3-' | awk '{print $1}')
if [ -z "$PODS" ]; then
    echo "没有找到目标POD!"
    exit 1
fi

# 3. 复制文件到每个POD并重启服务
for POD in $PODS; do
    echo "处理POD: $POD"
    
    # 复制文件
    for FILE in "${EXTRACTED_FILES[@]}"; do
        echo "复制文件 $FILE 到 $POD:$TARGET_DIR"
        kubectl cp "$FILE" "$POD_NAMESPACE/$POD:$TARGET_DIR/" -c "$POD"
        if [ $? -ne 0 ]; then
            echo "复制文件 $FILE 到 $POD 失败!"
            continue
        fi
    done
    
    # 重启服务
    echo "在POD $POD 中重启服务 $SERVICE_NAME..."
    kubectl exec -n "$POD_NAMESPACE" "$POD" -- systemctl restart "$SERVICE_NAME"
    if [ $? -ne 0 ]; then
        echo "重启服务 $SERVICE_NAME 失败!"
    else
        echo "服务 $SERVICE_NAME 重启成功!"
    fi
done

# 4. 清理工作目录
echo "清理临时目录 $TEMP_DIR..."
rm -rf "$TEMP_DIR"

echo "所有操作完成!"


https://github.com/ktbyers/netmiko/blob/v3.0.0/netmiko/hp/hp_comware.py