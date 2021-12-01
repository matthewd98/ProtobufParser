package main

type Documentation struct {
	Repo   Repo     `json:"repo"`
	Schema []Schema `json:"schema"`
}

type Repo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Schema struct {
	FilePath    string    `json:"filePath"`
	FileName    string    `json:"fileName"`
	Url         string    `json:"url"`
	PackageName string    `json:"packageName"`
	Enums       []Enum    `json:"enums"`
	Messages    []Message `json:"messages"`
	Services    []Service `json:"services"`
}

/* Enum */

type Enum struct {
	Name    string      `json:"name"`
	Comment string      `json:"comment"`
	Values  []EnumValue `json:"values"`
}

type EnumValue struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

/* Message */

type Message struct {
	Name           string          `json:"name"`
	Comment        string          `json:"comment"`
	Extensions     FieldExtensions `json:"extensions"`
	Fields         []MessageField  `json:"fields"`
	OneOfs         []OneOf         `json:"oneofs"`
	NestedMessages []Message       `json:"nestedMessages"`
}

type FieldExtensions struct {
	MinTag int `json:"minTag"`
	MaxTag int `json:"maxTag"`
}

type MessageField struct {
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	Type        string `json:"type"`
	Tag         int    `json:"tag"`
	IsRequired  bool   `json:"isRequired"`
	IsRepeated  bool   `json:"isRepeated"`
	IsExtension bool   `json:"isExtension"`
	Annotation  string `json:"annotation"` // e.g. [packed:true], [deprecated:true]
}

type OneOf struct {
	Name    string       `json:"name"`
	Comment string       `json:"comment"`
	Fields  []OneOfField `json:"fields"`
}

type OneOfField struct {
	Name       string `json:"name"`
	Comment    string `json:"comment"`
	Type       string `json:"type"`
	Tag        int    `json:"tag"`
	IsRepeated bool   `json:"isRepeated"`
	Annotation string `json:"annotation"` // e.g. [packed:true], [deprecated:true]
}

/* Service */

type Service struct {
	Name                 string `json:"name"`
	Comment              string `json:"comment"`
	RemoteProcedureCalls []Rpc  `json:"rpcs"`
}

type Rpc struct {
	Name      string  `json:"name"`
	Comment   string  `json:"comment"`
	RpcInput  RpcType `json:"rpcInput"`
	RpcOutput RpcType `json:"rpcOutput"`
}

type RpcType struct {
	Type     string `json:"type"`
	IsStream bool   `json:"isStream"`
}
