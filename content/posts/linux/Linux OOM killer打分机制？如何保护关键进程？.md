Linux 的 **OOM Killer（Out-of-Memory Killer）** 是内核在系统内存耗尽时主动终止进程以释放内存的机制。其核心是通过评分机制选择最“适合”终止的进程。以下是详细解析和关键进程保护方法：

---

### **1. OOM Killer 打分机制**
当系统内存不足时，内核会为每个进程计算一个 **OOM Score**（`oom_score` 值，范围 0~1000），分数越高越容易被杀死。计算依据包括：
- **内存占用比例**：进程消耗的内存占系统总内存的比例越高，分数越高。
- **进程优先级（Nice值）**：优先级越低（Nice值越高）的进程分数越高。
- **运行时间**：长时间运行的进程可能分数更低（倾向于保留）。
- **子进程内存**：父进程的分数会累加子进程的内存占用。
- **用户权限**：root 用户的进程可能获得轻微保护。
- **特殊标记**：如 `oom_score_adj` 或 `oom_adj`（旧版本）可手动调整分数。

#### **查看进程 OOM 分数**
```bash
cat /proc/<PID>/oom_score
```

---

### **2. 保护关键进程的方法**

#### **(1) 调整 `oom_score_adj`**
通过修改 `/proc/<PID>/oom_score_adj`（范围：-1000~1000），直接影响 `oom_score`：
- **降低分数**（减少被杀概率）：
  ```bash
  echo -500 > /proc/<PID>/oom_score_adj  # 关键进程设为负值
  ```
- **提高其他进程分数**（间接保护）：
  ```bash
  echo 1000 > /proc/<非关键PID>/oom_score_adj  # 让非关键进程更易被杀死
  ```

#### **(2) 使用 `systemd` 服务配置**
若进程由 `systemd` 管理，在单元文件中添加：
```ini
[Service]
OOMScoreAdjust=-500  # 等效于修改 oom_score_adj
```

#### **(3) 禁用 OOM Killer 对特定进程**
通过 `cgroup` 完全排除进程：
```bash
echo -17 > /proc/<PID>/oom_adj  # 旧内核版本（需确认支持）
```
> **注意**：新内核可能已移除此功能，建议优先用 `oom_score_adj`。

#### **(4) 限制非关键进程内存**
通过 `cgroups` 或 `ulimit` 限制非关键进程的内存使用，避免它们触发 OOM：
```bash
ulimit -v 500000  # 限制进程虚拟内存为 500MB
```

#### **(5) 内核参数调整**
- **`vm.overcommit_memory`**：
  - `0`（默认）：启发式内存分配，可能触发 OOM Killer。
  - `1`：总是允许超分配（风险：可能崩溃）。
  - `2`：严格拒绝超分配（需配合 `vm.overcommit_ratio`）。
  
  修改方式：
  ```bash
  sysctl vm.overcommit_memory=1
  ```

- **`vm.panic_on_oom`**：
  - `0`：启用 OOM Killer（默认）。
  - `1`：直接触发内核 panic（极端情况使用）。

---

### **3. 监控与调试**
- **查看 OOM 日志**：
  ```bash
  dmesg | grep -i "oom"
  journalctl --dmesg | grep -i "oom"
  ```
- **模拟 OOM 测试**：
  ```bash
  stress-ng --vm 1 --vm-bytes 90% --timeout 10s  # 谨慎使用！
  ```

---

### **最佳实践**
1. **优先保护**：数据库（如 MySQL）、关键服务（如 SSH）通过 `oom_score_adj` 设为负值。
2. **限制非关键进程**：避免单个进程耗尽内存。
3. **监控告警**：通过 Prometheus+Alertmanager 监控 `oom_kill` 事件。

通过合理配置，可以显著降低关键进程被误杀的风险。