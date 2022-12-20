// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syntax

// Syntax ={ Production }
// go 语法是产生式的集合
// Production = production_name "=" "[ Expression]" "."
// 产生式格式:  xxx(production_name) =  表达式
// Expresion = Term{"|" Term}
// 一个或多个终结符的集合组成
// Term = Factor { Factor}
// 一个或多个因子组成
//Factor= production_name | token [ "…" token ] | Group | Option | Repetition .
// Group       = "(" Expression ")" .
// Option      = "[" Expression "]" .
// Repetition  = "{" Expression "}" .

// ----------------------------------------------------------------------------
// Nodes

type Node interface {
    //抽象语法树节点
	// Pos() returns the position associated with the node as follows:\
    //节点关联位置
	// 1) The position of a node representing a terminal syntax production
	//    (Name, BasicLit, etc.) is the position of the respective production
	//    in the source.
    //  产生式在源代码中的位置, 代表的是 terminal syntax production 

	// 2) The position of a node representing a non-terminal production
	//    (IndexExpr, IfStmt, etc.) is the position of a token uniquely
	//    associated with that production; usually the left-most one
	//    ('[' for IndexExpr, 'if' for IfStmt, etc.)
    //  非终端产生式 的话则是最走边括号的位置.

	Pos() Pos
	aNode()
}

type node struct {
	// commented out for now since not yet used
    //这个注释的是不再使用的了
	// doc  *Comment // nil means no comment(s) attached
    // doc *Comment 代表没有注释出现
	pos Pos

}

func (n *node) Pos() Pos { return n.pos }
func (*node) aNode()     {}

// ----------------------------------------------------------------------------
// Files

// package PkgName; DeclList[0], DeclList[1], ...
type File struct {
    // ast 解析树的根节点

    //指令
	Pragma   Pragma
    //包名
	PkgName  *Name
    // 声明列表
	DeclList []Decl
    //结束位置
	EOF      Pos
    //组合了节点,就像继承才不多的用法.
	node
}

// ----------------------------------------------------------------------------
// Declarations
//一堆结构体
//Declaration   = ConstDecl | TypeDecl | VarDecl .
//TopLevelDecl  = Declaration | FunctionDecl | MethodDecl .

type (
    // Decl 声明语法节点
	Decl interface {
		Node
		aDecl()
	}

	//              Path
	// LocalPkgName Path
    //ImportDecl 语法节点
	ImportDecl struct {
        // 
		Group        *Group // nil means not part of a group
		//
        Pragma       Pragma
		// 本地包名词
        LocalPkgName *Name     // including "."; nil means no rename present
		//路径
        Path         *BasicLit // Path.Bad || Path.Kind == StringLit; nil means no path
		//实现了声明
        decl
	}

	// NameList
	// NameList      = Values
	// NameList Type = Values

    //常量声明
	ConstDecl struct {
        
		Group    *Group // nil means not part of a group
		
        Pragma   Pragma
		//名称列表
        NameList []*Name
		// 类型表达式
        Type     Expr // nil means no type
		// 值表达公式
        Values   Expr // nil means no values
		//实现了声明
        decl
	}

	// Name Type
	TypeDecl struct {

		Group      *Group // nil means not part of a group
		Pragma     Pragma
		//类型名
        Name       *Name
		//类型参数 
        TParamList []*Field // nil means no type parameters
		//别名
        Alias      bool
		//类型表达式
        Type       Expr
		//实现了声明
        decl
	}

	// NameList Type
	// NameList Type = Values
	// NameList      = Values
	VarDecl struct {

		Group    *Group // nil means not part of a group
		
        Pragma   Pragma
		
        NameList []*Name
		
        Type     Expr // nil means no type
		
        Values   Expr // nil means no values
		
        decl
	}

	// func          Name Type { Body }
	// func          Name Type
	// func Receiver Name Type { Body }
	// func Receiver Name Type
	FuncDecl struct {
		Pragma     Pragma
		Recv       *Field // nil means regular function
		Name       *Name
		TParamList []*Field // nil means no type parameters
		Type       *FuncType
		Body       *BlockStmt // nil means no body (forward declaration)
		decl
	}
)

type decl struct{ node }

func (*decl) aDecl() {}

// All declarations belonging to the same group point to the same Group node.
type Group struct {
	_ int // not empty so we are guaranteed different Group instances
}

// ----------------------------------------------------------------------------
// Expressions

func NewName(pos Pos, value string) *Name {
	n := new(Name)
	n.pos = pos
	n.Value = value
	return n
}

type (
	Expr interface {
		Node
		aExpr()
	}

	// Placeholder for an expression that failed to parse
	// correctly and where we can't provide a better node.
	BadExpr struct {
		expr
	}

	// Value
	Name struct {
		Value string
		expr
	}

	// Value
	BasicLit struct {
		Value string
		Kind  LitKind
		Bad   bool // true means the literal Value has syntax errors
		expr
	}

	// Type { ElemList[0], ElemList[1], ... }
	CompositeLit struct {
		Type     Expr // nil means no literal type
		ElemList []Expr
		NKeys    int // number of elements with keys
		Rbrace   Pos
		expr
	}

	// Key: Value
	KeyValueExpr struct {
		Key, Value Expr
		expr
	}

	// func Type { Body }
	FuncLit struct {
		Type *FuncType
		Body *BlockStmt
		expr
	}

	// (X)
	ParenExpr struct {
		X Expr
		expr
	}

	// X.Sel
	SelectorExpr struct {
		X   Expr
		Sel *Name
		expr
	}

	// X[Index]
	// X[T1, T2, ...] (with Ti = Index.(*ListExpr).ElemList[i])
	IndexExpr struct {
		X     Expr
		Index Expr
		expr
	}

	// X[Index[0] : Index[1] : Index[2]]
	SliceExpr struct {
		X     Expr
		Index [3]Expr
		// Full indicates whether this is a simple or full slice expression.
		// In a valid AST, this is equivalent to Index[2] != nil.
		// TODO(mdempsky): This is only needed to report the "3-index
		// slice of string" error when Index[2] is missing.
		Full bool
		expr
	}

	// X.(Type)
	AssertExpr struct {
		X    Expr
		Type Expr
		expr
	}

	// X.(type)
	// Lhs := X.(type)
	TypeSwitchGuard struct {
		Lhs *Name // nil means no Lhs :=
		X   Expr  // X.(type)
		expr
	}

	Operation struct {
		Op   Operator
		X, Y Expr // Y == nil means unary expression
		expr
	}

	// Fun(ArgList[0], ArgList[1], ...)
	CallExpr struct {
		Fun     Expr
		ArgList []Expr // nil means no arguments
		HasDots bool   // last argument is followed by ...
		expr
	}

	// ElemList[0], ElemList[1], ...
	ListExpr struct {
		ElemList []Expr
		expr
	}

	// [Len]Elem
	ArrayType struct {
		// TODO(gri) consider using Name{"..."} instead of nil (permits attaching of comments)
		Len  Expr // nil means Len is ...
		Elem Expr
		expr
	}

	// []Elem
	SliceType struct {
		Elem Expr
		expr
	}

	// ...Elem
	DotsType struct {
		Elem Expr
		expr
	}

	// struct { FieldList[0] TagList[0]; FieldList[1] TagList[1]; ... }
	StructType struct {
		FieldList []*Field
		TagList   []*BasicLit // i >= len(TagList) || TagList[i] == nil means no tag for field i
		expr
	}

	// Name Type
	//      Type
	Field struct {
		Name *Name // nil means anonymous field/parameter (structs/parameters), or embedded element (interfaces)
		Type Expr  // field names declared in a list share the same Type (identical pointers)
		node
	}

	// interface { MethodList[0]; MethodList[1]; ... }
	InterfaceType struct {
		MethodList []*Field
		expr
	}

	FuncType struct {
		ParamList  []*Field
		ResultList []*Field
		expr
	}

	// map[Key]Value
	MapType struct {
		Key, Value Expr
		expr
	}

	//   chan Elem
	// <-chan Elem
	// chan<- Elem
	ChanType struct {
		Dir  ChanDir // 0 means no direction
		Elem Expr
		expr
	}
)

type expr struct{ node }

func (*expr) aExpr() {}

type ChanDir uint

const (
	_ ChanDir = iota
	SendOnly
	RecvOnly
)

// ----------------------------------------------------------------------------
// Statements

type (
	Stmt interface {
		Node
		aStmt()
	}

	SimpleStmt interface {
		Stmt
		aSimpleStmt()
	}

	EmptyStmt struct {
		simpleStmt
	}

	LabeledStmt struct {
		Label *Name
		Stmt  Stmt
		stmt
	}

	BlockStmt struct {
		List   []Stmt
		Rbrace Pos
		stmt
	}

	ExprStmt struct {
		X Expr
		simpleStmt
	}

	SendStmt struct {
		Chan, Value Expr // Chan <- Value
		simpleStmt
	}

	DeclStmt struct {
		DeclList []Decl
		stmt
	}

	AssignStmt struct {
		Op       Operator // 0 means no operation
		Lhs, Rhs Expr     // Rhs == nil means Lhs++ (Op == Add) or Lhs-- (Op == Sub)
		simpleStmt
	}

	BranchStmt struct {
		Tok   token // Break, Continue, Fallthrough, or Goto
		Label *Name
		// Target is the continuation of the control flow after executing
		// the branch; it is computed by the parser if CheckBranches is set.
		// Target is a *LabeledStmt for gotos, and a *SwitchStmt, *SelectStmt,
		// or *ForStmt for breaks and continues, depending on the context of
		// the branch. Target is not set for fallthroughs.
		Target Stmt
		stmt
	}

	CallStmt struct {
		Tok  token // Go or Defer
		Call Expr
		stmt
	}

	ReturnStmt struct {
		Results Expr // nil means no explicit return values
		stmt
	}

	IfStmt struct {
		Init SimpleStmt
		Cond Expr
		Then *BlockStmt
		Else Stmt // either nil, *IfStmt, or *BlockStmt
		stmt
	}

	ForStmt struct {

		Init SimpleStmt // incl. *RangeClause
		//表达式
        Cond Expr
		//简单句
        Post SimpleStmt
		//块语句
        Body *BlockStmt
		//实现了语句
        stmt
	}

	SwitchStmt struct {
		Init   SimpleStmt
		Tag    Expr // incl. *TypeSwitchGuard
		Body   []*CaseClause
		Rbrace Pos
		stmt
	}

	SelectStmt struct {
		Body   []*CommClause
		Rbrace Pos
		stmt
	}
)

type (
	RangeClause struct {
		Lhs Expr // nil means no Lhs = or Lhs :=
		Def bool // means :=
		X   Expr // range X
		simpleStmt
	}

	CaseClause struct {
		Cases Expr // nil means default clause
		Body  []Stmt
		Colon Pos
		node
	}

	CommClause struct {
		Comm  SimpleStmt // send or receive stmt; nil means default clause
		Body  []Stmt
		Colon Pos
		node
	}
)

type stmt struct{ node }

func (stmt) aStmt() {}

type simpleStmt struct {
	stmt
}

func (simpleStmt) aSimpleStmt() {}

// ----------------------------------------------------------------------------
// Comments

// TODO(gri) Consider renaming to CommentPos, CommentPlacement, etc.
// Kind = Above doesn't make much sense.
type CommentKind uint

const (
	Above CommentKind = iota
	Below
	Left
	Right
)

type Comment struct {
    //注释类型
	Kind CommentKind
    //文本
	Text string
	//下个注释
    Next *Comment
}
