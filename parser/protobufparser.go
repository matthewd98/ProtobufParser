package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parse(protoFileContent string) *Schema {
	schema := &Schema{}

	reader := strings.NewReader(protoFileContent)
	currentWord, comment := "", ""

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		// Ignore whitespaces and new lines, and comments without a space between the keyword and body (e.g. //MyComment).
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' || currentWord == "//" || currentWord == "/*" {
			if currentWord == "" {
				continue
			} else {
				// Match token type
				switch currentWord {
				case "syntax":
					{
						statement := extractStatement(reader)
						parseAndVerifySyntax(statement)
						comment = ""
					}
				case "package":
					{
						statement := extractStatement(reader)
						schema.PackageName = parsePackage(statement)
						comment = ""
					}
				case "enum":
					{
						blockName, blockBody := extractBlock(reader)
						enum := parseEnum(blockName, blockBody, comment)
						schema.Enums = append(schema.Enums, *enum)
						comment = ""
					}
				case "message":
					{
						blockName, blockBody := extractBlock(reader)
						msg, enums := parseMessage(blockName, blockBody, comment)

						schema.Messages = append(schema.Messages, *msg)
						schema.Enums = append(schema.Enums, enums...)
						comment = ""
					}
				case "service":
					{
						blockName, blockBody := extractBlock(reader)
						service := parseService(blockName, blockBody, comment)
						schema.Services = append(schema.Services, *service)
						comment = ""
					}
				case "//":
					{
						if comment == "" {
							comment = readToEndOfLine(reader)
						} else {
							comment += fmt.Sprintf("\n%s", readToEndOfLine(reader))
						}
					}
				case "/*":
					{
						if comment == "" {
							comment = readMultiLineComment(reader)
						} else {
							comment += fmt.Sprintf("\n%s", readMultiLineComment(reader))
						}
					}
				case "import":
					{
						// Discard
						extractStatement(reader)
					}
				case "option":
					{
						// Discard
						extractStatement(reader)
					}
				default:
					{
						//fmt.Printf("Keyword %s not recognized\n", currentWord)
						//os.Exit(1)
					}
				}

				currentWord = ""
				continue
			}
		}

		currentWord += string(r)
	}

	return schema
}

// Read to end of line for a single-line comment
func readToEndOfLine(r *strings.Reader) string {
	restOfLine := ""

	for {
		r, _, err := r.ReadRune()
		if err == io.EOF || r == '\r' || r == '\n' {
			break
		}
		restOfLine += string(r)
	}

	return restOfLine
}

// Read multi-line comment, i.e. /* This is a comment\n  */
func readMultiLineComment(r *strings.Reader) string {
	comment := ""
	sequence := ""

	for {
		r, _, err := r.ReadRune()
		if err == io.EOF {
			break
		} else if r == '*' || r == '/' {
			sequence += string(r)
			if sequence == "*/" {
				break
			}
		} else {
			sequence = ""
		}

		comment += string(r)
	}

	return strings.TrimRight(comment[0:len(comment)-1], " ")
}

// Statement is a one-liner that ends with ';'
func extractStatement(r *strings.Reader) string {
	statement := ""

	for {
		r, _, err := r.ReadRune()
		if err == io.EOF || r == ';' {
			break
		}
		statement += string(r)
	}

	return statement
}

// Extract RPC statement which can end with a semicolon or {}
func extractRpcStatement(r *strings.Reader) string {
	rpcStatement := ""
	sequence := ""

	for {
		r, _, err := r.ReadRune()
		if err == io.EOF {
			break
		} else if r == ';' {
			break
		} else if r == '{' || r == '}' {
			sequence += string(r)
			if sequence == "{}" {
				break
			}
		} else {
			sequence = ""
		}

		rpcStatement += string(r)
	}

	return strings.TrimRight(rpcStatement[0:len(rpcStatement)-1], " ")
}

// Block is encapsulated by {}
func extractBlock(reader *strings.Reader) (blockName string, blockBody string) {
	blockName, blockBody = "", ""
	bracketCount, blockStarted := 0, false

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		switch r {
		case '{':
			blockStarted = true
			bracketCount++
		case '}':
			bracketCount--
		}

		if blockStarted {
			blockBody += string(r)
		} else if r != ' ' && r != '\r' && r != '\n' && r != '\t' {
			blockName += string(r)
		}

		if blockStarted && bracketCount == 0 {
			break
		}
	}

	return blockName, blockBody[1 : len(blockBody)-1]
}

func parseAndVerifySyntax(statement string) {
	if !strings.Contains(statement, "proto3") && !strings.Contains(statement, "proto2") {
		fmt.Println(`Only "proto3" syntax is supported`)
		os.Exit(1)
	}
}

func parsePackage(statement string) string {
	r, _ := regexp.Compile(`\s*(?P<PACKAGENAME>\w+)`)
	packageName := r.FindStringSubmatch(statement)[1]
	return packageName
}

func parseEnum(blockName string, blockBody string, comment string) *Enum {
	enum := &Enum{
		Name:    blockName,
		Comment: comment,
	}

	r, _ := regexp.Compile(`(?P<KEY>\w+)\s*=\s*(?P<VALUE>\d+)`)
	matches := r.FindAllStringSubmatch(blockBody, -1)
	for _, m := range matches {
		value, _ := strconv.Atoi(m[2])
		enumValue := &EnumValue{
			Name:  m[1],
			Value: value,
		}
		enum.Values = append(enum.Values, *enumValue)
	}

	return enum
}

/* parseMessage returns a message and a slice of enums because the enums aren't nested objects
in the Message struct, they are nested under the Enum field in the Documentation struct */
func parseMessage(msgName string, msgBlockBody string, comment string) (*Message, []Enum) {
	enums := make([]Enum, 0)
	msg := &Message{
		Name:    msgName,
		Comment: comment,
	}

	// TODO: handle extends
	// if blockName == "" {
	// }

	reader := strings.NewReader(msgBlockBody)
	currentWord, comment := "", ""

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		// Ignore whitespaces and new lines, and comments without a space between the keyword and body (e.g. //MyComment).
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' || currentWord == "//" || currentWord == "/*" {
			if currentWord == "" {
				continue
			} else {
				// Match token type
				switch currentWord {
				case "enum":
					{
						blockName, blockBody := extractBlock(reader)
						nestedEnumName := fmt.Sprintf("%s.%s", msgName, blockName)
						nestedEnum := parseEnum(nestedEnumName, blockBody, comment)
						enums = append(enums, *nestedEnum)
						comment = ""
					}
				case "message":
					{
						blockName, blockBody := extractBlock(reader)
						nestedMsgName := fmt.Sprintf("%s.%s", msgName, blockName)
						nestedMessage, nestedEnums := parseMessage(nestedMsgName, blockBody, comment)

						msg.NestedMessages = append(msg.NestedMessages, *nestedMessage)
						enums = append(enums, nestedEnums...)
						comment = ""
					}
				case "oneof":
					{
						blockName, blockBody := extractBlock(reader)
						oneOf := parseMessageOneOf(blockName, blockBody, comment)
						msg.OneOfs = append(msg.OneOfs, *oneOf)
						comment = ""
					}
				case "reserved":
					{
						// Discard
						extractStatement(reader)
					}
				case "extensions":
					{
						statement := extractStatement(reader)
						parseExtensions(statement)
					}
				case "option":
					{
						// Discard
						extractStatement(reader)
					}
				case "//":
					{
						if comment == "" {
							comment = readToEndOfLine(reader)
						} else {
							comment += fmt.Sprintf("\n%s", readToEndOfLine(reader))
						}
					}
				case "/*":
					{
						if comment == "" {
							comment = readMultiLineComment(reader)
						} else {
							comment += fmt.Sprintf("\n%s", readMultiLineComment(reader))
						}
					}
				case "required":
					{
						statement := extractStatement(reader)
						field := parseMessageField(statement, true, false, comment)
						msg.Fields = append(msg.Fields, *field)
						comment = ""
					}
				case "optional":
					{
						statement := extractStatement(reader)
						field := parseMessageField(statement, false, false, comment)
						msg.Fields = append(msg.Fields, *field)
						comment = ""
					}
				// Repeated field is inherently optional
				case "repeated":
					{
						statement := extractStatement(reader)
						field := parseMessageField(statement, false, true, comment)
						msg.Fields = append(msg.Fields, *field)
						comment = ""
					}
				// Optional, required and repeated are omitted
				default:
					{
						statement := fmt.Sprintf("%s %s", currentWord, extractStatement(reader))
						field := parseMessageField(statement, false, false, comment)
						msg.Fields = append(msg.Fields, *field)
						comment = ""
					}
				}

				currentWord = ""
				continue
			}
		}

		currentWord += string(r)
	}

	return msg, enums
}

func parseService(blockName string, blockBody string, comment string) *Service {
	service := &Service{
		Name:    blockName,
		Comment: comment,
	}

	reader := strings.NewReader(blockBody)
	currentWord, comment := "", ""

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		// Ignore whitespaces and new lines, and comments without a space between the keyword and body (e.g. //MyComment).
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' || currentWord == "//" || currentWord == "/*" {
			if currentWord == "" {
				continue
			} else {
				// Match token type
				switch currentWord {
				case "rpc":
					{
						statement := extractRpcStatement(reader)
						rpc := parseRpc(statement)
						service.RemoteProcedureCalls = append(service.RemoteProcedureCalls, *rpc)
						comment = ""
					}
				case "//":
					{
						if comment == "" {
							comment = readToEndOfLine(reader)
						} else {
							comment += fmt.Sprintf("\n%s", readToEndOfLine(reader))
						}
					}
				case "/*":
					{
						if comment == "" {
							comment = readMultiLineComment(reader)
						} else {
							comment += fmt.Sprintf("\n%s", readMultiLineComment(reader))
						}
					}
				}

				currentWord = ""
				continue
			}
		}

		currentWord += string(r)
	}

	return service
}

func parseRpc(statement string) *Rpc {
	r, _ := regexp.Compile(`(?P<RPCNAME>\w+)\s*\((?P<RPCINPUT>[\w|\s|\.]+)\)\s+returns\s+\((?P<RPCOUTPUT>[\w|\s|\.]+)\)`)
	matches := r.FindStringSubmatch(statement)

	rpc := &Rpc{
		Name: matches[1],
		RpcInput: RpcType{
			Type:     strings.Trim(strings.Trim(matches[2], "stream"), " "),
			IsStream: strings.Contains(matches[2], "stream"),
		},
		RpcOutput: RpcType{
			Type:     strings.Trim(strings.Trim(matches[3], "stream"), " "),
			IsStream: strings.Contains(matches[3], "stream"),
		},
	}

	return rpc
}

func parseMessageField(statement string, isRequired bool, isRepeated bool, comment string) *MessageField {
	field := &MessageField{
		IsRequired: isRequired,
		IsRepeated: isRepeated,
		Comment:    comment,
	}

	r, _ := regexp.Compile(`(?P<FIELDTYPE>[\w|\<|\>|,]+)\s+(?P<FIELDNAME>\w+)\s*=\s*(?P<FIELDTAG>\d+)`)
	matches := r.FindStringSubmatch(statement)
	field.Type, field.Name = matches[1], matches[2]
	field.Tag, _ = strconv.Atoi(matches[3])

	return field
}

func parseMessageOneOf(oneOfName string, oneBlockBody string, comment string) *OneOf {
	oneOf := &OneOf{
		Name:    oneOfName,
		Comment: comment,
	}

	parseOneOfFieldFunction := func(statement string, isRepeated bool, comment string) *OneOfField {
		field := &OneOfField{
			IsRepeated: isRepeated,
			Comment:    comment,
		}

		r, _ := regexp.Compile(`(?P<FIELDTYPE>\w+)\s+(?P<FIELDNAME>\w+)\s*=\s*(?P<FIELDTAG>\d+)`)
		matches := r.FindStringSubmatch(statement)
		field.Type, field.Name = matches[1], matches[2]
		field.Tag, _ = strconv.Atoi(matches[3])

		return field
	}

	reader := strings.NewReader(oneBlockBody)
	currentWord, comment := "", ""

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		// Ignore whitespaces and new lines, and comments without a space between the keyword and body (e.g. //MyComment).
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' || currentWord == "//" || currentWord == "/*" {
			if currentWord == "" {
				continue
			} else {
				// Match token type
				switch currentWord {
				case "//":
					{
						if comment == "" {
							comment = readToEndOfLine(reader)
						} else {
							comment += fmt.Sprintf("\n%s", readToEndOfLine(reader))
						}
					}
				case "/*":
					{
						if comment == "" {
							comment = readMultiLineComment(reader)
						} else {
							comment += fmt.Sprintf("\n%s", readMultiLineComment(reader))
						}
					}
				//TODO: is this even allowed in OneOf block?
				case "required":
					{
						statement := extractStatement(reader)
						field := parseOneOfFieldFunction(statement, false, comment)
						oneOf.Fields = append(oneOf.Fields, *field)
						comment = ""
					}
				//TODO: is this even allowed in OneOf block?
				case "optional":
					{
						statement := extractStatement(reader)
						field := parseOneOfFieldFunction(statement, false, comment)
						oneOf.Fields = append(oneOf.Fields, *field)
						comment = ""
					}
				// Repeated field is inherently optional
				case "repeated":
					{
						statement := extractStatement(reader)
						field := parseOneOfFieldFunction(statement, true, comment)
						oneOf.Fields = append(oneOf.Fields, *field)
						comment = ""
					}
				// Optional, required and repeated are omitted
				default:
					{
						statement := fmt.Sprintf("%s %s", currentWord, extractStatement(reader))
						field := parseOneOfFieldFunction(statement, false, comment)
						oneOf.Fields = append(oneOf.Fields, *field)
						comment = ""
					}
				}

				currentWord = ""
				continue
			}
		}

		currentWord += string(r)
	}

	return oneOf
}

func parseExtensions(statement string) *FieldExtensions {
	fieldExtensions := &FieldExtensions{}

	r, _ := regexp.Compile(`(?P<MINTAG>\d+)\s*to\s*(?P<MAXTAG>\d+)`)
	matches := r.FindStringSubmatch(statement)
	fieldExtensions.MinTag, _ = strconv.Atoi(matches[1])
	fieldExtensions.MaxTag, _ = strconv.Atoi(matches[2])

	return fieldExtensions
}
