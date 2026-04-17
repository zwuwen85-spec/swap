package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message WebSocket消息
type Message struct {
	Type      string `json:"type"`       // message/status/pong
	UserID    int64  `json:"user_id,omitempty"`    // 状态消息的用户ID
	SenderID  int64  `json:"sender_id"`
	ReceiverID int64 `json:"receiver_id"`
	Content   string `json:"content"`
	GoodsID   int64  `json:"goods_id,omitempty"`
	Online    bool   `json:"online"`                 // 在线状态
	Timestamp int64  `json:"timestamp"`
}

// Client WebSocket客户端
type Client struct {
	ID     int64                  // 用户ID
	Send   chan *Message          // 发送通道
	Hub    *Hub                   // Hub引用
	Conn   *websocket.Conn        // WebSocket连接
}

// Hub WebSocket连接中心
type Hub struct {
	Clients    map[int64]*Client    // 在线用户
	Register   chan *Client          // 注册通道
	Unregister chan *Client          // 注销通道
	Broadcast  chan *Message         // 广播通道
	mu         sync.RWMutex          // 读写锁
}

// NewHub 创建Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int64]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 256),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.register(client)

		case client := <-h.Unregister:
			h.unregister(client)

		case message := <-h.Broadcast:
			h.broadcast(message)
		}
	}
}

// register 注册客户端
func (h *Hub) register(client *Client) {
	h.mu.Lock()
	h.Clients[client.ID] = client
	log.Printf("[WebSocket] 用户 %d 上线，当前在线人数: %d", client.ID, len(h.Clients))
	h.mu.Unlock()

	// 广播上线消息（在解锁后进行，避免死锁）
	h.broadcast(&Message{
		Type:      "status",
		UserID:    client.ID,
		Online:    true,
		Timestamp: getCurrentTimestamp(),
	})
}

// unregister 注销客户端
func (h *Hub) unregister(client *Client) {
	h.mu.Lock()
	if _, ok := h.Clients[client.ID]; ok {
		delete(h.Clients, client.ID)
		close(client.Send)
		log.Printf("[WebSocket] 用户 %d 下线，当前在线人数: %d", client.ID, len(h.Clients))
		h.mu.Unlock()

		// 广播下线消息（在解锁后进行，避免死锁）
		h.broadcast(&Message{
			Type:      "status",
			UserID:    client.ID,
			Online:    false,
			Timestamp: getCurrentTimestamp(),
		})
	} else {
		h.mu.Unlock()
	}
}

// broadcast 广播消息
func (h *Hub) broadcast(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.Clients {
		// 如果是消息类型，只发给发送者和接收者
		if message.Type == "message" {
			if client.ID != message.SenderID && client.ID != message.ReceiverID {
				continue
			}
		}

		select {
		case client.Send <- message:
		default:
			// 发送通道满，跳过
			log.Printf("[WebSocket] 用户 %d 发送通道满", client.ID)
		}
	}
}

// GetClient 获取客户端
func (h *Hub) GetClient(userID int64) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	client, ok := h.Clients[userID]
	return client, ok
}

// GetOnlineUsers 获取在线用户
func (h *Hub) GetOnlineUsers() []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]int64, 0, len(h.Clients))
	for userID := range h.Clients {
		users = append(users, userID)
	}
	return users
}

// IsOnline 检查用户是否在线
func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.Clients[userID]
	return ok
}

// SendMessageToUser 发送消息给指定用户
func (h *Hub) SendMessageToUser(userID int64, message *Message) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.Clients[userID]; ok {
		select {
		case client.Send <- message:
			return true
		default:
			return false
		}
	}
	return false
}

// SendPong 发送心跳响应
func (c *Client) SendPong() {
	c.Send <- &Message{
		Type:      "pong",
		Timestamp: getCurrentTimestamp(),
	}
}

// ReadPump 读取泵（从WebSocket读取数据）
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Time{})
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket] 用户 %d 连接异常关闭: %v", c.ID, err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("[WebSocket] 用户 %d 消息解析失败: %v", c.ID, err)
			continue
		}

		// 处理心跳
		if msg.Type == "ping" {
			c.SendPong()
			continue
		}

		// 设置发送者ID
		msg.SenderID = c.ID
		msg.Timestamp = getCurrentTimestamp()

		// 广播消息
		c.Hub.Broadcast <- &msg
	}
}

// WritePump 写入泵（向WebSocket写入数据）
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			// 发送心跳
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"ping"}`)); err != nil {
				return
			}

		case message, ok := <-c.Send:
			if !ok {
				return
			}

			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("[WebSocket] 消息序列化失败: %v", err)
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("[WebSocket] 写入消息失败: %v", err)
				return
			}
		}
	}
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// Upgrader WebSocket升级器
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 生产环境应该检查Origin
		return true
	},
}

// WebSocketManager WebSocket管理器
type WebSocketManager struct {
	Hub *Hub
}

// NewWebSocketManager 创建WebSocket管理器
func NewWebSocketManager() *WebSocketManager {
	hub := NewHub()
	go hub.Run()

	return &WebSocketManager{
		Hub: hub,
	}
}

// HandleWebSocket 处理WebSocket连接
func (m *WebSocketManager) HandleWebSocket(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(200, gin.H{"code": 10002, "message": "未登录"})
		return
	}

	// 升级连接
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket] 升级连接失败: %v", err)
		return
	}

	// 创建客户端
	client := &Client{
		ID:   userID.(int64),
		Send: make(chan *Message, 256),
		Hub:  m.Hub,
		Conn: conn,
	}

	// 注册客户端
	m.Hub.Register <- client

	// 启动读写泵
	go client.ReadPump()
	go client.WritePump()
}

// BroadcastMessage 广播消息
func (m *WebSocketManager) BroadcastMessage(message *Message) {
	m.Hub.Broadcast <- message
}

// SendToUser 发送消息给指定用户
func (m *WebSocketManager) SendToUser(userID int64, message *Message) bool {
	return m.Hub.SendMessageToUser(userID, message)
}

// IsOnline 检查用户是否在线
func (m *WebSocketManager) IsOnline(userID int64) bool {
	return m.Hub.IsOnline(userID)
}

// GetOnlineUsers 获取在线用户列表
func (m *WebSocketManager) GetOnlineUsers() []int64 {
	return m.Hub.GetOnlineUsers()
}
