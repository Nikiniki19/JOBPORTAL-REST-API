package services

import (
	"context"
	"job-portal-api/internal/models"
)

func (s *Service) AddCompanyDetails(ctx context.Context, companyData models.Company) (models.Company, error) {
	companyData, err := s.UserRepo.CreateCom(companyData)
	if err != nil {
		return models.Company{}, err
	}
	return companyData, nil

}
func (s *Service) ViewAllCompanies(ctx context.Context) ([]models.Company, error) {
	companyDetails, err := s.UserRepo.GetAllTheCompanies()
	if err != nil {
		return nil, err
	}
	return companyDetails, nil
}

func (s *Service) ViewCompanyDetails(ctx context.Context, id uint64) (models.Company, error) {
	companyData, err := s.UserRepo.GetCompany(id)
	if err != nil {
		return models.Company{}, err
	}
	return companyData, nil
}
