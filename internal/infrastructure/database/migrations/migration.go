package migrations

import (
	"fmt"
	"log"
	"sort"
	"time"

	"gorm.io/gorm"
)

type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Batch     int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (MigrationRecord) TableName() string {
	return "migrations"
}

type MigrationEntry struct {
	Name string
	Up   func(*gorm.DB) error
	Down func(*gorm.DB) error
}

var registry []MigrationEntry

func RegisterMigration(name string, up func(*gorm.DB) error, down func(*gorm.DB) error) {
	registry = append(registry, MigrationEntry{
		Name: name,
		Up:   up,
		Down: down,
	})
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return err
	}

	var lastBatch int
	if err := db.Model(&MigrationRecord{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch).Error; err != nil {
		return err
	}

	sort.Slice(registry, func(i, j int) bool {
		return registry[i].Name < registry[j].Name
	})

	nextBatch := lastBatch + 1
	ran := 0

	for _, m := range registry {
		var count int64
		if err := db.Model(&MigrationRecord{}).Where("name = ?", m.Name).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		if err := m.Up(db); err != nil {
			return fmt.Errorf("failed migrating %s: %w", m.Name, err)
		}

		if err := db.Create(&MigrationRecord{Name: m.Name, Batch: nextBatch}).Error; err != nil {
			return err
		}

		log.Printf("✅ Migrated: %s (batch %d)\n", m.Name, nextBatch)
		ran++
	}

	if ran == 0 {
		log.Println("ℹ️  Nothing to migrate")
	}
	return nil
}

func Rollback(db *gorm.DB) error {
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return err
	}

	var lastBatch int
	if err := db.Model(&MigrationRecord{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch).Error; err != nil {
		return err
	}
	if lastBatch == 0 {
		log.Println("ℹ️  No migrations to rollback")
		return nil
	}

	var records []MigrationRecord
	if err := db.Where("batch = ?", lastBatch).Order("id DESC").Find(&records).Error; err != nil {
		return err
	}

	for _, rec := range records {
		for _, m := range registry {
			if m.Name == rec.Name {
				if err := m.Down(db); err != nil {
					return fmt.Errorf("failed rolling back %s: %w", m.Name, err)
				}
				break
			}
		}
		if err := db.Delete(&rec).Error; err != nil {
			return err
		}
		log.Printf("⏪ Rolled back: %s\n", rec.Name)
	}

	return nil
}

func DropAll(db *gorm.DB) error {
	sorted := make([]MigrationEntry, len(registry))
	copy(sorted, registry)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name > sorted[j].Name
	})

	for _, m := range sorted {
		if err := m.Down(db); err != nil {
			return fmt.Errorf("failed dropping %s: %w", m.Name, err)
		}
		log.Printf("🗑️  Dropped: %s\n", m.Name)
	}

	if db.Migrator().HasTable(&MigrationRecord{}) {
		if err := db.Migrator().DropTable(&MigrationRecord{}); err != nil {
			return err
		}
	}

	return nil
}

func Fresh(db *gorm.DB) error {
	if err := DropAll(db); err != nil {
		return err
	}
	return Migrate(db)
}
