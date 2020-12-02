// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package vacancy

import (
	json "encoding/json"
	models "github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
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

func easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy(in *jlexer.Lexer, out *VacRequest) {
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
		case "vac_id":
			out.Id = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "gender":
			out.Gender = string(in.String())
		case "salary_min":
			out.SalaryMin = int(in.Int())
		case "salary_max":
			out.SalaryMax = int(in.Int())
		case "description":
			out.Description = string(in.String())
		case "requirements":
			out.Requirements = string(in.String())
		case "duties":
			out.Duties = string(in.String())
		case "skills":
			out.Skills = string(in.String())
		case "sphere":
			if in.IsNull() {
				in.Skip()
				out.Sphere = nil
			} else {
				if out.Sphere == nil {
					out.Sphere = new(int)
				}
				*out.Sphere = int(in.Int())
			}
		case "employment":
			out.Employment = string(in.String())
		case "experience_month":
			out.ExperienceMonth = int(in.Int())
		case "location":
			out.Location = string(in.String())
		case "area_search":
			out.AreaSearch = string(in.String())
		case "career_level":
			out.CareerLevel = string(in.String())
		case "education_level":
			out.EducationLevel = string(in.String())
		case "email":
			out.EmpEmail = string(in.String())
		case "phone":
			out.EmpPhone = string(in.String())
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
func easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy(out *jwriter.Writer, in VacRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"vac_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"gender\":"
		out.RawString(prefix)
		out.String(string(in.Gender))
	}
	{
		const prefix string = ",\"salary_min\":"
		out.RawString(prefix)
		out.Int(int(in.SalaryMin))
	}
	{
		const prefix string = ",\"salary_max\":"
		out.RawString(prefix)
		out.Int(int(in.SalaryMax))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"requirements\":"
		out.RawString(prefix)
		out.String(string(in.Requirements))
	}
	{
		const prefix string = ",\"duties\":"
		out.RawString(prefix)
		out.String(string(in.Duties))
	}
	{
		const prefix string = ",\"skills\":"
		out.RawString(prefix)
		out.String(string(in.Skills))
	}
	{
		const prefix string = ",\"sphere\":"
		out.RawString(prefix)
		if in.Sphere == nil {
			out.RawString("null")
		} else {
			out.Int(int(*in.Sphere))
		}
	}
	{
		const prefix string = ",\"employment\":"
		out.RawString(prefix)
		out.String(string(in.Employment))
	}
	{
		const prefix string = ",\"experience_month\":"
		out.RawString(prefix)
		out.Int(int(in.ExperienceMonth))
	}
	{
		const prefix string = ",\"location\":"
		out.RawString(prefix)
		out.String(string(in.Location))
	}
	{
		const prefix string = ",\"area_search\":"
		out.RawString(prefix)
		out.String(string(in.AreaSearch))
	}
	{
		const prefix string = ",\"career_level\":"
		out.RawString(prefix)
		out.String(string(in.CareerLevel))
	}
	{
		const prefix string = ",\"education_level\":"
		out.RawString(prefix)
		out.String(string(in.EducationLevel))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.EmpEmail))
	}
	{
		const prefix string = ",\"phone\":"
		out.RawString(prefix)
		out.String(string(in.EmpPhone))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v VacRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v VacRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *VacRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *VacRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy(l, v)
}
func easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy1(in *jlexer.Lexer, out *VacListRequest) {
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
		case "Start":
			out.Start = uint(in.Uint())
		case "Limit":
			out.Limit = uint(in.Uint())
		case "CompId":
			out.CompId = string(in.String())
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
func easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy1(out *jwriter.Writer, in VacListRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Start\":"
		out.RawString(prefix[1:])
		out.Uint(uint(in.Start))
	}
	{
		const prefix string = ",\"Limit\":"
		out.RawString(prefix)
		out.Uint(uint(in.Limit))
	}
	{
		const prefix string = ",\"CompId\":"
		out.RawString(prefix)
		out.String(string(in.CompId))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v VacListRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v VacListRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *VacListRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *VacListRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy1(l, v)
}
func easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy2(in *jlexer.Lexer, out *RespList) {
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
		case "vacancyList":
			if in.IsNull() {
				in.Skip()
				out.Vacancies = nil
			} else {
				in.Delim('[')
				if out.Vacancies == nil {
					if !in.IsDelim(']') {
						out.Vacancies = make([]models.Vacancy, 0, 0)
					} else {
						out.Vacancies = []models.Vacancy{}
					}
				} else {
					out.Vacancies = (out.Vacancies)[:0]
				}
				for !in.IsDelim(']') {
					var v1 models.Vacancy
					(v1).UnmarshalEasyJSON(in)
					out.Vacancies = append(out.Vacancies, v1)
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
func easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy2(out *jwriter.Writer, in RespList) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"vacancyList\":"
		out.RawString(prefix[1:])
		if in.Vacancies == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Vacancies {
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
func (v RespList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RespList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RespList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RespList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy2(l, v)
}
func easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy3(in *jlexer.Lexer, out *Resp) {
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
		case "vacancy":
			if in.IsNull() {
				in.Skip()
				out.Vacancy = nil
			} else {
				if out.Vacancy == nil {
					out.Vacancy = new(models.Vacancy)
				}
				(*out.Vacancy).UnmarshalEasyJSON(in)
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
func easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy3(out *jwriter.Writer, in Resp) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"vacancy\":"
		out.RawString(prefix[1:])
		if in.Vacancy == nil {
			out.RawString("null")
		} else {
			(*in.Vacancy).MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Resp) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Resp) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Resp) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Resp) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy3(l, v)
}
func easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy4(in *jlexer.Lexer, out *Pair) {
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
		case "SphereInd":
			out.SphereInd = int(in.Int())
		case "Score":
			out.Score = int(in.Int())
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
func easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy4(out *jwriter.Writer, in Pair) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"SphereInd\":"
		out.RawString(prefix[1:])
		out.Int(int(in.SphereInd))
	}
	{
		const prefix string = ",\"Score\":"
		out.RawString(prefix)
		out.Int(int(in.Score))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Pair) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Pair) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson56de76c1EncodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Pair) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Pair) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson56de76c1DecodeGithubComGoParkMailRu20202MVVMGitApplicationVacancy4(l, v)
}