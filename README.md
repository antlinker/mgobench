# mgobench

> mongodb测试工具

## 获取

``` bash
$ go get github.com/antlinker/mgobench
```

## 使用说明

``` bash
$ mgobench -h
```

```
Usage of ./mgobench:
  -alsologtostderr
       	log to standard error as well as files
  -c string
       	集合名 (default "test1")
  -db string
       	数据库 (default "test")
  -doc_num int
       	写入文档数量 (default 100)
  -doc_size int
       	文档大小(单位B) (default 64)
  -go_num int
       	goroutine数量 (default 10)
  -log_backtrace_at value
       	when logging hits line file:N, emit a stack trace
  -log_dir string
       	If non-empty, write log files in this directory
  -logtostderr
       	log to standard error instead of files
  -mgo_url string
       	mongodb链接串 (default "mongodb://127.0.0.1:27017")
  -sleep int
       	写入文档的休眠间隔(单位毫秒，0则不休眠)
  -stderrthreshold value
       	logs at or above this threshold go to stderr
  -v value
       	log level for V logs
  -vmodule value
       	comma-separated list of pattern=N settings for file-filtered logging
```

## MIT License

```
Copyright (c) 2016 Lyric
```