package services

import (
	"context"
	"encoding/json"
	"job-portal-api/internal/models"
	"sync"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (s *Service) AddJobDetails(ctx context.Context, cj models.NewJobRequest, cid uint64) (models.Response, error) {
	// cj.CompanyId = uint64(cid)
	app := models.Job{
		CompanyId:           cid,
		JobTitle:            cj.JobTitle,
		Salary:              cj.Salary,
		MinimumNoticePeriod: cj.MinimumNoticePeriod,
		MaximumNoticePeriod: cj.MaximumNoticePeriod,
		Budget:              cj.Budget,
		JobDescription:      cj.JobDescription,
		MinExperience:       cj.MinExperience,
	}
	for _, v := range cj.QualificationIDs {
		tempData := models.Qualification{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.Qualifications = append(app.Qualifications, tempData)
	}
	for _, v := range cj.LocationIDs {
		tempData := models.Location{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.Locations = append(app.Locations, tempData)
	}
	for _, v := range cj.SkillIDs {
		tempData := models.Skill{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.Skills = append(app.Skills, tempData)
	}
	for _, v := range cj.WorkModeIDs {
		tempData := models.WorkMode{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.WorkModes = append(app.WorkModes, tempData)
	}
	for _, v := range cj.ShiftIDs {
		tempData := models.Shift{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.Shifts = append(app.Shifts, tempData)
	}
	for _, v := range cj.JobTypeIDs {
		tempData := models.JobType{
			Model: gorm.Model{
				ID: v,
			},
		}
		app.JobTypes = append(app.JobTypes, tempData)
	}
	jobData, err := s.UserRepo.PostJob(app)
	if err != nil {
		return models.Response{}, err
	}
	return jobData, nil
}

func (s *Service) ViewJobFromCompany(cid uint64) ([]models.Job, error) {
	jobData, err := s.UserRepo.GetJobsFromCompany(cid)
	if err != nil {
		return []models.Job{}, err
	}
	return jobData, nil
}

func (s *Service) ViewAllJobs(ctx context.Context) ([]models.Job, error) {
	jobDatas, err := s.UserRepo.GetAllJobs()
	if err != nil {
		return nil, err
	}
	return jobDatas, nil

}
func (s *Service) ViewJobById(ctx context.Context, jid uint64) ([]models.Job, error) {
	jobData, err := s.UserRepo.GetOneJob(jid)
	if err != nil {
		return []models.Job{}, nil
	}
	return jobData, nil
}

// func (s *Service) ProcessJobApplications(appData []models.NewUserApplication) ([]models.NewUserApplication, error) {
// 	ctx := context.Background()
// 	var wg = new(sync.WaitGroup)
// 	ch := make(chan models.NewUserApplication)
// 	var finalApplications []models.NewUserApplication

// 	for _, v := range appData {
// 		wg.Add(1)
// 		go func(v models.NewUserApplication) {
// 			defer wg.Done()
// 			var jobData models.Job
// 			val, err := s.rdb.GetCache(ctx, uint(v.ID))
// 			if err == redis.Nil {
// 				dbData, err := s.UserRepo.FetchJobData(v.ID)
// 				if err != nil {
// 					return
// 				}
// 				err = s.rdb.AddCache(ctx, uint(v.ID), dbData)
// 				if err != nil {
// 					return
// 				}
// 				jobData = dbData
// 			} else {
// 				err = json.Unmarshal([]byte(val), &jobData)
// 				if err == redis.Nil {
// 					return
// 				}
// 				if err != nil {
// 					return
// 				}

// 				check, val, err := s.MatchingCriteria(v, jobData)
// 				if err != nil {
// 					return
// 				}
// 				if check {
// 					ch <- val
// 				}
// 			}
// 		}(v)

// 	}

// 	go func() {
// 		wg.Wait()
// 		close(ch)
// 	}()

// 	for v := range ch {
// 		finalApplications = append(finalApplications, v)
// 	}

// 	return finalApplications, nil
// }

// func (s *Service) MatchingCriteria(appData models.NewUserApplication, jobData models.Job) (bool, models.NewUserApplication, error) {

// 	MatchedConditions := 0

// 	if appData.Jobs.Experience >= jobData.MinExperience {
// 		MatchedConditions++
// 	}

// 	if appData.Jobs.NoticePeriod >= jobData.MinimumNoticePeriod && appData.Jobs.NoticePeriod <= int(jobData.MaximumNoticePeriod) {
// 		MatchedConditions++
// 	}
// 	for _, v := range appData.Jobs.WorkModeIDs {
// 		for _, v1 := range jobData.WorkModes {
// 			if v == v1.ID {
// 				MatchedConditions++
// 				break
// 			}
// 		}
// 	}

// 	for _, v := range appData.Jobs.JobType {
// 		for _, v1 := range jobData.JobTypes {
// 			if v == v1.ID {
// 				MatchedConditions++
// 				break
// 			}
// 		}
// 	}

// 	for _, v := range appData.Jobs.Location {
// 		for _, v1 := range jobData.Locations {
// 			if v == v1.ID {
// 				MatchedConditions++
// 				break
// 			}
// 		}
// 	}

// 	for _, v := range appData.Jobs.Qualifications {
// 		for _, v1 := range jobData.Qualifications {
// 			if v == v1.ID {
// 				MatchedConditions++
// 				break
// 			}
// 		}
// 	}

// 	for _, v := range appData.Jobs.Skills {
// 		for _, v1 := range jobData.Skills {
// 			if v == v1.ID {
// 				MatchedConditions++
// 				break
// 			}
// 		}
// 	}

// 	for _, v := range appData.Jobs.Shift {
// 		for _, v1 := range jobData.Shifts {
// 			if v == v1.ID {
// 				MatchedConditions++
// 				break
// 			}
// 		}
// 	}

// 	totalConditions := 8
// 	if MatchedConditions*2 >= totalConditions {
// 		return true, appData, nil
// 	}

// 	return false, models.NewUserApplication{}, nil
// }



func(s *Service)ProcessJobApplications(applications []models.NewUserApplication)([]models.NewUserApplication,error){
	ctx:=context.Background()
	wg := new(sync.WaitGroup)
	ch := make(chan models.NewUserApplication)
	var finalData []models.NewUserApplication

	for _,v := range applications{
		wg.Add(1)
		go func (application models.NewUserApplication)  {
			defer wg.Done()

			var jobData models.Job

			val,err := s.rdb.GetCache(ctx,uint(application.ID))

			if err!=nil {
				jobDataFromDB, err := s.UserRepo.FetchJobData(application.ID)
				if err!=nil{
					log.Error().Err(err).Msg("invalid application job id does not exists")
					return
				}
				err = s.rdb.AddCache(ctx,uint(application.ID),jobDataFromDB)
				if err!=nil{
					return
				}
				jobData = jobDataFromDB
			}else{
				err = json.Unmarshal([]byte(val),&jobData)
				if err!=nil{
					return
				}
			}
			check := s.compareData(application,jobData)

			if check {
				ch <- application
			}

		}(v)
	}

	go func ()  {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		finalData = append(finalData, v)
	}

	return finalData, nil
}

func(s *Service) compareData(application models.NewUserApplication, jobData models.Job)bool{
	totalFields := 0
	matchedFields := 0

	totalFields++
	if application.Jobs.NoticePeriod>=jobData.MinimumNoticePeriod && application.Jobs.NoticePeriod<=int(jobData.MaximumNoticePeriod){
		matchedFields++
	}

	totalFields++
	if application.Jobs.Experience>=jobData.MinExperience && application.Jobs.Experience<=(jobData.MaxExperience){
		matchedFields++
	}

	count := 0
	totalFields++
	for _,v := range application.Jobs.Location{
		for _,v1 := range jobData.Locations{
			if v == v1.ID{
				count++
			}
		}
	}
	if count!=0 {
		matchedFields++
	}

	count = 0
	totalFields++
	for _,v := range application.Jobs.Skills{
		for _,v1 := range jobData.Skills{
			if v == v1.ID{
				count++
			}
		}
	}
	if count!=0 {
		matchedFields++
	}

	count = 0
	totalFields++
	for _,v := range application.Jobs.Qualifications{
		for _,v1 := range jobData.Qualifications{
			if v == v1.ID{
				count++
			}
		}
	}
	if count!=0 {
		matchedFields++
	}

	count = 0
	totalFields++
	for _,v := range application.Jobs.Shift{
		for _,v1 := range jobData.Shifts{
			if v == v1.ID{
				count++
			}
		}
	}
	if count!=0 {
		matchedFields++
	}

	count = 0
	totalFields++
	for _,v := range application.Jobs.JobType{
		for _,v1 := range jobData.JobTypes{
			if v == v1.ID{
				count++
			}
		}
	}
	if count!=0 {
		matchedFields++
	}

	if matchedFields*2 >= totalFields {
		return true
	}

	return false
}