# multi-version-api-sample

多版本 API 编程示例。

示例说明：

一开始，我们创建了 v1 版本的 User API 资源，Spec 里面包含两个字段：FullName 和 Age。

后来根据新的需求，我们需要将 FullName 字段拆分成 FirstName 和 LastName 字段，这时候我们就需要创建 v2 版本的 User API 资源，Spec 里面包含三个字段：FirstName、LastName 和 Age。

我们升级了 API 之后，之前的使用 FullName 创建 User 的方式我们仍然需要支持，做到 API 的向下兼容。

## 背景

API 不断的迭代，会出现多个版本的 API，需要集群同时支持多个版本的 API 以向下兼容。

## 开始

### 生成项目

参照脚本 `gen_project.sh` 生成项目。

### 编写不同版本的 API

例如 编写 `sampleapis.poneding.com/v1` 和 `sampleapis.poneding.com/v2` 版本的 `user` API。

存在多版本 API 时，需要使用 `// +kubebuilder:storageversion` 注释指定默认的存储版本。例如在 `v2.User` 的结构体上添加注释 `// +kubebuilder:storageversion`，则 `v2.User` 为默认的存储版本。

### 编写转换函数

1、`v1.User` 结构体添加函数 `ConvertTo` 和 `ConvertFrom`，用于转换到 v2.User 结构体。

2、`v2.User` 结构体添加函数 `Hub`，用于指定内部版本（Hub）的结构体。

### 编写 Webhook

最终 etcd 中存储 `v2.User` 版本的 API 数据，所以需要在存储到 etcd 之前，将 `v1.User` 版本的 API 数据转换为 `v2.User` 版本的 API 数据，这其中可能需要对数据做修改（Mutating）和验证（Validating）准入控制。

1、为 `v2.User` 结构体实现 `webhook.Defaulter` 和 `webhook.Validator` 接口函数，用于设置默认值。

2、为 `v2.User` 添加 `// +kubebuilder:webhook:path=/mutate-sampleapis-poneding-com-v2-user,mutating=true,failurePolicy=fail,groups=sampleapis.poneding.com,resources=users,verbs=create;update,versions=v2,name=muser.dp.io,sideEffects=None,admissionReviewVersions=v1
` 和 `// +kubebuilder:webhook:verbs=create;update;delete,path=/validate-sampleapis-poneding-com-v2-user,mutating=false,failurePolicy=fail,groups=sampleapis.poneding.com,resources=users,versions=v2,name=vuser.dp.io,sideEffects=None,admissionReviewVersions=v1` 注释，用于生成 Admission Webhook。

### 生成资源清单

```bash
make generate
make manifest
```

### Kustomize 资源清单

1、取消 `config/crd/kustomization.yaml` 文件中 `patches/cainjection_in_sampleapis_users.yaml` 和 `patches/webhook_in_sampleapis_users.yaml` 注释；

2、取消 `config/default/kustomization.yaml` 文件中 `../certmanager`、`../webhook`、`manager_webhook_patch.yaml`、`webhookcainjection_patch.yaml` 的注释；

3、取消 `config/default/kustomization.yaml` 文件中 `CERTMANAGER` 段 `replacements` 内容的注释。

### 镜像

1、修改 `Makefile` 中 `IMG` 变量，指定镜像地址，打多架构镜像并推送到镜像仓库中;

2、修改 `Makefile` 中 `PLATFORMS` 变量，只保留最常见的架构：`linux/amd64`、`linux/arm64`。

3、编译镜像并推送

```bash
make docker-buildx
```

## 部署和验证

```bash
make install deploy
```
> `kube-rbac-proxy` 默认会拉取 `gcr.io/kubebuilder/kube-rbac-proxy` 下的镜像，但是国内网络无法拉取，直接去除 `gcr.io` 镜像前缀即可，将会直接拉取 Docker Hub 的镜像。

配置 `config/samples/sampleapis_v1_user.yaml` 和 `config/samples/sampleapis_v2_user.yaml` 文件，分别创建 `v1.User` 和 `v2.User` 资源。

```bash
kubectl apply -f config/samples/sampleapis_v1_user.yaml
``` 

此时查看资源
```bash
kubectl get user.sampleapis.poneding.com user-sample -o yaml
```

能看到存储的是 `v2.User` 结构的数据。

## 卸载和清理

```bash
make undeploy uninstall
```
