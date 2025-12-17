// package config

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	_ "github.com/lib/pq" // Driver PostgreSQL
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// var (
// 	// DB adalah koneksi untuk PostgreSQL (Relational Data)
// 	DB *sql.DB

// 	// MongoDB adalah koneksi untuk MongoDB (Document Data)
// 	MongoDB *mongo.Database
// )

// // ConnectDB menginisialisasi semua koneksi database
// func ConnectDB() {
// 	connectPostgres()
// 	connectMongo()
// }

// func connectPostgres() {
// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
// 		os.Getenv("DB_HOST"),
// 		os.Getenv("DB_PORT"),
// 		os.Getenv("DB_USER"),
// 		os.Getenv("DB_PASSWORD"),
// 		os.Getenv("DB_NAME"),
// 		os.Getenv("DB_SSLMODE"),
// 	)

// 	var err error
// 	DB, err = sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Fatalf("❌ Gagal membuka driver PostgreSQL: %v", err)
// 	}

// 	if err = DB.Ping(); err != nil {
// 		log.Fatalf("❌ Gagal ping PostgreSQL: %v", err)
// 	}

// 	fmt.Println("✅ Berhasil terhubung ke PostgreSQL!")
// }

// func connectMongo() {
// 	// Timeout 10 detik untuk koneksi
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	mongoURI := os.Getenv("MONGO_URI")
// 	if mongoURI == "" {
// 		log.Fatal("❌ MONGO_URI tidak ditemukan di .env")
// 	}

// 	clientOptions := options.Client().ApplyURI(mongoURI)
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatalf("❌ Gagal koneksi awal ke MongoDB: %v", err)
// 	}

// 	// Verifikasi koneksi dengan Ping
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatalf("❌ Gagal ping MongoDB: %v", err)
// 	}

// 	dbName := os.Getenv("MONGO_DB_NAME")
// 	if dbName == "" {
// 		dbName = "prestasi_mhs_mongo" // Default fallback
// 	}

// 	MongoDB = client.Database(dbName)
// 	fmt.Println("✅ Berhasil terhubung ke MongoDB!")
// }


package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Driver PostgreSQL
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// DB adalah koneksi untuk PostgreSQL (Relational Data)
	DB *sql.DB

	// MongoDB adalah koneksi untuk MongoDB (Document Data)
	MongoDB *mongo.Database
)

// ConnectDB menginisialisasi koneksi PostgreSQL
// Dipanggil di main.go
func ConnectDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Gagal membuka driver PostgreSQL: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ Gagal ping PostgreSQL: %v", err)
	}

	fmt.Println("✅ Berhasil terhubung ke PostgreSQL!")
}

// ConnectMongo menginisialisasi koneksi MongoDB
// Dipanggil di main.go secara terpisah (Exported Function)
func ConnectMongo() {
	// Timeout 10 detik untuk koneksi
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		// Fallback default jika di .env kosong (untuk safety dev local)
		mongoURI = "mongodb://localhost:27017"
		log.Println("⚠️ MONGO_URI tidak ditemukan di .env, menggunakan default localhost.")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi awal ke MongoDB: %v", err)
	}

	// Verifikasi koneksi dengan Ping
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Gagal ping MongoDB: %v", err)
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "prestasi_mhs_mongo" // Default fallback
	}

	MongoDB = client.Database(dbName)
	fmt.Println("✅ Berhasil terhubung ke MongoDB!")
}