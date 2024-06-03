

## 1. android文本使用span来描述

在Android开发中，`TextView` 是用于显示文本内容的控件，而 `span` 则是一种机制，用于在 `TextView` 中对文本的某些部分应用特定的样式或行为。`Spannable` 是一个接口，它扩展了 `CharSequence` 接口，允许对文本内容中的部分文本进行样式处理。通过使用 `span`，开发者可以在 `TextView` 中实现富文本显示。

### 1.1 常见的 Span 类型

以下是一些常见的 `Span` 类型及其用途：

1. **`ForegroundColorSpan`**：改变文本的前景色（文本颜色）。
    ```java
    SpannableString spannableString = new SpannableString("Hello, World!");
    spannableString.setSpan(new ForegroundColorSpan(Color.RED), 0, 5, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
    textView.setText(spannableString);
    ```

2. **`BackgroundColorSpan`**：改变文本的背景色。
    ```java
    SpannableString spannableString = new SpannableString("Hello, World!");
    spannableString.setSpan(new BackgroundColorSpan(Color.YELLOW), 0, 5, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
    textView.setText(spannableString);
    ```

3. **`StyleSpan`**：应用文本样式，例如粗体（Typeface.BOLD）或斜体（Typeface.ITALIC）。
    ```java
    SpannableString spannableString = new SpannableString("Hello, World!");
    spannableString.setSpan(new StyleSpan(Typeface.BOLD), 0, 5, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
    textView.setText(spannableString);
    ```

4. **`UnderlineSpan`**：为文本添加下划线。
    ```java
    SpannableString spannableString = new SpannableString("Hello, World!");
    spannableString.setSpan(new UnderlineSpan(), 0, 5, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
    textView.setText(spannableString);
    ```

5. **`StrikethroughSpan`**：为文本添加删除线。
    ```java
    SpannableString spannableString = new SpannableString("Hello, World!");
    spannableString.setSpan(new StrikethroughSpan(), 0, 5, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
    textView.setText(spannableString);
    ```

6. **`ClickableSpan`**：使文本可点击，并定义点击行为。
    ```java
    SpannableString spannableString = new SpannableString("Click here to learn more!");
    spannableString.setSpan(new ClickableSpan() {
        @Override
        public void onClick(View widget) {
            // Handle click event
        }
    }, 0, 10, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
    textView.setText(spannableString);
    textView.setMovementMethod(LinkMovementMethod.getInstance()); // Make links clickable
    ```

7. **`ImageSpan`**：在文本中嵌入图像。
    ```java
    SpannableString spannableString = new SpannableString("Image: ");
    Drawable drawable = ContextCompat.getDrawable(context, R.drawable.image);
    drawable.setBounds(0, 0, drawable.getIntrinsicWidth(), drawable.getIntrinsicHeight());
    ImageSpan imageSpan = new ImageSpan(drawable, ImageSpan.ALIGN_BASELINE);
    spannableString.setSpan(imageSpan, 7, 8, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
    textView.setText(spannableString);
    ```

### 1.2 使用 Span 的基本步骤

1. **创建一个 `Spannable` 对象**：可以使用 `SpannableString` 或 `SpannableStringBuilder` 来包裹原始字符串。
2. **设置 Span**：使用 `setSpan` 方法将 `Span` 应用到 `Spannable` 对象的指定范围内。
3. **应用到 `TextView`**：将 `Spannable` 对象设置为 `TextView` 的文本。

### 1.3 代码示例

```java
TextView textView = findViewById(R.id.textView);
SpannableString spannableString = new SpannableString("Hello, World!");

// 设置前景色
ForegroundColorSpan colorSpan = new ForegroundColorSpan(Color.RED);
spannableString.setSpan(colorSpan, 0, 5, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);

// 设置背景色
BackgroundColorSpan backgroundColorSpan = new BackgroundColorSpan(Color.YELLOW);
spannableString.setSpan(backgroundColorSpan, 7, 12, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);

// 设置点击事件
ClickableSpan clickableSpan = new ClickableSpan() {
    @Override
    public void onClick(View widget) {
        Toast.makeText(widget.getContext(), "Clicked!", Toast.LENGTH_SHORT).show();
    }
};
spannableString.setSpan(clickableSpan, 7, 12, Spanned.SPAN_EXCLUSIVE_EXCLUSIVE);
textView.setMovementMethod(LinkMovementMethod.getInstance());

textView.setText(spannableString);
```

### 1.4 结论

通过使用 `Span`，开发者可以在 `TextView` 中创建丰富多样的文本样式和交互行为，从而提升应用的用户体验。



## 2. ios中文本使用span

在iOS开发中，类似于Android中的`Span`机制，`NSAttributedString`和`NSMutableAttributedString`提供了对文本的富文本处理功能。通过这些类，开发者可以对文本的不同部分应用特定的样式和属性，从而实现类似的效果。

### 2.1常见的属性和使用方式

以下是一些常见的属性以及在iOS中使用`NSAttributedString`和`NSMutableAttributedString`进行文本样式处理的示例：

1. **字体和字体大小**：使用`NSFontAttributeName`属性。
    ```swift
    let text = "Hello, World!"
    let attributedString = NSMutableAttributedString(string: text)
    attributedString.addAttribute(.font, value: UIFont.boldSystemFont(ofSize: 18), range: NSRange(location: 0, length: 5))
    ```

2. **文本颜色**：使用`NSForegroundColorAttributeName`属性。
    ```swift
    attributedString.addAttribute(.foregroundColor, value: UIColor.red, range: NSRange(location: 0, length: 5))
    ```

3. **背景颜色**：使用`NSBackgroundColorAttributeName`属性。
    ```swift
    attributedString.addAttribute(.backgroundColor, value: UIColor.yellow, range: NSRange(location: 7, length: 5))
    ```

4. **下划线**：使用`NSUnderlineStyleAttributeName`属性。
    ```swift
    attributedString.addAttribute(.underlineStyle, value: NSUnderlineStyle.single.rawValue, range: NSRange(location: 0, length: 5))
    ```

5. **删除线**：使用`NSStrikethroughStyleAttributeName`属性。
    ```swift
    attributedString.addAttribute(.strikethroughStyle, value: NSUnderlineStyle.single.rawValue, range: NSRange(location: 7, length: 5))
    ```

6. **超链接**：使用`NSLinkAttributeName`属性。
    ```swift
    attributedString.addAttribute(.link, value: URL(string: "https://www.example.com")!, range: NSRange(location: 7, length: 5))
    ```

7. **图像**：使用`NSTextAttachment`插入图像。
    ```swift
    let attachment = NSTextAttachment()
    attachment.image = UIImage(named: "image.png")
    let attachmentString = NSAttributedString(attachment: attachment)
    attributedString.append(attachmentString)
    ```

### 2.2 使用示例

以下是一个完整的示例，展示如何在`UILabel`中使用富文本：

```swift
import UIKit

class ViewController: UIViewController {
    override func viewDidLoad() {
        super.viewDidLoad()

        let label = UILabel()
        label.frame = CGRect(x: 20, y: 50, width: 300, height: 100)
        label.numberOfLines = 0
        self.view.addSubview(label)

        let text = "Hello, World!"
        let attributedString = NSMutableAttributedString(string: text)

        // 设置字体和颜色
        attributedString.addAttribute(.font, value: UIFont.boldSystemFont(ofSize: 18), range: NSRange(location: 0, length: 5))
        attributedString.addAttribute(.foregroundColor, value: UIColor.red, range: NSRange(location: 0, length: 5))

        // 设置背景色
        attributedString.addAttribute(.backgroundColor, value: UIColor.yellow, range: NSRange(location: 7, length: 5))

        // 添加下划线
        attributedString.addAttribute(.underlineStyle, value: NSUnderlineStyle.single.rawValue, range: NSRange(location: 0, length: 5))

        // 添加删除线
        attributedString.addAttribute(.strikethroughStyle, value: NSUnderlineStyle.single.rawValue, range: NSRange(location: 7, length: 5))

        // 添加点击链接（需启用交互和添加手势识别器）
        attributedString.addAttribute(.link, value: URL(string: "https://www.example.com")!, range: NSRange(location: 7, length: 5))

        // 添加图像
        let attachment = NSTextAttachment()
        attachment.image = UIImage(named: "image.png")
        let attachmentString = NSAttributedString(attachment: attachment)
        attributedString.append(attachmentString)

        label.attributedText = attributedString

        // 使链接可点击
        label.isUserInteractionEnabled = true
        label.addGestureRecognizer(UITapGestureRecognizer(target: self, action: #selector(handleTapOnLabel(_:))))
    }

    @objc func handleTapOnLabel(_ recognizer: UITapGestureRecognizer) {
        // 处理点击事件
        if let text = (recognizer.view as? UILabel)?.attributedText?.string {
            print("Tapped on label with text: \(text)")
        }
    }
}
```

### 2.3 结论

在iOS中，使用`NSAttributedString`和`NSMutableAttributedString`可以实现类似于Android中`Span`的富文本处理功能。通过这些类，开发者可以在`UILabel`、`UITextView`等控件中应用多种文本样式，创建丰富的用户界面效果。



## 3. drafty结构分析

### 3.1 示例

> this is **bold**, `code` and _italic_, ~~strike~~<br/>
> combined **bold and _italic_**<br/>
> an url: https://www.example.com/abc#fragment and another _[https://web.tinode.co](https://web.tinode.co)_<br/>
> this is a [@mention](#) and a [#hashtag](#) in a string<br/>
> second [#hashtag](#)<br/>

文档:
```js
{
	 "txt": "this is bold, code and italic, strike combined bold and italic an url: https://www.example.com/abc#fragment and another www.tinode.co this is a @mention and a #hashtag in a string second #hashtag,
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
      "ent":[
       { "tp":"LN", "data":{ "url":"https://www.example.com/abc#fragment" } },
       { "tp":"LN", "data":{ "url":"http://www.tinode.co" } },
       { "tp":"MN", "data":{ "val":"mention" } },
       { "tp":"HT", "data":{ "val":"hashtag" } }
      ]
  
}
```

文档的结构大概如下：

![](/image/drafty结构.png)

- fmt是一个列表；
- ent也是一个列表；
- fmt通过序号引用ent；
- fmt有嵌套关系，所以构成父子节点；



### 3.2 使用GO代码解析

我们用几个结构来描述这个文本：

文档中有三个部分：`txt`, `fmt`,  `ent`,

- 一个是文本的内容，支持多语言，所以`document`需要`txt []rune`来描述；
- `fmt`是一个格式列表，也就是描述格式的一个个`span`；同时有些`fmt`元素并不是格式，而是附件或者文件时，会应用后面的ent内容；
- `ent`数组用于描述各种附件或者链接；事实上是通过`fmt`数组的元素来引用的；

```go
type document struct {
	Txt string `json:"txt,omitempty"`
	txt []rune   // 支持多语言，使用字符数组来描述
	Fmt []style  `json:"fmt,omitempty"`
	Ent []entity `json:"ent,omitempty"`
}

type style struct {
	Tp     string `json:"tp,omitempty"`
	At     int    `json:"at,omitempty"`
	Length int    `json:"len,omitempty"`
	Key    int    `json:"key,omitempty"`
}

type entity struct {
	Tp   string                 `json:"tp,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}
```

**这三个机构主要的作用就是反序列化；**

这个文档在解析前就是一段`string`，那么就是一个无格式的文本；

如果这个文档使用`json`反序列化之后就是一个`map[string]interface{}`



使用上述的结构解析之后，fmt也可以转为对应的span，这些span 是一个一个罗列的小片段，

span之间是不能交叉的重叠，因为这样格式比较混乱；

span直接可以完全的从属，比如：

> combined **bold and _italic_**<br/>
>
> 这里就有一个嵌套，对应就是：
```js
{ "at":56, "len":6, "tp":"EM" },  // 第56字符开始，6个字符长度为斜体：italic
{ "at":47, "len":15, "tp":"ST" }, // 第47字符开始，15个字符长度为粗体： bold and _italic
{ "tp":"BR", "len":1, "at":62 },  // 换行
```

那么这里的span 不仅仅是并列的关系，而且还是一个树状的关系，这里**斜体部分**就是**粗体**的一个**子节点**;



那么需要一个`span`描述`fmt`内容以及`fmt`引用了`ent`中的内容，所以需要添加一个`data`元素；

但是为了描述各个span之间的这种嵌套关系，所以就是需要一个树状结构，这个树状结构`node`只是为了描述span之间的关系；

而`toTree`解析之后返回的的根节点就是包含了一级的各个`span`：这个toTree函数是一个深度递归的调用；

```go
type span struct {
	tp   string
	at   int
	end  int
	key  int
	data map[string]interface{}
}

type node struct {
	txt      []rune
	sp       *span
	children []*node
}
```

最后使用了一个`previewFormatter`函数执行格式化，因为GO格式化后的内容，毕竟不是在android 和 ios 中可以显示的控件中使用，所以和其他版本的格式化不太一样；

这个函数简单的总结就是一个递归函数，对其中的文本进行格式化处理，不过GO确实没有啥可处理的。