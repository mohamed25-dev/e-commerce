CREATE TABLE IF NOT EXISTS analytics (
    id VARCHAR(50) PRIMARY KEY,
    top_customers INT NOT NULL,
    total_sales DECIMAL(9, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(50) PRIMARY KEY,
    customer_id VARCHAR(50) NOT NULL,
    customer_name VARCHAR(50) NOT NULL,
    product_id VARCHAR(50) NOT NULL,
    product_name VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    total_price DECIMAL NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
