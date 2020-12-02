// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson7bf29b90DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels(in *jlexer.Lexer, out *ReqExperienceCustomComp) {
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
		case "name_job":
			out.NameJob = string(in.String())
		case "position":
			if in.IsNull() {
				in.Skip()
				out.Position = nil
			} else {
				if out.Position == nil {
					out.Position = new(string)
				}
				*out.Position = string(in.String())
			}
		case "begin":
			out.Begin = string(in.String())
		case "finish":
			if in.IsNull() {
				in.Skip()
				out.Finish = nil
			} else {
				if out.Finish == nil {
					out.Finish = new(string)
				}
				*out.Finish = string(in.String())
			}
		case "duties":
			if in.IsNull() {
				in.Skip()
				out.Duties = nil
			} else {
				if out.Duties == nil {
					out.Duties = new(string)
				}
				*out.Duties = string(in.String())
			}
		case "continue_to_today":
			out.ContinueToToday = bool(in.Bool())
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
func easyjson7bf29b90EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels(out *jwriter.Writer, in ReqExperienceCustomComp) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name_job\":"
		out.RawString(prefix[1:])
		out.String(string(in.NameJob))
	}
	{
		const prefix string = ",\"position\":"
		out.RawString(prefix)
		if in.Position == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Position))
		}
	}
	{
		const prefix string = ",\"begin\":"
		out.RawString(prefix)
		out.String(string(in.Begin))
	}
	{
		const prefix string = ",\"finish\":"
		out.RawString(prefix)
		if in.Finish == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Finish))
		}
	}
	{
		const prefix string = ",\"duties\":"
		out.RawString(prefix)
		if in.Duties == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Duties))
		}
	}
	{
		const prefix string = ",\"continue_to_today\":"
		out.RawString(prefix)
		out.Bool(bool(in.ContinueToToday))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ReqExperienceCustomComp) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7bf29b90EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ReqExperienceCustomComp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7bf29b90EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ReqExperienceCustomComp) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7bf29b90DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ReqExperienceCustomComp) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7bf29b90DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels(l, v)
}
func easyjson7bf29b90DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels1(in *jlexer.Lexer, out *ExperienceOfficialComp) {
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
		case "cand_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CandID).UnmarshalText(data))
			}
		case "resume_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ResumeID).UnmarshalText(data))
			}
		case "company_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CompanyID).UnmarshalText(data))
			}
		case "position":
			if in.IsNull() {
				in.Skip()
				out.Position = nil
			} else {
				if out.Position == nil {
					out.Position = new(string)
				}
				*out.Position = string(in.String())
			}
		case "begin":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Begin).UnmarshalJSON(data))
			}
		case "finish":
			if in.IsNull() {
				in.Skip()
				out.Finish = nil
			} else {
				if out.Finish == nil {
					out.Finish = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.Finish).UnmarshalJSON(data))
				}
			}
		case "duties":
			if in.IsNull() {
				in.Skip()
				out.Duties = nil
			} else {
				if out.Duties == nil {
					out.Duties = new(string)
				}
				*out.Duties = string(in.String())
			}
		case "continue_to_today":
			if in.IsNull() {
				in.Skip()
				out.ContinueToToday = nil
			} else {
				if out.ContinueToToday == nil {
					out.ContinueToToday = new(string)
				}
				*out.ContinueToToday = string(in.String())
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
func easyjson7bf29b90EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels1(out *jwriter.Writer, in ExperienceOfficialComp) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	{
		const prefix string = ",\"cand_id\":"
		out.RawString(prefix)
		out.RawText((in.CandID).MarshalText())
	}
	{
		const prefix string = ",\"resume_id\":"
		out.RawString(prefix)
		out.RawText((in.ResumeID).MarshalText())
	}
	{
		const prefix string = ",\"company_id\":"
		out.RawString(prefix)
		out.RawText((in.CompanyID).MarshalText())
	}
	{
		const prefix string = ",\"position\":"
		out.RawString(prefix)
		if in.Position == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Position))
		}
	}
	{
		const prefix string = ",\"begin\":"
		out.RawString(prefix)
		out.Raw((in.Begin).MarshalJSON())
	}
	{
		const prefix string = ",\"finish\":"
		out.RawString(prefix)
		if in.Finish == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.Finish).MarshalJSON())
		}
	}
	{
		const prefix string = ",\"duties\":"
		out.RawString(prefix)
		if in.Duties == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Duties))
		}
	}
	{
		const prefix string = ",\"continue_to_today\":"
		out.RawString(prefix)
		if in.ContinueToToday == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.ContinueToToday))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExperienceOfficialComp) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7bf29b90EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExperienceOfficialComp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7bf29b90EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExperienceOfficialComp) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7bf29b90DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExperienceOfficialComp) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7bf29b90DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels1(l, v)
}
func easyjson7bf29b90DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels2(in *jlexer.Lexer, out *ExperienceCustomComp) {
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
		case "cand_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CandID).UnmarshalText(data))
			}
		case "resume_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ResumeID).UnmarshalText(data))
			}
		case "name_job":
			out.NameJob = string(in.String())
		case "position":
			if in.IsNull() {
				in.Skip()
				out.Position = nil
			} else {
				if out.Position == nil {
					out.Position = new(string)
				}
				*out.Position = string(in.String())
			}
		case "begin":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Begin).UnmarshalJSON(data))
			}
		case "finish":
			if in.IsNull() {
				in.Skip()
				out.Finish = nil
			} else {
				if out.Finish == nil {
					out.Finish = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.Finish).UnmarshalJSON(data))
				}
			}
		case "duties":
			if in.IsNull() {
				in.Skip()
				out.Duties = nil
			} else {
				if out.Duties == nil {
					out.Duties = new(string)
				}
				*out.Duties = string(in.String())
			}
		case "continue_to_today":
			if in.IsNull() {
				in.Skip()
				out.ContinueToToday = nil
			} else {
				if out.ContinueToToday == nil {
					out.ContinueToToday = new(bool)
				}
				*out.ContinueToToday = bool(in.Bool())
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
func easyjson7bf29b90EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels2(out *jwriter.Writer, in ExperienceCustomComp) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	{
		const prefix string = ",\"cand_id\":"
		out.RawString(prefix)
		out.RawText((in.CandID).MarshalText())
	}
	{
		const prefix string = ",\"resume_id\":"
		out.RawString(prefix)
		out.RawText((in.ResumeID).MarshalText())
	}
	{
		const prefix string = ",\"name_job\":"
		out.RawString(prefix)
		out.String(string(in.NameJob))
	}
	{
		const prefix string = ",\"position\":"
		out.RawString(prefix)
		if in.Position == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Position))
		}
	}
	{
		const prefix string = ",\"begin\":"
		out.RawString(prefix)
		out.Raw((in.Begin).MarshalJSON())
	}
	{
		const prefix string = ",\"finish\":"
		out.RawString(prefix)
		if in.Finish == nil {
			out.RawString("null")
		} else {
			out.Raw((*in.Finish).MarshalJSON())
		}
	}
	{
		const prefix string = ",\"duties\":"
		out.RawString(prefix)
		if in.Duties == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Duties))
		}
	}
	{
		const prefix string = ",\"continue_to_today\":"
		out.RawString(prefix)
		if in.ContinueToToday == nil {
			out.RawString("null")
		} else {
			out.Bool(bool(*in.ContinueToToday))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExperienceCustomComp) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7bf29b90EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExperienceCustomComp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7bf29b90EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExperienceCustomComp) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7bf29b90DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExperienceCustomComp) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7bf29b90DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels2(l, v)
}