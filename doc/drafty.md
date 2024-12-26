# Drafty: 一种富文本格式


这份文档介绍了 Drafty，一种由 Tinode 使用的文本格式，用于为消息添加样式。Drafty 的目标是在表达能力足够的同时，不会开放太多的可能性以避免安全问题。你可以将它视为将 JSON 封装了一下的 [markdown](https://en.wikipedia.org/wiki/Markdown)。Drafty 受到了 Facebook 的 [draft.js](https://draftjs.org/) 规范的影响。截至撰写本文时，已经存在 [JavaScript](https://github.com/tinode/tinode-js/blob/master/src/drafty.js)、[Java](https://github.com/tinode/tindroid/blob/master/tinodesdk/src/main/java/co/tinode/tinodesdk/model/Drafty.java) 和 [Swift](https://github.com/tinode/ios/blob/master/TinodeSDK/model/Drafty.swift) 的实现。还有一个 [Go 实现](https://github.com/tinode/chat/blob/master/server/drafty/drafty.go)，可以将 Drafty 转换为纯文本和预览。

这样的好处就是在机器人交互时候可以发送一些格式化的文本，而不仅仅是微信和QQ那种无格式文本；

## 示例

> this is **bold**, `code` and _italic_, ~~strike~~<br/>
>  combined **bold and _italic_**<br/>
>  an url: https://www.example.com/abc#fragment and another _[https://web.tinode.co](https://web.tinode.co)_<br/>
>  this is a [@mention](#) and a [#hashtag](#) in a string<br/>
> second [#hashtag](#)<br/>

上面的格式可以用下面的方式描述:

```
draft: {txt: ' ',fmt: [{tp: 'null', at: 0, len: 1, key: 0}],ent: [{tp: 'IM', data: {ref=content://com.miui.gallery.open/raw/%2Fstorage%2Femulated%2F0%2FDCIM%2FCamera_XHS_1733128209586.jpg, size=844681, mime=image/jpeg, width=1270, height=1920}}]}
```



```js
{
   "txt":  "this is bold, code and italic, strike combined bold and italic an url: https://www.example.com/abc#fragment and another www.tinode.co this is a @mention and a #hashtag in a string second #hashtag",
   "fmt": [
       { "at":8, "len":4,"tp":"ST" },
       { "at":14, "len":4, "tp":"CO" },
       { "at":23, "len":6, "tp":"EM"},
       { "at":31, "len":6, "tp":"DL" },
       { "tp":"BR", "len":1, "at":37 },
       
       { "at":56, "len":6, "tp":"EM" },
       { "at":47, "len":15, "tp":"ST" },
       { "tp":"BR", "len":1, "at":62 },
       
       { "at":120, "len":13, "tp":"EM" },
       { "at":71, "len":36, "key":0 },
       { "at":120, "len":13, "key":1 },
       { "tp":"BR", "len":1, "at":133 },
       
       { "at":144, "len":8, "key":2 },
       { "at":159, "len":8, "key":3 },
       { "tp":"BR", "len":1, "at":179 },
       
       { "at":187, "len":8, "key":3 },
       { "tp":"BR", "len":1, "at":195 }
   ],
   "ent": [
       { "tp":"LN", "data":{ "url":"https://www.example.com/abc#fragment" } },
       { "tp":"LN", "data":{ "url":"http://www.tinode.co" } },
       { "tp":"MN", "data":{ "val":"mention" } },
       { "tp":"HT", "data":{ "val":"hashtag" } }
   ]
}
```

## 结构

Drafty 对象有三个字段：纯文本 `txt`、行内标记 `fmt` 和实体 `ent`。

- **纯文本**：消息被转换为纯 Unicode 文本，并且去除了所有标记，保存在 `txt` 字段中。

- **行内标记 fmt**：行内标记是一个包含了每种样式的数组。每种样式由一个对象表示，至少包含 `at` 和 `len` 字段。`at` 表示在 `txt` 中的偏移量，`len` 是要应用样式的字符数。

  样式的第三个值可以是 `tp` 或 `key`。如果是 `tp`，则表示样式是基本的文本装饰；如果是 `key`，则是对 `ent` 数组的索引，其中包含了更多的样式参数，比如图片或 URL。

- **实体 ent**：实体是一种需要额外数据的文本装饰，例如链接、提及或者标签等。实体由一个包含了两个字段的对象表示，`tp` 表示实体的类型，`data` 是类型相关的样式信息。

### 纯文本 `txt`

解释在 Drafty 中如何处理纯文本消息。

消息被发送时会被转换成纯 Unicode 文本，并且所有的标记都会被剥离掉，然后保存在 `txt` 字段中。通常情况下，一个有效的 Drafty 消息可能只包含 `txt` 字段，即纯文本内容。这意味着如果消息不需要任何样式或实体，那么只需要填充 `txt` 字段即可。

### 内联格式 `fmt`

这段说明了在 Drafty 中的行内格式化样式。具体来说：

- **行内格式化** 是一个包含了各种样式的数组，存储在 `fmt` 字段中。每种样式都由一个对象表示，至少包含 `at` 和 `len` 字段。`at` 表示相对于 `txt` 的偏移量（从 0 开始），`len` 表示要应用样式的字符数。每种样式的第三个值可能是 `tp` 或 `key`。

- 如果提供了 `tp`，则表示样式是基本的文本装饰，可以是以下值之一：
  - `BR`：换行。
  - `CO`：代码或等宽字体文本，可能带有不同的背景色。
  - `DL`：删除或删除线文本。
  - `EM`：强调文本，通常表示为斜体。
  - `FM`：表单/字段集合；也可以表示为实体。
  - `HD`：隐藏内容。
  - `HL`：高亮文本，例如不同颜色或背景的文本；颜色不能指定。
  - `RW`：格式的逻辑分组，一行；也可以表示为实体。
  - `ST`：粗体文本。

- 如果提供了 `key`，则是一个基于 `ent` 数组的从 0 开始的索引，其中包含了扩展样式参数，例如图片或 URL。可能的键值包括：
  - `AU`：嵌入式音频。
  - `BN`：交互式按钮。
  - `EX`：通用附件。
  - `FM`：表单/字段集合；也可以表示为基本装饰。
  - `HT`：标签，例如 [#hashtag](#)。
  - `IM`：内联图像。
  - `LN`：链接（URL）[https://api.tinode.co](https://api.tinode.co)。
  - `MN`：提及，例如 [@tinode](#)。
  - `RW`：格式的逻辑分组，一行；也可以表示为基本装饰。
  - `VC`：视频（和音频）通话。
  - `VD`：内联视频。



**示例：**

解释如何使用 Drafty 中的样式对象。具体来说：

- `{ "at":8, "len":4, "tp":"ST"}`：表示从 `txt` 中的偏移量 8 处开始的 4 个字符应用**粗体样式** (`ST`)。
- `{ "at":144, "len":8, "key":2 }`：表示将实体 `ent[2]` 插入到位置 144 处，该实体跨越 8 个字符。
- `{ "at":-1, "len":0, "key":4 }`：表示将实体 `ent[4]` 显示为**文件附件**，但不对文本应用任何样式。

这些示例中的样式对象可以缺少 `at`、`key` 和 `len` 值，客户端应该能够处理这些缺失的情况。如果缺少这些值，则默认将它们视为 `0`。

需要注意的是，`at` 和 `len` 的索引是以 [Unicode 代码点](https://developer.mozilla.org/en-US/docs/Glossary/Code_point) 衡量的，而不是字节或字符。目前对于包含 Fitzpatrick 皮肤色调修饰符、变体选择器或与 `ZWJ` 组合的多代码点字符（例如表情符号）的行为是未定义的。


#### `FM`: 表单、订单以及一组输入

表单可以提供对一组元素进行图形化的展示. 

<table>
<tr><th>是否同意?</th></tr>
<tr><td><a href="">是</a></td></tr>
<tr><td><a href="">否</a></td></tr>
</table>


```js
{
 "txt": "是否同意? 是，否",
 "fmt": [
   {"len": 20, "tp": "FM"}, // missing 'at' is zero: "at": 0
   {"len": 13, "tp": "ST"}
   {"at": 13, "len": 1, "tp": "BR"},
   {"at": 14, "len": 3}, // missing 'key' is zero: "key": 0
   {"at": 17, "len": 1, "tp": "BR"},
   {"at": 18, "len": 2, "key": 1},
 ],
 "ent": [
   {"tp": "BN", "data": {"name": "yes", "act": "pub", "val": "是!"}},
   {"tp": "BN", "data": {"name": "no", "act": "pub"}}
 ]
}
```
如果点击了按钮 `是` , 客户端应该向对方（群组）发送如下的信息:
```js
{
 "txt": "Yes",
 "fmt": [{
   "at":-1
 }],
 "ent": [{
   "tp": "EX",
   "data": {
     "mime": "application/json",
     "val": {
       "seq": 15, // seq id of the message containing the form.
       "resp": {"yes": "oh yes!"}
     }
   }
 }]
}
```

The form may be optionally represented as an entity:
```js
{
  "tp": "FM",
  "data": {
    "su": true
  }
}
```


`data.su` 描述了交互式表单元素在点击后的行为方式。当 `"su": true` 时，表示表单是 `single use` 的：即表单在第一次交互后应更改，以显示它不再接受输入。

这意味着一旦用户与表单交互，表单将被禁用或隐藏，以表示它已经被使用过，不再接受进一步的输入。通常，这是用于在用户提交表单后，防止用户再次进行相同的操作。

### 实体嵌入 `ent`

一般来说，实体是一种需要额外（可能很大）数据的文本装饰。一个实体由两个字段组成：`tp` 表示实体的类型，`data` 是依赖于类型的样式信息。未知的字段会被忽略。

这意味着在 Drafty 中，实体可以是多种类型的装饰，例如链接、提及、附件等，并且每种类型都有其自己的样式和数据。当客户端解析 Drafty 消息时，它会根据 `tp` 字段识别实体的类型，并根据该类型的样式信息对文本进行相应的装饰。

#### `AU`: 嵌入音频
`AU` 是一个音频。 `data` 包括下面的字段：

```js
{
  "tp": "AU",
  "data": {
    "mime": "audio/aac",
    "val": "Rt53jUU...iVBORw0KGgoA==",
    "ref": "/v0/file/s/e769gvt1ILE.m4v",
    "preview": "Aw4JKBkAAAAKMSM...vHxgcJhsgESAY"
    "duration": 180000,
    "name": "ding_dong.m4a",
    "size": 595496
  }
}
```
 * ``mime`: 数据格式，例如 'audio/ogg'，指定音频文件的 MIME 类型。
* `val`: 可选的内联音频数据，以 base64 编码形式提供。这是音频文件的实际数据。
* `ref`: 可选的外部音频引用，可以是音频文件的 URL 或文件路径。`val` 和 `ref` 两者之一必须存在，但不能同时存在。
* `preview`: 预览图像的 base64 编码数组，用于生成视觉预览。每个字节代表一个振幅条。
* `duration`: 音频记录的持续时间，以毫秒为单位。
* `name`: 原始文件的可选名称。
* `size`: 文件的可选大小，以字节为单位。

这些字段描述了音频实体的各个方面，例如其数据、类型、持续时间、大小等。这些信息可以用于在客户端显示音频消息，并提供相应的控制和交互功能。



要创建一个仅包含单个音频记录而没有文本的消息，可以使用以下 Drafty 格式：
```js
{
  "txt": " ",
  "fmt": [{ "len": 1 }],
  "ent": [{
    "tp": "AU",
    "data": {
      "mime": "audio/aac",
      "val": "Rt53jUU...iVBORw0KGgoA==",
      "ref": "/v0/file/s/e769gvt1ILE.m4v",
      "preview": "Aw4JKBkAAAAKMSM...vHxgcJhsgESAY",
      "duration": 180000,
      "name": "ding_dong.m4a",
      "size": 595496
    }
  }]
}

```

重要的安全注意事项：`val` 和 `ref` 字段可能包含恶意负载。客户端应该限制 `ref` 字段中的 URL 方案仅限于 `http` 和 `https`。客户端应该只有在将 `val` 字段正确转换为音频后才向用户展示其内容。

这意味着客户端在处理音频实体时应该采取必要的安全措施来防止恶意行为。在处理 `ref` 字段时，客户端应该验证 URL 的方案，确保它是安全的，防止恶意链接导致的安全问题。同时，在使用 `val` 字段的音频数据时，客户端应该进行正确的转换和验证，确保其是有效的音频数据，以防止恶意内容被播放或执行。


#### `BN`: 交互按钮
``BN` 提供了向服务器发送数据的选项，可以是原始服务器或其他服务器。 `data` 包含以下字段：

```js
{
  "tp": "BN",
  "data": {
    "name": "confirmation",
    "act": "url",
    "val": "some-value",
    "ref": "https://www.example.com/path/?foo=bar"
  }
}
```

* `act`: 按钮点击后的操作类型：
  * `pub`: 发送一个 Drafty 格式的 `{pub}` 消息到当前聊天会话中（私聊或者群），其中包含表单数据作为附件：
  ```js
  { "tp":"EX", "data":{ "mime":"application/json", "val": { "seq": 3, "resp": { "confirmation": "some-value" } } } }
  ```
  * `url`: 发送一个 HTTP GET 请求到 `data.ref` 字段中的 URL。以下查询参数将附加到 URL 中：`<name>=<val>`、`uid=<current-user-ID>`、`topic=<topic name>`、`seq=<message sequence ID>`。
  * `note`: 发送一个 `{note}` 消息到当前主题，`what` 设置为 `data`。
  ```js
  { "what": "data", "data": { "mime": "application/json", "val": { "seq": 3, "resp": { "confirmation": "some-value" } } } }
  ```
* `name`: 按钮的可选名称，将报告给服务器。
* `val`: 附加的不透明数据。
* `ref`: `url` 操作的 URL。

如果提供了 `name` 但未提供 `val`，则假定 `val` 为 `1`。如果未定义 `name`，则不会发送 `name` 或 `val`。

上面示例中的按钮将向 
```
https://www.example.com/path/?foo=bar&confirmation=some-value&uid=usrFsk73jYRR&topic=grpnG99YhENiQU&seq=3
```

发送一个 HTTP GET 请求，假设当前用户 ID 为 `usrFsk73jYRR`，主题为 `grpnG99YhENiQU`，带有按钮的消息的序列 ID 为 `3`。

**重要安全注意事项**：客户端应该将 `ref` 字段中的 URL 方案限制为仅允许 `http` 和 `https`。




#### `EX`: 文件附件
`EX` 是一个附件，示例如下：

```js
{
  "tp": "EX",
  "data": {
    "mime", "text/plain",
    "val", "Q3l0aG9uPT0w...PT00LjAuMAo=",
    "ref": "/v0/file/s/abcdef12345.txt",
    "name", "requirements.txt",
    "size": 1234
  }
}
```
- `mime`: 数据类型，例如 'application/octet-stream'。
- `val`: 可选的内联 base64 编码的文件数据。
- `ref`: 可选的外部文件数据引用。`val` 和 `ref` 必须二选一。
- `name`: 原始文件的可选名称。
- `size`: 文件的可选大小（以字节为单位）。

To generate a message with the file attachment shown as a downloadable file, use the following format:
```js
{
  at: -1,
  len: 0,
  key: <EX entity reference>
}
```




#### `IM`: 内联图或者引用图
`IM` 是一个图片；`data` 包括如下字段：

```js
{
  "tp": "IM",
  "data": {
    "mime": "image/png",
    "val": "Rt53jUU...iVBORw0KGgoA==",
    "ref": "/v0/file/s/abcdef12345.jpg",
    "width": 512,
    "height": 512,
    "name": "sample_image.png",
    "size": 123456
  }
}
```
- `mime`: 数据类型，例如 'image/jpeg'。
- `val`: 可选的内联图像数据：base64 编码的图像位。
- `ref`: 可选的外部图像数据引用。`val` 和 `ref` 必须二选一。
- `width`、`height`: 图像的线性尺寸，以像素为单位。
- `name`: 原始文件的可选名称。
- `size`: 文件的可选大小（以字节为单位）。

如果发送一个没有文本的纯图片，格式如下：
```js
{
  txt: " ",
  fmt: [{len: 1}],
  ent: [{tp: "IM", data: {<your image data here>}]}
}
```

_**重要安全注意事项**：`val` 和 `ref` 字段可能包含恶意有效负载。客户端应该限制 `ref` 字段中的 URL 方案仅限于 `http` 和 `https`。客户端应该仅当 `val` 字段正确转换为图像时才向用户呈现其内容。

#### `LN`: 链接 (URL)

`LN` 表示一个URL链接。`data` 包含一个字段 `url`，其值为URL链接的地址，例如 `"https://www.example.com/abc#fragment"`。

`MN` 表示提及某人，比如 `@alice`。`data` 包含一个字段 `val`，其值为被提及用户的ID，例如 
```js
{ "tp":"MN", "data":{ "val":"usrFsk73jYRR" } }
```

`HT` 表示标签，例如 `#tinode`。

`data` 包含一个字段 `val`，其值为标签的值，例如 
```js
{ "tp":"HT", "data":{ "val":"tinode" } }
```


#### `VC`: 视频电话控制信息
视频通话 `data` 包含当前的视频电话的状态和时长等:
```js
{
  "tp": "VC",
  "data": {
    "duration": 10000,
    "state": "disconnected",
    "incoming": false,
    "aonly": true
  }
}
```

`VC` 表示视频通话控制消息。`data` 包含以下字段：

- `duration`：通话持续时间，单位为毫秒。
- `state`：当前通话状态；支持的状态有：
  - `accepted`：通话已建立（进行中）。
  - `busy`：由于被呼叫方已经在另一个通话中，因此无法建立通话。
  - `finished`：之前建立的通话已成功结束。
  - `disconnected`：通话被中断，例如因为错误。
  - `missed`：通话未接听，即被呼叫方未接听电话。
  - `declined`：通话被拒绝，即被呼叫方在接听之前挂断电话。
- `incoming`：如果通话是呼入的，则为 true；否则为 false。
- `aonly`：如果这是一个仅音频通话（无视频），则为 true。

`VC` 也可以表示为格式 `"fmt": [{"len": 1, "tp": "VC"}]`，在这种情况下，所有通话信息都包含在外部消息的 `head` 字段中。



#### `VD`: 带预览的视频
``VD` 表示视频录制。`data` 包含以下字段：

```js
{
  "tp": "VD",
  "data": {
    "mime": "video/webm",
    "ref": "/v0/file/s/abcdef12345.webm",
    "preview": "AsTrsU...k86n00Ggo=="
    "preref": "/v0/file/s/abcdef54321.jpeg",
    "premime": "image/jpeg",
    "width": 640,
    "height": 360,
    "duration": 32000,
    "name": " bigbuckbunny.webm",
    "size": 1234567
  }
}
```
- `mime`: 视频数据类型，例如 'video/webm'。
- `val`: 可选的内联视频数据：base64编码的视频数据，通常不存在（null）。
- `ref`: 可选的外部视频数据的引用。`val` 或 `ref` 必须至少存在一个。
- `preview`: 可选的视频截图的 base64 编码图像（封面图）。
- `preref`: 可选的视频截图的外部引用（封面图）。
- `premime`: 可选的视频截图的数据类型（封面图）；如果缺少，默认为'image/jpeg'。
- `width`、`height`: 视频和封面图的线性尺寸（像素）。
- `duration`: 视频的时长（毫秒）。
- `name`: 原始文件的可选名称。
- `size`: 文件的可选大小（字节）。

一个没有文本的视频消息如下：
```js
{
  txt: " ",
  fmt: [{len: 1}],
  ent: [{tp: "VD", data: {<your video data here>}]}
}
```


