CREATE INDEX idx_email_actor ON medistock_db.auth_token(email,actor);
CREATE INDEX idx_vendor ON medistock_db.vendor_supply_prices(vendor_id);

CREATE TABLE IF NOT EXISTS hospitals (
    id VARCHAR(36) PRIMARY KEY, -- Using UUID for unique Hospital ID
    name VARCHAR(255) NOT NULL UNIQUE,
    address TEXT,
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS supplies (
    id VARCHAR(36) PRIMARY KEY, -- Using UUID for unique supply ID
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(100) UNIQUE,
    unit_of_measure VARCHAR(50) NOT NULL,
    category VARCHAR(100),
    is_vital BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS vendors (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE,
    contact_person VARCHAR(255),
    phone VARCHAR(50),
    email VARCHAR(255),
    address TEXT,
    overall_quality_rating DECIMAL(3,2),
    avg_delivery_time_days DECIMAL(5,2),
    score DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS vendor_supply_prices (
    vendor_id INT NOT NULL,
    supply_id VARCHAR(36) NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    quality_rating DECIMAL(3, 2),
    avg_delivery_days DECIMAL(5, 2),
    PRIMARY KEY (vendor_id, supply_id),
    FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON DELETE CASCADE on UPDATE CASCADE,
    FOREIGN KEY (supply_id) REFERENCES supplies(id) ON DELETE CASCADE on UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(36) PRIMARY KEY, -- Using UUID for unique order ID
    hospital_id VARCHAR(36) NOT NULL, -- NEW: Which hospital placed this order
    vendor_id INT NOT NULL,           -- Which vendor this order was placed with
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'Pending',
    total_amount DECIMAL(10, 2) NOT NULL,
    estimated_delivery_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (hospital_id) REFERENCES hospitals(id) on DELETE CASCADE, -- Link to hospitals table
    FOREIGN KEY (vendor_id) REFERENCES vendors(id) on DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS order_items (
    order_id VARCHAR(36) NOT NULL,
    supply_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL,
    unit_price_at_order DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (order_id, supply_id),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (supply_id) REFERENCES supplies(id) ON DELETE RESTRICT
);


CREATE TABLE IF NOT EXISTS inventory (
    hospital_id VARCHAR(36) NOT NULL, -- NEW: Which hospital this inventory record belongs to
    supply_id VARCHAR(36) NOT NULL,   -- Which supply item this inventory record is for
    current_stock INT NOT NULL DEFAULT 0,
    reorder_threshold INT NOT NULL DEFAULT 0,
    PRIMARY KEY (hospital_id, supply_id), -- Composite primary key: ensures each hospital has one record per supply
    FOREIGN KEY (hospital_id) REFERENCES hospitals(id) ON DELETE CASCADE on UPDATE CASCADE, -- If a hospital is deleted, its inventory is too
    FOREIGN KEY (supply_id) REFERENCES supplies(id) ON DELETE CASCADE on UPDATE CASCADE-- If a supply is deleted, its inventory record is too.
);