package seeder

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"os"
)

type seeder struct {
	container contract.IContainer
	store     []func() error
}

func (s *seeder) AddSeeder() {
	s.store = []func() error{s.provinceSeeder, s.citySeed}
}
func (s *seeder) Seed() {
	for _, f := range s.store {
		err := f()
		if err != nil {
			s.container.GetLogger().ErrorWithCategory(sflogger.Category.Database.Database, sflogger.SubCategory.Database.Migration, "province seed error", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

}

func InitializeSeeder(c contract.IContainer) {
	for _, arg := range os.Args[1:] {
		if arg == "--reset" || arg == "-reset" {
			seed := &seeder{
				container: c,
			}
			seed.AddSeeder()
			seed.Seed()
			break
		}
	}
}
