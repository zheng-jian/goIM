### A Instance Message web application with go websocket

##### Architecture
![image](https://raw.githubusercontent.com/zheng-jian/goIM/master/architecture.png)

##### run this application
- git clone this project
- cd goIM
- go run main.go

##### explanation
- I make two servers in the same area but listen two different ports
- IM server provides front pages for users to connect the websocket server
- Websocket server deals with the messages, when a user sends one message, 
all the other users will receive the same message.