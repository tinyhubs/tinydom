package tinydom

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func expect(t *testing.T, message string, result bool) {
	if result {
		return
	}

	fmt.Println(message)
	t.Fail()
}

func Test_example1(t *testing.T) {
	xmlstr := `
	<books>
	    <book><name>The Moon</name><author>Tom</author></book>
	    <book><name>Go west</name><author>Suny</author></book>
	</books>
	`

	doc, _ := LoadDocument(strings.NewReader(xmlstr))
	elem1 := doc.FirstChildElement("books").FirstChildElement("book").FirstChildElement("name")
	fmt.Println(elem1.Text()) //	The Moon

	elem2 := doc.FirstChildElement("books").FirstChildElement("book").LastChildElement("author")
	fmt.Println(elem2.Text()) //	Suny

}

func walk(m int, rootNode XMLNode) {
	if nil == rootNode {
		return
	}

	space := strings.Repeat("  ", m)
	for child := rootNode.FirstChild(); nil != child; child = child.Next() {
		fmt.Println(space, child.Value())
		walk(m+1, child)
	}
}

func Test_example2(t *testing.T) {
	doc := NewDocument()
	doc.InsertEndChild(NewProcInst("xml", `version="1.0" encoding="UTF-8"`))
	books := doc.InsertEndChild(NewElement("books"))
	book := books.InsertEndChild(NewElement("book"))
	name := book.InsertEndChild(NewElement("name"))
	name.InsertEndChild(NewText("The Moon"))

	doc.Accept(NewSimplePrinter(os.Stdout, PrintPretty))

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
	doc, _ := LoadDocument(strings.NewReader(xmlstr))
	talk := doc.FirstChildElement("talks").FirstChildElement("talk").Text()
	fmt.Print(talk)
}

func Test_example4(t *testing.T) {
	xmlstr := `<content><![CDATA[<example>This is ok in cdata text</example>]]></content>`
	doc, _ := LoadDocument(strings.NewReader(xmlstr))
	content := doc.FirstChildElement("content")
	fmt.Println("\nRead CDATA:", content.Text())
	fmt.Println("\nNormal Print:")
	doc.Accept(NewSimplePrinter(os.Stdout, PrintPretty))
	text := content.FirstChild().ToText()
	text.SetCDATA(true)
	fmt.Println("\nSpecial as CDATA:")
	doc.Accept(NewSimplePrinter(os.Stdout, PrintPretty))
}

func Test_Document_空文档_加载失败(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader(""))
	expect(t, "返回值检测", nil == doc)
	expect(t, "返回值检测", nil != err)
}

func Test_Document_格式错误_节点未关闭(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader("<node><elem></node>"))
	expect(t, "返回值检测", nil == doc)
	expect(t, "返回值检测", nil != err)
}

func Test_Document_格式错误_关闭节点多余(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader("<node><elem></elem></elem></node>"))
	expect(t, "返回值检测", nil == doc)
	expect(t, "返回值检测", nil != err)
}

func Test_Document_格式错误_多余的节点(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader("<node><elem></elem></node><hello/>"))
	expect(t, "返回值检测", nil == doc)
	expect(t, "返回值检测", nil != err)
}

func Test_Document_输出_各种元素遍历(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
	<!--comment1-->
	<!DOCTYPE poem>
	<node attr1="value1" attr2="value2"><elem><!--comment2--></elem><str>Hello world</str><hello/></node>`
	doc, err := LoadDocument(strings.NewReader(xml))
	doc.Accept(NewSimplePrinter(os.Stdout, PrintStream))
	expect(t, "返回值检测1", nil != doc)
	expect(t, "返回值检测2", nil == err)

	result1 := `<?xml version="1.0" encoding="UTF-8"?>` +
		`<!--comment1--><!DOCTYPE poem><node attr1="value1" attr2="value2"><elem><!--comment2--></elem><str>Hello world</str><hello/></node>`
	result2 := `<?xml version="1.0" encoding="UTF-8"?>` +
		`<!--comment1--><!DOCTYPE poem><node attr2="value2" attr1="value1"><elem><!--comment2--></elem><str>Hello world</str><hello/></node>`
	buf := bytes.NewBufferString("")
	doc.Accept(NewSimplePrinter(buf, PrintStream))
	expect(t, "检查输出3", (result1 == buf.String()) || (result2 == buf.String()))
}

func Test_Node_正常的XML文档_特殊场景_只有一个根节点(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader("<node></node>"))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)

	node := doc.FirstChild()
	expect(t, "检查节点名:node", node.Value() == "node")

	expect(t, "topo结构检查", doc == node.Parent())
	expect(t, "topo结构检查", doc == node.Document())

	expect(t, "topo结构检查", nil == node.FirstChild())
	expect(t, "topo结构检查", nil == node.LastChild())
	expect(t, "topo结构检查", nil == node.Prev())
	expect(t, "topo结构检查", nil == node.Next())

	expect(t, "topo结构检查", nil == node.FirstChildElement(""))
	expect(t, "topo结构检查", nil == node.LastChildElement(""))
	expect(t, "topo结构检查", nil == node.PrevElement(""))
	expect(t, "topo结构检查", nil == node.NextElement(""))

	expect(t, "转换检查", nil != node.ToElement())
	expect(t, "转换检查", nil == node.ToComment())
	expect(t, "转换检查", nil == node.ToDirective())
	expect(t, "转换检查", nil == node.ToDocument())
	expect(t, "转换检查", nil == node.ToProcInst())
	expect(t, "转换检查", nil == node.ToText())
}

func Test_Node_正常的XML文档_特殊场景_只有一个子节点(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader("<node><elem></elem></node>"))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)

	node := doc.FirstChild()
	elem := node.FirstChild()

	//  node
	expect(t, "检查节点名:node", "node" == node.Value())

	expect(t, "topo结构检查", doc == node.Parent())
	expect(t, "topo结构检查", doc == node.Document())

	expect(t, "topo结构检查", elem == node.FirstChild())
	expect(t, "topo结构检查", elem == node.LastChild())
	expect(t, "topo结构检查", nil == node.Prev())
	expect(t, "topo结构检查", nil == node.Next())

	expect(t, "topo结构检查", elem.ToElement() == node.FirstChildElement(""))
	expect(t, "topo结构检查", elem.ToElement() == node.LastChildElement(""))
	expect(t, "topo结构检查", nil == node.PrevElement(""))
	expect(t, "topo结构检查", nil == node.NextElement(""))

	expect(t, "topo结构检查", elem.ToElement() == node.FirstChildElement("elem"))
	expect(t, "topo结构检查", elem.ToElement() == node.LastChildElement("elem"))
	expect(t, "topo结构检查", nil == node.PrevElement(""))
	expect(t, "topo结构检查", nil == node.NextElement(""))

	//  elem
	expect(t, "检查节点名:elem", "elem" == elem.Value())

	expect(t, "topo结构检查", node == elem.Parent())
	expect(t, "topo结构检查", doc == elem.Document())

	expect(t, "topo结构检查", nil == elem.FirstChild())
	expect(t, "topo结构检查", nil == elem.LastChild())
	expect(t, "topo结构检查", nil == elem.Prev())
	expect(t, "topo结构检查", nil == elem.Next())

	expect(t, "topo结构检查", nil == elem.FirstChildElement(""))
	expect(t, "topo结构检查", nil == elem.LastChildElement(""))
	expect(t, "topo结构检查", nil == elem.PrevElement(""))
	expect(t, "topo结构检查", nil == elem.NextElement(""))
}

func Test_Node_正常的XML文档_丰富的文档结构(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader("<node><elem1></elem1><elem2></elem2><elem3></elem3></node>"))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)

	node := doc.FirstChild()
	elem1 := node.FirstChild()
	elem2 := elem1.Next()
	elem3 := elem2.Next()

	expect(t, "检查节点顺序关系", "node" == node.Value())
	expect(t, "检查节点顺序关系", "elem1" == elem1.Value())
	expect(t, "检查节点顺序关系", "elem2" == elem2.Value())
	expect(t, "检查节点顺序关系", "elem3" == elem3.Value())

	expect(t, "topo结构检查", node == elem1.Parent())
	expect(t, "topo结构检查", node == elem2.Parent())
	expect(t, "topo结构检查", node == elem3.Parent())

	expect(t, "topo结构检查", nil == elem1.Prev())
	expect(t, "topo结构检查", nil == elem3.Next())

	expect(t, "topo结构检查", elem1 == elem2.Prev())
	expect(t, "topo结构检查", elem3 == elem2.Next())

	expect(t, "topo结构检查", elem1.ToElement() == elem2.PrevElement(""))
	expect(t, "topo结构检查", elem3.ToElement() == elem2.NextElement(""))

	expect(t, "topo结构检查", elem1.ToElement() == elem2.PrevElement("elem1"))
	expect(t, "topo结构检查", elem3.ToElement() == elem2.NextElement("elem3"))

	expect(t, "topo结构检查", false == node.NoChildren())
	expect(t, "topo结构检查", true == elem1.NoChildren())
	expect(t, "topo结构检查", true == elem2.NoChildren())
	expect(t, "topo结构检查", true == elem3.NoChildren())
}

func Test_Node_修改文档_节点层次的增删改(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader("<node><elem1></elem1><elem2></elem2><elem3></elem3></node>"))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)

	node := doc.FirstChildElement("node")
	elem1 := node.FirstChildElement("elem1")
	//elem2 := node.FirstChildElement("elem2")

	new1 := elem1.InsertEndChild(NewElement("new1"))
	new2 := node.InsertEndChild(NewElement("new2"))
	//    new3 := node.InsertAfterChild(elem2, NewElement(doc, "new3"))
	new4 := elem1.InsertFirstChild(NewElement("new4"))
	expect(t, "添加成功", nil != new1)
	expect(t, "添加成功", nil != new2)
	//    expect(v, "添加成功", nil != new3)
	expect(t, "添加成功", nil != new4)

	doc.DeleteChild(new1)
	doc.DeleteChild(new2)
	//    doc.DeleteChild(new3)
	doc.DeleteChild(new4)
	node.DeleteChildren()

	node.SetName("NewNode")
	expect(t, "添加成功", "NewNode" == node.Value())
}

func Test_Element_属性同名错误(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader(`<node attr="value1" attr="value2"></node>`))
	expect(t, "返回值检测", nil == doc)
	expect(t, "返回值检测", nil != err)
}

func Test_Element_属性基本功能(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader(`<node attr1="value1" attr2="value2"></node>`))
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
	node.ForeachAttribute(func(attribute XMLAttribute) int {

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
	retult := node.ForeachAttribute(func(attribute XMLAttribute) int {
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
	expect(t, "遍历空属性列表总是返回0", 0 == node.ForeachAttribute(func(attribute XMLAttribute) int {
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
	doc, err := LoadDocument(strings.NewReader(xml))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)

	expect(t, "申明头转换测试", nil != doc.FirstChild().ToProcInst())
	expect(t, "申明头转换测试", nil == doc.FirstChild().ToElement())
	expect(t, "申明头转换测试", nil == doc.FirstChild().ToDocument())
	expect(t, "申明头转换测试", nil == doc.FirstChild().ToDirective())
	expect(t, "申明头转换测试", nil == doc.FirstChild().ToComment())
	expect(t, "申明头转换测试", nil == doc.FirstChild().ToText())

	procInst := doc.FirstChild().ToProcInst()
	expect(t, "有申明头的xml文档，第一个子节点是申明头", "xml" == procInst.Value())
	expect(t, "有申明头的xml文档，第一个子节点是申明头", "xml" == procInst.Target())
	expect(t, "有申明头的xml文档，第一个子节点是申明头", `version="1.0" encoding="UTF-8"` == procInst.Instruction())
	expect(t, "申明头下面一个是node", "node" == procInst.Next().Value())
}

func Test_Comment_基本功能测试(t *testing.T) {
	xml := `<!--comment1--><node><elem1><!--comment2--></elem1></node>`

	//  加载
	doc, err := LoadDocument(strings.NewReader(xml))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)

	//  转换
	expect(t, "转换测试", nil != doc.FirstChild().ToComment())
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
	comment3 := NewComment("Comment3")
	doc.FirstChildElement("").InsertEndChild(comment3)
	expect(t, "向文档添加注释", "Comment3" == doc.FirstChildElement("").LastChild().Value())
}

func Test_Comment_含有错误注释的XML流(t *testing.T) {
	xml := `<node><elem1><!--comment2</elem1>--></node>
    `
	doc, err := LoadDocument(strings.NewReader(xml))
	expect(t, "返回值检测", nil == doc)
	expect(t, "返回值检测", nil != err)
}

func Test_Text_基本功能测试(t *testing.T) {
	xml := "<node>text1<elem1>text2</elem1>\ttext3\n<elem2>\t\n      </elem2></node>"
	doc, err := LoadDocument(strings.NewReader(xml))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)

	//  转换
	text1node := doc.FirstChildElement("node").FirstChild()
	expect(t, "转换测试", nil != text1node.ToText())
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
	text3 := elem1.Next()
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
	doc, err := LoadDocument(strings.NewReader(xml))
	expect(t, "返回值检测", nil == doc)
	expect(t, "返回值检测", nil != err)
}

func Test_Directive_基本功能测试(t *testing.T) {
	xml := `
    <?xml version="1.0" encoding="UTF-8"?>
    <!DOCTYPE poem [
        <!ELEMENT poem (author, title, content)>
        <!ELEMENT author (#PCDATA)>
        <!ELEMENT title (#PCDATA)>
        <!ELEMENT content (#PCDATA)>
    ]>
    <!--为元素poem定义了三个子元素author title content，
    这三个元素必须要出现并且必须按照这个顺序
    少元素不行，多元素也不行
    -->
    <!--指明author,title,content里面的内容是字符串类型-->
    <poem>
        <author>王维</author>
        <title>鹿柴</title>
        <content>空山不见人，但闻人语声。返景入深林，复照青苔上。</content>
    </poem>
    `
	doc, err := LoadDocument(strings.NewReader(xml))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)

	//  转换
	doctype := doc.FirstChild().Next()
	expect(t, "转换测试", nil != doctype.ToDirective())
	expect(t, "转换测试", nil == doctype.ToText())
	expect(t, "转换测试", nil == doctype.ToComment())
	expect(t, "转换测试", nil == doctype.ToProcInst())
	expect(t, "转换测试", nil == doctype.ToDocument())
	expect(t, "转换测试", nil == doctype.ToElement())

	cmp := `DOCTYPE poem [
        <!ELEMENT poem (author, title, content)>
        <!ELEMENT author (#PCDATA)>
        <!ELEMENT title (#PCDATA)>
        <!ELEMENT content (#PCDATA)>
    ]`
	expect(t, "转换测试", cmp == doctype.Value())
}

func Test_Handle_空腹测试(t *testing.T) {
	handle := NewHandle(nil)
	expect(t, "空转换测试", nil == handle.ToDirective())
	expect(t, "空转换测试", nil == handle.ToText())
	expect(t, "空转换测试", nil == handle.ToComment())
	expect(t, "空转换测试", nil == handle.ToProcInst())
	expect(t, "空转换测试", nil == handle.ToDocument())
	expect(t, "空转换测试", nil == handle.ToElement())

	expect(t, "空周游测试", nil == handle.FirstChild().ToNode())
	expect(t, "空周游测试", nil == handle.LastChild().ToNode())
	expect(t, "空周游测试", nil == handle.Prev().ToNode())
	expect(t, "空周游测试", nil == handle.Next().ToNode())
	expect(t, "空周游测试", nil == handle.FirstChildElement("").ToNode())
	expect(t, "空周游测试", nil == handle.LastChildElement("").ToNode())
	expect(t, "空周游测试", nil == handle.PrevElement("").ToNode())
	expect(t, "空周游测试", nil == handle.NextElement("").ToNode())
}

func Test_Handle_基本功能测试(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
	<!--comment1-->
	<!DOCTYPE poem [
        <!ELEMENT poem (author, title, content)>
        <!ELEMENT author (#PCDATA)>
        <!ELEMENT title (#PCDATA)>
        <!ELEMENT content (#PCDATA)>
    	]>
	<node attr1="value1" attr2="value2"><elem><!--comment2--></elem><elem>126</elem><str>Hello world</str></node>`
	doc, err := LoadDocument(strings.NewReader(xml))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)

	handle := NewHandle(doc)
	expect(t, "周游测试", nil != handle.ToDocument())
	expect(t, "周游测试", "xml" == handle.FirstChild().ToProcInst().Value())
	expect(t, "周游测试", "comment1" == handle.FirstChild().Next().ToComment().Value())

	node := handle.FirstChildElement("node")
	expect(t, "周游测试", nil != node.Parent().ToDocument())
	expect(t, "周游测试", "" == node.FirstChildElement("elem").ToElement().Text())
	expect(t, "周游测试", "126" == node.LastChildElement("elem").ToElement().Text())
	expect(t, "周游测试", "" == node.LastChildElement("elem").PrevElement("elem").ToElement().Text())
	expect(t, "周游测试", "126" == node.FirstChildElement("elem").NextElement("elem").ToElement().Text())
	expect(t, "周游测试", "126" == node.LastChildElement("elem").FirstChild().ToText().Value())
	expect(t, "周游测试", "126" == node.LastChildElement("elem").LastChild().ToText().Value())

	str := node.FirstChildElement("str")
	expect(t, "周游测试", "elem" == str.Prev().ToElement().Value())
	expect(t, "周游测试", nil != handle.FirstChildElement("node").Prev().ToDirective())
}

func Test_Node_修改文档_对新添加的节点进行修改(t *testing.T) {
	doc, err := LoadDocument(strings.NewReader("<node><elem1/></node>"))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)

	node := doc.FirstChildElement("node")
	elem1 := node.FirstChildElement("elem1")
	elem1.SetName("newelem")

	node.DeleteChild(elem1)
	node.InsertEndChild(elem1)
}

func Test_Handle_基本功能测试_Parent(t *testing.T) {
	xml := `<node attr1="value1" attr2="value2"></node>`
	doc, err := LoadDocument(strings.NewReader(xml))
	expect(t, "返回值检测", nil != doc)
	expect(t, "返回值检测", nil == err)
	expect(t, "返回值检测", nil != doc.FirstChild().Parent().ToDocument())
}

func Test_TODO_Document_通过修改文档破坏xml文档的有效性(t *testing.T) {
}

func Test_TODO_Document_各种dom树输出(t *testing.T) {
}

func Test_TODO_Node_将另外一个文档的node添加到本文档(t *testing.T) {
}

func Test_Print(t *testing.T) {

	str := `<?xml version="1.0" encoding="UTF-8"?><module><changes><change sequence="0000&amp;00000"/><change >@properties.insert "${ITPAAS_HONE}/cloudcontrol/conf-itpaas/database/ccc.properties" &amp; "key" &quot; "value"</change></changes></module><!--  ddd  -->`
	exp := `<?xml version="1.0" encoding="UTF-8"?><module><changes><change sequence="0000&amp;00000"/><change>@properties.insert "${ITPAAS_HONE}/cloudcontrol/conf-itpaas/database/ccc.properties" &amp; "key" " "value"</change></changes></module><!--  ddd  -->`
	doc, _ := LoadDocument(strings.NewReader(str))
	buf := bytes.NewBufferString("")
	doc.Accept(NewSimplePrinter(buf, PrintOptions{}))
	expect(t, "检查格式化xml文档的结果", buf.String() == exp)
}

func Test_Version(t *testing.T) {
	datas, err := ioutil.ReadFile("version")
	if nil != err {
		t.Fail()
		return
	}

	ver := strings.TrimSpace(string(datas))
	if ver != Version() {
		t.Fail()
		return
	}
}

func Test_EscapeAttribute(t *testing.T) {

	tester := func(str string, esc string) {
		doc := NewDocument()
		if nil == doc {
			t.Fail()
			return
		}

		elem := NewElement("elem")
		elem.SetAttribute("attr", str)
		doc.InsertEndChild(elem)

		buf := bytes.NewBufferString("")
		doc.Accept(NewSimplePrinter(buf, PrintStream))

		compare := fmt.Sprintf(`<elem attr="%s"/>`, esc)
		if compare != buf.String() {
			t.Fail()
			return
		}
	}

	tester(`tinyhubs@126.com&tinyhubs@126.com`, `tinyhubs@126.com&amp;tinyhubs@126.com`)
	tester(`"tinyhubs@126.com`, `&quot;tinyhubs@126.com`)
	tester(`tinyhubs@126.com<`, `tinyhubs@126.com&lt;`)
	tester("aaa\n", `aaa&#xA;`)
	tester("\raaa", `&#xD;aaa`)
	tester(`aaa'`, `aaa'`)
	tester(`aaa>aaa`, `aaa>aaa`)
}

func Test_EscapeText(t *testing.T) {

	tester := func(str string, esc string) {
		doc := NewDocument()
		if nil == doc {
			t.Fail()
			return
		}

		elem := NewElement("elem")
		doc.InsertEndChild(elem)
		elem.SetText(str)

		buf := bytes.NewBufferString("")
		doc.Accept(NewSimplePrinter(buf, PrintStream))

		compare := fmt.Sprintf(`<elem>%s</elem>`, esc)
		if compare != buf.String() {
			t.Fail()
			return
		}
	}

	tester(`tinyhubs@126.com&tinyhubs@126.com`, `tinyhubs@126.com&amp;tinyhubs@126.com`)
	tester(`tinyhubs@126.com<`, `tinyhubs@126.com&lt;`)
	tester(`"tinyhubs@126.com`, `"tinyhubs@126.com`)
	tester("aaa\naaa", "aaa\naaa")
	tester("\raaa", "\raaa")
	tester(`aaa'`, `aaa'`)
	tester(`aaa>aaa`, `aaa>aaa`)
}

func Test_Inserts(t *testing.T) {
	doc := NewDocument()
	doc.InsertEndChild(NewElement("elem1")). //  <elem1></elem1>
		InsertFirstChild(NewElement("elem2")). //  <elem1><elem2></elem2></elem1>
		InsertFront(NewElement("elem3")). //  <elem1><elem3></elem3><elem2></elem2></elem1>
		InsertBack(NewElement("elem4")). //  <elem1><elem3></elem3><elem4></elem4><elem2></elem2></elem1>
		InsertElementFront("elem5"). //  <elem1><elem3></elem3><elem5></elem5><elem4></elem4><elem2></elem2></elem1>
		InsertElementBack("elem6"). //  <elem1><elem3></elem3><elem5></elem5><elem6></elem6><elem4></elem4><elem2></elem2></elem1>
		InsertElementEndChild("elem7"). //  <elem1><elem3></elem3><elem5></elem5><elem6><elem7><elem8></elem8></elem7></elem6><elem4></elem4><elem2></elem2></elem1>
		InsertElementFirstChild("elem8").
		InsertEndChild(NewElement("elem9")).
		InsertBack(NewElement("elem10"))

	exp := `<elem1><elem3/><elem5/><elem6><elem7><elem8><elem9/><elem10/></elem8></elem7></elem6><elem4/><elem2/></elem1>`

	buf := bytes.NewBufferString("")
	doc.Accept(NewSimplePrinter(buf, PrintStream))
	fmt.Println("pppp", buf.String())
	expect(t, "检查格式化输出结果", exp == buf.String())
}

//func Test_Inserts2(v *testing.T) {
//    doc := NewDocument()
//    elem1 := doc.InsertEndChild(NewElement("elem1"))     //  <elem1></elem1>
//    elem2 := elem1.InsertFirstChild(NewElement("elem2")) //  <elem1><elem2></elem2></elem1>
//    elem3 := elem2.InsertFront(NewElement("elem3"))      //  <elem1><elem3></elem3><elem2></elem2></elem1>
//    elem4 := elem3.InsertBack(NewElement("elem4"))       //  <elem1><elem3></elem3><elem4></elem4><elem2></elem2></elem1>
//    elem5 := elem4.InsertElementFront("elem5")           //  <elem1><elem3></elem3><elem5></elem5><elem4></elem4><elem2></elem2></elem1>
//
//    doc.Accept(NewSimplePrinter(os.Stdout, PrintStream))
//    fmt.Println("------------")
//    fmt.Println("elem5", elem5)
//}
//
//func Test_Inserts3(v *testing.T) {
//    doc := NewDocument()
//    elem1 := doc.InsertEndChild(NewElement("elem1"))     //  <elem1></elem1>
//    elem2 := elem1.InsertFirstChild(NewElement("elem2")) //  <elem1><elem2></elem2></elem1>
//    elem3 := elem2.InsertElementFront("elem3")            //  <elem1><elem3></elem3><elem2></elem2></elem1>
//    elem4 := elem2.InsertElementFront("elem4")            //  <elem1><elem3></elem3><elem4></elem4><elem2></elem2></elem1>
//
//    doc.Accept(NewSimplePrinter(os.Stdout, PrintStream))
//    fmt.Println("------------")
//    fmt.Println("elem5", elem4, elem3)
//}

func Test_XmlHandle_Parent_NULL(t *testing.T) {
	doc, _ := LoadDocument(bytes.NewBufferString(`<elem1/>`))
	nd := NewHandle(doc).FirstChild().Parent().Parent().Parent().ToNode()
	expect(t, "Document的上层只能为null", nil == nd)
}

func Test_Accept_Terminate(t *testing.T) {
	str := `
        <?xml version="1.0" encoding="UTF-8"?>
	<!--comment1-->
	<!DOCTYPE poem>
	<node attr1="value1" attr2="value2">
	        <elem>
	                <!--comment2-->
	        </elem>
	        <str>Hello world</str>
	        <hello/>
        </node>
        `

	doc, _ := LoadDocument(bytes.NewBufferString(str))
	doc.Accept(&DefaultVisitor{})

	// 各种提前中断遍历
	doc.Accept(&DefaultVisitor{EnterDocument: func(XMLDocument) bool {
		return false
	}})

	// 各种提前中断遍历
	doc.Accept(&DefaultVisitor{ExitDocument: func(XMLDocument) bool {
		return false
	}})

	// 各种提前中断遍历
	doc.Accept(&DefaultVisitor{EnterElement: func(element XMLElement) bool {
		return false
	}})

	// 各种提前中断遍历
	doc.Accept(&DefaultVisitor{ExitElement: func(element XMLElement) bool {
		return false
	}})

	// 各种提前中断遍历
	doc.Accept(&DefaultVisitor{Text: func(text XMLText) bool {
		return false
	}})

	// 各种提前中断遍历
	doc.Accept(&DefaultVisitor{ProcInst: func(text XMLProcInst) bool {
		return false
	}})

	// 各种提前中断遍历
	doc.Accept(&DefaultVisitor{Comment: func(text XMLComment) bool {
		return false
	}})

	// 各种提前中断遍历
	doc.Accept(&DefaultVisitor{Directive: func(text XMLDirective) bool {
		return false
	}})

}

func Test_Document_Cannot_Add_Sibling(t *testing.T) {
	doc := NewDocument()
	elem1 := doc.InsertBack(NewElement("elem1"))
	expect(t, "document 不能添加兄弟节点", nil == elem1)
	elem2 := doc.InsertFront(NewElement("elem2"))
	expect(t, "document 不能添加兄弟节点", nil == elem2)
}

func Test_Node_Split(t *testing.T) {
	s := `
	<!-- comment1 -->
	<node1>
		<node2/>
	</node1>
	`

	r := `
	<node3>
	<!-- comment3 -->
	</node3>
	<!-- yyy -->
	`

	sdoc, _ := LoadDocument(bytes.NewBufferString(s))
	rdoc, _ := LoadDocument(bytes.NewBufferString(r))
	elem := sdoc.FirstChildElement("")
	elem.Split()
	//sdoc.InsertEndChild(rdoc.FirstChildElement(""))
	sdoc.InsertEndChild(rdoc.LastChildElement(""))

	buf := bytes.NewBufferString("")
	sdoc.Accept(NewSimplePrinter(buf, PrintStream))
	expect(t, "检查输出结果,本用例用于检测tinydom的继承机制是否完善", buf.String() == `<!-- comment1 --><node3><!-- comment3 --></node3>`)
}

type Interface interface {
	Echo(string) string
}

type Base struct {
	impl *Base
}

func (b*Base) Echo(s string) string {
	return fmt.Sprintf("Base : %s", s)
}

type Div struct {
	Base
}

func (b*Div) Echo(s string) string {
	return fmt.Sprintf("Div : %s", s)
}

func Call(b Interface) {
	b.Echo("hello")
}

func Test_inhert(t *testing.T) {
	a := new(Base)
	b := new(Div)
	var c Interface = new(Div)
	//d := b.(*Base) 编译失败
	Call(a)
	Call(b)
	Call(c)
	//Call(d)
}

func Test_Attr_Order(t *testing.T) {
	s := `<node attr5="55"/>`
	doc, _ := LoadDocument(bytes.NewBufferString(s))
	node := doc.FirstChildElement("node")
	node.SetAttribute("attr2", "22")
	node.SetAttribute("attr3", "33")
	node.SetAttribute("attr4", "44")
	node.SetAttribute("attr6", "66")
	node.SetAttribute("attr9", "99")
	node.SetAttribute("attr", "")
	buf := bytes.NewBufferString("")
	doc.Accept(NewSimplePrinter(buf, PrintStream))
	expect(t, "属性的顺序就是添加的顺序,不会应为key的不断变化而导致属性输出时,属性间的相对位置发生不断变化",
	buf.String() == `<node attr5="55" attr2="22" attr3="33" attr4="44" attr6="66" attr9="99" attr=""/>`)
}
