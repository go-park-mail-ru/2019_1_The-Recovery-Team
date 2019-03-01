// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonE34310f8DecodeApiModels(in *jlexer.Lexer, out *HandlerError) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Description":
			out.Description = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE34310f8EncodeApiModels(out *jwriter.Writer, in HandlerError) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Description))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v HandlerError) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE34310f8EncodeApiModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v HandlerError) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE34310f8EncodeApiModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *HandlerError) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE34310f8DecodeApiModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *HandlerError) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE34310f8DecodeApiModels(l, v)
}
