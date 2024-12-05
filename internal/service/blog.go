package service

import (
	"time"
)

// BlogPost blog post
type BlogPost struct {
	ID              int64  `json:"id" gorm:"size:255;primaryKey;autoIncrement"`
	Code            string `json:"code" gorm:"size:255;uniqueIndex"`
	Title           string `json:"title,omitempty" gorm:"size:255"`
	ContentMarkdown string `json:"content_markdown,omitempty" gorm:"size:32767"`
	// ContentHTML string    `json:"content_html,omitempty" gorm:"size:32767"` // 2^16-1
	Status    string    `json:"status" gorm:"size:255"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (x *BlogPost) Fill() {
	// if x.ID == "" {
	// 	x.ID = uuid.New().String()
	// }
}

type BlogPostDAO struct {
	appService AppService
}

type BlogService interface {
	Posts() *BlogPostDAO
}

type defaultBlogService struct {
	appService AppService
	blogPost   BlogPostDAO
}

func newBlogService(appService AppService) BlogService {

	res := &defaultBlogService{

		appService: appService,
		blogPost: BlogPostDAO{
			appService: appService,
		},
	}

	return res
}

func (x *defaultBlogService) Posts() *BlogPostDAO {
	return &x.blogPost
}

func (x *BlogPostDAO) FindByID(id int64) (*BlogPost, error) {
	if id == 0 {
		return nil, nil // fmt.Errorf("id cannot be empty")
	}

	user := new(BlogPost)

	result := x.appService.Repository().Find(user, "id = ?", id)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return user, nil
}
func (x *BlogPostDAO) FindByCode(code string) (*BlogPost, error) {
	if code == "" {
		return nil, nil // fmt.Errorf("id cannot be empty")
	}

	user := new(BlogPost)

	result := x.appService.Repository().Find(user, "code = ?", code)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return user, nil
}
func (x *BlogPostDAO) ID(id int64) (int64, error) {
	if id == 0 {
		return 0, nil // fmt.Errorf("id cannot be empty")
	}

	user := new(BlogPost)

	result := x.appService.Repository().Select("id").Limit(1).Find(user, "id = ? ", id)

	if result.Error != nil || result.RowsAffected == 0 {
		return 0, result.Error
	}

	return user.ID, nil
}

func (x *BlogPostDAO) Code(code string) (int64, error) {
	if code == "" {
		return 0, nil // fmt.Errorf("id cannot be empty")
	}

	user := new(BlogPost)

	result := x.appService.Repository().Select("id").Find(user, "code = ?", code)

	if result.Error != nil || result.RowsAffected == 0 {
		return 0, result.Error
	}

	return user.ID, nil
}

func (x *BlogPostDAO) Create(data *BlogPost) error {

	repo := x.appService.Repository()
	data.Fill()
	res := repo.Create(data)
	return res.Error

}
func (x *BlogPostDAO) Update(data *BlogPost) error {
	repo := x.appService.Repository()
	res := repo.Model(data).Select("*" /*over all columns*/).Updates(data)
	return res.Error
}
func (x *BlogPostDAO) Delete(id int64) error {

	if id == 0 {
		return nil
	}

	repo := x.appService.Repository()
	res := repo.Delete(&BlogPost{ID: id})
	return res.Error
}
