package database

import (
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"savory-ai-server/app/storage"
	"savory-ai-server/utils/config"
)

// setup storage with gorm
type Database struct {
	DB  *gorm.DB
	Log zerolog.Logger
	Cfg *config.Config
}

type Seeder interface {
	Seed(*gorm.DB) error
	Count() (int, error)
}

func NewDatabase(cfg *config.Config, log zerolog.Logger) *Database {
	db := &Database{
		Cfg: cfg,
		Log: log,
	}

	return db
}

// connect storage
func (db *Database) ConnectDatabase() {
	conn, err := gorm.Open(postgres.Open(db.Cfg.DB.Postgres.DSN), &gorm.Config{})
	if err != nil {
		db.Log.Error().Err(err).Msg("An unknown error occurred when to connect the storage!")
	} else {
		db.Log.Info().Msg("Connected the storage succesfully!")
	}

	db.DB = conn
}

// shutdown storage
func (db *Database) ShutdownDatabase() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		db.Log.Error().Err(err).Msg("An unknown error occurred when to shutdown the storage!")
	} else {
		db.Log.Info().Msg("Shutdown the storage succesfully!")
	}
	if err = sqlDB.Close(); err != nil {
		db.Log.Error().Err(err).Msg("An unknown error occurred when to shutdown the storage!")
	}
}

// migrate models
func (db *Database) MigrateModels() {
	if err := db.DB.AutoMigrate(
		Models()...,
	); err != nil {
		db.Log.Error().Err(err).Msg("An unknown error occurred when to migrate the storage!")
	}
}

// list of models for migration
func Models() []any {
	return []any{
		&storage.User{},
		&storage.MenuCategory{},
		&storage.Dish{},
		&storage.Ingredient{},
		&storage.Table{},
		&storage.Question{},
		&storage.Restaurant{},
		&storage.WorkingHour{},
		&storage.Organization{},
		&storage.TableChatSessions{},
		&storage.TableChatMessage{},
		&storage.RestaurantChatSessions{},
		&storage.RestaurantChatMessage{},
		&storage.PasswordResetCode{},
	}
}

// seed data
func (db *Database) SeedModels(seeder ...Seeder) {
	for _, seed := range seeder {
		count, err := seed.Count()
		if err != nil {
			db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the storage!")
		}

		if count == 0 {
			if err := seed.Seed(db.DB); err != nil {
				db.Log.Error().Err(err).Msg("An unknown error occurred when to seed the storage!")
			}

			db.Log.Info().Msg("Seeded the storage succesfully!")
		} else {
			db.Log.Info().Msg("Database is already seeded!")
		}
	}

	db.Log.Info().Msg("Seeded the storage succesfully!")
}
