package models

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	JobTitle            string          `json:"job_title" validate:"required"`
	Salary              string          `json:"sal" validate:"required"`
	CompanyId           uint64          `json:"cid" validate:"required"`
	Comp                Company         `gorm:"ForeignKey:CompanyId"`
	MinimumNoticePeriod int             `json:"min_np" validate:"required"`
	MaximumNoticePeriod uint64          `json:"max_np" validate:"required"`
	Budget              float64         `json:"budget" validate:"required"`
	JobDescription      string          `json:"job_desc" validate:"required"`
	MinExperience       float64         `json:"min_exp" validate:"required"`
	MaxExperience       float64         `json:"max_exp" validate:"required"`
	Locations           []Location      `gorm:"many2many:job_locations;"`
	Skills              []Skill         `gorm:"many2many:job_skills;"`
	WorkModes           []WorkMode      `gorm:"many2many:job_work_modes;"`
	Qualifications      []Qualification `gorm:"many2many:job_qualifications;"`
	Shifts              []Shift         `gorm:"many2many:job_shifts;"`
	JobTypes            []JobType       `gorm:"many2many:job_jobtypes;"`
}
type Location struct {
	gorm.Model
	State string `json:"state" validate:"required"`
}
type Skill struct {
	gorm.Model
	Skillsets string `json:"skillsets" validate:"required"`
}
type WorkMode struct {
	gorm.Model
	Mode string `json:"work_mode" validate:"required"`
}
type Qualification struct {
	gorm.Model
	Degree string `json:"degree" validate:"required"`
}
type Shift struct {
	gorm.Model
	ShiftType string `json:"shift_type" validate:"required"`
}
type JobType struct {
	gorm.Model
	Typeofjob string `json:"job_type" validate:"required"`
}

type NewJobRequest struct {
	JobTitle            string  `json:"jobTitle" validate:"required"`
	Salary              string  `json:"sal" validate:"required"`
	MinimumNoticePeriod int     `json:"minNp" validate:"required"`
	MaximumNoticePeriod uint64  `json:"maxNp" validate:"required"`
	Budget              float64 `json:"budget" validate:"required"`
	JobDescription      string  `json:"jobDesc" validate:"required"`
	MinExperience       float64 `json:"minExp" validate:"required"`
	MaxExperience       float64 `json:"maxExp" validate:"required"`
	LocationIDs         []uint
	SkillIDs            []uint
	WorkModeIDs         []uint
	QualificationIDs    []uint
	ShiftIDs            []uint
	JobTypeIDs          []uint
}
type Response struct {
	ID uint64
}

type RequestFromUser struct {
	NoticePeriod   int     `json:"noticePeriod" validate:"required"`
	Location       []uint  `json:"location" `
	Skills         []uint  `json:"technologyStack" `
	Experience     float64 `json:"experience" validate:"required"`
	Qualifications []uint  `json:"qualifications"`
	Shift          []uint  `json:"shifts"`
	WorkModeIDs    []uint  `json:"work_modes"`
	JobType        []uint  `json:"job_type"`
}
type NewUserApplication struct {
	Name string          `json:"name"`
	Age  string          `json:"age"`
	ID   uint64          `json:"jid"`
	Jobs RequestFromUser `json:"job_application"`
}
