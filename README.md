# tel.0sn.net-ssh

```bash
$ ssh localhost -p 8022

Press 'q' to exit.
--------------------------
tel.0sn.netへようこそ！
_       _   ___                         _   
| |_ ___| | / _ \ ___ _ __    _ __   ___| |_ 
| __/ _ \ || | | / __| '_ \  | '_ \ / _ \ __|
| ||  __/ || |_| \__ \ | | |_| | | |  __/ |_ 
 \__\___|_(_)___/|___/_| |_(_)_| |_|\___|\__|
--------------------------
あなたは 1 人目の訪問者です。
--------------------------
Web: https://0sn.net
--------------------------
Connection to localhost closed.
```

## 使い方
```bash
$ go build -o main

#必ず最初にホストキーの生成をする
$ ssh-keygen -f host.key

#デフォルトでは8022で待ち受け
$ ./main
$ PORT=8000 ./main
```

### Docker
```bash
$ docker compose build

#必ず最初にホストキーの生成をする
$ ssh-keygen -f host.key

$ docker compose up -d
```