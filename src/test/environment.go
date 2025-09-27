package test

import (
	"avito-test/internal/handler"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Environment struct {
	t            *testing.T
	ShopClient   *ShopClient
	ShopDatabase *ShopServiceDatabase
}

func NewEnvironment(t *testing.T) *Environment {
	// TODO: test config
	database_uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		"postgres", "password", "localhost", 5432, "shop",
	)

	env := &Environment{
		t:            t,
		ShopClient:   NewShopClient("http://localhost:8080"),
		ShopDatabase: NewShopClientDatabase(database_uri, "./../../../migrations/init.sql"),
	}

	env.ShopDatabase.Clear()
	env.ShopDatabase.RunMigrations()

	return env
}

func (env *Environment) Close() {
	env.ShopDatabase.Close()
}

type ShopServiceDatabase struct {
	migrations_file string // TODO: real migration tool like goose
	pool            *pgxpool.Pool
}

func NewShopClientDatabase(uri string, migrations_file_path string) *ShopServiceDatabase {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, uri)
	if err != nil {
		log.Fatal(err)
	}

	return &ShopServiceDatabase{
		migrations_file: migrations_file_path,
		pool:            pool,
	}
}

func (d *ShopServiceDatabase) Close() {
	d.pool.Close()
}

func (d *ShopServiceDatabase) RunMigrations() {
	ctx := context.Background()

	sqlBytes, err := os.ReadFile(d.migrations_file)
	if err != nil {
		log.Fatal("Error reading SQL file: ", err)
	}

	_, err = d.pool.Exec(ctx, string(sqlBytes))
	if err != nil {
		log.Fatal("Error executing SQL: ", err)
	}
}

func (d *ShopServiceDatabase) Clear() {
	ctx := context.Background()

	_, err := d.pool.Exec(ctx, `
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
		GRANT ALL ON SCHEMA public TO PUBLIC;
		GRANT ALL ON SCHEMA public TO postgres;
	`)
	if err != nil {
		log.Fatal("Failed to reset database: ", err)
	}
}

func (env *Environment) EnsureAuthorized() {
	code, token := env.ShopClient.Auth("username", "password")

	require.Equal(env.t, 200, code)

	env.ShopClient.SetAuthToken(token.Token)
}

func (env *Environment) EnsureHaveAnotherUser() (string, handler.InfoResponse) {
	const username = "another_username"

	// NOTE: could be register or login, no difference
	code, another_user_token := env.ShopClient.Auth(username, "another_password")
	require.Equal(env.t, 200, code)

	another_client := *env.ShopClient
	another_client.SetAuthToken(another_user_token.Token)

	code, info := another_client.GetInfo()
	require.Equal(env.t, 200, code)

	return username, info
}

func (env *Environment) GetBalance() int {
	env.EnsureAuthorized()

	code, info := env.ShopClient.GetInfo()

	assert.Equal(env.t, 200, code)

	return info.Coins
}
