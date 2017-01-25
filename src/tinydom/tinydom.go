/*
Package tinydom	实现了一个简单的XML的DOM模型.

tidydom使用encoding/xml作为底层XML解析库，实现对XML文件的解析.使用tinydom提供的接口可以实现简单的XML文件的读取和生成。
tinydom借鉴了tinyxml和的tinyxml2的接口设计技巧，tinydom的接口和tinyxml类似，都提供了丰富的查找XML元素的查找手段。

一个XML文档由	XMLElement、XMLText、XMLComment、XMLDocument、XMLProcInst、XMLDirective者几种类型的节点组成。
XNLNode是所有这些节点的共同基础，XMLNode提供了丰富的节点元素遍历手段。
XMLVisitor提供了一种XML对象的元素遍历机制。
XMLHandle的所用是简化代码编写工作，使用XMLHandle将减少很多判空代码(if nil == xxx {}),活用XMLHandle将会让XML文件的元素事半功倍。

加载文档：

LoadDocument用于从一个文件流或者字符流读取XML数据，并构建出XMLDocument对象，一般用于读取XML文件的场景。
    import "tinydom"
    doc, err := tinydom.LoadDocument(strings.NewReader(s))

从文档中找到我们需要的元素：
FirstChildElement、LastChildElement、PreviousSiblingElement、NextSiblingElement这些接口主要是为了方便查找XMLElement元素，
大部分情况下我们建立XML文档的DMO模型就是为了对XMLElement进行访问。
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

新建文档：

NewDocument用于在内存中生成DOM，一般用于生成XML文件。
InsertEndChild、InsertFirstChild、InsertAfterChild、DeleteChildren、DeleteChild用于对XMLDocument进行修改。
下面的代码创建了一个XML文档：
    doc := tinydom.NewDocument()
    books := doc.InsertEndChild(tinydom.NewElement(doc, "books"))
    book := books.InsertEndChild(tinydom.NewElement(doc, "book"))
    name := book.InsertEndChild(tinydom.NewElement(doc, "name"))
    name.InsertEndChild(tinydom.NewText(doc, "The Moon"))
    doc.InsertEndChild(tinydom.NewProcInst(doc, "xml", `version="1.0" encoding="UTF-8"`))

我们可以使用XMLDocument.Accept方法来将这个XML文档输出：
    doc.Accept(tinydom.NewSimplePrinter(os.Stdout))

文档的遍历：

Parent、FirstChild、LastChild、PreviousSibling、NextSibling用于使我们可以方便地在XML的DOM树中游走。
下面这个函数可以用于对一个doc进行遍历，可以这样使用walk(doc)。
还有一个更好的替代方式是使用XMLVisitor接口对文档中的元素进行遍历。
    func walk(m int , rootNode tinydom.XMLNode) {
        if nil == rootNode {
            return
        }

        for child := rootNode.FirstChild(); nil != child; child = child.NextSibling() {
            fmt.Println(strings.Repeat(" ", m), child.Value())
            walk(m + 1, child)
        }
    }

XML字符转义：

受益于go的xml库，tinydom也支持XML字符转义，使用tinydom在读写xml的数据的时候不需要关注XML转义字符，tinydom自动会处理好，可参考下面的例子。
如果您需要自定义输出格式，那么文本雷荣时，需要通过xml.ExcapeText函数进行转义。
    xmlstr :=
        `<talks>
            <talk from="bill" to="tom">[&amp;&apos;&quot;&gt;&lt;] are the xml escape chars? </talk>
            <talk from="tom" to="bill">yes， that is right</talk>
         </talks>
        `
    doc, _ := tinydom.LoadDocument(strings.NewReader(xmlstr))
    talk := doc.FirstChildElement("talks").FirstChildElement("talk").Text()
    fmt.Print(talk) //  [&'"><] are the xml escape chars?

CDATA：

只有XMLText对象才涉及到CDATA，可以通过XMLText，tinydom能够自动识别CDATA，但是将DOM对象序列化成字符串时，除非节点指定了CDATA属性，否则会直接转义。
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


*/
package tinydom

import (
    "encoding/xml"
    "errors"
    "io"
    "bytes"
)

//  XMLAttribute    是一个元素的属性的接口
type XMLAttribute interface {
    Name() string
    Value() string
    SetValue(string)
}

//  XMLNode 定义了XML所有节点的基础设施，提供了基本的元素遍历、增删等操作,也提供了逆向转换能力.
type XMLNode interface {
    ToElement() XMLElement
    ToText() XMLText
    ToComment() XMLComment
    ToDocument() XMLDocument
    ToProcInst() XMLProcInst
    ToDirective() XMLDirective
    
    Value() string
    SetValue(newValue string)
    
    GetDocument() XMLDocument
    
    NoChildren() bool
    Parent() XMLNode
    FirstChild() XMLNode
    LastChild() XMLNode
    PreviousSibling() XMLNode
    NextSibling() XMLNode
    FirstChildElement(name string) XMLElement
    LastChildElement(name string) XMLElement
    PreviousSiblingElement(name string) XMLElement
    NextSiblingElement(name string) XMLElement
    
    InsertEndChild(node XMLNode) XMLNode
    InsertFirstChild(node XMLNode) XMLNode
    InsertAfterChild(afterThis XMLNode, addThis XMLNode) XMLNode
    DeleteChildren()
    DeleteChild(node XMLNode)
    Accept(visitor XMLVisitor) bool
    
    //  被迫入侵的接口
    getDocument() XMLDocument
    setDocument(document XMLDocument)
    getParent() XMLNode
    setParent(node XMLNode)
    getValue() string
    setValue(value string)
    getFirstChild() XMLNode
    setFirstChild(node XMLNode)
    getLastChild() XMLNode
    setLastChild(node XMLNode)
    getPrev() XMLNode
    setPrev(node XMLNode)
    getNext() XMLNode
    setNext(node XMLNode)
    
    unlink(child XMLNode)
}

//  XMLElement  提供了访问XML基本节点元素的能力
//
//  Name、SetName其实是Value和SetValue的别名，目的是为了使得接口更加符合直观理解。
//
//  Text、SetText的作用是设置<node>与</node>之间的文字，虽然文字都是有XMLText对象来承载的，但是通常来说直接在XMLElement中访问会更加方便。
//
//  FindAttribute和ForeachAttribute分别用于查找特定的XML节点的属性和遍历XML属性列表。
//
//  Attribute、SetAttribute、DeleteAttribute用于读取和删除属性。
type XMLElement interface {
    XMLNode
    
    Name() string
    SetName(name string)
    
    FindAttribute(name string) XMLAttribute
    ForeachAttribute(callback func(attribute XMLAttribute) int) int
    
    AttributeCount() int
    
    Attribute(name string, def string) string
    SetAttribute(name string, value string)
    DeleteAttribute(name string)
    
    Text() string
    SetText(text string)
}

//  XMLText 提供了对XML元素间文本的封装
type XMLText interface {
    XMLNode
    SetCDATA(isCData bool)
    CDATA() bool
}

type XMLComment interface {
    XMLNode
}

type XMLProcInst interface {
    XMLNode
    Target() string
    Instruction() string
}

type XMLDirective interface {
    XMLNode
}

type XMLDocument interface {
    XMLNode
}

type XMLVisitor interface {
    VisitEnterDocument(XMLDocument) bool
    VisitExitDocument(XMLDocument) bool
    
    VisitEnterElement(XMLElement) bool
    VisitExitElement(XMLElement) bool
    
    VisitProcInst(XMLProcInst) bool
    VisitText(XMLText) bool
    VisitComment(XMLComment) bool
    VisitDirective(XMLDirective) bool
}

type XMLHandle interface {
    Parent() XMLHandle
    FirstChild() XMLHandle
    LastChild() XMLHandle
    PreviousSibling() XMLHandle
    NextSibling() XMLHandle
    FirstChildElement(name string) XMLHandle
    LastChildElement(name string) XMLHandle
    PreviousSiblingElement(name string) XMLHandle
    NextSiblingElement(name string) XMLHandle
    
    ToElement() XMLElement
    ToText() XMLText
    ToComment() XMLComment
    ToDocument() XMLDocument
    ToProcInst() XMLProcInst
    ToDirective() XMLDirective
}

//=========================================================
type xmlAttributeImpl struct {
    name  string
    value string
}

func (this *xmlAttributeImpl) Name() string {
    return this.name
}

func (this *xmlAttributeImpl) Value() string {
    return this.value
}

func (this *xmlAttributeImpl) SetValue(newValue string) {
    this.value = newValue
}

//==================================================================

type xmlNodeImpl struct {
    document XMLDocument
    parent   XMLNode
    value    string
    
    firstChild XMLNode
    lastChild  XMLNode
    
    prev XMLNode
    next XMLNode
}

func (this *xmlNodeImpl) getDocument() XMLDocument {
    return this.document
}

func (this *xmlNodeImpl) setDocument(document XMLDocument) {
    this.document = document
}

func (this *xmlNodeImpl) getParent() XMLNode {
    return this.parent
}

func (this *xmlNodeImpl) setParent(node XMLNode) {
    this.parent = node
}

func (this *xmlNodeImpl) getValue() string {
    return this.value
}

func (this *xmlNodeImpl) setValue(value string) {
    this.value = value
}

func (this *xmlNodeImpl) getFirstChild() XMLNode {
    return this.firstChild
}

func (this *xmlNodeImpl) setFirstChild(node XMLNode) {
    this.firstChild = node
}

func (this *xmlNodeImpl) getLastChild() XMLNode {
    return this.lastChild
}

func (this *xmlNodeImpl) setLastChild(node XMLNode) {
    this.lastChild = node
}

func (this *xmlNodeImpl) getPrev() XMLNode {
    return this.prev
}

func (this *xmlNodeImpl) setPrev(node XMLNode) {
    this.prev = node
}

func (this *xmlNodeImpl) getNext() XMLNode {
    return this.next
}

func (this *xmlNodeImpl) setNext(node XMLNode) {
    this.next = node
}

func (this *xmlNodeImpl) ToElement() XMLElement {
    return nil
}

func (this *xmlNodeImpl) ToText() XMLText {
    return nil
}

func (this *xmlNodeImpl) ToComment() XMLComment {
    return nil
}

func (this *xmlNodeImpl) ToDocument() XMLDocument {
    return nil
}

func (this *xmlNodeImpl) ToProcInst() XMLProcInst {
    return nil
}

func (this *xmlNodeImpl) ToDirective() XMLDirective {
    return nil
}

func (this *xmlNodeImpl) Value() string {
    return this.value
}

func (this *xmlNodeImpl) SetValue(newValue string) {
    this.value = newValue
}

func (this *xmlNodeImpl) GetDocument() XMLDocument {
    return this.document
}

func (this *xmlNodeImpl) Parent() XMLNode {
    return this.parent
}

func (this *xmlNodeImpl) NoChildren() bool {
    return nil == this.firstChild
}

func (this *xmlNodeImpl) FirstChild() XMLNode {
    return this.firstChild
}

func (this *xmlNodeImpl) LastChild() XMLNode {
    return this.lastChild
}

func (this *xmlNodeImpl) PreviousSibling() XMLNode {
    return this.prev
}

func (this *xmlNodeImpl) NextSibling() XMLNode {
    return this.next
}

func (this *xmlNodeImpl) FirstChildElement(name string) XMLElement {
    for item := this.firstChild; nil != item; item = item.getNext() {
        elem := item.ToElement()
        if nil == elem {
            continue
        }
        
        if ("" == name) || (elem.Name() == name) {
            return elem
        }
    }
    
    return nil
}

func (this *xmlNodeImpl) LastChildElement(name string) XMLElement {
    
    for item := this.lastChild; nil != item; item = item.getPrev() {
        elem := item.ToElement()
        if nil == elem {
            continue
        }
        
        if ("" == name) || (elem.Name() == name) {
            return elem
        }
    }
    
    return nil
}

func (this *xmlNodeImpl) PreviousSiblingElement(name string) XMLElement {
    for item := this.prev; nil != item; item = item.getPrev() {
        elem := item.ToElement()
        if nil == elem {
            continue
        }
        
        if ("" == name) || (elem.Name() == name) {
            return elem
        }
    }
    
    return nil
}

func (this *xmlNodeImpl) NextSiblingElement(name string) XMLElement {
    
    for item := this.next; nil != item; item = item.getNext() {
        elem := item.ToElement()
        if nil == elem {
            continue
        }
        
        if ("" == name) || (elem.Name() == name) {
            return elem
        }
    }
    
    return nil
}

func (this *xmlNodeImpl) unlink(child XMLNode) {
    if child == this.firstChild {
        this.firstChild = this.firstChild.getNext()
    }
    
    if child == this.lastChild {
        this.lastChild = this.lastChild.getPrev()
    }
    
    if nil != child.getPrev() {
        child.getPrev().setNext(child.getNext())
    }
    
    if nil != child.getNext() {
        child.getNext().setPrev(child.getPrev())
    }
    
    child.setParent(nil)
}

func (this *xmlNodeImpl) InsertEndChild(addThis XMLNode) XMLNode {
    if addThis.getDocument() != this.document {
        return nil
    }
    
    if nil != addThis.getParent() {
        addThis.getParent().unlink(addThis)
    }
    
    if nil != this.lastChild {
        this.lastChild.setNext(addThis)
        addThis.setPrev(this.lastChild)
        this.lastChild = addThis
        addThis.setNext(nil)
    } else {
        this.firstChild = addThis
        this.lastChild = addThis
        
        addThis.setPrev(nil)
        addThis.setNext(nil)
    }
    
    addThis.setParent(this)
    return addThis
}

func (this *xmlNodeImpl) InsertFirstChild(addThis XMLNode) XMLNode {
    if addThis.getDocument() != this.document {
        return nil
    }
    
    if nil != addThis.getParent() {
        addThis.getParent().unlink(addThis)
    }
    
    if nil != this.firstChild {
        this.firstChild.setPrev(addThis)
        addThis.setNext(this.firstChild)
        this.firstChild = addThis
        addThis.setPrev(nil)
    } else {
        this.firstChild = addThis
        this.lastChild = addThis
        
        addThis.setPrev(nil)
        addThis.setNext(nil)
    }
    
    addThis.setParent(this)
    return addThis
}

func (this *xmlNodeImpl) InsertAfterChild(afterThis XMLNode, addThis XMLNode) XMLNode {
    if addThis.getDocument() != this.document {
        return nil
    }
    
    if afterThis.getParent() != this {
        return nil
    }
    
    if afterThis.getNext() == nil {
        return this.InsertEndChild(addThis)
    }
    
    if nil != addThis.getParent() {
        addThis.getParent().unlink(addThis)
    }
    
    addThis.setPrev(afterThis)
    addThis.setNext(afterThis.getNext())
    afterThis.getNext().setPrev(addThis)
    afterThis.setNext(addThis)
    addThis.setParent(this)
    return addThis
}

func (this *xmlNodeImpl) DeleteChildren() {
    for nil != this.firstChild {
        this.DeleteChild(this.firstChild)
    }
    
    this.firstChild = nil
    this.lastChild = nil
}

func (this *xmlNodeImpl) DeleteChild(node XMLNode) {
    this.unlink(node)
}

func (this *xmlNodeImpl) Accept(visitor XMLVisitor) bool {
    return false
}

//------------------------------------------------------------------

type xmlElementImpl struct {
    xmlNodeImpl
    
    //rootAttribute XMLAttribute
    attributes map[string]XMLAttribute
}

func (this *xmlElementImpl) ToElement() XMLElement {
    return this
}

func (this *xmlElementImpl) Accept(visitor XMLVisitor) bool {
    
    if visitor.VisitEnterElement(this) {
        for node := this.FirstChild(); nil != node; node = node.NextSibling() {
            if !node.Accept(visitor) {
                break
            }
        }
    }
    
    return visitor.VisitExitElement(this)
}

func (this *xmlElementImpl) Name() string {
    return this.Value()
}

func (this *xmlElementImpl) SetName(name string) {
    this.SetValue(name)
}

func (this *xmlElementImpl) FindAttribute(name string) XMLAttribute {
    attr, ok := this.attributes[name]
    if !ok {
        return nil
    }
    
    return attr
}

func (this *xmlElementImpl) AttributeCount() int {
    return len(this.attributes)
}

func (this *xmlElementImpl) Attribute(name string, def string) string {
    attr, ok := this.attributes[name]
    if !ok {
        return def
    }
    
    return attr.Value()
}

func (this *xmlElementImpl) SetAttribute(name string, value string) {
    
    attr, ok := this.attributes[name]
    if ok {
        attr.SetValue(value)
        return
    }
    
    this.attributes[name] = NewAttribute(name, value)
}

func (this *xmlElementImpl) DeleteAttribute(name string) {
    delete(this.attributes, name)
}

func (this *xmlElementImpl) Text() string {
    if text := this.FirstChild(); (nil != text) && (nil != text.ToText()) {
        return text.Value()
    }
    
    return ""
}

func (this *xmlElementImpl) SetText(inText string) {
    if node := this.FirstChild(); (nil != node) && (nil != node.ToText()) {
        node.SetValue(inText)
    } else {
        theText := NewText(this.getDocument(), inText)
        this.InsertFirstChild(theText)
    }
}

func (this *xmlElementImpl) ForeachAttribute(callback func(attribute XMLAttribute) int) int {
    for _, value := range this.attributes {
        if ret := callback(value); 0 != ret {
            return ret
        }
    }
    
    return 0
}

//------------------------------------------------------------------

type xmlCommentImpl struct {
    xmlNodeImpl
}

func (this *xmlCommentImpl) ToComment() XMLComment {
    return this
}

func (this *xmlCommentImpl) Accept(visitor XMLVisitor) bool {
    return visitor.VisitComment(this)
}

//------------------------------------------------------------------

type xmlProcInstImpl struct {
    xmlNodeImpl
    instruction string
}

func (this *xmlProcInstImpl) ToProcInst() XMLProcInst {
    return this
}

func (this *xmlProcInstImpl) Accept(visitor XMLVisitor) bool {
    return visitor.VisitProcInst(this)
}

func (this *xmlProcInstImpl) Target() string {
    return this.value
}

func (this *xmlProcInstImpl) Instruction() string {
    return this.instruction
}

//------------------------------------------------------------------

type xmlDocumentImpl struct {
    xmlNodeImpl
}

func (this *xmlDocumentImpl) ToDocument() XMLDocument {
    return this
}
func (this *xmlDocumentImpl) Accept(visitor XMLVisitor) bool {
    
    if visitor.VisitEnterDocument(this) {
        for node := this.FirstChild(); nil != node; node = node.NextSibling() {
            if !node.Accept(visitor) {
                break
            }
        }
    }
    
    return visitor.VisitExitDocument(this)
}

//------------------------------------------------------------------

type xmlTextImpl struct {
    xmlNodeImpl
    cdata bool
}

func (this *xmlTextImpl) ToText() XMLText {
    return this
}
func (this *xmlTextImpl) Accept(visitor XMLVisitor) bool {
    return visitor.VisitText(this)
}
func (this *xmlTextImpl) SetCDATA(isCData bool) {
    this.cdata = isCData
}
func (this *xmlTextImpl) CDATA() bool {
    return this.cdata
}

//------------------------------------------------------------------

type xmlDirectiveImpl struct {
    xmlNodeImpl
}

func (this *xmlDirectiveImpl) ToDirective() XMLDirective {
    return this
}

func (this *xmlDirectiveImpl) Accept(visitor XMLVisitor) bool {
    return visitor.VisitDirective(this)
}

//------------------------------------------------------------------

//	NewText	创建一个新的XMLText对象
func NewText(document XMLDocument, text string) XMLText {
    node := new(xmlTextImpl)
    node.document = document
    node.value = text
    return node
}

//	XMLComment	创建一个新的XMLComment对象
func NewComment(document XMLDocument, comment string) XMLComment {
    node := new(xmlCommentImpl)
    node.document = document
    node.value = comment
    return node
}

//	NewElement	创建一个新的XMLElement对象
func NewElement(document XMLDocument, name string) XMLElement {
    node := new(xmlElementImpl)
    node.document = document
    node.value = name
    node.attributes = make(map[string]XMLAttribute)
    return node
}

//	NewProcInst	创建一个新的XMLProcInst对象
func NewProcInst(document XMLDocument, target string, inst string) XMLProcInst {
    node := new(xmlProcInstImpl)
    node.document = document
    node.value = target
    node.instruction = inst
    return node
}

//	NewDirective	创建一个新的XMLDirective对象
func NewDirective(document XMLDocument, directive string) XMLDirective {
    node := new(xmlDirectiveImpl)
    node.document = document
    node.value = directive
    return node
}

//	NewAttribute	创建一个新的XMLAttribute对象.
//	name和value分别用于指定属性的名称和值
func NewAttribute(name string, value string) XMLAttribute {
    attr := new(xmlAttributeImpl)
    attr.name = name
    attr.value = value
    return attr
}

//	NewDocument	创建一个全新的XMLDocument对象
func NewDocument() XMLDocument {
    node := new(xmlDocumentImpl)
    node.document = node
    return node
}

//	LoadDocument	从rd流中读取XML码流并构建成XMLDocument对象
func LoadDocument(rd io.Reader) (XMLDocument, error) {
    doc := NewDocument()
    var parent XMLNode = doc
    decoder := xml.NewDecoder(rd)
    var token xml.Token
    var err error
    for token, err = decoder.Token(); nil == err; token, err = decoder.Token() {
        switch token.(type) {
        case xml.StartElement:
            startElement := token.(xml.StartElement)
            node := NewElement(doc, startElement.Name.Local)
            for _, item := range startElement.Attr {
                node.SetAttribute(item.Name.Local, item.Value)
            }
            parent.InsertEndChild(node)
            parent = node
        case xml.EndElement:
            //endElement := token.(xml.EndElement)
            parent = parent.Parent()
        case xml.Comment:
            comment := token.(xml.Comment)
            node := NewComment(doc, string(comment))
            parent.InsertEndChild(node)
        case xml.Directive:
            directive := token.(xml.Directive)
            node := NewDirective(doc, string(directive))
            parent.InsertEndChild(node)
        case xml.ProcInst:
            procInst := token.(xml.ProcInst)
            node := NewProcInst(doc, procInst.Target, string(procInst.Inst))
            parent.InsertEndChild(node)
        case xml.CharData:
            charData := token.(xml.CharData)
            if len(bytes.TrimSpace(charData)) > 0 {
                node := NewText(doc, string(charData))
                parent.InsertEndChild(node)
            }
        default:
            return nil, errors.New("Unsupported token type")
        }
    }
    
    if (nil == err) || (io.EOF == err) {
        
        //  不能是空文档
        if nil == doc.FirstChildElement(""){
            return nil, errors.New("XML document missing the root element")
        }
        
        return doc, nil
    }
    
    return nil, err
}

//------------------------------------------------------------------
type xmlSimplePrinter struct {
    writer io.Writer
}

func NewSimplePrinter(writer io.Writer) XMLVisitor {
    visitor := new(xmlSimplePrinter)
    visitor.writer = writer
    return visitor
}

func (this *xmlSimplePrinter) VisitEnterDocument(node XMLDocument) bool {
    return true
}

func (this *xmlSimplePrinter) VisitExitDocument(node XMLDocument) bool {
    return true
}

func (this *xmlSimplePrinter) VisitEnterElement(node XMLElement) bool {
    io.WriteString(this.writer, "<")
    io.WriteString(this.writer, node.Name())
    
    node.ForeachAttribute(func(attribute XMLAttribute) int {
        io.WriteString(this.writer, attribute.Name())
        io.WriteString(this.writer, `="`)
        xml.EscapeText(this.writer, []byte(attribute.Value()))
        io.WriteString(this.writer, `"`)
        return 0
    })
    
    io.WriteString(this.writer, ">")
    
    return true
}

func (this *xmlSimplePrinter) VisitExitElement(node XMLElement) bool {
    io.WriteString(this.writer, "</")
    io.WriteString(this.writer, node.Name())
    io.WriteString(this.writer, ">")
    return true
}

func (this *xmlSimplePrinter) VisitProcInst(node XMLProcInst) bool {
    io.WriteString(this.writer, "<?")
    io.WriteString(this.writer, node.Target())
    io.WriteString(this.writer, " ")
    io.WriteString(this.writer, node.Instruction())
    io.WriteString(this.writer, "?>")
    io.WriteString(this.writer, "\n")
    return true
}

func (this *xmlSimplePrinter) VisitText(node XMLText) bool {
    if node.CDATA() {
        io.WriteString(this.writer, "<![CDATA[")
        io.WriteString(this.writer, node.Value())
        io.WriteString(this.writer, "]]")
        return true
    }
    
    xml.EscapeText(this.writer, []byte(node.Value()))
    return true
}

func (this *xmlSimplePrinter) VisitComment(node XMLComment) bool {
    io.WriteString(this.writer, "<!--")
    xml.EscapeText(this.writer, []byte(node.Value()))
    io.WriteString(this.writer, "-->")
    return true
}

func (this *xmlSimplePrinter) VisitDirective(node XMLDirective) bool {
    io.WriteString(this.writer, "<!")
    xml.EscapeText(this.writer, []byte(node.Value()))
    io.WriteString(this.writer, ">")
    return true
}

//------------------------------------------------------------------
type XMLHandleImpl struct {
    node XMLNode
}

func NewHandle(node XMLNode) XMLHandle {
    handle := new(XMLHandleImpl)
    handle.node = node
    return handle
}

func (this*XMLHandleImpl) Parent() XMLHandle {
    if nil != this.node {
        this.node = this.node.Parent()
    }
    
    return this
}

func (this*XMLHandleImpl) FirstChild() XMLHandle {
    if nil != this.node {
        this.node = this.node.FirstChild()
    }
    
    return this
}

func (this*XMLHandleImpl) LastChild() XMLHandle {
    if nil != this.node {
        this.node = this.node.LastChild()
    }
    
    return this
}

func (this*XMLHandleImpl) PreviousSibling() XMLHandle {
    if nil != this.node {
        this.node = this.node.PreviousSibling()
    }
    
    return this
}

func (this*XMLHandleImpl) NextSibling() XMLHandle {
    if nil != this.node {
        this.node = this.node.NextSibling()
    }
    
    return this
}

func (this*XMLHandleImpl) FirstChildElement(name string) XMLHandle {
    if nil != this.node {
        this.node = this.node.FirstChildElement(name)
    }
    
    return this
}

func (this*XMLHandleImpl) LastChildElement(name string) XMLHandle {
    if nil != this.node {
        this.node = this.node.LastChildElement(name)
    }
    
    return this
}

func (this*XMLHandleImpl) PreviousSiblingElement(name string) XMLHandle {
    if nil != this.node {
        this.node = this.node.PreviousSiblingElement(name)
    }
    
    return this
}

func (this*XMLHandleImpl) NextSiblingElement(name string) XMLHandle {
    if nil != this.node {
        this.node = this.node.NextSiblingElement(name)
    }
    
    return this
}

func (this*XMLHandleImpl) ToElement() XMLElement {
    if nil != this.node {
        return this.node.ToElement()
    }
    
    return nil
}

func (this*XMLHandleImpl) ToText() XMLText {
    if nil != this.node {
        return this.node.ToText()
    }
    
    return nil
}

func (this*XMLHandleImpl) ToComment() XMLComment {
    if nil != this.node {
        return this.node.ToComment()
    }
    
    return nil
}

func (this*XMLHandleImpl) ToDocument() XMLDocument {
    if nil != this.node {
        return this.node.ToDocument()
    }
    
    return nil
}

func (this*XMLHandleImpl) ToProcInst() XMLProcInst {
    if nil != this.node {
        return this.node.ToProcInst()
    }
    
    return nil
}

func (this*XMLHandleImpl) ToDirective() XMLDirective {
    if nil != this.node {
        return this.node.ToDirective()
    }
    
    return nil
}
