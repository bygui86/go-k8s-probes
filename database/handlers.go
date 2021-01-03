package database

import (
	"context"
	"database/sql"

	"github.com/opentracing/opentracing-go"
)

func GetProducts(db *sql.DB, start, count int, ctx context.Context) ([]*Product, error) {
	span := opentracing.StartSpan(
		"get-products-db",
		opentracing.ChildOf(opentracing.SpanFromContext(ctx).Context()))
	defer span.Finish()

	span.SetTag("query", getProductsQuery)
	span.SetTag("count", count)
	span.SetTag("start", start)
	span.LogKV(
		"query", getProductsQuery,
		"count", count,
		"start", start,
	)

	rows, err := db.QueryContext(ctx, getProductsQuery, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*Product, 0)
	for rows.Next() {
		var prod Product
		if err := rows.Scan(&prod.ID, &prod.Name, &prod.Price); err != nil {
			return nil, err
		}
		products = append(products, &prod)
	}

	span.SetTag("products-found", len(products))
	span.LogKV("products-found", len(products))

	return products, nil
}

func GetProduct(db *sql.DB, product *Product, ctx context.Context) error {
	span := opentracing.StartSpan(
		"get-product-db",
		opentracing.ChildOf(opentracing.SpanFromContext(ctx).Context()))
	defer span.Finish()

	span.SetTag("product-id", product.ID)
	span.LogKV("product-id", product.ID)

	return db.QueryRowContext(ctx, getProductQuery, product.ID).
		Scan(&product.Name, &product.Price)
}

func CreateProduct(db *sql.DB, product *Product, ctx context.Context) error {
	span := opentracing.StartSpan(
		"create-product-db",
		opentracing.ChildOf(opentracing.SpanFromContext(ctx).Context()))
	defer span.Finish()

	span.SetTag("product", product.String())
	span.LogKV("product", product.String())

	err := db.QueryRowContext(ctx, createProductQuery, product.Name, product.Price).Scan(&product.ID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProduct(db *sql.DB, product *Product, ctx context.Context) error {
	span := opentracing.StartSpan(
		"update-product-db",
		opentracing.ChildOf(opentracing.SpanFromContext(ctx).Context()))
	defer span.Finish()

	span.SetTag("product", product.String())
	span.LogKV("product", product.String())

	_, err := db.ExecContext(ctx, updateProductQuery, product.Name, product.Price, product.ID)
	return err
}

func DeleteProduct(db *sql.DB, productId int, ctx context.Context) error {
	span := opentracing.StartSpan(
		"delete-product-db",
		opentracing.ChildOf(opentracing.SpanFromContext(ctx).Context()))
	defer span.Finish()

	span.SetTag("product-id", productId)
	span.LogKV("product-id", productId)

	_, err := db.ExecContext(ctx, deleteProductQuery, productId)
	return err
}
