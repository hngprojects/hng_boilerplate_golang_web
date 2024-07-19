package seed

import "gorm.io/gorm"

type Seeder struct {
	Run func(*gorm.DB) error
}

func Run() {

}
