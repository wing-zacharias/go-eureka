# go-eureka说明

## 结构说明

结构图
![avatar](https://github.com/wing-zacharias/go-eureka/blob/master/1.png)

```text
1.config获取配置
  Eureka Config中可配置自动探测指定网络中eureka服务或者指定eureka地址列表
  Global Config中可做一些常规配置,例如本地微服务、logger等
2.erpc
  cluster：eureka集群实体,并提供一些实体服务,例如选举出leader的方法
  eurekabaseservice：eureka server实体的基础服务,例如eureka自动探测eureka服务器并加入集群
  eurekaservice：提供eureka的一些具体api服务,例如向eureka cluster注册、注销、查询应用、查询实例等等
  feignclient：feign client会选取cluster中的leader instance提供远程调用服务,具有一定的负载均衡能力
3.client
  本地web服务,可作为微服务注册到eureka
4.其他为工具集
```

## 使用说明

例如:
从配置文件中获取eureka的相关信息,并将本地微服务注册上去,再远程调用

```text
    erkSvr := erpc.GetEurekaServerFromConfig()
    es := erpc.NewEurekaService(erkSvr)
    fc := es.GetFeignClient("appName", "/contextPath")
    res, _ := fc.GetForEntity("/endpoint")
    fmt.Println(string(res))
```

## 其他

```text
1.日志写入文件和https功能尚未加入
2.eureka config中
如果配置自动检测为true,则eurekaNodes配置将无效;
但如果eurekaNodes配置中的服务器在自动检测网段并且服务正常的话,依然可以被检测到并作为cluster的一个节点;
如果network配置为空,则会检测本机ip所在网段中eureka服务;
```
