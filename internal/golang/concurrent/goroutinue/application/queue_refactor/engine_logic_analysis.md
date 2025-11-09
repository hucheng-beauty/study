# Engine 并发任务引擎逻辑分析

## 📋 整体架构

这是一个基于**工作池模式（Worker Pool Pattern）**的并发任务调度引擎，支持优雅关闭和任务排空。

---

## 🏗️ 核心组件

### 1. **数据结构**

```
Engine {
    wg         sync.WaitGroup           // 等待所有worker完成
    once       sync.Once               // 确保Stop只执行一次
    opts       *Options                // 配置选项
    workerPool chan chan interface{}   // worker池（容量=WorkersCount）
    jobChannel []chan interface{}      // 每个worker的专属任务通道
    jobQueue   chan interface{}        // 全局任务队列（容量=QueueLength）
    done       chan struct{}           // 关闭信号
}
```

### 2. **配置选项**

- `Ctx`: 上下文控制，用于取消操作
- `Func`: 任务处理函数
- `WorkersCount`: worker数量（默认10个）
- `QueueLength`: 任务队列长度（默认10000）

---

## 🔄 工作流程

### **阶段1: 初始化（NewEngine）**

```
1. 参数校验和默认值设置
2. 创建 workerPool（缓冲通道，大小=worker数量）
3. 创建 jobChannel 数组（每个worker一个专属通道）
4. 创建 jobQueue（全局任务队列）
5. 创建 done 信号通道
```

### **阶段2: 启动调度（Schedule）**

```
启动 N 个 worker goroutine
    └─> 调用 worker(i) 方法

启动 1 个调度器 goroutine
    └─> 调用 schedule() 方法
```

---

## 👷 Worker 逻辑详解

```go
func (e *Engine) worker(i int) {
for {
// 步骤1: 注册到工作池（告诉调度器：我空闲了）
case e.workerPool <- e.jobChannel[i]:

// 步骤2: 等待任务到达
select {
case job := <-e.jobChannel[i]:
e.opts.Func(job) // 执行任务
case <-e.opts.Ctx.Done():
return // 上下文取消，退出
}

// 随时监听上下文取消
case <-e.opts.Ctx.Done():
return
}
}
```

**工作流程：**

1. Worker将自己的 `jobChannel[i]` 发送到 `workerPool`（表示：我空闲）
2. 阻塞等待调度器通过 `jobChannel[i]` 分配任务
3. 收到任务后执行 `Func(job)`
4. 完成后回到步骤1（循环）

---

## 📮 调度器逻辑详解

```go
func (e *Engine) schedule() {
for {
select {
// 场景1: 从任务队列获取任务
case job := <-e.jobQueue:
// 从worker池获取一个空闲worker
case jobChannel := <-e.workerPool:
select {
case jobChannel <- job: // 分配任务给worker
case <-e.opts.Ctx.Done():
e.opts.Func(job) // 上下文取消时直接处理当前任务
return
}
}

// 场景2: 收到停止信号
case <-e.done:
e.drainQueue() // 排空剩余任务
return
}
}
}
```

**工作流程：**

1. 从 `jobQueue` 获取待处理任务
2. 从 `workerPool` 获取空闲worker的通道
3. 将任务发送到该worker的专属通道
4. 如果收到 `done` 信号，进入排空模式

---

## 🚰 任务排空逻辑（drainQueue）

```go
func (e *Engine) drainQueue() {
for {
select {
case job, ok := <-e.jobQueue:
if !ok {
return // 队列已关闭
}
e.opts.Func(job) // 直接处理，不分配给worker
drained++

default:
return // 队列为空，退出
}
}
}
```

**关键点：**

- **不再分配给worker**：因为worker可能正在退出
- **直接执行任务**：在调度器goroutine中同步处理
- **快速排空**：使用 `default` 分支立即检测空队列

---

## 🛑 优雅关闭流程

### **用户调用顺序：**

```go
engine.Stop() // 1. 停止接收新任务
cancel()      // 2. 取消上下文（通知worker退出）
engine.Wait() // 3. 等待所有worker退出
```

### **内部执行流程：**

```
Stop() 被调用
  └─> close(done)  // 关闭done通道
       └─> schedule() 收到 <-done 信号
            └─> 调用 drainQueue()
                 └─> 处理 jobQueue 中所有剩余任务
                      └─> schedule() 退出

同时，上下文被取消
  └─> 所有 worker 收到 <-Ctx.Done() 信号
       └─> worker goroutine 退出
            └─> wg.Done() 被调用（defer中）

Wait() 阻塞等待
  └─> 直到所有 worker 的 wg.Done() 都被调用
       └─> Wait() 返回，关闭完成
```

---

## 📊 时序图

```
用户        Engine        Scheduler       Worker-1    Worker-2
 |            |              |              |           |
 |--Submit--->|              |              |           |
 |            |--job-------->|              |           |
 |            |              |              |           |
 |            |              |<--register---|           |
 |            |              |              |           |
 |            |              |---job------->|           |
 |            |              |              |--Func()   |
 |            |              |              |           |
 |            |              |<----------register-------|
 |            |              |              |           |
 |--Submit--->|              |              |           |
 |            |--job-------->|              |           |
 |            |              |---job----------------->|
 |            |              |              |           |--Func()
 |--Stop()--->|              |              |           |
 |            |--close(done)->|              |           |
 |            |              |--drainQueue()->          |
 |--cancel()-->--Ctx.Done()---------------->|---------->|
 |            |              |              |exit       |exit
 |--Wait()--->|              |              |           |
 |            |<--wg.Wait()------------------           |
 |<--return---|              |              |           |
```

---

## ✅ 优点

1. ✨ **Worker池复用**：避免频繁创建goroutine
2. 🛡️ **Panic恢复**：worker和调度器都有恢复机制
3. 🚪 **优雅关闭**：支持排空任务队列
4. 🔒 **并发安全**：使用通道通信，避免共享内存

## 🎯 使用示例

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    // 创建引擎
    engine := NewEngine(&Options{
        Ctx: ctx,
        Func: func(job interface{}) {
            fmt.Printf("Processing: %v\n", job)
            time.Sleep(100 * time.Millisecond)
        },
        WorkersCount: 5,
        QueueLength:  100,
    })

    // 启动调度
    engine.Schedule()

    // 提交任务
    for i := 0; i < 20; i++ {
        engine.Submit(fmt.Sprintf("task-%d", i))
    }

    // 优雅关闭
    time.Sleep(2 * time.Second)
    engine.Stop() // 停止接收新任务
    cancel()      // 通知worker退出
    engine.Wait() // 等待完成

    fmt.Println("All done!")
}
```

