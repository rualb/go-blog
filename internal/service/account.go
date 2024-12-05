package service

// UserAccount Username,Email,NormalizedEmail are uniqueIndex with condition "not empty"
type UserAccount struct {
	ID       string `gorm:"primaryKey"`
	Username string // `gorm:"uniqueIndex:,where:username != ''"`
	// PhoneNumber string
	// Email string // use this on emailing and show
	// NormalizedEmail string // `gorm:"uniqueIndex:,where:normalized_email != ''"` // use this on search
	// // SecurityStamp   string // Key := Base32(Random(32))  HMACSHA1(Key)  Key == VTOQQ2PQKD7A2KTSXU7OFLKUNI7QEZRJ
	// PasswordHash string
	// CreatedAt    time.Time
	Roles string
}

// // HasAllRoles checks if the user has all the specified roles.
// func (x *UserAccount) HasAllRoles(roles ...string) bool {
// 	return utilaccess.HasAllRoles(x.Roles, roles...)
// }

// AccountService is a service for managing user account.
type AccountService interface {
	FindByID(id string) (*UserAccount, error)
}

type defaultAccountService struct {
	appService AppService
}

// NewAccountService is constructor.
func newAccountService(appService AppService) AccountService {

	return &defaultAccountService{
		appService: appService,
	}
}

func (x defaultAccountService) FindByID(id string) (*UserAccount, error) {

	if id == "" {
		return nil, nil // fmt.Errorf("id cannot be empty")
	}

	user := new(UserAccount)

	result := x.appService.Repository().Find(user, "id = ?", id)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return user, nil
}
