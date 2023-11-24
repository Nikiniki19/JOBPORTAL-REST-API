package repository

import (
	"errors"
	"job-portal-api/internal/models"

	"github.com/rs/zerolog/log"
)

func (r *Repo) CreateCom(nc models.Company) (models.Company, error) {
	err := r.DB.Create(&nc).Error
	if err != nil {
		log.Info().Err(err).Send()
		return models.Company{}, errors.New("company cannot be created")
	}
	return nc, nil
}
func (r *Repo) GetAllTheCompanies() ([]models.Company, error) {
	var f []models.Company
	err := r.DB.Find(&f).Error
	if err != nil {
		log.Info().Err(err).Send()
		return []models.Company{}, err
	}
	return f, nil
}

func (r *Repo) GetCompany(id uint64) (models.Company, error) {
	var z models.Company
	ax := r.DB.Where("id=?", id)
	err := ax.First(&z).Error
	if err != nil {
		log.Info().Err(err).Send()
		return models.Company{}, err
	}
	return z, nil
}
