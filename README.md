go_push_server
=========

go push demo

  特性：1)  自动负载均衡
        2)  消息可离线
        3)  自定义不同产品线
        4)  优雅关闭（会把所有消息发送完再关闭）。
        
server端：
  go build push/push.go
  
client端：
  go build client/client.go
