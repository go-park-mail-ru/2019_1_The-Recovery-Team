// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package chat

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

func easyjsonC83507a6DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(in *jlexer.Lexer, out *ActionRaw) {
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
		case "type":
			out.Type = string(in.String())
		case "payload":
			out.Payload = string(in.String())
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
func easyjsonC83507a6EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(out *jwriter.Writer, in ActionRaw) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	if in.Payload != "" {
		const prefix string = ",\"payload\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Payload))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ActionRaw) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC83507a6EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ActionRaw) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC83507a6EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ActionRaw) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC83507a6DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ActionRaw) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC83507a6DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(l, v)
}
func easyjsonC83507a6DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(in *jlexer.Lexer, out *Action) {
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
		case "type":
			out.Type = string(in.String())
		case "payload":
			if m, ok := out.Payload.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Payload.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Payload = in.Interface()
			}
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
func easyjsonC83507a6EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(out *jwriter.Writer, in Action) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	if in.Payload != nil {
		const prefix string = ",\"payload\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if m, ok := in.Payload.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Payload.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Payload))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Action) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC83507a6EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Action) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC83507a6EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Action) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC83507a6DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Action) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC83507a6DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(l, v)
}
