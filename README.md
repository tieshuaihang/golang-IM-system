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