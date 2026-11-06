package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Nama collection MongoDB sesuai dokument database-design-platform-produk-digital.md
const (
	ColProjectDocuments = "project_documents"
	ColInvitationGuests = "invitation_guests"
	ColRSVPResponses    = "rsvp_responses"
	ColProjectRevisions = "project_revisions"
)

// EnsureIndexes membuat indeks MongoDB yang dibutuhkan.
// Idempoten: aman dipanggil berulang karena indeks sudah ada tidak dibuat lagi.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	for _, params := range indexDefinitions() {
		coll := db.Collection(params.collection)
		if _, err := coll.Indexes().CreateMany(ctx, params.models); err != nil {
			return &IndexError{Collection: params.collection, Err: err}
		}
	}
	return nil
}

// IndexError error pembuatan index MongoDB
type IndexError struct {
	Collection string
	Err        error
}

func (e *IndexError) Error() string {
	return "mongodb index gagal dibuat pada collection " + e.Collection + ": " + e.Err.Error()
}

func (e *IndexError) Unwrap() error { return e.Err }

type collectionIndex struct {
	collection string
	models     []mongo.IndexModel
}

// indexDefinitions daftar indeks sesuai section 4 dokument rancangan:
// - project_documents: UNIQUE (projectId, productType) karena setiap proyek
//   punya satu dokumen konten aktif per jenis produk
// - invitation_guests: UNIQUE (projectId, slug) agar slug tamu unik per proyek
// - rsvp_responses: (projectId, guestId) untuk pencarian respons tamu
// - project_revisions: (projectId, version) untuk riwayat konten
func indexDefinitions() []collectionIndex {
	return []collectionIndex{
		{
			collection: ColProjectDocuments,
			models: []mongo.IndexModel{
				{
					Keys: bson.D{
						{Key: "projectId", Value: 1},
						{Key: "productType", Value: 1},
					},
					Options: options.Index().
						SetName("uniq_projectId_productType").
						SetUnique(true),
				},
				{
					Keys:    bson.D{{Key: "projectId", Value: 1}},
					Options: options.Index().SetName("idx_projectId"),
				},
			},
		},
		{
			collection: ColInvitationGuests,
			models: []mongo.IndexModel{
				{
					Keys: bson.D{
						{Key: "projectId", Value: 1},
						{Key: "slug", Value: 1},
					},
					Options: options.Index().
						SetName("uniq_projectId_slug").
						SetUnique(true),
				},
				{
					Keys:    bson.D{{Key: "projectId", Value: 1}},
					Options: options.Index().SetName("idx_projectId"),
				},
			},
		},
		{
			collection: ColRSVPResponses,
			models: []mongo.IndexModel{
				{
					Keys: bson.D{
						{Key: "projectId", Value: 1},
						{Key: "guestId", Value: 1},
					},
					Options: options.Index().SetName("idx_projectId_guestId"),
				},
				{
					Keys:    bson.D{{Key: "guestId", Value: 1}},
					Options: options.Index().SetName("idx_guestId"),
				},
			},
		},
		{
			collection: ColProjectRevisions,
			models: []mongo.IndexModel{
				{
					Keys: bson.D{
						{Key: "projectId", Value: 1},
						{Key: "version", Value: -1},
					},
					Options: options.Index().SetName("idx_projectId_version"),
				},
			},
		},
	}
}
