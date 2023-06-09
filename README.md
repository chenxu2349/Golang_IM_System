# Golang_IM_System
学习简易Go网络编程，实现一个简单即时通讯系统
win可以安装netcat包，配置好后可以使用nc 127.0.0.1 8888测试连接

## 一、构建基础server  
1、构建server类型  
2、创建一个server对象  
3、启动server服务  
4、处理链接业务  
![img.png](images/img1.png)


## 二、用户上线功能
1、创建user类型  
2、创建user对象，并创建方法监听user对应的channel消息  
3、改写server：新增OnlineMap和Message属性，在处理客户端上线的Handler中创建并
添加用户，新增广播消息方法，新增监听广播消息channel的方法，然后启用协程去监听server中
的Message
![img.png](images/img2.png)

## 三、将用户消息进行广播
![img.png](images/img3.png)

## 四、用户业务封装
![img.png](images/img4.png)

## 五、在线用户查询
![img.png](images/img5.png)

## 六、修改用户名
![img.png](images/img6.png)

## 七、超时强踢功能
![img.png](images/img7.png)

## 八、私聊功能
![img.png](images/img8.png)

## 九、客户端实现
新建一个目录用于存放client.go，因为里面有另外一个main函数与主包冲突
先go build出可执行文件client.exe，然后./client -h(-ip, -port)使用
![img.png](images/img9_1.png)
<br>
<br>
![img.png](images/img9_2.png)
<br>
<br>
菜单显示：新增flag属性，新增menu()方法获取用户输入，新增Run()方法主业务循环，最后在
main函数中调用Run()方法

<br>
<br>
客户端更改用户名功能：新增UpdateName()方法更新用户名，然后加入到Run业务分支中，
添加用于处理服务器回执消息的方法DealRsponse()，开启一个go协程去处理消息

<br>
<br>
客户端公聊功能：新增公聊方法，加入到Run分支业务中

<br>
<br>
客户端私聊功能：1、查询当前在线用户
2、选择一个用户进入私聊
3、私聊方法