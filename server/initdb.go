package main

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func InitDB() error {
	return DB.BeginFunc(
		context.Background(),
		func(t pgx.Tx) error {
			_, err := t.Exec(
				context.Background(),
				`CREATE TABLE IF NOT EXISTS totp_tokens (
					id BIGSERIAL PRIMARY KEY,
					secret TEXT NOT NULL,
					created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
				)`,
			)
			if err != nil {
				return err
			}

			_, err = t.Exec(
				context.Background(),
				`CREATE TABLE IF NOT EXISTS auth_methods (
					id BIGSERIAL PRIMARY KEY,
					user_id BIGINT NOT NULL,
					type SMALLINT NOT NULL,
					data_id BIGSERIAL NOT NULL,
					created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
				)`,
			)
			if err != nil {
				return err
			}
			return nil
		},
	)

}
