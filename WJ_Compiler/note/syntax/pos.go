// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syntax

import "fmt"

// PosMax 是可以无损表示的最大行或列值。
// 大于 PosMax 的传入值（参数）将被设置为 PosMax。
// PosMax is the largest line or column value that can be represented without loss.
// Incoming values (arguments) larger than PosMax will be set to PosMax.
const PosMax = 1 << 30
// 1262485504 十亿级别

// A Pos represents an absolute (line, col) source position
// with a reference to position base for computing relative
// (to a file, or line directive) position information.
// Pos 表示绝对（行，列）源位置，并引用位置基来计算相对（到文件或行指令）位置信息。

// Pos values are intentionally light-weight so that they
// can be created without too much concern about space use.
//占用内存少

type Pos struct {
    //在文件中的开始位置
	base      *PosBase
    // 行列,uint 可以表示的范围0~4294967295 四十多个亿.
    // 通常一本书也就是百万级别,而UTF-8 三个字节表示 表示一个中文字,
    //  要表示一个纯文本的百万字级别的需要的磁盘空间大概要多少?
    // 一百万级别 大概就是 千*千 1MB就是 1024*1024 字节*3 就是 3M就能存下百万个中文字符. 
    // 如果换成数据库 ,一条数据 假设是一百个字符,也就是 300B , 四条数据就是 1K左右, 
    // 百万条数据就是 1K * 1MB 也就是 1 GB 左右,数据库就是个吃磁盘的怪兽.
    // 一块2T 的磁盘 可以为一百万人,提供 400字符*1024*2( 80万字 左右) ,在忽略数据库本身占用的字节数的情况下.

	line, col uint32
}

// MakePos returns a new Pos for the given PosBase, line and column.
// 
func MakePos(base *PosBase, line, col uint) Pos { return Pos{base, sat32(line), sat32(col)} }

// TODO(gri) IsKnown makes an assumption about linebase < 1.
// Maybe we should check for Base() != nil instead.

//返回结构体位置
func (pos Pos) Pos() Pos       { return pos }
//返回 位置是否被识别
func (pos Pos) IsKnown() bool  { return pos.line > 0 }
//返回位置的基准点
func (pos Pos) Base() *PosBase { return pos.base }
// 返回位置的行
func (pos Pos) Line() uint     { return uint(pos.line) }
// 返回位置的列
func (pos Pos) Col() uint      { return uint(pos.col) }
// 返回位置引用的文件名
func (pos Pos) RelFilename() string { return pos.base.Filename() }

//
func (pos Pos) RelLine() uint {
	b := pos.base
	if b.Line() == 0 {
		// base line is unknown => relative line is unknown
		return 0
	}
	return b.Line() + (pos.Line() - b.Pos().Line())
}

func (pos Pos) RelCol() uint {
	b := pos.base
	if b.Col() == 0 {
		// base column is unknown => relative column is unknown
		// (the current specification for line directives requires
		// this to apply until the next PosBase/line directive,
		// not just until the new newline)
		return 0
	}
	if pos.Line() == b.Pos().Line() {
		// pos on same line as pos base => column is relative to pos base
		return b.Col() + (pos.Col() - b.Pos().Col())
	}
	return pos.Col()
}

// 比较两个位置的大小
// Cmp compares the positions p and q and returns a result r as follows:
//
//	r <  0: p is before q
//	r == 0: p and q are the same position (but may not be identical)
//	r >  0: p is after q
//
// If p and q are in different files, p is before q if the filename
// of p sorts lexicographically before the filename of q.
func (p Pos) Cmp(q Pos) int {
	pname := p.RelFilename()
	qname := q.RelFilename()
	switch {
	case pname < qname:
		return -1
	case pname > qname:
		return +1
	}

	pline := p.Line()
	qline := q.Line()
	switch {
	case pline < qline:
		return -1
	case pline > qline:
		return +1
	}

	pcol := p.Col()
	qcol := q.Col()
	switch {
	case pcol < qcol:
		return -1
	case pcol > qcol:
		return +1
	}

	return 0
}


//-------------------------多态用用法,不同结构可以拥有私有的同名函数-------------------------------------
func (pos Pos) String() string {
	rel := position_{pos.RelFilename(), pos.RelLine(), pos.RelCol()}
	abs := position_{pos.Base().Pos().RelFilename(), pos.Line(), pos.Col()}
	s := rel.String()
	if rel != abs {
		s += "[" + abs.String() + "]"
	}
	return s
}

// TODO(gri) cleanup: find better name, avoid conflict with position in error_test.go
type position_ struct {
	filename  string
	line, col uint
}

func (p position_) String() string {
	if p.line == 0 {
		if p.filename == "" {
			return "<unknown position>"
		}
		return p.filename
	}
	if p.col == 0 {
		return fmt.Sprintf("%s:%d", p.filename, p.line)
	}
	return fmt.Sprintf("%s:%d:%d", p.filename, p.line, p.col)
}
//-------------------------------------------------------------------


// A PosBase represents the base for relative position information:
// PosBase 代表了 基本的相对位置信息
// At position pos, the relative position is filename:line:col.
// 在Position pos 中 ,相对位置是文件名:line:col
type PosBase struct {
    // 位置
	pos       Pos
    //文件名
	filename  string
    //行.列
	line, col uint32
    //
	trimmed   bool // whether -trimpath has been applied
}

// NewFileBase returns a new PosBase for the given filename.
// A file PosBase's position is relative to itself, with the
// position being filename:1:1.
// 为一个文件新建一个PosBase ,postion 将长成 filename:1:1 的样子
func NewFileBase(filename string) *PosBase {
    //
	return NewTrimmedFileBase(filename, false)
}


// NewTrimmedFileBase is like NewFileBase, but allows specifying Trimmed.
func NewTrimmedFileBase(filename string, trimmed bool) *PosBase {
	//新建一个base
    base := &PosBase{MakePos(nil, linebase, colbase), filename, linebase, colbase, trimmed}
	base.pos.base = base
	return base
}

// NewLineBase returns a new PosBase for a line directive "line filename:line:col"
// relative to pos, which is the position of the character immediately following
// the comment containing the line directive. For a directive in a line comment,
// that position is the beginning of the next line (i.e., the newline character
// belongs to the line comment).
func NewLineBase(pos Pos, filename string, trimmed bool, line, col uint) *PosBase {
	return &PosBase{pos, filename, sat32(line), sat32(col), trimmed}
}

func (base *PosBase) IsFileBase() bool {
	if base == nil {
		return false
	}
	return base.pos.base == base
}

func (base *PosBase) Pos() (_ Pos) {
	if base == nil {
		return
	}
	return base.pos
}

//返回 文件名
func (base *PosBase) Filename() string {
	if base == nil {
		return ""
	}
	return base.filename
}

//返回行
func (base *PosBase) Line() uint {
	if base == nil {
		return 0
	}
	return uint(base.line)
}

// 返回Col
// 为什么 uint32要强制转成uint?
// uint32 和 uint 的区别
//uint 看机器字长,64位机器是64位,32位机是32位

func (base *PosBase) Col() uint {
	if base == nil {
		return 0
	}
	return uint(base.col)
}

//
func (base *PosBase) Trimmed() bool {
	if base == nil {
		return false
	}
	return base.trimmed
}

// 输入 一个uint ,返回一个uint 
// 返回比较小的那一个

func sat32(x uint) uint32 {
	if x > PosMax {
		return PosMax
	}
	return uint32(x)
}
