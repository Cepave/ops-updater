# updater

updater只有一个功能，就是升级其他业务系统的agent。

## 功能描述

每隔指定周期（配置文件中的interval，单位是秒）调用一下ops-meta的http接口，上报当前管理的各个agent的状态和版本号，比如：

```
[
    {
        "name": "falcon-agent",
        "version": "1.0.0",
        "status": "started" 
    },
    {
        "name": "dinp-agent",
        "version": "2.0.0",
        "status": "stoped" 
    }
]
```

response中顺便带回服务端配置的meta信息，也就是各个agent的版本号、tarball地址、要求的操作等等，比如

```
[
    {
        "name": "falcon-agent",
        "version": "1.0.0",
        "tarball": "http://11.11.11.11:8888/falcon",
        "md5": "http://11.11.11.11:8888/falcon",
        "cmd": "start" 
    },
    {
        "name": "dinp-agent",
        "version": "2.0.0",
        "tarball": "http://11.11.11.11:8888/dinp",
        "md5": "",
        "cmd": "stop" 
    }
]
```

updater将这个信息与自我内存中的信息做对比，该升级的升级，该启动的启动，该停止的停止。sorry，不能与内存中的信息做对比，应该直接去各个agent目录调用control脚本查看，因为agent有可能crash，或者做了一些手工操作，这将导致updater内存中的信息并不能实时反应线上情况

此处tarball和md5的配置并不是一个全路径（类似这样： http://11.11.11.11:8888/falcon/falcon-agent-1.0.0.tar.gz ），只是一部分路径，因为我们已经知道name和version了，就可以拼接出全路径了，算是一种规范化吧，拼接规范是：`{$tarball}/{$name}-{$version}.tar.gz`

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

## 启停

```
./control start|stop|restart|tail
```

配置文件是cfg.json，我们提供了cfg.example.json作为配置模板，`mv cfg.example.json cfg.json`，然后修改成合适的配置

## 设计注意点

- updater访问ops-mata不能太集中，故而要加一个随机数，启动之后先随机sleep一下，然后再发送心跳请求
- updater只升级agent就可以了，事情做得少才不容易出错
- updater可能依赖一些Linux工具才能正常工作，比如tar、md5sum等，在启动之前要做验证，如果没有，退出

## agent约定

- 必须提供control脚本
- `./control start`可以启动agent
- `./control stop`可以停止agent
- `./control status`打印出状态信息，只能是started或者stoped
- control文件已经有可执行权限，并且在tarball根目录下


