package services

import (
	"proyecto-bd-final/internal/models"
	"proyecto-bd-final/internal/repository"
)

type BookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetBookByID(id int) (*models.Book, error) {
	return nil, nil
}
