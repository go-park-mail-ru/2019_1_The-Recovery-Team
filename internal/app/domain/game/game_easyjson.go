// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package game

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

func easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame(in *jlexer.Lexer, out *SetItemPayload) {
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
		case "playerId":
			out.PlayerId = uint64(in.Uint64())
		case "itemType":
			out.ItemType = string(in.String())
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
func easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame(out *jwriter.Writer, in SetItemPayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"playerId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.PlayerId))
	}
	{
		const prefix string = ",\"itemType\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ItemType))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SetItemPayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SetItemPayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SetItemPayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SetItemPayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame(l, v)
}
func easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame1(in *jlexer.Lexer, out *SetGameStartPayload) {
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
		case "Field":
			if in.IsNull() {
				in.Skip()
				out.Field = nil
			} else {
				if out.Field == nil {
					out.Field = new(Field)
				}
				easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame2(in, out.Field)
			}
		case "Players":
			if in.IsNull() {
				in.Skip()
				out.Players = nil
			} else {
				in.Delim('[')
				if out.Players == nil {
					if !in.IsDelim(']') {
						out.Players = make([]Player, 0, 1)
					} else {
						out.Players = []Player{}
					}
				} else {
					out.Players = (out.Players)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Player
					(v1).UnmarshalEasyJSON(in)
					out.Players = append(out.Players, v1)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame1(out *jwriter.Writer, in SetGameStartPayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Field\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Field == nil {
			out.RawString("null")
		} else {
			easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame2(out, *in.Field)
		}
	}
	{
		const prefix string = ",\"Players\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Players == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Players {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SetGameStartPayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SetGameStartPayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SetGameStartPayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SetGameStartPayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame1(l, v)
}
func easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame2(in *jlexer.Lexer, out *Field) {
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
		case "cells":
			if in.IsNull() {
				in.Skip()
				out.Cells = nil
			} else {
				in.Delim('[')
				if out.Cells == nil {
					if !in.IsDelim(']') {
						out.Cells = make([]Cell, 0, 1)
					} else {
						out.Cells = []Cell{}
					}
				} else {
					out.Cells = (out.Cells)[:0]
				}
				for !in.IsDelim(']') {
					var v4 Cell
					easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame3(in, &v4)
					out.Cells = append(out.Cells, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "width":
			out.Width = int(in.Int())
		case "height":
			out.Height = int(in.Int())
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
func easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame2(out *jwriter.Writer, in Field) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.Cells) != 0 {
		const prefix string = ",\"cells\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v5, v6 := range in.Cells {
				if v5 > 0 {
					out.RawByte(',')
				}
				easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame3(out, v6)
			}
			out.RawByte(']')
		}
	}
	if in.Width != 0 {
		const prefix string = ",\"width\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Width))
	}
	if in.Height != 0 {
		const prefix string = ",\"height\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Height))
	}
	out.RawByte('}')
}
func easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame3(in *jlexer.Lexer, out *Cell) {
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
		case "row":
			out.Row = int(in.Int())
		case "col":
			out.Col = int(in.Int())
		case "type":
			out.Type = string(in.String())
		case "hasBox":
			out.HasBox = bool(in.Bool())
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
func easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame3(out *jwriter.Writer, in Cell) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"row\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Row))
	}
	{
		const prefix string = ",\"col\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Col))
	}
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
	{
		const prefix string = ",\"hasBox\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.HasBox))
	}
	out.RawByte('}')
}
func easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame4(in *jlexer.Lexer, out *Players) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Players, 0, 1)
			} else {
				*out = Players{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v7 Player
			(v7).UnmarshalEasyJSON(in)
			*out = append(*out, v7)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame4(out *jwriter.Writer, in Players) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in {
			if v8 > 0 {
				out.RawByte(',')
			}
			(v9).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Players) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Players) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Players) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Players) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame4(l, v)
}
func easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame5(in *jlexer.Lexer, out *Player) {
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
			out.Id = uint64(in.Uint64())
		case "x":
			out.X = int(in.Int())
		case "y":
			out.Y = int(in.Int())
		case "items":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Items = make(map[string]uint64)
				} else {
					out.Items = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v10 uint64
					v10 = uint64(in.Uint64())
					(out.Items)[key] = v10
					in.WantComma()
				}
				in.Delim('}')
			}
		case "loseRound":
			if in.IsNull() {
				in.Skip()
				out.LoseRound = nil
			} else {
				if out.LoseRound == nil {
					out.LoseRound = new(int)
				}
				*out.LoseRound = int(in.Int())
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
func easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame5(out *jwriter.Writer, in Player) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.Id))
	}
	{
		const prefix string = ",\"x\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.X))
	}
	{
		const prefix string = ",\"y\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Y))
	}
	{
		const prefix string = ",\"items\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Items == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v11First := true
			for v11Name, v11Value := range in.Items {
				if v11First {
					v11First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v11Name))
				out.RawByte(':')
				out.Uint64(uint64(v11Value))
			}
			out.RawByte('}')
		}
	}
	if in.LoseRound != nil {
		const prefix string = ",\"loseRound\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(*in.LoseRound))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Player) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Player) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Player) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Player) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame5(l, v)
}
func easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame6(in *jlexer.Lexer, out *InitPlayersPayload) {
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
		case "playerIds":
			if in.IsNull() {
				in.Skip()
				out.PlayersId = nil
			} else {
				in.Delim('[')
				if out.PlayersId == nil {
					if !in.IsDelim(']') {
						out.PlayersId = make([]uint64, 0, 8)
					} else {
						out.PlayersId = []uint64{}
					}
				} else {
					out.PlayersId = (out.PlayersId)[:0]
				}
				for !in.IsDelim(']') {
					var v12 uint64
					v12 = uint64(in.Uint64())
					out.PlayersId = append(out.PlayersId, v12)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame6(out *jwriter.Writer, in InitPlayersPayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"playerIds\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.PlayersId == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v13, v14 := range in.PlayersId {
				if v13 > 0 {
					out.RawByte(',')
				}
				out.Uint64(uint64(v14))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InitPlayersPayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InitPlayersPayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InitPlayersPayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InitPlayersPayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame6(l, v)
}
func easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame7(in *jlexer.Lexer, out *InitPlayerReadyPayload) {
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
		case "playerId":
			out.PlayerId = uint64(in.Uint64())
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
func easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame7(out *jwriter.Writer, in InitPlayerReadyPayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"playerId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.PlayerId))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InitPlayerReadyPayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InitPlayerReadyPayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InitPlayerReadyPayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InitPlayerReadyPayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame7(l, v)
}
func easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame8(in *jlexer.Lexer, out *InitPlayerMovePayload) {
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
		case "playerId":
			out.PlayerId = uint64(in.Uint64())
		case "move":
			out.Move = string(in.String())
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
func easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame8(out *jwriter.Writer, in InitPlayerMovePayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"playerId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.PlayerId))
	}
	{
		const prefix string = ",\"move\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Move))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InitPlayerMovePayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InitPlayerMovePayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InitPlayerMovePayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InitPlayerMovePayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame8(l, v)
}
func easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame9(in *jlexer.Lexer, out *InitItemUsePayload) {
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
		case "playerId":
			out.PlayerId = uint64(in.Uint64())
		case "itemType":
			out.ItemType = string(in.String())
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
func easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame9(out *jwriter.Writer, in InitItemUsePayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"playerId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.PlayerId))
	}
	{
		const prefix string = ",\"itemType\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ItemType))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InitItemUsePayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InitItemUsePayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656EncodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InitItemUsePayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InitItemUsePayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656DecodeGithubComGoParkMailRu20191TheRecoveryTeamInternalAppDomainGame9(l, v)
}
