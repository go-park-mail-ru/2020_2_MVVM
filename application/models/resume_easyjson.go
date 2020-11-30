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

func easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels(in *jlexer.Lexer, out *Resume) {
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
				in.AddError((out.ResumeID).UnmarshalText(data))
			}
		case "cand_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CandID).UnmarshalText(data))
			}
		case "Candidate":
			easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels1(in, &out.Candidate)
		case "title":
			out.Title = string(in.String())
		case "salary_min":
			if in.IsNull() {
				in.Skip()
				out.SalaryMin = nil
			} else {
				if out.SalaryMin == nil {
					out.SalaryMin = new(int)
				}
				*out.SalaryMin = int(in.Int())
			}
		case "salary_max":
			if in.IsNull() {
				in.Skip()
				out.SalaryMax = nil
			} else {
				if out.SalaryMax == nil {
					out.SalaryMax = new(int)
				}
				*out.SalaryMax = int(in.Int())
			}
		case "description":
			out.Description = string(in.String())
		case "skills":
			out.Skills = string(in.String())
		case "gender":
			out.Gender = string(in.String())
		case "education_level":
			if in.IsNull() {
				in.Skip()
				out.EducationLevel = nil
			} else {
				if out.EducationLevel == nil {
					out.EducationLevel = new(string)
				}
				*out.EducationLevel = string(in.String())
			}
		case "career_level":
			if in.IsNull() {
				in.Skip()
				out.CareerLevel = nil
			} else {
				if out.CareerLevel == nil {
					out.CareerLevel = new(string)
				}
				*out.CareerLevel = string(in.String())
			}
		case "place":
			if in.IsNull() {
				in.Skip()
				out.Place = nil
			} else {
				if out.Place == nil {
					out.Place = new(string)
				}
				*out.Place = string(in.String())
			}
		case "experience_month":
			if in.IsNull() {
				in.Skip()
				out.ExperienceMonth = nil
			} else {
				if out.ExperienceMonth == nil {
					out.ExperienceMonth = new(int)
				}
				*out.ExperienceMonth = int(in.Int())
			}
		case "area_search":
			if in.IsNull() {
				in.Skip()
				out.AreaSearch = nil
			} else {
				if out.AreaSearch == nil {
					out.AreaSearch = new(string)
				}
				*out.AreaSearch = string(in.String())
			}
		case "date_create":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.DateCreate).UnmarshalJSON(data))
			}
		case "education":
			if in.IsNull() {
				in.Skip()
				out.Education = nil
			} else {
				in.Delim('[')
				if out.Education == nil {
					if !in.IsDelim(']') {
						out.Education = make([]Education, 0, 0)
					} else {
						out.Education = []Education{}
					}
				} else {
					out.Education = (out.Education)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Education
					(v1).UnmarshalEasyJSON(in)
					out.Education = append(out.Education, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "custom_experience":
			if in.IsNull() {
				in.Skip()
				out.ExperienceCustomComp = nil
			} else {
				in.Delim('[')
				if out.ExperienceCustomComp == nil {
					if !in.IsDelim(']') {
						out.ExperienceCustomComp = make([]ExperienceCustomComp, 0, 0)
					} else {
						out.ExperienceCustomComp = []ExperienceCustomComp{}
					}
				} else {
					out.ExperienceCustomComp = (out.ExperienceCustomComp)[:0]
				}
				for !in.IsDelim(']') {
					var v2 ExperienceCustomComp
					(v2).UnmarshalEasyJSON(in)
					out.ExperienceCustomComp = append(out.ExperienceCustomComp, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "avatar":
			out.Avatar = string(in.String())
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
func easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels(out *jwriter.Writer, in Resume) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.ResumeID).MarshalText())
	}
	{
		const prefix string = ",\"cand_id\":"
		out.RawString(prefix)
		out.RawText((in.CandID).MarshalText())
	}
	{
		const prefix string = ",\"Candidate\":"
		out.RawString(prefix)
		easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels1(out, in.Candidate)
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"salary_min\":"
		out.RawString(prefix)
		if in.SalaryMin == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.SalaryMin))
		}
	}
	{
		const prefix string = ",\"salary_max\":"
		out.RawString(prefix)
		if in.SalaryMax == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.SalaryMax))
		}
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"skills\":"
		out.RawString(prefix)
		out.String(string(in.Skills))
	}
	{
		const prefix string = ",\"gender\":"
		out.RawString(prefix)
		out.String(string(in.Gender))
	}
	{
		const prefix string = ",\"education_level\":"
		out.RawString(prefix)
		if in.EducationLevel == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.EducationLevel))
		}
	}
	{
		const prefix string = ",\"career_level\":"
		out.RawString(prefix)
		if in.CareerLevel == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.CareerLevel))
		}
	}
	{
		const prefix string = ",\"place\":"
		out.RawString(prefix)
		if in.Place == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Place))
		}
	}
	{
		const prefix string = ",\"experience_month\":"
		out.RawString(prefix)
		if in.ExperienceMonth == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.ExperienceMonth))
		}
	}
	{
		const prefix string = ",\"area_search\":"
		out.RawString(prefix)
		if in.AreaSearch == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.AreaSearch))
		}
	}
	{
		const prefix string = ",\"date_create\":"
		out.RawString(prefix)
		out.Raw((in.DateCreate).MarshalJSON())
	}
	{
		const prefix string = ",\"education\":"
		out.RawString(prefix)
		if in.Education == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Education {
				if v3 > 0 {
					out.RawByte(',')
				}
				(v4).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"custom_experience\":"
		out.RawString(prefix)
		if in.ExperienceCustomComp == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.ExperienceCustomComp {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Resume) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Resume) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Resume) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Resume) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels(l, v)
}
func easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels1(in *jlexer.Lexer, out *Candidate) {
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
			easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels2(in, &out.User)
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
func easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels1(out *jwriter.Writer, in Candidate) {
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
		easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels2(out, in.User)
	}
	out.RawByte('}')
}
func easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels2(in *jlexer.Lexer, out *User) {
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
func easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels2(out *jwriter.Writer, in User) {
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
	out.RawByte('}')
}
func easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels3(in *jlexer.Lexer, out *ListBriefResumeInfo) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(ListBriefResumeInfo, 0, 0)
			} else {
				*out = ListBriefResumeInfo{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v7 BriefResumeInfo
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
func easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels3(out *jwriter.Writer, in ListBriefResumeInfo) {
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
func (v ListBriefResumeInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ListBriefResumeInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ListBriefResumeInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ListBriefResumeInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels3(l, v)
}
func easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels4(in *jlexer.Lexer, out *BriefResumeInfo) {
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
		case "avatar":
			out.Avatar = string(in.String())
		case "resume_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ResumeID).UnmarshalText(data))
			}
		case "cand_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CandID).UnmarshalText(data))
			}
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "place":
			if in.IsNull() {
				in.Skip()
				out.Place = nil
			} else {
				if out.Place == nil {
					out.Place = new(string)
				}
				*out.Place = string(in.String())
			}
		case "location":
			if in.IsNull() {
				in.Skip()
				out.AreaSearch = nil
			} else {
				if out.AreaSearch == nil {
					out.AreaSearch = new(string)
				}
				*out.AreaSearch = string(in.String())
			}
		case "name":
			out.Name = string(in.String())
		case "surname":
			out.Surname = string(in.String())
		case "email":
			out.Email = string(in.String())
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
func easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels4(out *jwriter.Writer, in BriefResumeInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix[1:])
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"resume_id\":"
		out.RawString(prefix)
		out.RawText((in.ResumeID).MarshalText())
	}
	{
		const prefix string = ",\"cand_id\":"
		out.RawString(prefix)
		out.RawText((in.CandID).MarshalText())
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.RawText((in.UserID).MarshalText())
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"place\":"
		out.RawString(prefix)
		if in.Place == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Place))
		}
	}
	{
		const prefix string = ",\"location\":"
		out.RawString(prefix)
		if in.AreaSearch == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.AreaSearch))
		}
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
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BriefResumeInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BriefResumeInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson39b3a2f5EncodeGithubComGoParkMailRu20202MVVMGitApplicationModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BriefResumeInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BriefResumeInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson39b3a2f5DecodeGithubComGoParkMailRu20202MVVMGitApplicationModels4(l, v)
}
