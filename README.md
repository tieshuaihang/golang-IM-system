# golang-IM-system

### v1 搭建基础server

### v2 实现用户上线通知

- server 对象保存所有在线用户map,需要加读写锁，遍历用户及增删用户时使用
- 用户上线时将上线消息发给 server 对象的 Message channel
- server 启动时，启动协程监听 Message channel，如果有消息，遍历发给所有用户对象的 C channel
- user 对象有 C channel, 用于接受广播信息，创建user时，启动协程监听该channel,有数据时，将其发给客户端

### v3 用户消息广播
- 启动一个协程获取用户发来的消息，注意如何使用Conn读取数据，结尾处理以及断开时处理
- 将用户消息通过 server 的 Message channel 进行广播

### v4 查询在线用户  修改用户名

### v5 超时强踢
- handler 中定义 isActive channel，每次收到用户消息，就证明用户在线，往 channel中发数据
- 使用 select 评估 定时器 和  isActive channel，每次 isActive channel 命中就会重置定时器，定时器超时，关闭连接
- 手动关闭连接时，读取用户数据的协程会收到长度为0的数据，使用户下线

### v6 客户端实现