# autoacme

一个自动化申请证书并部署的工具

### 已完成

- [x] 阿里云
  - [x] 申请证书
    - [x] DNS
    - [x] OSS
    - [ ] HTTP
  - [x] 上传证书
    - [x] CDN
- [ ] 七牛云
  - [x] SSL上传

### 使用方法

```shell
docker run -it --rm \
  --name autoacme \
  -v ${PWD}/config/config.yaml:/etc/autoacme/config.yaml \
  -v ${PWD}/runtime:/etc/autoacme \
  ghcr.io/liasica/autoacme:master
```

### 配置文件

```yaml
account: ca@example.com # 邮箱
dns:
  - 223.5.5.5 # DNS服务器

domains:
  - domain: share.example.com # 待申请证书域名
    provider: DNS # 证书申请方式
    dnsProvider:
      accessKeyId: LAIT5... # 阿里云AccessKey
      accessKeySecret: rTMKy... # 阿里云AccessKey
    ossProvider:
      accessKeyId: LAIT5... # 阿里云AccessKey
      accessKeySecret: rTMKy... # 阿里云AccessKey
      bucket: share-example # OSS Bucket
      endpoint: oss-cn-beijing.aliyuncs.com # OSS Endpoint
      path: # 证书存储在oss上的路径
    hooks:
      - name: CDN # 部署到阿里云CDN
        cdnHook:
          accessKeyId: LAIT5... # 阿里云AccessKey
          accessKeySecret: rTMKy... # 阿里云AccessKey
```

### 阿里云权限

- AliyunOSSFullAccess
- AliyunYundunCertFullAccess
- AliyunDNSFullAccess <dns:DescribeDomains, dns:DeleteDomainRecord, dns:AddDomainRecord>
- AliyunCDNFullAccess <cdn:SetDomainServerCertificate>
