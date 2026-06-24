// Package mongo provee una implementación de domain.ProductRepository sobre
// MongoDB. Se activa con DB_DRIVER=mongo.
package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/francososa97/product-api/internal/domain"
)

// ProductRepository persiste productos en una colección de MongoDB.
type ProductRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewProductRepository conecta con MongoDB, verifica la conexión con un ping y
// devuelve un repositorio listo para usar. El context controla el timeout de
// conexión.
func NewProductRepository(ctx context.Context, uri, dbName, collectionName string) (*ProductRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar a MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("no se pudo hacer ping a MongoDB: %w", err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &ProductRepository{client: client, collection: collection}, nil
}

// Close cierra la conexión con MongoDB. Debe llamarse al apagar la aplicación.
func (r *ProductRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}

func (r *ProductRepository) GetAll(ctx context.Context, sortByPriceAsc bool) ([]domain.Product, error) {
	order := -1 // descendente por defecto
	if sortByPriceAsc {
		order = 1
	}

	opts := options.Find().SetSort(bson.D{{Key: "price", Value: order}})
	cursor, err := r.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	products := make([]domain.Product, 0)
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var product domain.Product
	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	filter := bson.D{{Key: "_id", Value: product.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: product.Name},
		{Key: "price", Value: product.Price},
	}}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	res, err := r.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}
