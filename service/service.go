package service

import (
	"context"
	"log"
	"time"
	"user-access-management/db"

	_ "github.com/lib/pq" // pq driver.
)

const (
	timeout time.Duration = 1 //seconds
)

//Service is an interface for other concrete service to implement sql DB operation
type Service interface {
	DBStatus() (bool, error)
}

//UserManagementService is a class with DB connection information. It also implements Service interface
type UserManagementService struct {
	db db.SQLDatabase
}

//NewUserManagementService is a constructor which can creates an object of the User Management Service class
func NewUserManagementService(db db.SQLDatabase) *UserManagementService {
	return &UserManagementService{
		db: db,
	}
}

//DBStatus returns the database connection status
func (d *UserManagementService) DBStatus() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	err := d.db.PingContext(ctx)
	if err != nil {
		log.Printf("ping context error: %v", err)
		return false, err
	}
	return true, nil
}
