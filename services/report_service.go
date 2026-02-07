package services

import (
	"go-kasir-api/models"
	"go-kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetReport(startDate, endDate string) (*models.Report, error) {
	return s.repo.GetReport(startDate, endDate)
}
