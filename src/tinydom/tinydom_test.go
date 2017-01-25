package tinydom_test

import (
    "fmt"
    "os"
    "strings"
    "tinydom"
    "testing"
)

func expect(t *testing.T, message string, result bool) {
    if result {
        fmt.Println(message)
        return
    }
    
    t.Fail()
}

func Test_example1(t *testing.T) {
    xmlstr := `
    <books>
        <book><name>The Moon</name><author>Tom</author></book>
        <book><name>Go west</name><author>Suny</author></book>
    </books>
    `
    doc, _ := tinydom.LoadDocument(strings.NewReader(xmlstr))
    elem1 := doc.FirstChildElement("books").FirstChildElement("book").FirstChildElement("name")
    fmt.Println(elem1.Text()) //	The Moon
    
    elem2 := doc.FirstChildElement("books").FirstChildElement("book").LastChildElement("author")
    fmt.Println(elem2.Text()) //	Suny
    
}

func walk(m int, rootNode tinydom.XMLNode) {
    if nil == rootNode {
        return
    }
    
    space := strings.Repeat("  ", m)
    for child := rootNode.FirstChild(); nil != child; child = child.NextSibling() {
        fmt.Println(space, child.Value())
        walk(m+1, child)
    }
}

func Test_example2(t *testing.T) {
    doc := tinydom.NewDocument()
    doc.InsertEndChild(tinydom.NewProcInst(doc, "xml", `version="1.0" encoding="UTF-8"`))
    books := doc.InsertEndChild(tinydom.NewElement(doc, "books"))
    book := books.InsertEndChild(tinydom.NewElement(doc, "book"))
    name := book.InsertEndChild(tinydom.NewElement(doc, "name"))
    name.InsertEndChild(tinydom.NewText(doc, "The Moon"))
    
    doc.Accept(tinydom.NewSimplePrinter(os.Stdout))
    
    fmt.Println()
    
    walk(0, doc)
}

func Test_example3(t *testing.T) { //Me
    xmlstr :=
        `<talks>
            <talk from="bill" to="tom">[&amp;&apos;&quot;&gt;&lt;] are the xml escape chars? </talk>
            <talk from="tom" to="bill">yes, that is right</talk>
         </talks>
        `
    doc, _ := tinydom.LoadDocument(strings.NewReader(xmlstr))
    talk := doc.FirstChildElement("talks").FirstChildElement("talk").Text()
    fmt.Print(talk)
}

func Test_example4(t *testing.T) {
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
}

func Test_空文档_加载失败(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader(""))
    expect(t, "返回值检测", nil == doc)
    expect(t, "返回值检测", nil != err)
}

func Test_格式错误_节点未关闭(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node><elem></node>"))
    expect(t, "返回值检测", nil == doc)
    expect(t, "返回值检测", nil != err)
}

func Test_格式错误_关闭节点多余(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node><elem></elem></elem></node>"))
    expect(t, "返回值检测", nil == doc)
    expect(t, "返回值检测", nil != err)
}

func Test_格式错误_多余的节点(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node><elem></elem></node><hello/>"))
    doc.Accept(tinydom.NewSimplePrinter(os.Stdout))
    expect(t, "返回值检测", nil == doc)
    expect(t, "返回值检测", nil != err)
}

func Test_正常的XML文档1(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node></node>"))
    expect(t, "返回值检测", nil != doc)
    expect(t, "返回值检测", nil == err)
    
    node := doc.FirstChild()
    expect(t, "检查节点名:node", node.Value() == "node")
    
    expect(t, "topo结构检查", doc.ToNode() == node.Parent())
    expect(t, "topo结构检查", doc == node.GetDocument())
    
    expect(t, "topo结构检查", nil == node.FirstChild())
    expect(t, "topo结构检查", nil == node.LastChild())
    expect(t, "topo结构检查", nil == node.PreviousSibling())
    expect(t, "topo结构检查", nil == node.NextSibling())
    
    expect(t, "topo结构检查", nil == node.FirstChildElement(""))
    expect(t, "topo结构检查", nil == node.LastChildElement(""))
    expect(t, "topo结构检查", nil == node.PreviousSiblingElement(""))
    expect(t, "topo结构检查", nil == node.NextSiblingElement(""))
    
    expect(t, "转换检查", nil != node.ToElement())
    expect(t, "转换检查", nil == node.ToComment())
    expect(t, "转换检查", nil == node.ToDirective())
    expect(t, "转换检查", nil == node.ToDocument())
    expect(t, "转换检查", nil == node.ToProcInst())
    expect(t, "转换检查", nil == node.ToText())
}

func Test_正常的XML文档2(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node><elem></elem></node>"))
    expect(t, "返回值检测", nil != doc)
    expect(t, "返回值检测", nil == err)
    
    node := doc.FirstChild()
    elem := node.FirstChild()
    
    //  node
    expect(t, "检查节点名:node", "node" == node.Value())
    
    expect(t, "topo结构检查", doc.ToNode() == node.Parent())
    expect(t, "topo结构检查", doc == node.GetDocument())
    
    expect(t, "topo结构检查", elem == node.FirstChild())
    expect(t, "topo结构检查", elem == node.LastChild())
    expect(t, "topo结构检查", nil == node.PreviousSibling())
    expect(t, "topo结构检查", nil == node.NextSibling())
    
    expect(t, "topo结构检查", elem.ToElement() == node.FirstChildElement(""))
    expect(t, "topo结构检查", elem.ToElement() == node.LastChildElement(""))
    expect(t, "topo结构检查", nil == node.PreviousSiblingElement(""))
    expect(t, "topo结构检查", nil == node.NextSiblingElement(""))
    
    expect(t, "topo结构检查", elem.ToElement() == node.FirstChildElement("elem"))
    expect(t, "topo结构检查", elem.ToElement() == node.LastChildElement("elem"))
    expect(t, "topo结构检查", nil == node.PreviousSiblingElement(""))
    expect(t, "topo结构检查", nil == node.NextSiblingElement(""))
    
    //  elem
    expect(t, "检查节点名:elem", "elem" == elem.Value())
    
    expect(t, "topo结构检查", node.ToNode() == elem.Parent())
    expect(t, "topo结构检查", doc == elem.GetDocument())
    
    expect(t, "topo结构检查", nil == elem.FirstChild())
    expect(t, "topo结构检查", nil == elem.LastChild())
    expect(t, "topo结构检查", nil == elem.PreviousSibling())
    expect(t, "topo结构检查", nil == elem.NextSibling())
    
    expect(t, "topo结构检查", nil == elem.FirstChildElement(""))
    expect(t, "topo结构检查", nil == elem.LastChildElement(""))
    expect(t, "topo结构检查", nil == elem.PreviousSiblingElement(""))
    expect(t, "topo结构检查", nil == elem.NextSiblingElement(""))
}

func Test_正常的XML文档3(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node><elem1></elem1><elem2></elem2><elem3></elem3></node>"))
    expect(t, "返回值检测", nil != doc)
    expect(t, "返回值检测", nil == err)
    
    node := doc.FirstChild()
    elem1 := node.FirstChild()
    elem2 := elem1.NextSibling()
    elem3 := elem2.NextSibling()
    
    expect(t, "检查节点顺序关系", "node" == node.Value())
    expect(t, "检查节点顺序关系", "elem1" == elem1.Value())
    expect(t, "检查节点顺序关系", "elem2" == elem2.Value())
    expect(t, "检查节点顺序关系", "elem3" == elem3.Value())
    
    expect(t, "topo结构检查", node.ToNode() == elem1.Parent())
    expect(t, "topo结构检查", node.ToNode() == elem2.Parent())
    expect(t, "topo结构检查", node.ToNode() == elem3.Parent())
    
    expect(t, "topo结构检查", nil == elem1.PreviousSibling())
    expect(t, "topo结构检查", nil == elem3.NextSibling())
    
    expect(t, "topo结构检查", elem1 == elem2.PreviousSibling())
    expect(t, "topo结构检查", elem3 == elem2.NextSibling())
    
    expect(t, "topo结构检查", elem1.ToElement() == elem2.PreviousSiblingElement(""))
    expect(t, "topo结构检查", elem3.ToElement() == elem2.NextSiblingElement(""))
    
    expect(t, "topo结构检查", elem1.ToElement() == elem2.PreviousSiblingElement("elem1"))
    expect(t, "topo结构检查", elem3.ToElement() == elem2.NextSiblingElement("elem3"))
    
    expect(t, "topo结构检查", false == node.NoChildren())
    expect(t, "topo结构检查", true == elem1.NoChildren())
    expect(t, "topo结构检查", true == elem2.NoChildren())
    expect(t, "topo结构检查", true == elem3.NoChildren())
}

func Test_Element_属性同名错误(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader(`<node attr="value1" attr="value2"></node>`))
    expect(t, "返回值检测", nil == doc)
    expect(t, "返回值检测", nil != err)
}

func Test_Element_属性基本功能(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader(`<node attr1="value1" attr2="value2"></node>`))
    expect(t, "返回值检测", nil != doc)
    expect(t, "返回值检测", nil == err)
    
    //  get
    node := doc.FirstChildElement("")
    expect(t, "元素个数", 2 == node.AttributeCount())
    expect(t, "元素的值检测", "value1" == node.Attribute("attr1", ""))
    expect(t, "元素的值检测", "value2" == node.Attribute("attr2", ""))
    expect(t, "不存在的元素", "(default)" == node.Attribute("attr3", "(default)"))
    
    //  modify
    node.SetAttribute("attr1", "(modified1)")
    node.DeleteAttribute("attr2") //  删除已经存在的属性
    node.SetAttribute("attr3", "(add3)")
    node.DeleteAttribute("attr4") //  删除不存在的属性
    expect(t, "元素个数", 2 == node.AttributeCount())
    attr3 := node.FindAttribute("attr3")
    expect(t, "属性名检测", "attr3" == attr3.Name())
    expect(t, "属性值检测", "(add3)" == attr3.Value())
    
    //  enum
    hitCount := 0
    node.ForeachAttribute(func(attribute tinydom.XMLAttribute) int {
        
        if "attr1" == attribute.Name() {
            expect(t, "检查元素值", "(modified1)" == attribute.Value())
            hitCount++
            return 0
        }
        
        if "attr3" == attribute.Name() {
            expect(t, "检查元素值", "(add3)" == attribute.Value())
            hitCount++
            return 0
        }
        
        return 0
    })
    expect(t, "检查属性遍历命中次数", 2 == hitCount)
    
    //  return vaue of callback
    hitCount = 0
    retult := node.ForeachAttribute(func(attribute tinydom.XMLAttribute) int {
        if "attr1" == attribute.Name() {
            expect(t, "检查元素值", "(modified1)" == attribute.Value())
            hitCount++
            return -44
        }
        
        if "attr3" == attribute.Name() {
            expect(t, "检查元素值", "(add3)" == attribute.Value())
            hitCount++
            return -55
        }
        return -66
    })
    expect(t, "提前返回那么只能命中一次", 1 == hitCount)
    expect(t, "提前返回的返回值", (-44 == retult) || (-55 == retult))
    
    //  clean all
    node.ClearAttributes()
    expect(t, "元素个数", 0 == node.AttributeCount())
    expect(t, "清除所有属性之后", "(default1)" == node.Attribute("attr1", "(default1)"))
    expect(t, "清除所有属性之后", "(default2)" == node.Attribute("attr2", "(default2)"))
    expect(t, "清除所有属性之后", "(default3)" == node.Attribute("attr3", "(default3)"))
    expect(t, "遍历空属性列表总是返回0", 0 == node.ForeachAttribute(func(attribute tinydom.XMLAttribute) int {
        return 11
    }))
    expect(t, "在空属性列表中查找，总是返回nil", nil == node.FindAttribute("attr1"))
    
    //  test attr's method
    attr1 := node.SetAttribute("attr1", "value1")
    expect(t, "修改属性值", "value1" == node.Attribute("attr1", "(default1)"))
    attr1.SetValue("NewValue")
    expect(t, "修改属性值", "NewValue" == node.Attribute("attr1", "(default1)"))
}

func Test_ProcInst_基本功能测试(t *testing.T) {
    xml := `<?xml version="1.0" encoding="UTF-8"?>
    <node attr1="value1" attr2="value2"></node>
    `
    doc, err := tinydom.LoadDocument(strings.NewReader(xml))
    expect(t, "返回值检测", nil != doc)
    expect(t, "返回值检测", nil == err)
    
    expect(t, "申明头转换测试", nil != doc.FirstChild().ToProcInst())
    expect(t, "申明头转换测试", nil != doc.FirstChild().ToNode())
    expect(t, "申明头转换测试", nil == doc.FirstChild().ToElement())
    expect(t, "申明头转换测试", nil == doc.FirstChild().ToDocument())
    expect(t, "申明头转换测试", nil == doc.FirstChild().ToDirective())
    expect(t, "申明头转换测试", nil == doc.FirstChild().ToComment())
    expect(t, "申明头转换测试", nil == doc.FirstChild().ToText())
    
    procInst := doc.FirstChild().ToProcInst()
    expect(t, "有申明头的xml文档，第一个子节点是申明头", "xml" == procInst.Value())
    expect(t, "有申明头的xml文档，第一个子节点是申明头", "xml" == procInst.Target())
    expect(t, "有申明头的xml文档，第一个子节点是申明头", `version="1.0" encoding="UTF-8"` == procInst.Instruction())
    expect(t, "申明头下面一个是node", "node" == procInst.NextSibling().Value())
}

func Test_Comment_基本功能测试(t *testing.T) {
    xml := `<!--comment1--><node><elem1><!--comment2--></elem1></node>`
    
    //  加载
    doc, err := tinydom.LoadDocument(strings.NewReader(xml))
    expect(t, "返回值检测", nil != doc)
    expect(t, "返回值检测", nil == err)
    
    //  转换
    expect(t, "转换测试", nil != doc.FirstChild().ToComment())
    expect(t, "转换测试", nil != doc.FirstChild().ToNode())
    expect(t, "转换测试", nil == doc.FirstChild().ToDirective())
    expect(t, "转换测试", nil == doc.FirstChild().ToProcInst())
    expect(t, "转换测试", nil == doc.FirstChild().ToDocument())
    expect(t, "转换测试", nil == doc.FirstChild().ToText())
    expect(t, "转换测试", nil == doc.FirstChild().ToElement())
    
    //  获取注释内容
    comment1 := doc.FirstChild().ToComment()
    comment2 := doc.FirstChildElement("node").FirstChildElement("elem1").FirstChild().ToComment()
    expect(t, "返回值检测", nil != comment1)
    expect(t, "返回值检测", nil != comment2)
    
    //  修改注释
    comment1.SetComment("New\nComment")
    expect(t, "修改注释内容", "New\nComment" == comment1.Comment())
    
    //  添加注释
    comment3 := tinydom.NewComment(doc, "Comment3")
    doc.FirstChildElement("").InsertEndChild(comment3)
    expect(t, "向文档添加注释", "Comment3" == doc.FirstChildElement("").LastChild().Value())
}

func Test_Comment_含有错误注释的XML流(t *testing.T) {
    xml := `<node><elem1><!--comment2</elem1>--></node>
    `
    doc, err := tinydom.LoadDocument(strings.NewReader(xml))
    expect(t, "返回值检测", nil == doc)
    expect(t, "返回值检测", nil != err)
}

func Test_Text_基本功能测试(t *testing.T) {
    xml := "<node>text1<elem1>text2</elem1>\ttext3\n<elem2>\t\n      </elem2></node>"
    doc, err := tinydom.LoadDocument(strings.NewReader(xml))
    expect(t, "返回值检测", nil != doc)
    expect(t, "返回值检测", nil == err)
    
    //  转换
    text1node := doc.FirstChildElement("node").FirstChild()
    expect(t, "转换测试", nil != text1node.ToText())
    expect(t, "转换测试", nil != text1node.ToNode())
    expect(t, "转换测试", nil == text1node.ToComment())
    expect(t, "转换测试", nil == text1node.ToDirective())
    expect(t, "转换测试", nil == text1node.ToProcInst())
    expect(t, "转换测试", nil == text1node.ToDocument())
    expect(t, "转换测试", nil == text1node.ToElement())
    
    //  获得节点
    node := doc.FirstChildElement("node")
    elem1 := node.FirstChildElement("elem1")
    elem2 := node.FirstChildElement("elem2")
    expect(t, "获得节点对象", nil != node)
    expect(t, "获得节点对象", nil != elem1)
    expect(t, "获得节点对象", nil != elem2)
    
    //  获得Text对象
    text1 := node.FirstChild()
    text2 := elem1.FirstChild()
    text3 := elem1.NextSibling()
    text4 := elem2.FirstChild()
    expect(t, "获得Text对象", nil != text1)
    expect(t, "获得Text对象", nil != text2)
    expect(t, "获得Text对象", nil != text3)
    expect(t, "全空白的Text不会被读取", nil == text4)
    
    //  Text的数据获取
    expect(t, "Text的数据获取", "text1" == text1.Value())
    expect(t, "Text的数据获取", "text2" == text2.Value())
    expect(t, "如果不是全部为空白，那么空白部分也属于Text", "\ttext3\n" == text3.Value())

    //  父的Element可以通过Text直接获得第一个子节点的文本
    expect(t, "Text的数据获取", "text1" == node.Text())
    
    //  修改Text
    text1.SetValue("Hello World")
    expect(t, "修改Text的内容", "Hello World" == text1.Value())
    
    node.SetText("<TEXT>TextListXML</TEXT>")
    expect(t, "node可以直接获取第一个Text节点的值", "<TEXT>TextListXML</TEXT>" == node.Text())
    
    elem2.SetText("NewText")
    expect(t, "当通过Element.Text设置一个没有Text节点的值时，自动新建Text子节点", "NewText" == elem2.Text())
}

func Test_Text_Text出现在跟节点之外(t *testing.T) {
    xml := `<node></node>texterror`
    doc, err := tinydom.LoadDocument(strings.NewReader(xml))
    expect(t, "返回值检测", nil == doc)
    expect(t, "返回值检测", nil != err)
}


