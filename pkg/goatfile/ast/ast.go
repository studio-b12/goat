package ast

type Pos int

type Goatfile struct {
	Comments   []Comment
	Imports    []Import
	Separators []Separator
	Actions    []Action
	Sections   []Section
}

type Comment struct {
	Pos     Pos
	Content string
}

type Import struct {
	Pos  Pos
	Path string
}

type Separator struct {
	Pos Pos
	Len int
}

type Action interface{}

type Section interface{}

type SectionSetup struct {
	Pos     Pos
	Actions []Action
}

type SectionTests struct {
	Pos     Pos
	Actions []Action
}

type SectionTeardown struct {
	Pos     Pos
	Actions []Action
}

type SectionDefaults struct {
	Pos     Pos
	Request PartialRequest
}

type LogSection struct {
	Pos     Pos
	Content string
}

type Execute struct {
	Pos        Pos
	Path       string
	Parameters KV
	Returns    Assignments
}

type KV map[string]any

type Assignments map[string]string

type Request struct {
	Pos    Pos
	Head   RequestHead
	Blocks []RequestBlock
}

type PartialRequest struct {
	Pos    Pos
	Blocks []RequestBlock
}

type HeaderEntries map[string]string

type TextBlock string

type RequestHead struct {
	Method string
	Url    string
}

type RequestBlock interface{}

type RequestOptions KV

type RequestHeaders HeaderEntries

type RequestQueryParams KV

type RequestAuth KV

type RequestPreScript TextBlock

type RequestScript TextBlock

type FileDescriptor string
