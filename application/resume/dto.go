package resume

type SearchParams struct {
	KeyWords        *string  `json:"keywords"`
	SalaryMin       *int     `json:"salary_min"`
	SalaryMax       *int     `json:"salary_max"`
	Gender          []string `json:"gender"`
	EducationLevel  []string `json:"education_level"`
	CareerLevel     []string `json:"career_level"`
	ExperienceMonth []int    `json:"experience_month"`
	AreaSearch      []string `json:"area_search"`
}
