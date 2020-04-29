# SockerMQTTWatcher
## SockerMQTTWatcher是什么  
  SockerMQTTWatcher是Socker的MQTT模块，支持双向订阅，支持在线/离线认证，数据上传和指令下发
## 使用
  编译项目，将编译好的mqttWatcher文件放在Socker目录下，运行Socker时会调用mqttWatcher  
### 运行  
```
  sudo ./mqttWatcher start ContainerNAME
```
### 停止  
```
  sudo ./mqttWatcher stop ContainerNAME
```
