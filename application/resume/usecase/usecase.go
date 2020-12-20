package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/makePdf/pdfGenerator"
	resume2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/resume"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"strings"
	"time"
)

type ResumeUseCase struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
	userUseCase user.UseCase
	//educationUseCase education.UseCase
	customExpUseCase custom_experience.UseCase
	strg             resume.Repository
}

func NewUseCase(infoLogger *logger.Logger,
	errorLogger *logger.Logger,
	userUseCase user.UseCase,
//educationUseCase education.UseCase,
	customExpUseCase custom_experience.UseCase,
	strg resume.Repository) resume.UseCase {
	usecase := ResumeUseCase{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		userUseCase: userUseCase,
		//educationUseCase: educationUseCase,
		customExpUseCase: customExpUseCase,
		strg:             strg,
	}
	return &usecase
}

func (u *ResumeUseCase) Create(template models.Resume) (*models.Resume, error) {
	// create resume
	template.DateCreate = time.Now().Format(time.RFC3339)
	if template.Sphere == nil {
		template.Sphere = &common.DefaultSphere
	}
	for i := range template.ExperienceCustomComp {
		if *(template.ExperienceCustomComp[i].ContinueToToday) {
			dateFinish := time.Now()
			template.ExperienceCustomComp[i].Finish = &dateFinish
		}
		template.ExperienceCustomComp[i].ResumeID = template.ResumeID
		template.ExperienceCustomComp[i].CandID = template.CandID

	}
	result, err := u.strg.Create(template)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *ResumeUseCase) Update(resume models.Resume) (*models.Resume, error) {
	oldResume, err := u.strg.GetById(resume.ResumeID)
	if err != nil {
		err = fmt.Errorf("error in get resume by id: %w", err)
		return nil, err
	}
	if resume.CandID != oldResume.CandID {
		err = fmt.Errorf("this user cannot update this resume")
		return nil, err
	}
	resume.DateCreate = oldResume.DateCreate
	if resume.Sphere == nil {
		resume.Sphere = &common.DefaultSphere
	}
	//err = u.customExpUseCase.DropAllFromResume(resume.ResumeID)
	//if err != nil {
	//	return nil, err
	//}
	return u.strg.Update(resume)
}

func (u *ResumeUseCase) GetAllUserResume(userid uuid.UUID) ([]models.BriefResumeInfo, error) {
	r, err := u.strg.GetAllUserResume(userid)
	if err != nil {
		err = fmt.Errorf("error in get my resume: %w", err)
		return nil, err
	}

	return DoBriefRespResume(r)
}

func (u *ResumeUseCase) Search(searchParams resume2.SearchParams) ([]models.BriefResumeInfo, error) {
	if searchParams.KeyWords != nil {
		*searchParams.KeyWords = strings.ToLower(*searchParams.KeyWords)
		searchParams.KeywordsGeo = strings.ToLower(searchParams.KeywordsGeo)
	}
	r, err := u.strg.Search(&searchParams)
	if err != nil {
		err = fmt.Errorf("error in resume search: %w", err)
		return nil, err
	}

	return DoBriefRespResume(r)
}

func (u *ResumeUseCase) GetById(id uuid.UUID) (*models.Resume, error) {
	r, err := u.strg.GetById(id)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return nil, err
	}
	return r, nil
}

func (u *ResumeUseCase) List(start, limit uint) ([]models.BriefResumeInfo, error) {
	if limit >= 200 {
		return nil, fmt.Errorf("Limit is too high. ")
	}
	r, err := u.strg.List(start, limit)
	if err != nil {
		err = fmt.Errorf("error in resume get list from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	return DoBriefRespResume(r)
}

func (u *ResumeUseCase) AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoriteID, error) {
	return u.strg.AddFavorite(favoriteForEmpl)
}

func (u *ResumeUseCase) RemoveFavorite(favoriteForEmpl models.FavoritesForEmpl) error {
	oldFavorite, err := u.strg.GetFavoriteByID(favoriteForEmpl.FavoriteID)
	if err != nil {
		return err
	}
	if oldFavorite.EmplID != favoriteForEmpl.EmplID {
		err = fmt.Errorf("error in remove favorite: method not allowed")
		return err
	}
	return u.strg.RemoveFavorite(favoriteForEmpl.FavoriteID)
}

func (u *ResumeUseCase) GetAllEmplFavoriteResume(userID uuid.UUID) ([]models.BriefResumeInfo, error) {
	r, err := u.strg.GetAllEmplFavoriteResume(userID)
	if err != nil {
		err = fmt.Errorf("error in get list favorite resume: %w", err)
		return nil, err
	}
	return DoBriefRespResume(r)
}

func (u *ResumeUseCase) GetFavoriteByResume(userID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error) {
	return u.strg.GetFavoriteForResume(userID, resumeID)
}

func (u *ResumeUseCase) GetFavoriteByID(favoriteID uuid.UUID) (*models.FavoritesForEmpl, error) {
	return u.strg.GetFavoriteByID(favoriteID)
}

func DoBriefRespResume(resumes []models.Resume) ([]models.BriefResumeInfo, error) {
	var briefRespResumes []models.BriefResumeInfo
	for i := range resumes {
		brief, err := resumes[i].Brief()
		if err != nil {
			return nil, err
		}
		briefRespResumes = append(briefRespResumes, *brief)
	}
	return briefRespResumes, nil
}

func (u *ResumeUseCase) DeleteResume(resId uuid.UUID, candId uuid.UUID) error {
	return u.strg.Delete(resId, candId)
}

func (u *ResumeUseCase) MakePdf(id uuid.UUID) error {
	r, err := u.strg.GetById(id)
	if err != nil {
		err = fmt.Errorf("error in resume get by id func : %w", err)
		return err
	}

	pgf := pdfGenerator.NewRequestPdf("")

	//html template path
	templatePath := "makePdf/templates/resume.html"

	//path for download pdf
	outputPath := fmt.Sprintf(`static/pdf/%s.pdf`, id.String())

	//html template data
	template := convertToTemplate(*r)

	if err := pgf.ParseTemplate(templatePath, template); err == nil {
		_, err := pgf.GeneratePDF(outputPath)
		if err != nil {
			err = fmt.Errorf("error in resume pdf generated: %w", err)
			return err
		}
	} else {
		err = fmt.Errorf("error in resume pdf generated: %w", err)
		return err
	}

	return nil
}

type ExperienceTemplate struct {
	NameJob  string
	Position string
	Begin    string
	Finish   string
	Duties   string
}

type resumeTemplate struct {
	Avatar          string
	Name            string
	Surname         string
	AreaSearch      string
	Phone           string
	Email           string
	Sphere          string
	Experience      string
	Education       string
	Salary          string
	Description     string
	ExperienceItems []ExperienceTemplate
	Skills          string
}

func convertToTemplate(r models.Resume) resumeTemplate {
	t := resumeTemplate{
		Avatar:      r.Avatar,
		Name:        r.CandName,
		Surname:     r.CandSurname,
		Email:       r.Candidate.User.Email,
		Description: r.Description,
		Skills:      r.Skills,
	}
	if r.AreaSearch != nil {
		t.AreaSearch = *(r.AreaSearch)
	} else {
		t.AreaSearch = ""
	}
	if r.Candidate.User.Phone != nil {
		t.Phone = *(r.Candidate.User.Phone)
	} else {
		t.Phone = ""
	}
	if r.Sphere != nil {
		t.Sphere = strconv.Itoa(*r.Sphere)
	} else {
		t.Sphere = ""
	}
	if r.ExperienceMonth != nil {
		t.Experience = strconv.Itoa(*r.ExperienceMonth)
	} else {
		t.Experience = ""
	}
	if r.EducationLevel != nil {
		t.Education = *(r.EducationLevel)
	} else {
		t.Education = ""
	}
	if r.SalaryMin != nil {
		if r.SalaryMax != nil {
			t.Salary = fmt.Sprintf(`%v - %v`, *r.SalaryMin, *r.SalaryMax)
		}
		t.Salary = strconv.Itoa(*r.SalaryMin)
	} else if r.SalaryMax != nil {
		t.Salary = strconv.Itoa(*r.SalaryMax)
	} else {
		t.Salary = ""
	}

	var exp []ExperienceTemplate
	for _, val := range r.ExperienceCustomComp {
		expItem := ExperienceTemplate{NameJob: val.NameJob,
			Begin: fmt.Sprintf(val.Begin.Format("2006-01-02"))}
		if val.Position != nil {
			expItem.Position = *(val.Position)
		} else {
			expItem.Position = ""
		}
		if val.Duties != nil {
			expItem.Duties = *(val.Duties)
		} else {
			expItem.Duties = ""
		}
		if val.Finish != nil {
			expItem.Finish = fmt.Sprintf(val.Finish.Format("2006-01-02"))
		} else {
			expItem.Finish = "по настоящее время"
		}
		exp = append(exp, expItem)
	}
	t.ExperienceItems = exp

	return t
}
