package main

import (
    "fmt"
    "log"
    "sync"
    "time"
)

// ============ 流式接口完整示例 ============

// StreamSession 代表一个流式连接会话
type StreamSession struct {
    ID       string
    MsgChan  chan string // 流式消息通道
    RespChan chan string // 响应通道
    closed   bool
    mu       sync.Mutex
}

func NewStreamSession(id string) *StreamSession {
    return &StreamSession{
        ID:       id,
        MsgChan:  make(chan string, 100), // 缓冲100条消息
        RespChan: make(chan string, 10),
    }
}

func (s *StreamSession) Close() {
    s.mu.Lock()
    defer s.mu.Unlock()
    if !s.closed {
        close(s.MsgChan)
        close(s.RespChan)
        s.closed = true
    }
}

func (s *StreamSession) SendMessage(msg string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.closed {
        return fmt.Errorf("session closed")
    }
    select {
    case s.MsgChan <- msg:
        return nil
    default:
        return fmt.Errorf("message buffer full")
    }
}

// StreamWorker 处理流式连接
type StreamWorker struct {
    sess      *StreamSession
    done      bool
    startTime time.Time
    timeout   time.Duration
}

func NewStreamWorker(sess *StreamSession, timeout time.Duration) *StreamWorker {
    return &StreamWorker{
        sess:      sess,
        startTime: time.Now(),
        timeout:   timeout,
    }
}

func (w *StreamWorker) handle(message string) error {
    // 模拟消息处理
    log.Printf("[Worker %s] 处理消息: %s", w.sess.ID, message)
    time.Sleep(50 * time.Millisecond) // 模拟处理耗时

    // 发送响应
    response := fmt.Sprintf("已处理: %s", message)
    select {
    case w.sess.RespChan <- response:
    default:
        log.Printf("[Worker %s] 响应通道已满，丢弃响应", w.sess.ID)
    }

    return nil
}

func (w *StreamWorker) checkTimeout() error {
    if time.Since(w.startTime) > w.timeout {
        return fmt.Errorf("worker timeout after %v", w.timeout)
    }
    return nil
}

func (w *StreamWorker) clean() {
    log.Printf("[Worker %s] 清理资源", w.sess.ID)
    w.sess.Close()
}

// Start 启动 Worker，持续处理流式消息
func (w *StreamWorker) Start() {
    log.Printf("[Worker %s] 启动", w.sess.ID)
    tick := time.NewTicker(time.Second)
    defer tick.Stop()

    var err error
    msgCount := 0

    for !w.done && err == nil {
        select {
        case msg, ok := <-w.sess.MsgChan:
            if !ok {
                log.Printf("[Worker %s] 消息通道已关闭，处理了 %d 条消息", w.sess.ID, msgCount)
                w.done = true
                break
            }
            msgCount++
            err = w.handle(msg)

        case <-tick.C:
            err = w.checkTimeout()
            if err != nil {
                log.Printf("[Worker %s] 超时检查失败: %v", w.sess.ID, err)
            }
        }
    }

    if err != nil {
        log.Printf("[Worker %s] 错误: %v", w.sess.ID, err)
    }
    w.clean()
}

// StreamServer 流式服务器
type StreamServer struct {
    limitSize    int
    workerSize   int
    workerBuffer chan *StreamWorker
    limitChan    chan struct{}
    wg           sync.WaitGroup
}

func NewStreamServer(limitSize, workerSize int) *StreamServer {
    s := &StreamServer{
        limitSize:    limitSize,
        workerSize:   workerSize,
        workerBuffer: make(chan *StreamWorker, workerSize),
        limitChan:    make(chan struct{}, limitSize),
    }
    s.startDispatcher()
    return s
}

// startDispatcher 启动 Worker 调度器
func (s *StreamServer) startDispatcher() {
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        log.Println("[Dispatcher] 启动")

        for worker := range s.workerBuffer {
            // 获取限流槽位（阻塞直到有可用槽位）
            s.limitChan <- struct{}{}
            log.Printf("[Dispatcher] 分配槽位给 Worker %s (当前活跃: %d/%d)",
                worker.sess.ID, len(s.limitChan), s.limitSize)

            s.wg.Add(1)
            go func(w *StreamWorker) {
                defer func() {
                    <-s.limitChan // 释放槽位
                    log.Printf("[Dispatcher] 释放槽位 (当前活跃: %d/%d)",
                        len(s.limitChan), s.limitSize)
                    s.wg.Done()
                }()
                w.Start()
            }(worker)
        }

        log.Println("[Dispatcher] 停止")
    }()
}

// AcceptConnection 接受新连接（提交 Worker）
func (s *StreamServer) AcceptConnection(sessionID string, timeout time.Duration) (*StreamSession, error) {
    session := NewStreamSession(sessionID)
    worker := NewStreamWorker(session, timeout)

    select {
    case s.workerBuffer <- worker:
        log.Printf("[Server] 接受连接 %s (队列中: %d/%d)",
            sessionID, len(s.workerBuffer), s.workerSize)
        return session, nil
    default:
        return nil, fmt.Errorf("server busy, queue full")
    }
}

// Shutdown 优雅关闭
func (s *StreamServer) Shutdown() {
    log.Println("[Server] 开始关闭...")
    close(s.workerBuffer)
    s.wg.Wait()
    log.Println("[Server] 关闭完成")
}

// ============ 示例：模拟客户端 ============

// simulateClient 模拟一个流式客户端
func simulateClient(server *StreamServer, clientID string, messageCount int) {
    // 1. 建立连接
    session, err := server.AcceptConnection(clientID, 10*time.Second)
    if err != nil {
        log.Printf("[Client %s] 连接失败: %v", clientID, err)
        return
    }

    log.Printf("[Client %s] 连接建立", clientID)

    // 2. 启动响应接收器
    go func() {
        for resp := range session.RespChan {
            log.Printf("[Client %s] 收到响应: %s", clientID, resp)
        }
    }()

    // 3. 发送流式消息
    for i := 0; i < messageCount; i++ {
        msg := fmt.Sprintf("消息-%d", i+1)
        err := session.SendMessage(msg)
        if err != nil {
            log.Printf("[Client %s] 发送失败: %v", clientID, err)
            break
        }
        log.Printf("[Client %s] 发送: %s", clientID, msg)
        time.Sleep(100 * time.Millisecond) // 模拟消息间隔
    }

    // 4. 关闭连接
    time.Sleep(500 * time.Millisecond) // 等待处理完成
    session.Close()
    log.Printf("[Client %s] 连接关闭", clientID)
}

// ============ 主函数 ============

func RunStreamExample() {
    log.Println("========== 流式接口示例 ==========")

    // 创建服务器：最多5个并发连接，队列容量20
    server := NewStreamServer(5, 20)

    // 模拟10个客户端，每个发送5条消息
    var wg sync.WaitGroup
    for i := 1; i <= 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            clientID := fmt.Sprintf("Client-%d", id)
            simulateClient(server, clientID, 5)
        }(i)
        time.Sleep(200 * time.Millisecond) // 错开连接时间
    }

    // 等待所有客户端完成
    wg.Wait()

    // 关闭服务器
    server.Shutdown()

    log.Println("========== 示例完成 ==========")
}

// ============ 对比示例：错误的单层设计 ============

// WrongDesign 错误的设计：只用一个消息通道
func WrongDesignExample() {
    log.Println("========== 错误设计示例 ==========")

    // ❌ 问题：所有连接共享一个消息通道
    // 无法区分消息来自哪个连接
    // 无法限制连接数
    globalMsgChan := make(chan string, 10000)

    // Worker池：固定5个Worker
    for i := 0; i < 5; i++ {
        go func(id int) {
            for msg := range globalMsgChan {
                log.Printf("[Worker %d] 处理消息: %s", id, msg)
                time.Sleep(50 * time.Millisecond)
            }
        }(i)
    }

    // 问题演示
    // 1. 无法关联消息和连接
    globalMsgChan <- "Client1-Msg1"
    globalMsgChan <- "Client2-Msg1"
    // Worker 无法知道这两条消息是否来自同一个客户端

    // 2. 无法保证同一连接的消息顺序
    // Client1 的消息可能被不同的 Worker 处理

    // 3. 无法独立管理每个连接的生命周期
    // 无法针对某个连接做超时控制

    time.Sleep(time.Second)
    close(globalMsgChan)

    log.Println("========== 错误设计示例结束 ==========")
}
