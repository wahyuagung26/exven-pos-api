-- Schema Database untuk Aplikasi POS Multi-Tenant (PostgreSQL)
-- Created for retail POS with multi-outlet support

-- =============================================
-- TENANT & SUBSCRIPTION MANAGEMENT
-- =============================================

-- Tabel untuk paket langganan
CREATE TABLE subscription_plans (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    max_outlets INTEGER NOT NULL DEFAULT 1,
    max_users INTEGER NOT NULL DEFAULT 1,
    max_products INTEGER DEFAULT NULL, -- NULL = unlimited
    max_transactions_per_month INTEGER DEFAULT NULL,
    features JSONB, -- Array fitur yang tersedia
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);



-- Tabel tenant (perusahaan/bisnis)
CREATE TABLE tenants (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    business_type VARCHAR(100),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    tax_number VARCHAR(50),
    logo_url VARCHAR(500),
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    currency VARCHAR(3) DEFAULT 'IDR',
    is_active BOOLEAN DEFAULT TRUE,
    trial_ends_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenants_email ON tenants(email);
CREATE INDEX idx_tenants_is_active ON tenants(is_active);



-- Tabel langganan tenant
CREATE TYPE subscription_status AS ENUM ('active', 'cancelled', 'expired', 'pending');

CREATE TABLE tenant_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    subscription_plan_id BIGINT NOT NULL,
    status subscription_status DEFAULT 'pending',
    starts_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ends_at TIMESTAMP WITH TIME ZONE NOT NULL,
    auto_renew BOOLEAN DEFAULT TRUE,
    payment_method VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (subscription_plan_id) REFERENCES subscription_plans(id)
);

CREATE INDEX idx_tenant_subscriptions_tenant_status ON tenant_subscriptions(tenant_id, status);
CREATE INDEX idx_tenant_subscriptions_ends_at ON tenant_subscriptions(ends_at);



-- =============================================
-- USER MANAGEMENT & ROLES
-- =============================================

-- Tabel roles
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    permissions JSONB, -- Array permission strings
    is_system BOOLEAN DEFAULT FALSE, -- System roles tidak bisa dihapus
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabel users
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar_url VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE NULL,
    email_verified_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    UNIQUE (tenant_id, email)
);

CREATE INDEX idx_users_tenant_active ON users(tenant_id, is_active);



-- =============================================
-- OUTLET MANAGEMENT
-- =============================================

-- Tabel outlets
CREATE TABLE outlets (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL, -- Kode unik outlet dalam tenant
    description TEXT,
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    phone VARCHAR(20),
    email VARCHAR(255),
    manager_id BIGINT, -- User yang bertanggung jawab
    is_active BOOLEAN DEFAULT TRUE,
    settings JSONB, -- Pengaturan khusus outlet (printer, tax, etc)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE (tenant_id, code)
);

CREATE INDEX idx_outlets_tenant_active ON outlets(tenant_id, is_active);



-- Tabel untuk assign user ke outlet
CREATE TABLE user_outlets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE,
    UNIQUE (user_id, outlet_id)
);

CREATE INDEX idx_user_outlets_outlet_active ON user_outlets(outlet_id, is_active);

-- =============================================
-- PRODUCT MANAGEMENT
-- =============================================

-- Tabel kategori produk
CREATE TABLE product_categories (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    parent_id BIGINT NULL, -- Untuk sub-kategori
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(500),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES product_categories(id) ON DELETE SET NULL
);

CREATE INDEX idx_product_categories_tenant_parent ON product_categories(tenant_id, parent_id);
CREATE INDEX idx_product_categories_sort_order ON product_categories(sort_order);



-- Tabel produk
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    category_id BIGINT,
    sku VARCHAR(100) NOT NULL,
    barcode VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    unit VARCHAR(50) DEFAULT 'pcs', -- satuan (pcs, kg, liter, dll)
    cost_price DECIMAL(12,2) DEFAULT 0.00,
    selling_price DECIMAL(12,2) NOT NULL,
    min_stock INTEGER DEFAULT 0,
    track_stock BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    images JSONB, -- Array URL gambar
    variants JSONB, -- Array varian produk (size, color, dll)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES product_categories(id) ON DELETE SET NULL,
    UNIQUE (tenant_id, sku)
);

CREATE INDEX idx_products_tenant_category ON products(tenant_id, category_id);
CREATE INDEX idx_products_barcode ON products(barcode);
CREATE INDEX idx_products_name ON products(name);



-- Tabel stock per outlet
CREATE TABLE product_stocks (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER DEFAULT 0, -- Stock yang di-reserve untuk order
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE,
    UNIQUE (product_id, outlet_id)
);

CREATE INDEX idx_product_stocks_outlet_quantity ON product_stocks(outlet_id, quantity);



-- =============================================
-- CUSTOMER MANAGEMENT
-- =============================================

-- Tabel pelanggan
CREATE TYPE gender_type AS ENUM ('male', 'female');

CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    code VARCHAR(50), -- Kode pelanggan
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    birth_date DATE,
    gender gender_type,
    loyalty_points INTEGER DEFAULT 0,
    total_spent DECIMAL(15,2) DEFAULT 0.00,
    visit_count INTEGER DEFAULT 0,
    last_visit_at TIMESTAMP WITH TIME ZONE NULL,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE (tenant_id, code)
);

CREATE INDEX idx_customers_tenant_phone ON customers(tenant_id, phone);
CREATE INDEX idx_customers_tenant_email ON customers(tenant_id, email);



-- =============================================
-- SALES & TRANSACTIONS
-- =============================================

-- Types untuk payment method dan status
CREATE TYPE payment_method_type AS ENUM ('cash', 'card', 'transfer', 'ewallet', 'multiple');
CREATE TYPE transaction_status_type AS ENUM ('completed', 'cancelled', 'refunded');

-- Tabel transaksi penjualan dengan denormalisasi customer data
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    cashier_id BIGINT NOT NULL,
    customer_id BIGINT NULL,
    -- Snapshot data customer saat transaksi (denormalisasi)
    customer_name_snapshot VARCHAR(255),
    customer_phone_snapshot VARCHAR(20),
    customer_email_snapshot VARCHAR(255),
    -- Snapshot data cashier untuk historical reference
    cashier_name_snapshot VARCHAR(255) NOT NULL,
    -- Snapshot data outlet untuk historical reference  
    outlet_name_snapshot VARCHAR(255) NOT NULL,
    outlet_code_snapshot VARCHAR(50) NOT NULL,
    -- Data transaksi
    transaction_number VARCHAR(100) NOT NULL UNIQUE,
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    subtotal DECIMAL(15,2) NOT NULL,
    discount_amount DECIMAL(15,2) DEFAULT 0.00,
    tax_amount DECIMAL(15,2) DEFAULT 0.00,
    total_amount DECIMAL(15,2) NOT NULL,
    paid_amount DECIMAL(15,2) NOT NULL,
    change_amount DECIMAL(15,2) DEFAULT 0.00,
    payment_method payment_method_type NOT NULL,
    status transaction_status_type DEFAULT 'completed',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (cashier_id) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL
);

CREATE INDEX idx_transactions_tenant_outlet_date ON transactions(tenant_id, outlet_id, transaction_date);
CREATE INDEX idx_transactions_cashier_date ON transactions(cashier_id, transaction_date);
CREATE INDEX idx_transactions_customer_snapshot ON transactions(customer_name_snapshot, customer_phone_snapshot);
CREATE INDEX idx_transactions_outlet_snapshot ON transactions(outlet_code_snapshot);
CREATE INDEX idx_transactions_cashier_snapshot ON transactions(cashier_name_snapshot);



-- Tabel detail item penjualan dengan denormalisasi untuk historical accuracy
CREATE TABLE transaction_items (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    -- Snapshot data produk saat transaksi (denormalisasi)
    product_name_snapshot VARCHAR(255) NOT NULL,
    product_sku_snapshot VARCHAR(100) NOT NULL,
    product_category_snapshot VARCHAR(255),
    product_unit_snapshot VARCHAR(50) DEFAULT 'pcs',
    -- Data transaksi
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL,
    cost_price_snapshot DECIMAL(12,2) DEFAULT 0.00, -- Untuk profit calculation
    discount_amount DECIMAL(12,2) DEFAULT 0.00,
    total_price DECIMAL(15,2) NOT NULL,
    notes TEXT,
    
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX idx_transaction_items_transaction ON transaction_items(transaction_id);
CREATE INDEX idx_transaction_items_product ON transaction_items(product_id);
CREATE INDEX idx_transaction_items_product_sku_snapshot ON transaction_items(product_sku_snapshot);
CREATE INDEX idx_transaction_items_product_name_snapshot ON transaction_items(product_name_snapshot);

-- Tabel pembayaran (untuk multiple payment method)
CREATE TYPE payment_method_single AS ENUM ('cash', 'card', 'transfer', 'ewallet');

CREATE TABLE transaction_payments (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    payment_method payment_method_single NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    reference_number VARCHAR(100), -- Nomor referensi untuk non-cash
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE
);

CREATE INDEX idx_transaction_payments_transaction ON transaction_payments(transaction_id);

-- =============================================
-- STOCK MOVEMENTS & INVENTORY
-- =============================================

-- Types untuk stock movement
CREATE TYPE movement_type AS ENUM ('in', 'out', 'adjustment', 'transfer');
CREATE TYPE reference_type AS ENUM ('sale', 'purchase', 'adjustment', 'transfer', 'initial');

CREATE TABLE stock_movements (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    movement_type movement_type NOT NULL,
    quantity INTEGER NOT NULL, -- Bisa negatif untuk stock out
    reference_type reference_type NOT NULL,
    reference_id BIGINT, -- ID dari transaksi terkait
    notes TEXT,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX idx_stock_movements_product_outlet ON stock_movements(product_id, outlet_id);
CREATE INDEX idx_stock_movements_outlet_date ON stock_movements(outlet_id, created_at);
CREATE INDEX idx_stock_movements_reference ON stock_movements(reference_type, reference_id);

-- =============================================
-- BACKUP TABLES FOR DATA RETENTION
-- =============================================

-- Tabel backup transaksi untuk tenant gratis
CREATE TABLE archived_transactions (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    cashier_id BIGINT NOT NULL,
    customer_id BIGINT NULL,
    -- Snapshot data customer saat transaksi (denormalisasi)
    customer_name_snapshot VARCHAR(255),
    customer_phone_snapshot VARCHAR(20),
    customer_email_snapshot VARCHAR(255),
    -- Snapshot data cashier untuk historical reference
    cashier_name_snapshot VARCHAR(255) NOT NULL,
    -- Snapshot data outlet untuk historical reference  
    outlet_name_snapshot VARCHAR(255) NOT NULL,
    outlet_code_snapshot VARCHAR(50) NOT NULL,
    -- Data transaksi
    transaction_number VARCHAR(100) NOT NULL,
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL,
    discount_amount DECIMAL(15,2) DEFAULT 0.00,
    tax_amount DECIMAL(15,2) DEFAULT 0.00,
    total_amount DECIMAL(15,2) NOT NULL,
    paid_amount DECIMAL(15,2) NOT NULL,
    change_amount DECIMAL(15,2) DEFAULT 0.00,
    payment_method payment_method_type NOT NULL,
    status transaction_status_type DEFAULT 'completed',
    notes TEXT,
    original_created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    original_updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    archived_reason VARCHAR(100) DEFAULT 'data_retention_policy'
);

CREATE INDEX idx_archived_transactions_tenant_date ON archived_transactions(tenant_id, transaction_date);
CREATE INDEX idx_archived_transactions_number ON archived_transactions(transaction_number);
CREATE INDEX idx_archived_transactions_archived_date ON archived_transactions(archived_at);

-- Tabel backup detail item penjualan 
CREATE TABLE archived_transaction_items (
    id BIGINT PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    -- Snapshot data yang sudah tersimpan dari tabel asli
    product_name_snapshot VARCHAR(255) NOT NULL,
    product_sku_snapshot VARCHAR(100) NOT NULL,
    product_category_snapshot VARCHAR(255),
    product_unit_snapshot VARCHAR(50) DEFAULT 'pcs',
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL,
    cost_price_snapshot DECIMAL(12,2) DEFAULT 0.00,
    discount_amount DECIMAL(12,2) DEFAULT 0.00,
    total_price DECIMAL(15,2) NOT NULL,
    notes TEXT,
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_archived_transaction_items_transaction ON archived_transaction_items(transaction_id);
CREATE INDEX idx_archived_transaction_items_product_sku_snapshot ON archived_transaction_items(product_sku_snapshot);
CREATE INDEX idx_archived_transaction_items_archived_date ON archived_transaction_items(archived_at);

-- Tabel backup pembayaran
CREATE TABLE archived_transaction_payments (
    id BIGINT PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    payment_method payment_method_single NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    reference_number VARCHAR(100),
    notes TEXT,
    original_created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_archived_transaction_payments_transaction ON archived_transaction_payments(transaction_id);
CREATE INDEX idx_archived_transaction_payments_archived_date ON archived_transaction_payments(archived_at);

-- Types untuk data retention
CREATE TYPE retention_type AS ENUM ('transaction_archive', 'transaction_delete', 'audit_cleanup');
CREATE TYPE retention_status AS ENUM ('success', 'failed', 'partial');

-- Tabel untuk tracking data retention policy
CREATE TABLE data_retention_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    retention_type retention_type NOT NULL,
    records_affected INTEGER NOT NULL,
    date_from DATE NOT NULL,
    date_to DATE NOT NULL,
    execution_time DECIMAL(8,3), -- dalam detik
    status retention_status DEFAULT 'success',
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

CREATE INDEX idx_data_retention_logs_tenant_type_date ON data_retention_logs(tenant_id, retention_type, created_at);

-- =============================================
-- SYSTEM & AUDIT
-- =============================================

-- Tabel untuk audit log
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT,
    action VARCHAR(100) NOT NULL,
    table_name VARCHAR(100),
    record_id BIGINT,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_audit_logs_tenant_date ON audit_logs(tenant_id, created_at);
CREATE INDEX idx_audit_logs_user_date ON audit_logs(user_id, created_at);
CREATE INDEX idx_audit_logs_table_record ON audit_logs(table_name, record_id);

-- =============================================
-- INDEXES TAMBAHAN UNTUK PERFORMA
-- =============================================

-- Index untuk reporting dan analytics
CREATE INDEX idx_sales_tenant_date_status ON transactions(tenant_id, transaction_date, status);
CREATE INDEX idx_sales_outlet_date_total ON transactions(outlet_id, transaction_date, total_amount);
CREATE INDEX idx_customers_tenant_visit ON customers(tenant_id, last_visit_at);
CREATE INDEX idx_products_tenant_active ON products(tenant_id, is_active);

-- =============================================
-- SAMPLE DATA ROLES
-- =============================================

INSERT INTO roles (name, display_name, description, permissions, is_system) VALUES 
('super_admin', 'Super Admin', 'Full system access', '["*"]', TRUE),
('tenant_owner', 'Pemilik Bisnis', 'Full access dalam tenant', '["tenant.*"]', TRUE),
('manager', 'Manager', 'Mengelola outlet dan laporan', '["outlet.*", "reports.*", "products.*", "customers.*"]', TRUE),
('cashier', 'Kasir', 'Melakukan penjualan', '["sales.*", "customers.read", "products.read"]', TRUE);

-- =============================================
-- SAMPLE SUBSCRIPTION PLANS
-- =============================================

INSERT INTO subscription_plans (name, description, price, max_outlets, max_users, features) VALUES 
('Free', 'Paket gratis dengan batasan fitur dan retensi data 14 hari', 0.00, 1, 2, '["basic_pos", "basic_reports", "data_retention_14_days"]'),
('Starter', 'Paket untuk bisnis kecil', 99000.00, 2, 5, '["full_pos", "advanced_reports", "customer_management", "data_retention_unlimited"]'),
('Business', 'Paket untuk bisnis menengah', 299000.00, 5, 15, '["full_pos", "advanced_reports", "customer_management", "inventory_management", "multi_payment", "data_retention_unlimited"]'),
('Enterprise', 'Paket untuk bisnis besar', 599000.00, 999, 999, '["full_pos", "advanced_reports", "customer_management", "inventory_management", "multi_payment", "api_access", "custom_integration", "data_retention_unlimited"]');

-- =============================================
-- VIEWS FOR EASY ACCESS TO TRANSACTION DATA
-- =============================================

-- View untuk melihat semua transaksi (aktif + archived) dengan denormalized data
CREATE VIEW all_transactions_view AS
SELECT 
    id, tenant_id, outlet_id, cashier_id, customer_id, 
    customer_name_snapshot, customer_phone_snapshot, customer_email_snapshot,
    cashier_name_snapshot, outlet_name_snapshot, outlet_code_snapshot,
    transaction_number, transaction_date, subtotal, discount_amount, tax_amount, 
    total_amount, paid_amount, change_amount, payment_method, status, notes,
    original_created_at as created_at, original_updated_at as updated_at,
    'archived' as data_source
FROM archived_transactions
UNION ALL
SELECT 
    id, tenant_id, outlet_id, cashier_id, customer_id,
    customer_name_snapshot, customer_phone_snapshot, customer_email_snapshot,
    cashier_name_snapshot, outlet_name_snapshot, outlet_code_snapshot,
    transaction_number, transaction_date, subtotal, discount_amount, tax_amount,
    total_amount, paid_amount, change_amount, payment_method, status, notes,
    created_at, updated_at,
    'active' as data_source
FROM transactions;

-- =============================================
-- ADDITIONAL INDEXES FOR PERFORMANCE
-- =============================================

-- Index untuk archived tables performance
CREATE INDEX idx_archived_trans_tenant_date ON archived_transactions(tenant_id, transaction_date);
CREATE INDEX idx_archived_items_trans_product ON archived_transaction_items(transaction_id, product_id);