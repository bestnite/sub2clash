# sub2clash

将订阅链接转换为 Clash、Clash.Meta 配置  
[预览](https://clash.nite07.com/)

## 特性

- 开箱即用的规则、策略组配置
- 自动根据节点名称按国家划分策略组
- 多订阅合并
- 自定义 Rule Provider、Rule
- 支持多种协议
  - Shadowsocks
  - ShadowsocksR
  - Vmess
  - Vless （Clash.Meta）
  - Trojan
  - Hysteria （Clash.Meta）
  - Hysteria2 （Clash.Meta）
  - Socks5
  - Anytls （Clash.Meta）

## 使用

### 部署

- [docker compose](./compose.yml)
- 运行[二进制文件](https://github.com/bestnite/sub2clash/releases/latest)

### 配置

支持多种配置方式，按优先级排序：

1. **配置文件**：支持多种格式（YAML、JSON），按以下优先级搜索：
   - `config.yaml` / `config.yml`
   - `config.json`
   - `sub2clash.yaml` / `sub2clash.yml`
   - `sub2clash.json`
2. **环境变量**：使用 `SUB2CLASH_` 前缀，例如 `SUB2CLASH_ADDRESS=0.0.0.0:8011`
3. **默认值**：内置默认配置

| 配置项                | 环境变量                        | 说明                                    | 默认值                                                                                               |
| --------------------- | ------------------------------- | --------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| address               | SUB2CLASH_ADDRESS               | 服务监听地址                            | `0.0.0.0:8011`                                                                                       |
| meta_template         | SUB2CLASH_META_TEMPLATE         | 默认 meta 模板 URL                      | `https://raw.githubusercontent.com/bestnite/sub2clash/refs/heads/main/templates/template_meta.yaml`  |
| clash_template        | SUB2CLASH_CLASH_TEMPLATE        | 默认 clash 模板 URL                     | `https://raw.githubusercontent.com/bestnite/sub2clash/refs/heads/main/templates/template_clash.yaml` |
| request_retry_times   | SUB2CLASH_REQUEST_RETRY_TIMES   | 请求重试次数                            | `3`                                                                                                  |
| request_max_file_size | SUB2CLASH_REQUEST_MAX_FILE_SIZE | 请求文件最大大小（byte）                | `1048576`                                                                                            |
| cache_expire          | SUB2CLASH_CACHE_EXPIRE          | 订阅缓存时间（秒）                      | `300`                                                                                                |
| log_level             | SUB2CLASH_LOG_LEVEL             | 日志等级：`debug`,`info`,`warn`,`error` | `info`                                                                                               |
| short_link_length     | SUB2CLASH_SHORT_LINK_LENGTH     | 短链长度                                | `6`                                                                                                  |

#### 配置文件示例

参考示例文件：

- [config.example.yaml](./config.example.yaml) - YAML 格式
- [config.example.json](./config.example.json) - JSON 格式

### API

#### `GET /convert/:config`

获取 Clash/Clash.Meta 配置链接

| Path 参数 | 类型   | 说明                                           |
| --------- | ------ | ---------------------------------------------- |
| config    | string | Base64 URL Safe 编码后的 JSON 字符串，格式如下 |

##### `config` JSON 结构

| Query 参数         | 类型              | 是否必须                 | 默认值    | 说明                                                                                                     |
| ------------------ | ----------------- | ------------------------ | --------- | -------------------------------------------------------------------------------------------------------- |
| clashType          | int               | 是                       | 1         | 配置文件类型 (1: Clash, 2: Clash.Meta)                                                                   |
| subscriptions      | []string          | sub/proxy 至少有一项存在 | -         | 订阅链接（v2ray 或 clash 格式），可以在链接结尾加上`#名称`，来给订阅中的节点加上统一前缀（可以输入多个） |
| proxies            | []string          | sub/proxy 至少有一项存在 | -         | 节点分享链接（可以输入多个）                                                                             |
| refresh            | bool              | 否                       | `false`   | 强制刷新配置（默认缓存 5 分钟）                                                                          |
| template           | string            | 否                       | -         | 外部模板链接或内部模板名称                                                                               |
| ruleProviders      | []RuleProvider    | 否                       | -         | 规则                                                                                                     |
| rules              | []Rule            | 否                       | -         | 规则                                                                                                     |
| autoTest           | bool              | 否                       | `false`   | 国家策略组是否自动测速                                                                                   |
| lazy               | bool              | 否                       | `false`   | 自动测速是否启用 lazy                                                                                    |
| sort               | string            | 否                       | `nameasc` | 国家策略组排序策略，可选值 `nameasc`、`namedesc`、`sizeasc`、`sizedesc`                                  |
| replace            | map[string]string | 否                       | -         | 通过正则表达式重命名节点                                                                                 |
| remove             | string            | 否                       | -         | 通过正则表达式删除节点                                                                                   |
| nodeList           | bool              | 否                       | `false`   | 只输出节点                                                                                               |
| ignoreCountryGroup | bool              | 否                       | `false`   | 是否忽略国家分组                                                                                         |
| userAgent          | string            | 否                       | -         | 订阅 user-agent                                                                                          |
| useUDP             | bool              | 否                       | `false`   | 是否使用 UDP                                                                                             |

###### `RuleProvider` 结构

| 字段     | 类型   | 说明                                                             |
| -------- | ------ | ---------------------------------------------------------------- |
| behavior | string | rule-set 的 behavior                                             |
| url      | string | rule-set 的 url                                                  |
| group    | string | 该规则集使用的策略组名                                           |
| prepend  | bool   | 如果为 `true` 规则将被添加到规则列表顶部，否则添加到规则列表底部 |
| name     | string | 该 rule-provider 的名称，不能重复                                |

###### `Rule` 结构

| 字段    | 类型   | 说明                                                             |
| ------- | ------ | ---------------------------------------------------------------- |
| rule    | string | 规则                                                             |
| prepend | bool   | 如果为 `true` 规则将被添加到规则列表顶部，否则添加到规则列表底部 |

### 模板

可以通过变量自定义模板中的策略组代理节点  
具体参考下方默认模板

- `<all>` 为添加所有节点
- `<countries>` 为添加所有国家策略组
- `<地区二位字母代码>` 为添加指定地区所有节点，例如 `<hk>` 将添加所有香港节点

#### 默认模板

- [Clash](./templates/template_clash.yaml)
- [Clash.Meta](./templates/template_meta.yaml)

## 开发

### 添加新协议支持

添加新协议支持需要实现以下组件：

1. 在 `parser` 目录下实现协议解析器，用于解析节点链接
2. 在 `model/proxy` 目录下定义协议结构体

## 贡献者

[![](https://contrib.rocks/image?repo=bestnite/sub2clash)](https://github.com/bestnite/sub2clash/graphs/contributors)
