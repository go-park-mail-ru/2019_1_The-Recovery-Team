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

func easyjson521a5691DecodeApiModels(in *jlexer.Lexer, out *Score) {
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
		case "record":
			out.Record = int(in.Int())
		case "win":
			out.Win = int(in.Int())
		case "loss":
			out.Loss = int(in.Int())
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
func easyjson521a5691EncodeApiModels(out *jwriter.Writer, in Score) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Record != 0 {
		const prefix string = ",\"record\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Record))
	}
	if in.Win != 0 {
		const prefix string = ",\"win\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Win))
	}
	if in.Loss != 0 {
		const prefix string = ",\"loss\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Loss))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Score) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeApiModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Score) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeApiModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Score) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeApiModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Score) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeApiModels(l, v)
}
func easyjson521a5691DecodeApiModels1(in *jlexer.Lexer, out *Profiles) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Profiles, 0, 1)
			} else {
				*out = Profiles{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Profile
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson521a5691EncodeApiModels1(out *jwriter.Writer, in Profiles) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Profiles) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeApiModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Profiles) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeApiModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Profiles) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeApiModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Profiles) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeApiModels1(l, v)
}
func easyjson521a5691DecodeApiModels2(in *jlexer.Lexer, out *ProfileRegistration) {
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
		case "nickname":
			out.Nickname = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjson521a5691EncodeApiModels2(out *jwriter.Writer, in ProfileRegistration) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"nickname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Nickname))
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProfileRegistration) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeApiModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProfileRegistration) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeApiModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProfileRegistration) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeApiModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProfileRegistration) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeApiModels2(l, v)
}
func easyjson521a5691DecodeApiModels3(in *jlexer.Lexer, out *ProfileLogin) {
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
		case "email":
			out.Email = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjson521a5691EncodeApiModels3(out *jwriter.Writer, in ProfileLogin) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProfileLogin) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeApiModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProfileLogin) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeApiModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProfileLogin) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeApiModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProfileLogin) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeApiModels3(l, v)
}
func easyjson521a5691DecodeApiModels4(in *jlexer.Lexer, out *ProfileInfo) {
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
		case "email":
			out.Email = string(in.String())
		case "nickname":
			out.Nickname = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjson521a5691EncodeApiModels4(out *jwriter.Writer, in ProfileInfo) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Email != "" {
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	if in.Nickname != "" {
		const prefix string = ",\"nickname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Nickname))
	}
	if in.Password != "" {
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProfileInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeApiModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProfileInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeApiModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProfileInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeApiModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProfileInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeApiModels4(l, v)
}
func easyjson521a5691DecodeApiModels5(in *jlexer.Lexer, out *Profile) {
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
		case "id":
			out.ID = uint64(in.Uint64())
		case "avatar":
			out.Avatar = string(in.String())
		case "record":
			out.Record = int(in.Int())
		case "win":
			out.Win = int(in.Int())
		case "loss":
			out.Loss = int(in.Int())
		case "email":
			out.Email = string(in.String())
		case "nickname":
			out.Nickname = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjson521a5691EncodeApiModels5(out *jwriter.Writer, in Profile) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.ID))
	}
	if in.Avatar != "" {
		const prefix string = ",\"avatar\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Avatar))
	}
	if in.Record != 0 {
		const prefix string = ",\"record\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Record))
	}
	if in.Win != 0 {
		const prefix string = ",\"win\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Win))
	}
	if in.Loss != 0 {
		const prefix string = ",\"loss\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Loss))
	}
	if in.Email != "" {
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	if in.Nickname != "" {
		const prefix string = ",\"nickname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Nickname))
	}
	if in.Password != "" {
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Profile) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeApiModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Profile) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeApiModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Profile) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeApiModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Profile) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeApiModels5(l, v)
}
