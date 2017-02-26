#大体功能
Web 上的私信系统,能够实时收发消息，主要分成两部分：HTTP的restful API提供用户信息管理和历史消息管理，websocket部分提供消息的实时收发，用户信息与消息保存在mysql数据库，支持配置文件配置参数，
#后续改进点
目前是单服务器单进程结构，后期可以把http部分和websocket部分拆开，http部分单独部署提供服务，websocket部分可以作为IM系统的连接层，后面增加mq与转发逻辑来增加系统稳定性与并发量，还需要增加cache来缓存用户基本信息以及token。

#接口协议
##HTTP
###注册
curl http://localhost:8080/user -H "Action:register" -d "id=god4&pwd=123"

{"rc":0,"msg":"success"}
###登录
curl http://localhost:8080/user -H "Action:login" -d "id=god&pwd=123"

{"rc":0,"msg":"success","data":{"userid":6,"token":"T5GjdEEA7JcvcVp63uXScA=="}}
###获取好友
curl http://localhost:8080/user -H "Action:friendsList" -H "Token:sqUt7cZVy9QupUYS0uwFKg==" -d "userid=6"

{"rc":0,"msg":"success","data":[{"userid":3,"account":"lw"}]}
###添加好友
curl http://localhost:8080/user -H "Action:addFriend" -H "Token:sqUt7cZVy9QupUYS0uwFKg==" -d "userid=6&peerid=8"

{"rc":0,"msg":"success"}
###删除好友
curl http://localhost:8080/user -H "Action:delFriend" -H "Token:sqUt7cZVy9QupUYS0uwFKg==" -d "userid=6&peerid=8"

{"rc":0,"msg":"success"}

###todo
获取历史消息
删除消息
保存消息

##websocket
###发送消息

ws://localhost:8080/msg

{"fromid":6,"toid":3,"token":"sqUt7cZVy9QupUYS0uwFKg==","content":"hello"}

**注意：** token使用的是login返回的token

###收到消息
{"fromid":6,"toid":3,"content":"hello"}