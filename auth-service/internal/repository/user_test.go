package repository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestUserRepo_Integration(t *testing.T) {
	ctx := context.Background()

	// 1. Запускаем контейнер с PostgreSQL (временный)
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %s", err)
	}
	// Гарантируем, что контейнер остановится после теста
	defer postgres.Terminate(ctx)

	// 2. Получаем DSN для подключения к тестовой базе
	host, err := postgres.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get host: %s", err)
	}
	port, err := postgres.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("Failed to get port: %s", err)
	}

	dsn := "postgres://postgres:postgres@" + host + ":" + port.Port() + "/test?sslmode=disable"

	// 3. Подключаемся
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %s", err)
	}
	defer pool.Close()

	// Создаем таблицу для теста (в реальном проекте используем миграции)
	_, err = pool.Exec(ctx, `CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %s", err)
	}

	// 4. Инициализируем репозиторий
	repo := NewUserRepo(pool)

	// 5. Проверяем логику: Создаем пользователя
	email := "integration-test@example.com"
	passwordHash := "hashed_password_123"

	user, err := repo.Create(ctx, email, passwordHash)
	if err != nil {
		t.Fatalf("Create failed: %s", err)
	}

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}

	// 6. Проверяем логику: Получаем пользователя
	fetchedUser, err := repo.GetByEmail(ctx, email)
	if err != nil {
		t.Fatalf("GetByEmail failed: %s", err)
	}

	if fetchedUser == nil {
		t.Fatal("User not found")
	}

	if fetchedUser.Email != email {
		t.Errorf("Expected email %s, got %s", email, fetchedUser.Email)
	}

	t.Logf("✅ Integration test passed! User ID: %d", user.ID)
}
