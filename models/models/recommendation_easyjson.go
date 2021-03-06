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

func easyjson5c54f0e1DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels(in *jlexer.Lexer, out *Recommendation) {
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
		case "ID":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
			}
		case "Sphere0":
			out.Sphere0 = int(in.Int())
		case "Sphere1":
			out.Sphere1 = int(in.Int())
		case "Sphere2":
			out.Sphere2 = int(in.Int())
		case "Sphere3":
			out.Sphere3 = int(in.Int())
		case "Sphere4":
			out.Sphere4 = int(in.Int())
		case "Sphere5":
			out.Sphere5 = int(in.Int())
		case "Sphere6":
			out.Sphere6 = int(in.Int())
		case "Sphere7":
			out.Sphere7 = int(in.Int())
		case "Sphere8":
			out.Sphere8 = int(in.Int())
		case "Sphere9":
			out.Sphere9 = int(in.Int())
		case "Sphere10":
			out.Sphere10 = int(in.Int())
		case "Sphere11":
			out.Sphere11 = int(in.Int())
		case "Sphere12":
			out.Sphere12 = int(in.Int())
		case "Sphere13":
			out.Sphere13 = int(in.Int())
		case "Sphere14":
			out.Sphere14 = int(in.Int())
		case "Sphere15":
			out.Sphere15 = int(in.Int())
		case "Sphere16":
			out.Sphere16 = int(in.Int())
		case "Sphere17":
			out.Sphere17 = int(in.Int())
		case "Sphere18":
			out.Sphere18 = int(in.Int())
		case "Sphere19":
			out.Sphere19 = int(in.Int())
		case "Sphere20":
			out.Sphere20 = int(in.Int())
		case "Sphere21":
			out.Sphere21 = int(in.Int())
		case "Sphere22":
			out.Sphere22 = int(in.Int())
		case "Sphere23":
			out.Sphere23 = int(in.Int())
		case "Sphere24":
			out.Sphere24 = int(in.Int())
		case "Sphere25":
			out.Sphere25 = int(in.Int())
		case "Sphere26":
			out.Sphere26 = int(in.Int())
		case "Sphere27":
			out.Sphere27 = int(in.Int())
		case "Sphere28":
			out.Sphere28 = int(in.Int())
		case "Sphere29":
			out.Sphere29 = int(in.Int())
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
func easyjson5c54f0e1EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels(out *jwriter.Writer, in Recommendation) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ID\":"
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	{
		const prefix string = ",\"Sphere0\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere0))
	}
	{
		const prefix string = ",\"Sphere1\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere1))
	}
	{
		const prefix string = ",\"Sphere2\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere2))
	}
	{
		const prefix string = ",\"Sphere3\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere3))
	}
	{
		const prefix string = ",\"Sphere4\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere4))
	}
	{
		const prefix string = ",\"Sphere5\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere5))
	}
	{
		const prefix string = ",\"Sphere6\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere6))
	}
	{
		const prefix string = ",\"Sphere7\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere7))
	}
	{
		const prefix string = ",\"Sphere8\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere8))
	}
	{
		const prefix string = ",\"Sphere9\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere9))
	}
	{
		const prefix string = ",\"Sphere10\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere10))
	}
	{
		const prefix string = ",\"Sphere11\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere11))
	}
	{
		const prefix string = ",\"Sphere12\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere12))
	}
	{
		const prefix string = ",\"Sphere13\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere13))
	}
	{
		const prefix string = ",\"Sphere14\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere14))
	}
	{
		const prefix string = ",\"Sphere15\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere15))
	}
	{
		const prefix string = ",\"Sphere16\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere16))
	}
	{
		const prefix string = ",\"Sphere17\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere17))
	}
	{
		const prefix string = ",\"Sphere18\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere18))
	}
	{
		const prefix string = ",\"Sphere19\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere19))
	}
	{
		const prefix string = ",\"Sphere20\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere20))
	}
	{
		const prefix string = ",\"Sphere21\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere21))
	}
	{
		const prefix string = ",\"Sphere22\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere22))
	}
	{
		const prefix string = ",\"Sphere23\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere23))
	}
	{
		const prefix string = ",\"Sphere24\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere24))
	}
	{
		const prefix string = ",\"Sphere25\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere25))
	}
	{
		const prefix string = ",\"Sphere26\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere26))
	}
	{
		const prefix string = ",\"Sphere27\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere27))
	}
	{
		const prefix string = ",\"Sphere28\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere28))
	}
	{
		const prefix string = ",\"Sphere29\":"
		out.RawString(prefix)
		out.Int(int(in.Sphere29))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Recommendation) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5c54f0e1EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Recommendation) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5c54f0e1EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Recommendation) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5c54f0e1DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Recommendation) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5c54f0e1DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels(l, v)
}
