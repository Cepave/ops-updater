# updater

updater只有一个功能，就是升级其他业务系统的agent。

## 功能描述

每隔指定周期（启动的时候通过命令行传入）调用一下ops-meta的http接口，上报当前管理的各个agent的版本号，response中顺便带回服务端配置的meta信息，也就是各个agent的版本号、状态、tarball地址等等，比如

```
[
    {
        "name": "falcon-agent",
        "version": "1.0.0",
        "tarball": "http://11.11.11.11:8888/falcon",
        "md5": "http://11.11.11.11:8888/falcon",
        "status": "start" 
    },
    {
        "name": "dinp-agent",
        "version": "2.0.0",
        "tarball": "http://11.11.11.11:8888/dinp",
        "md5": "",
        "status": "stop" 
    }
]
```

updater将这个信息与自我内存中的信息做对比，该升级的升级，该启动的启动，该停止的停止。

此处tarball和md5的配置并不是一个全路径（类似这样： http://11.11.11.11:8888/falcon/falcon-agent-1.0.0.tar.gz ），只是一部分路径，因为我们已经知道name和version了，就可以拼接出全路径了，算是一种规范化吧

## 目录结构

举个例子：

```
/home/work/ops
  |- ops-updater
  |- ops-updater.pid
  |- ops-updater.log
  |- falcon-agent
    |- 1.0.0
      |- falcon-agent
      |- control
    |- 1.0.1
      |- falcon-agent
      |- control
    |- .version
  |- dinp-agent
    |- 2.0.0
    |- 2.0.1
    |- .version
```

`.version`文件是当前这个agent应该使用的版本

## 启动参数

```
nohup ./ops-updater --interval=300 --debug=false --server=127.0.0.1:2200 &> ops-updater.log &
```

以上启动参数通常都可以不用加，采用默认值即可，当然，`--server`一般还是要加的，配置成ops-meta的地址

## 设计注意点

- updater访问ops-mata不能太集中，故而要加一个随机数，可以使用hostname+timestamp作为随机数种子
- 如果updater正在执行一些命令，结果被kill，下次重启的时候如何保证agent目录干净？tarball解压完成之后要删掉，这将作为解压是否成功的标志，如果发现某个agent目录里还有一个tarball，重新解压缩
- updater只升级agent就可以了，事情做得少才不容易出错
- updater可能依赖一些Linux工具才能正常工作，比如tar、md5sum等，在启动之前要做验证，如果没有，退出
