package plugin

import (
	"gorm.io/gorm"
)

type MultiTenancyPlugin struct{}

func (p *MultiTenancyPlugin) Name() string {
	return "multiTenancyPlugin"
}

func (p *MultiTenancyPlugin) Initialize(db *gorm.DB) (err error) {

	db.Callback().Create().Before("gorm:create").Register("multiTenancy:before_create", p.beforeCreate)
	db.Callback().Query().Before("gorm:query").Register("multiTenancy:before_query", p.beforeQuery)
	return
}

func (p *MultiTenancyPlugin) beforeCreate(db *gorm.DB) {

	if scopeID, ok := db.Get("scope_id"); ok {
		if field := db.Statement.Schema.LookUpField("scope_id"); field != nil {
			field.Set(db.Statement.Context, db.Statement.ReflectValue, scopeID)
		}
	}
}

func TenantScope(scopeID interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("scope_id = ?", scopeID)
	}
}
func (p *MultiTenancyPlugin) beforeQuery(db *gorm.DB) {

	if field := db.Statement.Schema.LookUpField("scope_id"); field != nil {
		db.Scopes(TenantScope(db.Statement.Context.Value("scope_id").(string)))
	}

}
