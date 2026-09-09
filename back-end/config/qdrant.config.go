package config

import (
	"context"
	"os"
	"strings"

	"github.com/Nixon-Alexander/resolve_now.git/model"
	"github.com/qdrant/go-client/qdrant"
)

type QDrant struct {
	client *qdrant.Client
}

func isAlreadyExistsErr(err error) bool {
	return strings.Contains(err.Error(), "already exists")
}

func InitQDrant() (*QDrant, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: os.Getenv("QDRANT_HOST"),
		Port: 6334,
	})

	if err != nil {
		return nil, err
	}

	return &QDrant{
		client: client,
	}, nil
}

func (qd *QDrant) CreateChunksCollection() error {
	err := qd.client.CreateCollection(
		context.Background(),
		&qdrant.CreateCollection{
			CollectionName: "document_chunks",
			VectorsConfig: qdrant.NewVectorsConfig(
				&qdrant.VectorParams{
					Size:     768,
					Distance: qdrant.Distance_Cosine,
				},
			),
		},
	)

	// Kalau collection sudah ada, abaikan error "already exists" dan lanjut bikin index
	if err != nil && !isAlreadyExistsErr(err) {
		return err
	}

	// Buat index untuk field yang dipakai filtering (company_id)
	_, err = qd.client.CreateFieldIndex(
		context.Background(),
		&qdrant.CreateFieldIndexCollection{
			CollectionName: "document_chunks",
			FieldName:      "company_id",
			FieldType:      qdrant.FieldType_FieldTypeInteger.Enum(),
		},
	)

	if err != nil && !isAlreadyExistsErr(err) {
		return err
	}

	return nil
}

func (qd *QDrant) UpsertChunk(
	ctx context.Context,
	id uint64,
	embedding []float32,
	payload map[string]any,
) error {
	point := &qdrant.PointStruct{
		Id:      qdrant.NewIDNum(id),
		Vectors: qdrant.NewVectors(embedding...),
		Payload: qdrant.NewValueMap(payload),
	}

	_, err := qd.client.Upsert(
		ctx,
		&qdrant.UpsertPoints{
			CollectionName: "document_chunks",
			Points:         []*qdrant.PointStruct{point},
		},
	)

	return err
}

func (qd *QDrant) SearchChunks(
	ctx context.Context,
	embedding []float32,
	companyID int64,
	limit uint64,
) ([]model.ChunkSearchResult, error) {
	points, err := qd.client.Query(
		ctx,
		&qdrant.QueryPoints{
			CollectionName: "document_chunks",
			Query:          qdrant.NewQuery(embedding...),
			Filter: &qdrant.Filter{
				Must: []*qdrant.Condition{
					qdrant.NewMatchInt("company_id", companyID),
				},
			},
			Limit:       &limit,
			WithPayload: qdrant.NewWithPayload(true),
		},
	)

	if err != nil {
		return nil, err
	}

	results := make([]model.ChunkSearchResult, 0, len(points))

	for _, p := range points {
		payload := p.GetPayload()

		results = append(results, model.ChunkSearchResult{
			ChunkID:           p.GetId().GetNum(),
			DocumentID:        payload["document_id"].GetIntegerValue(),
			DocumentVersionID: payload["document_version_id"].GetIntegerValue(),
			ChunkIndex:        payload["chunk_index"].GetIntegerValue(),
			Content:           payload["content"].GetStringValue(),
			Score:             p.GetScore(),
		})
	}

	return results, nil
}
