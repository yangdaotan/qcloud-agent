# User Guide

## url
qcloud.bike.io:8360

关于黑石物理机规范，更多见：[腾讯云api](https://cloud.tencent.com/document/api/386/9308)

## bm（黑石物理机）

- GetDeviceInstanceIds （查询机器实例Id)


```
http://qcloud.bike.io:8360/bm?
	action=GetDeviceInstanceIds
	&lanIps.0=10.201.95.2
	&lanIps.1=10.201.54.2
```
参数说明

```
lanIps表示多个ip，分别以lanIps.0,lanIps.1...作为参数
```

## bmlb (黑石负载均衡)

- CreateBmLoadBalancer （创建）


在nodeIp的网段中创建一个lb

```
http://qcloud.bike.io:8360/bmlb?
	action=CreateBmLoadBalancer
	&nodeIp=10.201.95.2 //在nodeIp所在vpc的子网段中分配一个lb
```
 2. **修改lb名称**

 在nodeIp的网段中创建一个lb，并且修改lb名字

```
http://qcloud.bike.io:8360/bmlb?
	action=CreateBmLoadBalancer
	&nodeIp=10.201.95.2		//nodeIp所在vpc的子网段分配一个lb
	&loadBalancerName=ydt-test-lb // 将lb名称修改为ydt-test-lb
```

- ModifyBmLoadBalancerAttributes (修改lb名称)

```
http://qcloud.bike.io:8360/bmlb?
	action=ModifyBmLoadBalancerAttributes
	&loadBalancerId=lbbj4zk7oj //负责均衡器Id
	&loadBalancerName=ydt-test-lb
```

- CreateBmListeners （创建负载均衡监听器）

```
http://qcloud.bike.io:8360/bmlb?
	action=CreateBmListeners
	&loadBalancerId=lb-bj4zk7oj 	// 负责均衡器Id
	&listeners.0.loadBalancerPort=9876	//创建监听器0，监听端口为9876
	&listeners.0.protocol=tcp				// 四层
	&listeners.0.listenerName=NAMESRV		// 监听器名称
	&listeners.0.toaFlag=1					// 是否获取clienIP
	&listeners.0.healthSwitch=1			// 是否开启监控监测
	&listeners.1.loadBalancerPort=8080	// 创建监听器1，监听端口为8080
	&listeners.1.protocol=tcp				// 监听器1为四层
	&listeners.1.listenerName=CONSOLE		// 监听器1名称为CONSOLE
	&listeners.1.toaFlag=0					// 不获取clientIP
	&listeners.1.healthSwitch=1			// 开启监控监测

```

- BindBmL4ListenerRs （绑定物理机器到负载均衡监听器）

```
http://qcloud.bike.io:8360/bmlb?
	action=BindBmL4ListenerRs
	&loadBalancerId=lb-bj4zk7oj	// 负载均衡器Id
	&listenerId=lbl-8zuk81nl		// 监听器Id
	&backends.0.port=30959			// 第一台后端物理机的端口为30959
	&backends.0.instanceId=cpm-ayqlmx5v  // 第一台后端物理机Id
	&backends.0.weight=10			// 第一台负责均衡权重
```

## cns (域名解析)
- RecordCreate （添加解析）

```
http://qcloud.bike.io:8360/cns?
	action=RecordCreate
	&subDomain=mq.ydt	// 子域名
	&vip=10.201.116.2   // 负载均衡ip
```


## k8s (一键部署lb、lbs、record)
- Create (创建)

```
http://qcloud.bike.io:8360/k8s?
    &action=Create
	subDomain=zk.ydt    //域名，zk.ydt.bike.io格式
	&lbName=test-ydt-lb //负责均衡名称
	&ports.0=9092:30002 // lb listener 0端口映射关系
	&ports.1=9093:30003 // lb listener 1端口映射关系
	&nodeIps.0=10.201.95.2  // 后端机器0
	&nodeIps.1=10.201.53.2  // 后端机器1
```

- Delete (销毁lb，record...)

```
http://qcloud.bike.io:8360/k8s?
    &action=Delete
	subDomain=zk.ydt    //域名，zk.ydt.bike.io格式
	&lbName=test-ydt-lb //负责均衡名称
```

## 用法
- 替换mobike.io/infra/qcloud-agent/qcloud/consts.go中的key
```
const SecretId  = "ReplaceMeUseYourSecretId"
const SecretKey  = "ReplaceMeUseYourSecretKey"
```

## build

```
gox -os linux -arch amd64
operator-sdk build docker.bike.io/infra/qcloud-agent:v0.0.1

# qcloud-agent
