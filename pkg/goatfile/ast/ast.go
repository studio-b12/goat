package ast

type Pos struct {
	Pos     int
	Line    int
	LinePos int
}

type Goatfile struct {
	Dir string

	Comments   []Comment
	Imports    []Import
	Delimiters []Delimiter
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

type Delimiter struct {
	Pos      Pos
	ExtraLen int
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
	Parameters KVList[any]
	Returns    Assignments
}

type KVList[TVal any] []KV[TVal]

func (t KVList[TVal]) Get(key string) (val TVal, ok bool) {
	for _, kv := range t {
		if kv.Key == key {
			return kv.Value, true
		}
	}
	return val, false
}

func (t KVList[TVal]) GetUnchecked(key string) TVal {
	val, _ := t.Get(key)
	return val
}

func (t KVList[TVal]) ToMap() map[string]TVal {
	m := make(map[string]TVal)
	for _, kv := range t {
		m[kv.Key] = kv.Value
	}
	return m
}

func (t KVList[TVal]) ToMultiMap() map[string][]TVal {
	m := make(map[string][]TVal)
	for _, kv := range t {
		m[kv.Key] = append(m[kv.Key], kv.Value)
	}
	return m
}

type Assignments struct {
	KVList[string]
}

type KV[TVal any] struct {
	Pos   Pos
	Key   string
	Value TVal
}

type Request struct {
	Pos    Pos
	Head   RequestHead
	Blocks []RequestBlock
}

type PartialRequest struct {
	Pos    Pos
	Blocks []RequestBlock
}

type HeaderEntries struct {
	KVList[string]
}

type DataContent interface {
}

type TextBlock struct {
	Content string
}

type FileDescriptor struct {
	Path        string
	ContentType string
}

type VarDescriptor struct {
	VarName     string
	ContentType string
}

type NoContent struct{}

type RequestHead struct {
	Method string
	Url    string
}

type RequestBlock interface{}

type RequestOptions struct {
	KVList[any]
}

type RequestHeader struct {
	HeaderEntries
}

type RequestQueryParams struct {
	KVList[any]
}

type RequestAuth struct {
	KVList[any]
}

type RequestBody struct {
	DataContent
}

type RequestPreScript struct {
	DataContent
}

type RequestScript struct {
	DataContent
}

type FormData struct {
	KVList[any]
}
