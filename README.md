# GoLang WebSocket with binary communication with JavaScript
## Installation
make install
make run
## Communication protocol
Full packet
```
[                     Header                       ] [         body         ]
["om"][version:uint8][action:uint8][bodysize:uint32] [....:map[string]string]
```
Body packet - Key value pairs separated by [0xc0,0x80]
```
[key:string][0xc0,0x80][value:string][0xc0,0x80]
```

Look in ***ws.html*** file for JavaScript for ***converting from binary to js object.***

### Dependinces
 - github.com/gorilla/websocket
 - gopkg.in/ini.v1