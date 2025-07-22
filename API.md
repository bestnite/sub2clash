# `GET /convert/:config`

获取 Clash/Clash.Meta 配置链接

| Path 参数 | 类型   | 说明                                           |
| --------- | ------ | ---------------------------------------------- |
| config    | string | Base64 URL Safe 编码后的 JSON 字符串，格式如下 |

## `config` JSON 结构

| Query 参数         | 类型              | 是否必须                 | 默认值    | 说明                                                                              |
| ------------------ | ----------------- | ------------------------ | --------- | --------------------------------------------------------------------------------- |
| clashType          | int               | 是                       | 1         | 配置文件类型 (1: Clash, 2: Clash.Meta)                                            |
| subscriptions      | []string          | sub/proxy 至少有一项存在 | -         | 订阅链接，可以在链接结尾加上`#名称`，来给订阅中的节点加上统一前缀（可以输入多个） |
| proxies            | []string          | sub/proxy 至少有一项存在 | -         | 节点分享链接（可以输入多个）                                                      |
| refresh            | bool              | 否                       | `false`   | 强制刷新配置（默认缓存 5 分钟）                                                   |
| template           | string            | 否                       | -         | 外部模板链接或内部模板名称                                                        |
| ruleProviders      | []RuleProvider    | 否                       | -         | 规则                                                                              |
| rules              | []Rule            | 否                       | -         | 规则                                                                              |
| autoTest           | bool              | 否                       | `false`   | 国家策略组是否自动测速                                                            |
| lazy               | bool              | 否                       | `false`   | 自动测速是否启用 lazy                                                             |
| sort               | string            | 否                       | `nameasc` | 国家策略组排序策略，可选值 `nameasc`、`namedesc`、`sizeasc`、`sizedesc`           |
| replace            | map[string]string | 否                       | -         | 通过正则表达式重命名节点                                                          |
| remove             | string            | 否                       | -         | 通过正则表达式删除节点                                                            |
| nodeList           | bool              | 否                       | `false`   | 只输出节点                                                                        |
| ignoreCountryGroup | bool              | 否                       | `false`   | 是否忽略国家分组                                                                  |
| userAgent          | string            | 否                       | -         | 订阅 user-agent                                                                   |
| useUDP             | bool              | 否                       | `false`   | 是否使用 UDP                                                                      |

### `RuleProvider` 结构

| 字段     | 类型   | 说明                                                             |
| -------- | ------ | ---------------------------------------------------------------- |
| behavior | string | rule-set 的 behavior                                             |
| url      | string | rule-set 的 url                                                  |
| group    | string | 该规则集使用的策略组名                                           |
| prepend  | bool   | 如果为 `true` 规则将被添加到规则列表顶部，否则添加到规则列表底部 |
| name     | string | 该 rule-provider 的名称，不能重复                                |

### `Rule` 结构

| 字段    | 类型   | 说明                                                             |
| ------- | ------ | ---------------------------------------------------------------- |
| rule    | string | 规则                                                             |
| prepend | bool   | 如果为 `true` 规则将被添加到规则列表顶部，否则添加到规则列表底部 |

# `POST /short`

获取短链，Content-Type 为 `application/json`
具体参考使用可以参考 [api\templates\index.html](api/static/index.html)

| Body 参数 | 类型   | 是否必须 | 默认值 | 说明                      |
| --------- | ------ | -------- | ------ | ------------------------- |
| url       | string | 是       | -      | 需要转换的 Query 参数部分 |
| password  | string | 否       | -      | 短链密码                  |

# `GET /s/:hash`

短链跳转
`hash` 为动态路由参数，可以通过 `/short` 接口获取

| Query 参数 | 类型   | 是否必须 | 默认值 | 说明     |
| ---------- | ------ | -------- | ------ | -------- |
| password   | string | 否       | -      | 短链密码 |

# `PUT /short`

更新短链，Content-Type 为 `application/json`

| Body 参数 | 类型   | 是否必须 | 默认值 | 说明                      |
| --------- | ------ | -------- | ------ | ------------------------- |
| url       | string | 是       | -      | 需要转换的 Query 参数部分 |
| password  | string | 否       | -      | 短链密码                  |
| hash      | string | 是       | -      | 短链 hash                 |
