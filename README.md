# go-napcat

[![Go Reference](https://pkg.go.dev/badge/github.com/nekoite/go-napcat.svg)](https://pkg.go.dev/github.com/nekoite/go-napcat) [![Go Report Card](https://goreportcard.com/badge/github.com/nekoite/go-napcat)](https://goreportcard.com/report/github.com/nekoite/go-napcat) ![Build Status](https://github.com/nekoite/go-napcat/actions/workflows/build.yml/badge.svg) [![codecov](https://codecov.io/gh/nekoite/go-napcat/graph/badge.svg?token=ZW82R4ZV7F)](https://codecov.io/gh/nekoite/go-napcat)

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

## 消息链与消息片段

本框架只支持传入原始消息为 JSON 数组。它将会变成消息事件及其它带有消息的结构体中的消息字段，类型为 `message.Chain`（消息链）。其中包含多个 `message.Segment`（消息片段）。每个消息片段包含消息类型 `Type message.SegmentType` 以及消息内容 `Data any`。消息内容的类型根据消息类型不同，为 `message.*Data` 的指针的一种（例如 `Type` 为 `SegmentTypeText` 时，`Data` 的类型是 `*message.TextData`）。

使用 `New*()` 构建新的消息内容，或使用 `New*Segment()` 构建消息片段。也可以使用 `New*().Segment()` 构建消息片段。

使用 `Segment::AsChain()` 构建仅包含当前消息片段的消息链。使用 `Segment::Get*Data()` 获取对应类型的消息内容。如果消息类型不匹配，则返回 `nil`。

使用 `message.NewChain()` 构造空消息链。使用 `Chain::PrependSegment(Segment)`，`Chain::AddSegment(Segment)`，`Chain::AddSegments(...Segment)`，以及 `Chain::Add*()` 来操作消息链。

具体内容请看[文档](https://pkg.go.dev/github.com/nekoite/go-napcat/message#pkg-index)。

## 事件与指令

### 事件处理顺序

1. 指令（Command）- 仅对消息事件有效
2. 监听所有事件的处理器（HandlerAllTypes）
3. 监听各种事件的处理器

### 事件

包裹：`event`。

事件类型：`EventType*`，`MetaEvent[Subtype|Type]*`，`MessageEvent[Subtype|Type]*`， `NoticeEvent[Subtype|Type]*`，`RequestEventType*`，`GroupRequestSubtype*`，`HonorType*`。

所有事件都实现 `IEvent` 接口，使用 `GetEventType()` 查询事件类型后，将事件转为一个具体实现结构体。具体实现在 `event.*Event` 结构体。

可以使用 `event.GetAs` 函数做类型转换，失败返回 `nil`。使用 `event.GetAsUnsafe` 做类型转换，效果和 `e.(*T)` 一样，失败会 panic。使用 `event.GetAsOrError` 做类型转换，失败会返回 `nil, errors.ErrTypeAssertion`。

使用 `IEvent::SetContext(any)` 和 `IEvent::Context()` 来设置或获取你想使用的上下文。它可以是任意对象。建议使用一个 `map`。

### 指令

包裹：`event`。指令使用 [kong](https://github.com/alecthomas/kong) 处理。处理方式和命令行一样。

默认情况下，kong 的 stdout，stderr 和返回值将变成参数（字符串或整形）传入回调函数。这一行为可以被重载。

指令的字符串内容均为转义后的 CQ 码字符串。使用前请使用 `message.ParseCQString` 转换为 `*message.Chain`，或者使用 `message.UnescapeCQString` 对它进行反转义（如果你只需要用到字符串形式的话）。

#### 定义

想要定义一个指令，需要实现 `event.ICommand` 接口。其中，`GetName()` 返回指令名称或前缀，以及名称模式。返回的指令名称不需要经过转义。**指令名称或前缀中不能包含空格，且只能是普通字符串（不能有 CQ 码）。**如果想要绕过这一限制，请使用事件监听器。

首先，指令名称是从开头到*第一个非文本元素*或*第一个空格*为止的字符串。例如，`send x y z` 的指令名是 `send`，`send[CQ:image,file=./a.jpg]` 的指令名也是 `send`。

如果名称模式为 `CmdNameModePrefix`，则提供的名称将作为前缀与上述处理后的指令名称匹配。例如，如果消息是 `prefABC`，而 `GetName()` 返回 `("pref", CmdNameModePrefix)`，则匹配成功，并且指令名称是 `pref`，剩余部分是 `ABC`。消息前缀将先被反转义后作比较。

如果名称模式为 `CmdNameModeNormal`，则直接用处理后得到的指令名称反转义后进行匹配查找。

如果指令还实现了 `event.ICommandStopPropagation` 接口，并且函数返回 `true`，则处理完这个指令后，事件将不会继续传播（不会触发其它事件处理器）。

#### 全局指令激活前缀

全局激活前缀是指只有在消息开头有这个字符串时，指令才会被激活并检查。

调用 `bot.SetGlobalCommandPrefix` 函数来设置前缀。这个字符串会被从消息中移除后，剩下的部分将被传递给指令。若设置为空字符串，则表示没有激活前缀（默认情况，与不调用此函数一样）。与前缀指令一样，前缀中不能包含空格，且只能是普通字符串。

#### 指令的参数分割

在剩下的字符串中，有两种参数分割方式。

- 当 `event.ICommand::SplitBySpaceOnly()` 返回 `true` 时，只使用空格分割。例如，`x y z` 将会分割为 `[x, y, z]`，`a1[CQ:image,file=./a.jpg]a2 b c` 将会被分割为 `[a1[CQ:image,file=./a.jpg]a2, b, c]`。
- 当 `event.ICommand::SplitBySpaceOnly()` 返回 `false` 时，不同类型元素之间也会分割。例如，`a1[CQ:image,file=./a.jpg]a2 b c` 将会被分割为 `[a1, [CQ:image,file=./a.jpg], a2, b, c]`。

如果需要在参数中附带空格，需要将内容使用 `"` 括起来变为引号字符串。例如，`"a1[CQ:image,file=./a.jpg] a2" b c` 将会被分割为 `[a1[CQ:image,file=./a.jpg] a2, b, c]`（`SplitBySpaceOnly = true` 模式下）。在引号字符串中，需要使用以下转义字符：

- `"` -> `\"`
- `\` -> `\\`

> [!NOTE]
> CQ 码中的引号与转义字符不会有任何影响

`event.ICommandWithPreprocess::Preprocess(remaining string)` 接口函数用于对消息在分割之前进行预处理。这里传入的 `remaining` 参数将是除去指令名称的剩余的分割之前的字符串。如果指令实现这个接口，则会在分割前调用此函数。如果不实现这个接口，则直接进行分割。

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
