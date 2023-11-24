package repository

import (
	"errors"
	"job-portal-api/internal/models"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=jobRepo.go -destination=jobRepo_mock.go -package=repository

func (r *Repo) PostJob(nj models.Job) (models.Response, error) {

	res := r.DB.Create(&nj).Error
	if res != nil {
		log.Info().Err(res).Send()
		return models.Response{}, errors.New("job creation failed")
	}
	return models.Response{ID: uint64(nj.ID)}, nil
}
func (r *Repo) GetJobsFromCompany(comapny_id uint64) ([]models.Job, error) {
	var l []models.Job
	vx := r.DB.Where("company_id=?", comapny_id)
	err := vx.Find(&l).Error
	if err != nil {
		log.Info().Err(err).Send()
		return []models.Job{}, err
	}
	return l, nil
}
func (r *Repo) GetAllJobs() ([]models.Job, error) {
	var a []models.Job
	err := r.DB.Find(&a).Error
	if err != nil {
		return nil, err
	}
	return a, nil
}
func (r *Repo) GetOneJob(jid uint64) ([]models.Job, error) {
	var q []models.Job
	ax := r.DB.Where("id=?", jid)
	err := ax.Find(&jid).Error
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (r *Repo) FetchJobData(jid uint64) (models.Job, error) {
	var j models.Job
	result := r.DB.Preload("Comp").
		Preload("Locations").
		Preload("Skills").
		Preload("Qualifications").
		Preload("Shifts").
		Where("id = ?", jid).
		Find(&j)
	if result.Error != nil {

		log.Info().Err(result.Error).Send()
		return models.Job{}, result.Error
	}

	return j, nil
}
