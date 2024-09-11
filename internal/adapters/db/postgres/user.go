package postgres_repo

type userRepo struct {
	db *mongo.Database
}

func NewBookStorage(db *mongo.Database) *userRepo {
	return &userRepo{db: db}
}

func (bs *userRepo) GetOne(id string) *entity.Book {
	return nil
}
func (bs *userRepo) GetAll(limit, offset int) []*entity.Book {
	return nil
}
func (bs *userRepo) Create(book *entity.Book) *entity.Book {
	return nil
}
func (bs *userRepo) Delete(book *entity.Book) error {
	return nil
}
