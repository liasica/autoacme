# aliacme

一个自动化申请证书并部署到阿里云产品的工具

### 已完成

- [ ] 申请证书
  - [x] DNS
  - [ ] CDN
  - [ ] HTTP
- [ ] 上传证书
  - [x] CDN

### 使用方法

```shell
docker run -it --rm \
  --name aliacme \
  -v ${PWD}/config/config.yaml:/etc/aliacme/config.yaml \
  -v ${PWD}/runtime:/etc/aliacme \
  ghcr.io/liasica/aliacme:master
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
      - name: CDN # 部署到CDN
        cdnHook:
          accessKeyId: LAIT5... # 阿里云AccessKey
          accessKeySecret: rTMKy... # 阿里云AccessKey
```

### 阿里云权限

- AliyunOSSFullAccess
- AliyunYundunCertFullAccess
- AliyunDNSFullAccess <dns:DescribeDomains, dns:DeleteDomainRecord, dns:AddDomainRecord>
- AliyunCDNFullAccess <cdn:SetDomainServerCertificate>