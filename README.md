# sit-clock
久坐定时提示，可以内置MP3播放
使用translate.exe 可以将文件转为base64字符串
将转换的字符串添加到music.go文件中并修改main.go文件中audio相关的部分即可完成内置。
已经内置kun的三条语音

```
./translate cai1.mp3 
```

编译
```
go build -ldflags -H=windowsgui
```
<img width="386" alt="1672678158743" src="https://user-images.githubusercontent.com/53003567/210259535-8e43d02b-c9fc-44db-b6c1-b57ba4951df9.png">
