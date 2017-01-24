# tinydom

go语言的xml流的dom解析器。

# tinydom简介
tinydom	实现了一个简单的XML的DOM模型.

tidydom使用encoding/xml作为底层XML解析库，实现对XML文件的解析.使用tinydom提供的接口可以实现简单的XML文件的读取和生成。
tinydom借鉴了[tinyxml2](http://www.grinninglizard.com/tinyxml2/index.html)的接口设计技巧，提供了丰富的查找XML元素的查找手段。

# 接口设计
一个XML文档由`XMLDocument`、`XMLElement`、`XMLText`、`XMLComment`、`XMLProcInst`、`XMLDirective`者几种类型的节点组成。
`XMLDocument`是一个XML文档的根节点。
`XMLElement`是XML文档的基本节点元素，一个XMLElement可以含有多个XMLAttribute。
`XMLText`是XML的文本元素，支持CDATA和XML字符转义。
`XMLComment`表示的是XML的注释，是`<!--` 与 `-->`之间的部分。
`XMLProcInst`表示的是`<?`与`?>`之间的部分，一般出现在xml文档的声明部分。
`XMLDirective`表示的是`<!`与`>`之间的部分，一般为DTD。
`XNLNode`是所有这些节点的共同基础，XMLNode提供了丰富的节点元素遍历手段。
`XMLVisitor`提供了一种XML对象的元素遍历机制。
`XMLHandle`的所用是简化代码编写工作，使用XMLHandle将减少很多判空代码(if nil == xxx {}),活用XMLHandle将会让XML文件的元素事半功倍。

# 如何使用
##  加载文档
LoadDocument用于从一个文件流或者字符流读取XML数据，并构建出XMLDocument对象，一般用于读取XML文件的场景。
```go
  import "tinydom"
  doc, err := tinydom.LoadDocument(strings.NewReader(s))
```
从文档中找到我们需要的元素：
FirstChildElement、LastChildElement、PreviousSiblingElement、NextSiblingElement这些接口主要是为了方便查找XMLElement元素，
大部分情况下我们建立XML文档的DMO模型就是为了对XMLElement进行访问。
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

##  新建文档
NewDocument用于在内存中生成DOM，一般用于生成XML文件。
InsertEndChild、InsertFirstChild、InsertAfterChild、DeleteChildren、DeleteChild用于对XMLDocument进行修改。
下面的代码创建了一个XML文档：
```go
    doc := tinydom.NewDocument()
    books := doc.InsertEndChild(tinydom.NewElement(doc, "books"))
    book := books.InsertEndChild(tinydom.NewElement(doc, "book"))
    name := book.InsertEndChild(tinydom.NewElement(doc, "name"))
    name.InsertEndChild(tinydom.NewText(doc, "The Moon"))
    doc.InsertEndChild(tinydom.NewProcInst(doc, "xml", `version="1.0" encoding="UTF-8"`))
```

我们可以使用XMLDocument.Accept方法来将这个XML文档输出：
```go
    doc.Accept(tinydom.NewSimplePrinter(os.Stdout))
```

##  文档的遍历
`Parent`、`FirstChild`、`LastChild`、`PreviousSibling`、`NextSibling`用于使我们可以方便地在XML的DOM树中游走。
下面这个函数可以用于对一个doc进行遍历：
```go
    func walk(m int , rootNode tinydom.XMLNode) {
        if nil == rootNode {
            return
        }
        for child := rootNode.FirstChild(); nil != child; child = child.NextSibling() {
            fmt.Println(strings.Repeat(" ", m), child.Value())
            walk(m + 1, child)
        }
    }
```
您可以这样调用：
```go
walk(doc)。
```
还有一个更好的替代方式是使用XMLVisitor接口对文档中的元素进行遍历，可参见代码中XMLHandle的接口定义。

##  XML字符转义
受益于go的xml库，tinydom也支持XML字符转义，使用tinydom在读写xml的数据的时候不需要关注XML转义字符，tinydom自动会处理好，可参考下面的例子。
如果您需要自定义输出格式，那么文本雷荣时，需要通过xml.ExcapeText函数进行转义。
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

##  CDATA
只有XMLText对象才涉及到CDATA，可以通过XMLText，tinydom能够自动识别CDATA，但是将DOM对象序列化成字符串时，除非节点指定了CDATA属性，否则会直接转义。
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
