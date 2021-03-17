# Go基础课 -区块链

本次练习的主题是"挖矿"，由你来编写"挖矿client"，从老师编写的"矿池server"下载计算任务，按照要求编写你的"挖矿"逻辑，并将计算得到的"区块"提交给"矿池server"，验证成功的话你就通过了本次"考试"。

老师作为"矿池主"，对各位"矿工们"的劳动贡献一清二楚，所以大家不需要截图给老师，只要自己挖到矿就自然会通过"考试"了。


## 前置知识

简单了解区块链挖矿的原理：[《阮一峰-区块链入门教程》](https://www.ruanyifeng.com/blog/2017/12/blockchain-tutorial.html)

## 开发工作

你需要阅读整个代码结构，然后解决3个作业，最终才能挖到矿：
* 作业A：补全Json2Block函数（注意struct的tag是否正确）
* 作业B：补全Block2Json函数
* 作业C：补全MiningBlock函数（核心挖矿逻辑）

## 启动方法

```
cd course6/client
go run main.go -address 服务器地址 -content 你的名字
```

* address是"矿池server"的地址，老师会提供给大家
* content是写入到"区块"上的内容，你至少应该写上你的名字，因为老师是根据这个来统计各位的作业的
