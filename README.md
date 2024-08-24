# go-napcat

上啊！瞌睡猫猫！

带有 NapCat 支持的 OneBot 协议库。

> [!NOTE]
> 项目还处于 v0 阶段，测试不完善，接口不稳定！请谨慎使用。

## 导入

```sh
go get -u github.com/nekoite/go-napcat
```

## 机器人与初始化

首先，使用 `gonapcat.Init(*config.LogConfig)` 进行初始化并配置日志。不要忘记使用 `defer gonapcat.Finalize()` 或在最后调用这个函数来安全退出。

然后，使用 `gonapcat.NewBot(*config.BotConfig)` 创建新的机器人实例。Config 中使用的机器人 QQ 号需要与对应的 NapCat（或 OneBot）客户端上使用的一致。

然后，使用 `bot.RegisterHandler*` 系列方法配置监听器。要开始监听，使用 `bot.Start()`。要停止监听，使用 `bot.Close()`。停止后，无法再重新启动，需要创建新的机器人实例。

> [!CAUTION]
> `bot.Start()` 是非阻塞的。请使用管道或 `WaitGroup` 阻塞当前 Go 程。

## 例子

<https://github.com/nekoite/go-napcat/tree/master/examples>

## 事件与指令

### 事件处理顺序

1. 监听所有事件的处理器（HandlerAllTypes）
2. 指令（Command）- 仅对消息事件有效
3. 监听各种事件的处理器

### 事件

包裹：`event`。

事件类型：`EventType*`，`MetaEvent[Subtype|Type]*`，`MessageEvent[Subtype|Type]*`， `NoticeEvent[Subtype|Type]*`，`RequestEventType*`，`GroupRequestSubtype*`，`HonorType*`。

所有事件都实现 `IEvent` 接口，使用 `GetEventType()` 查询事件类型后，将事件转为一个具体实现结构体。具体实现在 `*Event`。

### 指令

包裹：`event`。指令使用 [kong](https://github.com/alecthomas/kong) 处理。处理方式和命令行一样。

默认情况下，kong 的 stdout，stderr 和返回值将变成参数（字符串或整形）传入回调函数。这一行为可以被重载。

指令的字符串内容均为转义后的 CQ 码字符串。使用前请使用 `message.ParseCQString` 转换为 `*message.Chain`，或者使用 `message.UnescapeCQString` 对它进行反转义（如果你只需要用到字符串形式的话）。

#### 定义

想要定义一个指令，需要实现 `event.ICommand` 接口。其中，`GetName()` 返回指令名称或前缀，以及名称模式。返回的指令名称不需要经过转义。**指令名称中不能包含空格。**

首先，指令名称是从开头到*第一个非文本元素*或*第一个空格*为止的字符串。例如，`send x y z` 的指令名是 `send`，`send[CQ:image,file=./a.jpg]` 的指令名也是 `send`。

如果名称模式为 `CmdNameModePrefix`，则提供的名称将作为前缀与上述处理后的指令名称匹配。例如，如果消息是 `prefABC`，而 `GetName()` 返回 `("pref", CmdNameModePrefix)`，则匹配成功，并且指令名称是 `pref`，剩余部分是 `ABC`。消息前缀将先被反转义后作比较。

如果名称模式为 `CmdNameModeNormal`，则直接用处理后得到的指令名称反转义后进行匹配查找。

#### 指令的参数分割

在剩下的字符串中，有两种参数分割方式。

- 当 `event.ICommand::SplitBySpaceOnly()` 返回 `true` 时，只使用空格分割。例如，`x y z` 将会分割为 `[x, y, z]`，`a1[CQ:image,file=./a.jpg]a2 b c` 将会被分割为 `[a1[CQ:image,file=./a.jpg]a2, b, c]`。
- 当 `event.ICommand::SplitBySpaceOnly()` 返回 `false` 时，不同类型元素之间也会分割。例如，`a1[CQ:image,file=./a.jpg]a2 b c` 将会被分割为 `[a1, [CQ:image,file=./a.jpg], a2, b, c]`。

如果需要在参数中附带空格，需要将内容使用 `"` 括起来变为引号字符串。例如，`"a1[CQ:image,file=./a.jpg] a2" b c` 将会被分割为 `[a1[CQ:image,file=./a.jpg] a2, b, c]`（`SplitBySpaceOnly = true` 模式下）。在引号字符串中，需要使用以下转义字符：

- `"` -> `\"`
- `\` -> `\\`

> [!NOTE]
> CQ 码中的引号与转义字符不会有任何影响

`event.ICommand::Preprocess(remaining string)` 接口函数用于对消息在分割之前进行预处理。这里传入的 `remaining` 参数将是除去指令名称的剩余的分割之前的字符串。

## OneBot WebSocket API 调用

API 集成于 `Bot` 对象。返回的是 `*api.Resp[T]`。

### 扩展接口

```go
ext := api.NewExtension("name").WithActions(map[api.Action]api.GetNewResultFunc{
    // ...
})
err := ext.Register()
```

#### 内置扩展

- NapCat 扩展：`napcat.Extension.Register()`

## 日志

请使用机器人对象自带的 `bot.[Log|Info|Debug|Warn|Error|Fatal]`。格式使用 [zap]。

## 自带库

包括但不限于：

- 用于打印日志的 [zap]
- 用于解析 JSON 的 [goccy/go-json](https://github.com/goccy/go-json) 和 [gjson](https://github.com/tidwall/gjson)
- 用于解析命令行参数的 [kong](https://github.com/alecthomas/kong)

可以直接使用以上推荐的库。

[zap]: https://github.com/uber-go/zap
