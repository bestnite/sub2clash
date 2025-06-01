# `GET /clash`, `GET /meta`

获取 Clash/Clash.Meta 配置链接

| Query 参数   | 类型   | 是否必须                 | 默认值    | 说明                                                                                                                                                                                                                                      |
| ------------ | ------ | ------------------------ | --------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| sub          | string | sub/proxy 至少有一项存在 | -         | 订阅链接，可以在链接结尾加上`#名称`，来给订阅中的节点加上统一前缀（可以输入多个，用 `,` 分隔）                                                                                                                                            |
| proxy        | string | sub/proxy 至少有一项存在 | -         | 节点分享链接（可以输入多个，用 `,` 分隔）                                                                                                                                                                                                 |
| refresh      | bool   | 否                       | `false`   | 强制刷新配置（默认缓存 5 分钟）
| skipErrors      | bool   | 否                       | `false`   | 跳过错误订阅                                                                                                                                                                                                            |
| template     | string | 否                       | -         | 外部模板链接或内部模板名称                                                                                                                                                                                                                |
| ruleProvider | string | 否                       | -         | 格式 `[Behavior,Url,Group,Prepend,Name],[Behavior,Url,Group,Prepend,Name]...`，其中 `Group` 是该规则集使用的策略组名，`Prepend` 为 bool 类型，如果为 `true` 规则将被添加到规则列表顶部，否则添加到规则列表底部（会调整到 MATCH 规则之前） |
| rule         | string | 否                       | -         | 格式 `[Rule,Prepend],[Rule,Prepend]...`，其中 `Prepend` 为 bool 类型，如果为 `true` 规则将被添加到规则列表顶部，否则添加到规则列表底部（会调整到 MATCH 规则之前）                                                                         |
| autoTest     | bool   | 否                       | `false`   | 国家策略组是否自动测速                                                                                                                                                                                                                    |
| lazy         | bool   | 否                       | `false`   | 自动测速是否启用 lazy                                                                                                                                                                                                                     |
| sort         | string | 否                       | `nameasc` | 国家策略组排序策略，可选值 `nameasc`、`namedesc`、`sizeasc`、`sizedesc`                                                                                                                                                                   |
| replace      | string | 否                       | -         | 通过正则表达式重命名节点，格式 `[<ReplaceKey>,<ReplaceTo>],[<ReplaceKey>,<ReplaceTo>]...`                                                                                                                                                 |
| remove       | string | 否                       | -         | 通过正则表达式删除节点                                                                                                                                                                                                                    |
| nodeList     | bool   | 否                       | `false`   | 只输出节点                                                                                                                                                                                                                                |

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
