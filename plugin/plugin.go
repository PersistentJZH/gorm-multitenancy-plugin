package plugin

import (
	"gorm.io/gorm"
)

type MultiTenancyPlugin struct{}

func (p *MultiTenancyPlugin) Name() string {
	return "multiTenancyPlugin"
}

func (p *MultiTenancyPlugin) Initialize(db *gorm.DB) (err error) {
	// 添加你的创建和查询钩子
	db.Callback().Create().Before("gorm:create").Register("multiTenancy:before_create", p.beforeCreate)
	db.Callback().Query().Before("gorm:query").Register("multiTenancy:before_query", p.beforeQuery)
	return
}

func (p *MultiTenancyPlugin) beforeCreate(db *gorm.DB) {
	// 在创建记录前添加 TenantID
	if tenantID, ok := db.Get("tenant_id"); ok {
		if field := db.Statement.Schema.LookUpField("TenantID"); field != nil {
			field.Set(db.Statement.Context, db.Statement.ReflectValue, tenantID)
		}
	}
}

func TenantScope(tenantID interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ?", tenantID)
	}
}
func (p *MultiTenancyPlugin) beforeQuery(db *gorm.DB) {
	// 在查询记录前添加 TenantID
	if tenantID, ok := db.Get("tenant_id"); ok {
		if field := db.Statement.Schema.LookUpField("TenantID"); field != nil {
			db.Scopes(TenantScope(tenantID))
		}
	}
}
