package store

import (
	"context"
	"log"
	"os"

	"server/errors"
	"server/objects"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type pg struct {
	db *gorm.DB
}

// NewPostgresPostStore returns a postgres implementation of Post store
func NewPostgresPostStore(conn string) IPostStore {
	// create database connection
	db, err := gorm.Open(postgres.Open(conn),
		&gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "", log.LstdFlags),
				logger.Config{
					LogLevel: logger.Info,
					Colorful: true,
				},
			),
		},
	)
	if err != nil {
		panic("Enable to connect to database: " + err.Error())
	}
	if err := db.AutoMigrate(&objects.Post{}); err != nil {
		panic("Enable to migrate database: " + err.Error())
	}
	// return store implementation
	return &pg{db: db}
}

func (p *pg) Get(ctx context.Context, in *objects.GetRequest) (*objects.Post, error) {
	post := &objects.Post{}
	// take Post where id == uid from database
	err := p.db.WithContext(ctx).Take(post, "id = ?", in.ID).Error
	if err == gorm.ErrRecordNotFound {
		// not found
		return nil, errors.ErrPostNotFound
	}
	return post, err
}

func (p *pg) List(ctx context.Context, in *objects.ListRequest) ([]*objects.Post, error) {
	if in.Limit == 0 || in.Limit > objects.MaxListLimit {
		in.Limit = objects.MaxListLimit
	}
	query := p.db.WithContext(ctx).Limit(in.Limit)
	if in.After != "" {
		query = query.Where("id > ?", in.After)
	}
	list := make([]*objects.Post, 0, in.Limit)
	err := query.Order("id").Find(&list).Error
	return list, err
}

func (p *pg) Create(ctx context.Context, in *objects.CreateRequest) error {
	if in.Post == nil {
		return errors.ErrObjectIsRequired
	}
	in.Post.ID = GenerateUniqueID()
	return p.db.WithContext(ctx).
		Create(in.Post).
		Error
}

// func (p *pg) UpdateDetails(ctx context.Context, in *objects.UpdateDetailsRequest) error {
// 	post := &objects.Post{
// 		ID:          in.ID,
// 		title:        in.title,
// 		Description: in.Description,
// 		Website:     in.Website,
// 		Address:     in.Address,
// 		PhoneNumber: in.PhoneNumber,
// 		UpdatedOn:   p.db.NowFunc(),
// 	}
// 	return p.db.WithContext(ctx).Model(post).
// 		Select("name", "description", "website", "address", "phone_number", "updated_on").
// 		Updates(post).
// 		Error
// }

func (p *pg) Delete(ctx context.Context, in *objects.DeleteRequest) error {
	post := &objects.Post{ID: in.ID}
	return p.db.WithContext(ctx).Model(post).
		Delete(post).
		Error
}
