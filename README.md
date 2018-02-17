# Overview 

[![Build Status](https://travis-ci.org/tinyhubs/tinydom.svg?branch=master)](https://travis-ci.org/tinyhubs/tinydom)
[![GoDoc](https://godoc.org/github.com/tinyhubs/tinydom?status.svg)](https://godoc.org/github.com/tinyhubs/tinydom)
[![Language](https://img.shields.io/badge/language-go-lightgrey.svg)](https://github.com/tinyhubs/tinydom)
[![License](https://img.shields.io/badge/license-New%20BSD-yellow.svg?style=flat)](LICENSE)
[![codecov](https://codecov.io/gh/tinyhubs/tinydom/branch/master/graph/badge.svg)](https://codecov.io/gh/tinyhubs/tinydom)
[![goreport](https://www.goreportcard.com/badge/github.com/tinyhubs/tinydom)](https://www.goreportcard.com/report/github.com/tinyhubs/tinydom)

tinydom是一个非验证的，轻量级的，经过充分测试的go语言(golang)xml流的dom构造器。

# tinydom简介

tidydom使用golang的encoding/xml标准库作为底层XML文本流的解析器。使用tinydom提供的接口可以实现简单的XML文件的读取和生成。
tinydom借鉴了[tinyxml2](http://www.grinninglizard.com/tinyxml2/index.html)的接口设计技巧，提供了丰富的XML元素的查找手段。



# 如何使用
# 接口定义
一个XML文档由`XMLDocument`、`XMLElement`、`XMLText`、`XMLComment`、`XMLProcInst`、`XMLDirective`这几种类型的节点组成。

- `XMLDocument`是一个XML文档的根节点。
- `XMLElement`是XML文档的基本节点元素，一个`XMLElement`可以含有多个`XMLAttribute`。
- `XMLText`是XML的文本元素，支持CDATA和XML字符转义。
- `XMLComment`表示的是XML的注释，是`<!--` 与 `-->`之间的部分。
- `XMLProcInst`表示的是`<?`与`?>`之间的部分，一般出现在xml文档的声明部分。
- `XMLDirective`表示的是`<!`与`>`之间的部分，一般为DTD。
- `XNLNode`是所有这些节点的共同基础，XMLNode提供了丰富的节点元素遍历手段。
- `XMLVisitor`提供了一种XML对象的元素遍历机制。
- `XMLHandle`的作用是简化代码编写工作，使用`XMLHandle`将减少很多判空处理的代码(if nil == xxx {}),活用`XMLHandle`可以让我们的编码工作事半功倍，代码也更加健壮。

##  加载文档
`tinydom.LoadDocument`用于从一个文件流或者字符流读取XML数据，并构建出`tinydom.XMLDocument`对象，一般用于读取XML文件的场景。

```go
import "tinydom"
doc, err := tinydom.LoadDocument(strings.NewReader(s))
```

`FirstChildElement`、`LastChildElement`、`PrevElement`、`NextElement`这几个函数，主要是为了方便查找`XMLElement`元素，
大部分情况下我们建立XML文档的DOM模型就是为了对XMLElement进行访问。

```go
xmlstr := `
<books>
    <book><name>The Moon</name><author>Tom</author></book>
    <book><name>Go west</name><author>Suny</author></book>
<books>
`
doc, _ := tinydom.LoadDocument(strings.NewReader(xmlstr))
elem1 := doc.FirstChildElement("books").FirstChildElement("book").FirstChildElement("name")
fmt.Println(elem1.Text()) //	The Moon
elem2 := doc.FirstChildElement("books").FirstChildElement("book").LastChildElement("author")
fmt.Println(elem2.Text()) //	Suny
```


##  查找节点

- 获取子节点

`FirstChild`、`LastChild`这两个函数用于获取某个节点的子节点,不管节点是什么类型都可以用这两个函数获取到.

`FirstChildElement`、`LastChildElement`这两个函数专用于查找指定的Element节点,如果指定了name参数,那么只查找指定名字的子Element节点.

- 获取相邻节点

`Prev`、`Next`这两个函数用于查找当前节点的前一个或者后一个兄弟节点.

`PrevElement`、`NextElement`这两个函数用于查找当前节点的上一个或者下一个Element节点.我们同样可以通过指定这两个函数的name参数来查找执行名字的Element节点.

- 获取父节点

`Parent`函数直接获取当前节点的父节点.

- 节点类型转换

我们提供了一堆的转换函数:

`ToElement() XMLElement`

`ToText() XMLText`

`ToComment() XMLComment`

`ToDocument() XMLDocument`

`ToProcInst() XMLProcInst`

`ToDirective() XMLDirective`

- 获取节点属性

`XMLElement`接口提供了多个获取属性的函数:

查找: `FindAttribute(name string) XMLAttribute`

遍历: `ForeachAttribute(callback func(attribute XMLAttribute) int) int`

属性个数统计: `AttributeCount() int`

直接获取属性字符串: `Attribute(name string, def string) string`


##  文档的遍历
`Parent`、`FirstChild`、`LastChild`、`Prev`、`Next`用于使我们可以方便地在XML的DOM树中游走。
下面这个函数可以用于对一个doc进行遍历：

```go
func walk(m int , rootNode tinydom.XMLNode) {
    if nil == rootNode {
        return
    }
    for child := rootNode.FirstChild(); nil != child; child = child.Next() {
        fmt.Println(strings.Repeat(" ", m), child.Value())
        walk(m + 1, child)
    }
}
```

您可以这样调用：

```go
walk(0, doc)
```

还有一个更好的替代方式是使用`tinydom.XMLVisitor`接口对文档中的元素进行遍历，可参见代码中`tinydom.XMLVisitor`的接口定义。

##  新建文档
`tinydom.NewDocument`用于在内存中生成DOM，一般用于生成XML文件。

tinydom提供了一系列的NewXXX方法用于创建各种不同类型的节点:

`tinydom.NewText(document XMLDocument, text string) XMLText`

`tinydom.NewComment(document XMLDocument, comment string) XMLComment`

`tinydom.NewElement(document XMLDocument, name string) XMLElement`

`tinydom.NewProcInst(document XMLDocument, target string, inst string) XMLProcInst`

`tinydom.NewDirective(document XMLDocument, directive string) XMLDirective`

而下面这些函数用于将任意类型的节点加入当前节点,或者对节点进行删除操作:

- 将node添加为本节点的最后一个子节点(最常用):`InsertEndChild(node XMLNode) XMLNode`

- 将node添加为本节点的第一个子节点:`InsertFirstChild(node XMLNode) XMLNode`

- 将addThis添加到本节点的后面:`InsertBack(addThis XMLNode) XMLNode`

- 将addThis添加到本节点的前面:`InsertFront(addThis XMLNode) XMLNode`

- 删除本节点所有的子节点:`DeleteChildren()`

- 删除本节点的指定的子节点:`DeleteChild(node XMLNode)`

- 将本节点从其所属的document中拆除:`Split()`

我们也可以对节点的属性进行操作:

- 新增属性或者修改属性值: `SetAttribute(name string, value string) XMLAttribute`

- 删除属性: `DeleteAttribute(name string) XMLAttribute`

- 删除所有属性: `ClearAttributes()`

下面的代码创建了一个XML文档：

```go
doc := tinydom.NewDocument()
books := doc.InsertEndChild(tinydom.NewElement(doc, "books"))
book := books.InsertEndChild(tinydom.NewElement(doc, "book"))
name := book.InsertEndChild(tinydom.NewElement(doc, "name"))
name.InsertEndChild(tinydom.NewText(doc, "The Moon"))
doc.InsertEndChild(tinydom.NewProcInst(doc, "xml", `version="1.0" encoding="UTF-8"`))
```

我们可以使用`tinydom.XMLDocument`的`Accept`方法来将这个XML文档输出：

```go
doc.Accept(tinydom.NewSimplePrinter(os.Stdout, tinydom.PrettyPrint))
```

##  输出

tinydom采用了访问者模式(参见`tinydom.XMLVisitor`接口)来对文档的所有节点进行遍历,`tinydom.XMLVisitor`和`tinydom.XMLDocument`的
`Accept`方法结合基本可以输出满足我们大多数场景的XML文档输出任务.我们完全可以使用该机制自己定制文档输出格式.

不过,为了方便大多数使用场景,tinydom仍然提供了一个专用于打印的visitor.下面这行代码用于直接向屏幕打印XML文档:

```go
doc.Accept(tinydom.NewSimplePrinter(os.Stdout, tinydom.PrettyPrint))
```

`tinydom.NewSimplePrinter`的接口如下:

```go
func NewSimplePrinter(writer io.Writer, options PrintOptions) XMLVisitor
```

`tinydom.NewSimplePrinter`的第二个参数用于控制输出格式:

```go
type PrintOptions struct {
    Indent        []byte //  缩进前缀,只允许填写tab或者空白,如果Indent长度为0表示折行但是不缩进,如果Indent为null表示不折行
    TextWrapWidth int    //  超过多长才强制换行
}
```

为简化编码tinydom也提供了两种缺省的`PrintOptions`:

- `tinydom.PrintPretty` 优美打印: 节点输出自动折行,并按4个空格缩进
- `tinydom.PrintStream` 流式打印: 节点输出不带换行,除非Text部分有换行

对于自定义XML文档输出模式而言,处理XML字符转义是个麻烦,因为你必须处理一些细节.但tinydom也可在这方面帮助你.tinydom提供了
`tinydom.EscapeAttribute`和`tinydom.EscapeText`来方便处理属性和`XMLText`中的转义字符.您也可以使用golang自带
的`xml.EscapeText`,只是这个函数做了更多的转义,会导致文档更难阅读和编辑.


##  XML字符转义
受益于go的xml库，tinydom也支持XML字符转义，使用tinydom在读写xml的数据的时候不需要关注XML转义字符，tinydom自动会处理好，可参考下面的例子：

```go
xmlstr :=
    `<talks>
        <talk from="bill" to="tom">[&amp;&apos;&quot;&gt;&lt;] are the xml escape chars? </talk>
        <talk from="tom" to="bill">yes， that is right</talk>
     </talks>
    `
doc, _ := tinydom.LoadDocument(strings.NewReader(xmlstr))
talk := doc.FirstChildElement("talks").FirstChildElement("talk").Text()
fmt.Print(talk) //  [&'"><] are the xml escape chars?
```

xml文档输出时,可使用`tinydom.EscapeAttribute`和`tinydom.EscapeText`来对字符进行转义.

##  CDATA
只有XMLText对象才涉及到CDATA，可以通过XMLText获取到CDATA对象的数据，tinydom能够自动识别CDATA，但是将DOM对象序列化成字符串时，除非节点指定了CDATA属性，否则会直接转义。

```go
xmlstr := `<content><![CDATA[<example>This is ok in cdata text</example>]]></content>`
doc, _ := tinydom.LoadDocument(strings.NewReader(xmlstr))
content := doc.FirstChildElement("content")
fmt.Println("\nRead CDATA:", content.Text())
fmt.Println("\nNormal Print:")
doc.Accept(tinydom.NewSimplePrinter(os.Stdout))
text := content.FirstChild().ToText()
text.SetCDATA(true)
fmt.Println("\nSpecial as CDATA:")
doc.Accept(tinydom.NewSimplePrinter(os.Stdout))
```

##  名字空间
不支持：
虽然golang标准库是能够正常处理名字空间的，但当前tinydom还无法正确处理xml的名字空间，所有带有名字空间前缀的节点或者属性都会被丢弃。后续计划将这块功能补齐。


##  BOM
golang的xml解析器自身还不支持BOM，所以本解析器还无法解析带BOM头的xml文件。

## Changelog

#### 1.0.0 初始版本

#### 1.1.0 小版本改进,能力增强,bug解决

- 文档输出增加打印选项控制,支持"优美打印" <font color="red">`接口变更`</font> `NewSimplePrinter`
- 优化字符转义切新增转义处理的接口 `tinydom.ExcapeText` `tinydom.ExcapeAttribute`
- 简化接口 `Previous` 缩写为 `Prev`,`Sibling`单词从所有接口中删除
- 删除接口 `InsertAfterChild` 建议使用 `InsertBack` 或者 `InsertElementBack` 代替
- 增加接口 `Split`,`InsertFront`,`InsertBack`
- 增加版本识别函数 `tinydom.Version`
- 解决用例稳定性问题
- 取消支持go-1.4.x版本 `因为该版本没有覆盖率统计工具`
- 完善文档
- github仓库增加了构建服务,文档服务
- 补充用例,增加代码覆盖率到90%以上

#### 1.2.0 小版本改进,能力增强,bug解决

- 只是直接指定文件名加载文档或者保存文档:因为这两种场景也比较常见,比较实用
- 发现1.1.0版本新增函数引入的bug,可能导致对象被切片
- 将两个全局量PreetyPrint改名为Print在前:因为发现Print在前更容易记忆
- 完善文档

