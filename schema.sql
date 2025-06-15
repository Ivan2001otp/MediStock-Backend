-- Create the medistock_db database if it doesn't exist
CREATE DATABASE IF NOT EXISTS medistock_db;

-- Use the medistock_db database
USE medistock_db;

CREATE TABLE IF NOT EXISTS supplies (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL, 
    sku VARCHAR(128) UNIQUE NOT NULL,
    current_stock INT NOT NULL DEFAULT 0,
    unit_of_measure  VARCHAR(32) NOT NULL,
    is_vital BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS vendors (
    id  INT PRIMARY KEY AUTO_INCREMENT ,
    name VARCHAR(128) NOT NULL, 
    contact_person VARCHAR(255) ,
    phone VARCHAR(16),
    email VARCHAR(64),
    address TEXT,
    overall_quality_rating DECIMAL(3,2),
    avg_delivery_time DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS vendor_supply_prices(
    id INT PRIMARY KEY AUTO_INCREMENT,
    vendor_id INT ,
    supply_id INT ,
    unit_price DECIMAL(10,2) NOT NULL,
    quality_rating DECIMAL(3,2) NOT NULL,
    estimated_delivery_days DECIMAL(5,2),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (supply_id) REFERENCES supplies(id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    supply_id INT,
    vendor_id INT,
    quantity_ordered  INT NOT NULL ,
    order_price DECIMAL(10,2),
    order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ,
    expected_delivery_date TIMESTAMP,
    status VARCHAR(16) NOT NULL, 
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (supply_id) REFERENCES supplies(id) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON DELETE SET NULL ON UPDATE CASCADE
)
