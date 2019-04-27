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

func easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(in *jlexer.Lexer, out *SetSessionPayload) {
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
		case "sessionId":
			out.SessionID = string(in.String())
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
func easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(out *jwriter.Writer, in SetSessionPayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"sessionId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.SessionID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SetSessionPayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SetSessionPayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SetSessionPayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SetSessionPayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat(l, v)
}
func easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(in *jlexer.Lexer, out *InitUpdateMessagePayload) {
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
		case "messageId":
			out.Id = uint64(in.Uint64())
		case "authorId":
			if in.IsNull() {
				in.Skip()
				out.Author = nil
			} else {
				if out.Author == nil {
					out.Author = new(uint64)
				}
				*out.Author = uint64(in.Uint64())
			}
		case "data":
			(out.Data).UnmarshalEasyJSON(in)
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
func easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(out *jwriter.Writer, in InitUpdateMessagePayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"messageId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.Id))
	}
	{
		const prefix string = ",\"authorId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Author == nil {
			out.RawString("null")
		} else {
			out.Uint64(uint64(*in.Author))
		}
	}
	{
		const prefix string = ",\"data\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Data).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InitUpdateMessagePayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InitUpdateMessagePayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InitUpdateMessagePayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InitUpdateMessagePayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat1(l, v)
}
func easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat2(in *jlexer.Lexer, out *InitMessagePayload) {
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
		case "sessionId":
			out.SessionID = string(in.String())
		case "author":
			if in.IsNull() {
				in.Skip()
				out.Author = nil
			} else {
				if out.Author == nil {
					out.Author = new(uint64)
				}
				*out.Author = uint64(in.Uint64())
			}
		case "toId":
			if in.IsNull() {
				in.Skip()
				out.Receiver = nil
			} else {
				if out.Receiver == nil {
					out.Receiver = new(uint64)
				}
				*out.Receiver = uint64(in.Uint64())
			}
		case "data":
			(out.Data).UnmarshalEasyJSON(in)
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
func easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat2(out *jwriter.Writer, in InitMessagePayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"sessionId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.SessionID))
	}
	{
		const prefix string = ",\"author\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Author == nil {
			out.RawString("null")
		} else {
			out.Uint64(uint64(*in.Author))
		}
	}
	{
		const prefix string = ",\"toId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Receiver == nil {
			out.RawString("null")
		} else {
			out.Uint64(uint64(*in.Receiver))
		}
	}
	{
		const prefix string = ",\"data\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Data).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InitMessagePayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InitMessagePayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InitMessagePayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InitMessagePayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat2(l, v)
}
func easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat3(in *jlexer.Lexer, out *InitGlobalMessagesPayload) {
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
		case "start":
			out.Start = int(in.Int())
		case "limit":
			out.Limit = int(in.Int())
		case "authorId":
			if in.IsNull() {
				in.Skip()
				out.Author = nil
			} else {
				if out.Author == nil {
					out.Author = new(uint64)
				}
				*out.Author = uint64(in.Uint64())
			}
		case "toId":
			if in.IsNull() {
				in.Skip()
				out.Receiver = nil
			} else {
				if out.Receiver == nil {
					out.Receiver = new(uint64)
				}
				*out.Receiver = uint64(in.Uint64())
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
func easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat3(out *jwriter.Writer, in InitGlobalMessagesPayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"start\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Start))
	}
	{
		const prefix string = ",\"limit\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Limit))
	}
	{
		const prefix string = ",\"authorId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Author == nil {
			out.RawString("null")
		} else {
			out.Uint64(uint64(*in.Author))
		}
	}
	{
		const prefix string = ",\"toId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Receiver == nil {
			out.RawString("null")
		} else {
			out.Uint64(uint64(*in.Receiver))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InitGlobalMessagesPayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InitGlobalMessagesPayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9b8f5552EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InitGlobalMessagesPayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InitGlobalMessagesPayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9b8f5552DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainChat3(l, v)
}
