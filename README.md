# updater

updater只有一个功能，就是升级其他业务系统的agent。

## 功能描述

每隔指定周期（启动的时候通过命令行传入）调用一下ops-meta的http接口，上报当前管理的各个agent的版本号，比如：

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

response中顺便带回服务端配置的meta信息，也就是各个agent的版本号、状态、tarball地址等等，比如

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

## agent约定

- 必须提供control脚本
- `./control start`可以启动agent
- `./control stop`可以停止agent
- `./control status`打印出状态信息，只能是started或者stoped，不过实现的时候，updater去获取的时候要用strings.contains

## updater处理流程

```
for agent in http.response:
    handle(agent)

def handle(agent):
    insure_desired_version(agent)
    real_version = cat {agent.name}/.version
    if agent.version == real_version:
        version_equal(agent)
    else:
        version_not_equal(agent)

def version_equal(agent):
    handle agent.cmd
    已经启动了就无需再启动
    已经stop了就无需再stop

def insure_desired_version(agent):
    dir = "{agent.name}/{agent.version}"
    if not dir.exists:
        mkdir
        deploy_files(agent)
        return
    if tarball.exists:
        untar

def deploy_files(agent):
    download(agent.tarball)
    download(agent.md5)
    md5sum -c
    untar

def untar(tarball, dir):
    tar zxvf tarball
    rm -f tarball

```

- 1. 遍历response中的各个agent信息
- 2. 对每一个agent，看现在线上运行的agent版本与要求的是否一致
- 3. 一致：要求stop就去`./control stop`，要求start就去`./control start`这些shell指令都要加超时


