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

func Test_Document_空文档_加载失败(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader(""))
    expect(t, "doc为空", nil == doc)
    expect(t, "doc为空", nil != err)
}

func Test_Document_格式错误_节点未关闭(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node><elem></node>"))
    expect(t, "doc为空", nil == doc)
    expect(t, "doc为空", nil != err)
}

func Test_Document_格式错误_关闭节点多余(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node><elem></elem></elem></node>"))
    expect(t, "doc为空", nil == doc)
    expect(t, "doc为空", nil != err)
}

func Test_Document_格式错误_多余的节点(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node><elem></elem></node><hello/>"))
    doc.Accept(tinydom.NewSimplePrinter(os.Stdout))
    expect(t, "doc为空", nil == doc)
    expect(t, "doc为空", nil != err)
}

func Test_Document_正常的XML文档1(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node></node>"))
    expect(t, "doc为空", nil != doc)
    expect(t, "doc为空", nil == err)
    
    node := doc.FirstChild()
    expect(t, "检查节点名:node", node.Value() == "node")
    
    expect(t, "topo结构检查", doc == node.Parent())
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

func Test_Document_正常的XML文档2(t *testing.T) {
    doc, err := tinydom.LoadDocument(strings.NewReader("<node><elem></elem></node>"))
    expect(t, "doc为空", nil != doc)
    expect(t, "doc为空", nil == err)
    
    node := doc.FirstChild()
    elem := node.FirstChild()
    
    //  node
    expect(t, "检查节点名:node", "node" == node.Value())
    
    expect(t, "topo结构检查", doc == node.Parent())
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
    
    expect(t, "topo结构检查", node == elem.Parent())
    expect(t, "topo结构检查", doc == elem.GetDocument())
    
    expect(t, "topo结构检查", nil == node.FirstChild())
    expect(t, "topo结构检查", nil == node.LastChild())
    expect(t, "topo结构检查", nil == node.PreviousSibling())
    expect(t, "topo结构检查", nil == node.NextSibling())
    
    expect(t, "topo结构检查", nil == node.FirstChildElement(""))
    expect(t, "topo结构检查", nil == node.LastChildElement(""))
    expect(t, "topo结构检查", nil == node.PreviousSiblingElement(""))
    expect(t, "topo结构检查", nil == node.NextSiblingElement(""))
}
