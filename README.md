### 下载代码
```shell
go get https://gitlab.com/kunzhangs/go-wiki-spider
```

### 新建数据库与修改数据库信息
-  进入engine/ simple.go根据Content结构的定义创建数据库，并在init函数中修正数据库的链接配置

### 编译软件
```shell
cd $GOPATH/src/go-wiki-spider/wiki
go install
```

### 编写种子配置文件
- 按照seed.conf格式填写几个种子，并存放在前面编译软件相同的目录 $GOPATH/bin

### 执行软件
```shell
cd  $GOPATH/bin
./wiki
```
