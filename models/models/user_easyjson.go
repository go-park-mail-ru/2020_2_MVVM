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

func easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels(in *jlexer.Lexer, out *UserLogin) {
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
		key := in.UnsafeFieldName(false)
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
func easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels(out *jwriter.Writer, in UserLogin) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix[1:])
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserLogin) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserLogin) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserLogin) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserLogin) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels(l, v)
}
func easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels1(in *jlexer.Lexer, out *User) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
			}
		case "user_type":
			out.UserType = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "surname":
			out.Surname = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "phone":
			if in.IsNull() {
				in.Skip()
				out.Phone = nil
			} else {
				if out.Phone == nil {
					out.Phone = new(string)
				}
				*out.Phone = string(in.String())
			}
		case "social_network":
			if in.IsNull() {
				in.Skip()
				out.SocialNetwork = nil
			} else {
				if out.SocialNetwork == nil {
					out.SocialNetwork = new(string)
				}
				*out.SocialNetwork = string(in.String())
			}
		case "avatar":
			out.AvatarPath = string(in.String())
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
func easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels1(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	{
		const prefix string = ",\"user_type\":"
		out.RawString(prefix)
		out.String(string(in.UserType))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"surname\":"
		out.RawString(prefix)
		out.String(string(in.Surname))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"phone\":"
		out.RawString(prefix)
		if in.Phone == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Phone))
		}
	}
	{
		const prefix string = ",\"social_network\":"
		out.RawString(prefix)
		if in.SocialNetwork == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.SocialNetwork))
		}
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.AvatarPath))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels1(l, v)
}
func easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels2(in *jlexer.Lexer, out *RespUser) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user":
			if in.IsNull() {
				in.Skip()
				out.User = nil
			} else {
				if out.User == nil {
					out.User = new(User)
				}
				(*out.User).UnmarshalEasyJSON(in)
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
func easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels2(out *jwriter.Writer, in RespUser) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix[1:])
		if in.User == nil {
			out.RawString("null")
		} else {
			(*in.User).MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RespUser) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RespUser) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RespUser) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RespUser) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels2(l, v)
}
func easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels3(in *jlexer.Lexer, out *Employer) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "empl_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
			}
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "comp_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CompanyID).UnmarshalText(data))
			}
		case "Favorites":
			if in.IsNull() {
				in.Skip()
				out.Favorites = nil
			} else {
				in.Delim('[')
				if out.Favorites == nil {
					if !in.IsDelim(']') {
						out.Favorites = make([]FavoritesForEmpl, 0, 0)
					} else {
						out.Favorites = []FavoritesForEmpl{}
					}
				} else {
					out.Favorites = (out.Favorites)[:0]
				}
				for !in.IsDelim(']') {
					var v1 FavoritesForEmpl
					(v1).UnmarshalEasyJSON(in)
					out.Favorites = append(out.Favorites, v1)
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
func easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels3(out *jwriter.Writer, in Employer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"empl_id\":"
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.RawText((in.UserID).MarshalText())
	}
	{
		const prefix string = ",\"comp_id\":"
		out.RawString(prefix)
		out.RawText((in.CompanyID).MarshalText())
	}
	{
		const prefix string = ",\"Favorites\":"
		out.RawString(prefix)
		if in.Favorites == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Favorites {
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
func (v Employer) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Employer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Employer) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Employer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels3(l, v)
}
func easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels4(in *jlexer.Lexer, out *Candidate) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "cand_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
			}
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "User":
			(out.User).UnmarshalEasyJSON(in)
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
func easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels4(out *jwriter.Writer, in Candidate) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"cand_id\":"
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.RawText((in.UserID).MarshalText())
	}
	{
		const prefix string = ",\"User\":"
		out.RawString(prefix)
		(in.User).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Candidate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Candidate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9e1087fdEncodeGithubComGoParkMailRu20202MVVMGitModelsModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Candidate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Candidate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9e1087fdDecodeGithubComGoParkMailRu20202MVVMGitModelsModels4(l, v)
}
