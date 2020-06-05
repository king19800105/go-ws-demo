# go-ws-demo

### 目录结构
- cmd ：启动命令
- configs ：配置文件
- pkg ：核心包
    - config ：配置对象
    - hardware ：核心业务处理
    - server ：websocket服务

### 测试（启动服务后）
- 方式一：使用 cmd/instrument/main_test.go 进行测试
- 方式二：打开go-ws-demo.html文件，（1）点"链接"，（2）输入框中输入：{"event":"login", "timestamp": 1591254792, "data":"{}"}