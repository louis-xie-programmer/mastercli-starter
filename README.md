# mastercli-starter

> 🚀 一个基于 **Golang** 的 **Master-Worker 框架 + CLI 应用** 脚手架，适用于构建并发任务调度系统。

该项目演示了如何使用 **Cobra** + **Viper** 搭建现代命令行程序，并实现一个可扩展的 **并发任务处理框架**。  
内置：
- **Worker 池**（可配置大小）
- **带缓冲任务队列**
- **上下文驱动的优雅退出**
- **指数退避重试机制**
- **配置文件 + 环境变量加载**
- **结构化日志**
- **多命令 CLI**

---

结合公众号博文，可快速了解本系统的核心功能和实现原理, 后续会持续更新，欢迎大家扫码关注“代码扳手”公众号，获取更多技术交流信息。

![img.png](img.png)


## 🏗 架构概览

```
+-------------------+       +-------------------+
|   Command Line    |       |   Config Loader   |
|  (Cobra Commands) |       |    (Viper)        |
+--------+----------+       +---------+---------+
         |                            |
         v                            v
+---------------------------------------------+
|                Master (Manager)            |
|  - Job Queue (chan job.Job)                 |
|  - Results Channel                          |
|  - Worker Pool Management                   |
+--------------------+------------------------+
                     |
        +------------+-------------+
        |                          |
+-------v-------+          +-------v-------+
|   Worker #1   |          |   Worker #N   |
|  Executes Job |  ...     |  Executes Job |
+---------------+          +---------------+
```

---

## ⚙️ 安装与构建

```bash
git clone https://github.com/louis-xie-programmer/mastercli-starter.git
cd mastercli-starter

# 安装依赖
go mod download

# 构建可执行文件
make build
```

编译完成后生成：
```
./mastercli
```

---

## 🚀 快速上手

### 查看命令
```bash
./mastercli --help
```

输出示例：
```
Master framework + CLI in Go

Usage:
  mastercli [command]

Available Commands:
  help        Help about any command
  run         Run a single job synchronously and print the result
  start       Start the master and process jobs
```

---

### 启动 Master（默认 4 个 Worker，20 个示例任务）
```bash
./mastercli start
```
示例输出：
```
{"level":"info","time":1734150515,"app":"mastercli","message":"starting"}
{"level":"info","job_id":"job-001","message":"completed"}
{"level":"warn","job_id":"job-005","error":"transient failure for job job-005","message":"failed"}
{"level":"info","message":"shutdown signal received"}
```

---

### 从文件加载任务
```bash
echo -e "task1\ntask2\ntask3" > tasks.txt
./mastercli start -f tasks.txt
```

此时每一行会被当作一个 Job 的 `payload`。

---

### 单任务执行（同步）
```bash
./mastercli run -p "hello" -d 500
```
参数说明：
- `-p` / `--payload`：任务数据
- `-d` / `--duration`：执行耗时（毫秒）
- `-F` / `--fail-once`：模拟一次性失败（触发重试）

---

## 🛠 配置系统

默认读取 `configs/config.yaml`：
```yaml
app:
  name: mastercli
  log_level: info   # debug, info, warn, error
master:
  workers: 4
  queue_size: 64
  max_retries: 2
  backoff_ms: 250
```

支持环境变量覆盖（前缀 `MASTERCLI_`）：
```bash
MASTERCLI_APP_LOG_LEVEL=debug ./mastercli start
```

---

## 🔄 Job 生命周期

1. **生成 / 读取**：`start` 命令会生成 demo 任务或从文件读取。
2. **入队**：调用 `Manager.Submit(job)` 放入 `jobs` channel。
3. **Worker 消费**：Worker 从队列取出任务，执行 `Do(ctx)` 方法。
4. **结果回传**：执行结果发送到 `results` channel。
5. **结果处理**：
   - 成功：记录日志
   - 失败：触发重试逻辑（指数退避）
6. **任务结束 / 超时 / 中断**：由 `context.Context` 控制。

---

## ⏱ 重试机制说明

- 每个失败任务会根据 `max_retries` 配置进行重试
- 每次重试前等待 `backoff_ms * 2^(attempt-1)` 毫秒
- 重试时会记录当前尝试次数和退避时间
- 演示实现中，并未持久化 Job，仅做日志输出

---
