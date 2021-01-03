package database

const (
	createTableQuery = `CREATE TABLE IF NOT EXISTS products(
	id SERIAL,
	name TEXT NOT NULL,
	price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
	CONSTRAINT products_pkey PRIMARY KEY (id)
)`

	getProductsQuery   = "SELECT id,name,price FROM products LIMIT $1 OFFSET $2"
	getProductQuery    = "SELECT name, price FROM products WHERE id=$1"
	createProductQuery = "INSERT INTO products(name, price) VALUES($1, $2) RETURNING id"
	updateProductQuery = "UPDATE products SET name=$1, price=$2 WHERE id=$3"
	deleteProductQuery = "DELETE FROM products WHERE id=$1"
)
