package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/exven/pos-system/internal/config"
	"github.com/exven/pos-system/shared/infrastructure/database"
	"github.com/exven/pos-system/shared/utils/crypto"
	"gorm.io/gorm"
)

func main() {
	var command string
	flag.StringVar(&command, "command", "", "Migration command: migrate, seed, rollback")
	flag.Parse()

	if command == "" {
		log.Fatal("Command is required. Use -command=migrate|seed|rollback")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.Initialize(database.Config{
		Host:               cfg.Database.Host,
		Port:               cfg.Database.Port,
		User:               cfg.Database.User,
		Password:           cfg.Database.Password,
		DBName:             cfg.Database.DBName,
		SSLMode:            cfg.Database.SSLMode,
		MaxConnections:     cfg.Database.MaxConnections,
		MaxIdleConnections: cfg.Database.MaxIdleConnections,
		ConnMaxLifetime:    cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	switch command {
	case "migrate":
		if err := runMigrations(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	case "seed":
		if err := runSeeds(db); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		fmt.Println("Demo data seeded successfully")
	case "rollback":
		fmt.Println("Rollback functionality not implemented yet")
		os.Exit(1)
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

func runMigrations(db *gorm.DB) error {
	// Auto migrate all schema models in proper order to handle foreign key dependencies
	return db.AutoMigrate(
		// Subscription and tenant management
		&database.SubscriptionPlan{},
		&database.Tenant{},
		&database.TenantSubscription{},
		
		// User management and roles
		&database.Role{},
		&database.User{},
		&database.Outlet{},
		&database.UserOutlet{},
		
		// Product management
		&database.ProductCategory{},
		&database.Product{},
		&database.ProductStock{},
		
		// Customer management
		&database.Customer{},
		
		// Sales and transactions
		&database.SalesTransaction{},
		&database.TransactionItem{},
		&database.TransactionPayment{},
		
		// Stock movements and inventory
		&database.StockMovement{},
		
		// Archive tables for data retention
		&database.ArchivedTransaction{},
		&database.ArchivedTransactionItem{},
		&database.ArchivedTransactionPayment{},
		
		// System and audit tables
		&database.AuditLog{},
		&database.DataRetentionLog{},
	)
}

func runSeeds(db *gorm.DB) error {
	seeder := NewSeeder(db)
	return seeder.Run()
}

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Run() error {
	fmt.Println("Starting seeding process...")

	// Seed subscription plans first
	if err := s.SeedSubscriptionPlans(); err != nil {
		return fmt.Errorf("failed to seed subscription plans: %w", err)
	}

	// Seed roles
	if err := s.SeedRoles(); err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	// Seed demo tenant
	tenantID, err := s.SeedDemoTenant()
	if err != nil {
		return fmt.Errorf("failed to seed demo tenant: %w", err)
	}

	// Seed demo users
	if err := s.SeedDemoUsers(tenantID); err != nil {
		return fmt.Errorf("failed to seed demo users: %w", err)
	}

	// Seed demo outlets
	if err := s.SeedDemoOutlets(tenantID); err != nil {
		return fmt.Errorf("failed to seed demo outlets: %w", err)
	}

	fmt.Println("Seeding completed successfully!")
	return nil
}

func (s *Seeder) SeedSubscriptionPlans() error {
	fmt.Println("Seeding subscription plans...")

	plans := []database.SubscriptionPlan{
		{
			Name:                    "Free",
			Description:             "Paket gratis dengan batasan fitur dan retensi data 14 hari",
			Price:                   0.00,
			MaxOutlets:              1,
			MaxUsers:                2,
			MaxProducts:             nil,
			MaxTransactionsPerMonth: nil,
			Features:                database.JSONFeatures{"basic_pos", "basic_reports", "data_retention_14_days"},
			IsActive:                true,
		},
		{
			Name:                    "Starter",
			Description:             "Paket untuk bisnis kecil",
			Price:                   99000.00,
			MaxOutlets:              2,
			MaxUsers:                5,
			MaxProducts:             nil,
			MaxTransactionsPerMonth: nil,
			Features:                database.JSONFeatures{"full_pos", "advanced_reports", "customer_management", "data_retention_unlimited"},
			IsActive:                true,
		},
		{
			Name:                    "Business",
			Description:             "Paket untuk bisnis menengah",
			Price:                   299000.00,
			MaxOutlets:              5,
			MaxUsers:                15,
			MaxProducts:             nil,
			MaxTransactionsPerMonth: nil,
			Features:                database.JSONFeatures{"full_pos", "advanced_reports", "customer_management", "inventory_management", "multi_payment", "data_retention_unlimited"},
			IsActive:                true,
		},
		{
			Name:                    "Enterprise",
			Description:             "Paket untuk bisnis besar",
			Price:                   599000.00,
			MaxOutlets:              999,
			MaxUsers:                999,
			MaxProducts:             nil,
			MaxTransactionsPerMonth: nil,
			Features:                database.JSONFeatures{"full_pos", "advanced_reports", "customer_management", "inventory_management", "multi_payment", "api_access", "custom_integration", "data_retention_unlimited"},
			IsActive:                true,
		},
	}

	for _, plan := range plans {
		var existingPlan database.SubscriptionPlan
		if err := s.db.Where("name = ?", plan.Name).First(&existingPlan).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := s.db.Create(&plan).Error; err != nil {
					return fmt.Errorf("failed to create subscription plan %s: %w", plan.Name, err)
				}
				fmt.Printf("Created subscription plan: %s\n", plan.Name)
			} else {
				return fmt.Errorf("failed to check existing subscription plan %s: %w", plan.Name, err)
			}
		} else {
			fmt.Printf("Subscription plan %s already exists, skipping...\n", plan.Name)
		}
	}

	return nil
}

func (s *Seeder) SeedRoles() error {
	fmt.Println("Seeding roles...")

	roles := []database.Role{
		{
			Name:        "super_admin",
			DisplayName: "Super Admin",
			Description: "Full system access",
			Permissions: database.JSONPermissions{"*"},
			IsSystem:    true,
		},
		{
			Name:        "tenant_owner",
			DisplayName: "Pemilik Bisnis",
			Description: "Full access dalam tenant",
			Permissions: database.JSONPermissions{"tenant.*"},
			IsSystem:    true,
		},
		{
			Name:        "manager",
			DisplayName: "Manager",
			Description: "Mengelola outlet dan laporan",
			Permissions: database.JSONPermissions{"outlet.*", "reports.*", "products.*", "customers.*"},
			IsSystem:    true,
		},
		{
			Name:        "cashier",
			DisplayName: "Kasir",
			Description: "Melakukan penjualan",
			Permissions: database.JSONPermissions{"sales.*", "customers.read", "products.read"},
			IsSystem:    true,
		},
	}

	for _, role := range roles {
		var existingRole database.Role
		if err := s.db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := s.db.Create(&role).Error; err != nil {
					return fmt.Errorf("failed to create role %s: %w", role.Name, err)
				}
				fmt.Printf("Created role: %s\n", role.Name)
			} else {
				return fmt.Errorf("failed to check existing role %s: %w", role.Name, err)
			}
		} else {
			fmt.Printf("Role %s already exists, skipping...\n", role.Name)
		}
	}

	return nil
}

func (s *Seeder) SeedDemoTenant() (uint64, error) {
	fmt.Println("Seeding demo tenant...")

	// Check if demo tenant already exists
	var existingTenant database.Tenant
	if err := s.db.Where("email = ?", "demo@example.com").First(&existingTenant).Error; err == nil {
		fmt.Printf("Demo tenant already exists with ID: %d\n", existingTenant.ID)
		return existingTenant.ID, nil
	}

	// Get the Free subscription plan
	var freePlan database.SubscriptionPlan
	if err := s.db.Where("name = ?", "Free").First(&freePlan).Error; err != nil {
		return 0, fmt.Errorf("failed to find Free subscription plan: %w", err)
	}

	// Create demo tenant
	tenant := database.Tenant{
		Name:         "Demo Coffee Shop",
		BusinessType: "Restaurant",
		Email:        "demo@example.com",
		Phone:        "+62812345678",
		Address:      "Jl. Demo No. 123",
		City:         "Jakarta",
		Province:     "DKI Jakarta",
		PostalCode:   "12345",
		TaxNumber:    "12.345.678.9-012.000",
		Timezone:     "Asia/Jakarta",
		Currency:     "IDR",
		IsActive:     true,
	}

	if err := s.db.Create(&tenant).Error; err != nil {
		return 0, fmt.Errorf("failed to create demo tenant: %w", err)
	}

	// Create tenant subscription
	subscription := database.TenantSubscription{
		TenantID:           tenant.ID,
		SubscriptionPlanID: freePlan.ID,
		Status:             database.SubscriptionStatusActive,
		StartsAt:           time.Now(),
		EndsAt:             time.Now().AddDate(1, 0, 0), // 1 year from now
		AutoRenew:          true,
		PaymentMethod:      "demo",
	}

	if err := s.db.Create(&subscription).Error; err != nil {
		return 0, fmt.Errorf("failed to create tenant subscription: %w", err)
	}

	fmt.Printf("Created demo tenant: %s with ID: %d\n", tenant.Name, tenant.ID)
	return tenant.ID, nil
}

func (s *Seeder) SeedDemoUsers(tenantID uint64) error {
	fmt.Println("Seeding demo users...")

	// Get roles
	var ownerRole, managerRole, cashierRole database.Role
	if err := s.db.Where("name = ?", "tenant_owner").First(&ownerRole).Error; err != nil {
		return fmt.Errorf("failed to find tenant_owner role: %w", err)
	}
	if err := s.db.Where("name = ?", "manager").First(&managerRole).Error; err != nil {
		return fmt.Errorf("failed to find manager role: %w", err)
	}
	if err := s.db.Where("name = ?", "cashier").First(&cashierRole).Error; err != nil {
		return fmt.Errorf("failed to find cashier role: %w", err)
	}

	// Hash default password
	hashedPassword, err := crypto.HashPassword("password123")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	users := []database.User{
		{
			TenantID:     tenantID,
			RoleID:       ownerRole.ID,
			Username:     "owner",
			Email:        "owner@demo.com",
			PasswordHash: hashedPassword,
			FullName:     "Demo Owner",
			Phone:        "+62812345001",
			IsActive:     true,
		},
		{
			TenantID:     tenantID,
			RoleID:       managerRole.ID,
			Username:     "manager",
			Email:        "manager@demo.com",
			PasswordHash: hashedPassword,
			FullName:     "Demo Manager",
			Phone:        "+62812345002",
			IsActive:     true,
		},
		{
			TenantID:     tenantID,
			RoleID:       cashierRole.ID,
			Username:     "cashier1",
			Email:        "cashier1@demo.com",
			PasswordHash: hashedPassword,
			FullName:     "Demo Cashier 1",
			Phone:        "+62812345003",
			IsActive:     true,
		},
		{
			TenantID:     tenantID,
			RoleID:       cashierRole.ID,
			Username:     "cashier2",
			Email:        "cashier2@demo.com",
			PasswordHash: hashedPassword,
			FullName:     "Demo Cashier 2",
			Phone:        "+62812345004",
			IsActive:     true,
		},
	}

	for _, user := range users {
		var existingUser database.User
		if err := s.db.Where("tenant_id = ? AND username = ?", user.TenantID, user.Username).First(&existingUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := s.db.Create(&user).Error; err != nil {
					return fmt.Errorf("failed to create user %s: %w", user.Username, err)
				}
				fmt.Printf("Created user: %s (%s)\n", user.Username, user.FullName)
			} else {
				return fmt.Errorf("failed to check existing user %s: %w", user.Username, err)
			}
		} else {
			fmt.Printf("User %s already exists, skipping...\n", user.Username)
		}
	}

	return nil
}

func (s *Seeder) SeedDemoOutlets(tenantID uint64) error {
	fmt.Println("Seeding demo outlets...")

	// Get manager user for outlet assignment
	var managerUser database.User
	if err := s.db.Where("tenant_id = ? AND username = ?", tenantID, "manager").First(&managerUser).Error; err != nil {
		return fmt.Errorf("failed to find manager user: %w", err)
	}

	outlets := []database.Outlet{
		{
			TenantID:    tenantID,
			Name:        "Main Store",
			Code:        "MAIN",
			Description: "Main outlet for Demo Coffee Shop",
			Address:     "Jl. Demo No. 123",
			City:        "Jakarta",
			Province:    "DKI Jakarta",
			PostalCode:  "12345",
			Phone:       "+62812345678",
			Email:       "main@demo.com",
			ManagerID:   &managerUser.ID,
			IsActive:    true,
			Settings:    database.JSONSettings{"tax_rate": 10, "currency": "IDR", "printer_enabled": true},
		},
	}

	// Get cashier users for outlet assignment
	var cashierUsers []database.User
	if err := s.db.Where("tenant_id = ?", tenantID).Joins("JOIN roles ON users.role_id = roles.id AND roles.name = ?", "cashier").Find(&cashierUsers).Error; err != nil {
		return fmt.Errorf("failed to find cashier users: %w", err)
	}

	for _, outlet := range outlets {
		var existingOutlet database.Outlet
		if err := s.db.Where("tenant_id = ? AND code = ?", outlet.TenantID, outlet.Code).First(&existingOutlet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := s.db.Create(&outlet).Error; err != nil {
					return fmt.Errorf("failed to create outlet %s: %w", outlet.Code, err)
				}
				fmt.Printf("Created outlet: %s (%s)\n", outlet.Code, outlet.Name)

				// Assign manager and cashiers to the outlet
				userOutlets := []database.UserOutlet{
					{
						UserID:   managerUser.ID,
						OutletID: outlet.ID,
						IsActive: true,
					},
				}

				// Add all cashiers to the outlet
				for _, cashier := range cashierUsers {
					userOutlets = append(userOutlets, database.UserOutlet{
						UserID:   cashier.ID,
						OutletID: outlet.ID,
						IsActive: true,
					})
				}

				for _, userOutlet := range userOutlets {
					if err := s.db.Create(&userOutlet).Error; err != nil {
						return fmt.Errorf("failed to assign user %d to outlet %d: %w", userOutlet.UserID, userOutlet.OutletID, err)
					}
				}

				fmt.Printf("Assigned %d users to outlet %s\n", len(userOutlets), outlet.Code)
			} else {
				return fmt.Errorf("failed to check existing outlet %s: %w", outlet.Code, err)
			}
		} else {
			fmt.Printf("Outlet %s already exists, skipping...\n", outlet.Code)
		}
	}

	return nil
}