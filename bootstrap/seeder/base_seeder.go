package seeder

import (
	"errors"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/internal/contract"
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
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
			configuration, err := c.GetConfigurationRepo().GetByKey("seed_executed")
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					c.GetLogger().ErrorWithCategory(sflogger.Category.Database.Database, sflogger.SubCategory.Database.Migration, "failed to get seed_executed configuration", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				configuration = &domain.Configuration{
					Key:   "seed_executed",
					Value: "false",
				}
			}
			if configuration.Value == "false" {
				seed.AddSeeder()
				seed.Seed()
			}
			if err := c.GetConfigurationRepo().SetKey("seed_executed", "true"); err != nil {
				c.GetLogger().ErrorWithCategory(sflogger.Category.Database.Database, sflogger.SubCategory.Database.Migration, "failed to update seed_executed configuration", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			break
		}
	}
}
