# Opay  [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/opay)

Opay 以订单为主线、面向接口开发的在线支付通用模块。

# 特点

- 完全面向接口开发

- 支持超时自动撤销处理订单

# 使用步骤

1. 注册资产账户操作接口实例

2. 实现订单接口

3. 注册订单类型对应的操作接口实例

4. 新建服务实例 var opay=NewOpay(db, 5000)

5. 开启服务协程 go opay.Serve()

6. 推送订单 err:=opay.Push(iOrd)

7. 使用 <-iOrd.Done() 等待订单处理结束
