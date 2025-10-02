# CatBot

一个 QQ 机器人，调用 napcat 的 websocket 接口

## 功能

| 功能 | 触发方式 |
| --- | --- |
| 获取天气 | .weather |
| ping | .ping <addr> |
| ai聊天 | .chat <聊天信息> |
| 你猜 | .nc |
| 羡慕死了 | xmsl |
| 找吃的 | .findfood <地址> |
| 查公交 | .bus <公交站点名> |
| 防炸群 | 此功能需要设置为管理员并向 i@mei.lv 发邮件或者找 qq:2181331836 开启(开启后会自动撤回违规消息/图片并自动禁言) |

## `pkg/napcat`

这个 package 封装了一些 napcat 的操作，具体用法可以看代码中的实现，比如：

```go
func main() {
 client := napcat.New(
  "ws://127.0.0.1:3001/?access_token=^l^}BOdE[8s<k@g@",
  internal.SendGroupMsg,
  napcat.WithRetryDelay(5*time.Second),
 )
 client.Start(internal.SendGroupMsg)
}
```

连接函数`client.Start`会自动过滤掉心跳包并自动进行超时重连  
其中 `internal.SendGroupMsg` 是一个函数,当收到消息时会将 `[]byte` 类型的原始数据发送给该函数处理，同时传递一个 `websocket.Conn` 对象，该对象可以用于发送消息  

`Marshal` 函数只能用来发送群聊消息，其他消息类型还没有实现  
`Parse` 函数只能用于解析群聊消息  
