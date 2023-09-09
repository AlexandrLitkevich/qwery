package store

import (
	"context"
	"fmt"
	"github.com/AlexandrLitkevich/qwery/internal/config"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/user"
	"github.com/AlexandrLitkevich/qwery/internal/lib/logger/sl"
	"github.com/AlexandrLitkevich/qwery/internal/storage"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log/slog"
)

type EtcdStore struct {
	Cli *clientv3.Client
}

func New(ctx context.Context, log *slog.Logger, cfg *config.Config) (*EtcdStore, error) {
	const op = "storage.etcd.New"

	cl3, err := clientv3.New(clientv3.Config{
		Context:     ctx,
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
	})
	if err != nil {
		log.Error("fail to connection etcd in func:", op, sl.Err(err))
		return nil, fmt.Errorf("fail to connection")
	}
	log.Info("created connection etcd")

	return &EtcdStore{Cli: cl3}, nil
}

func (s *EtcdStore) GetUser(ctx context.Context, log *slog.Logger, userId string) (string, error) {
	const op = "storage.etcd.GetUser"

	log.Info("GetUser")
	firstUserPath := storage.GetUserPath(userId)
	log.Info("run get user")
	kv := clientv3.NewKV(s.Cli)
	response, err := kv.Get(ctx, firstUserPath)
	if err != nil {
		log.Error("fail to get user(Store) in func:", op, sl.Err(err))
		return "", err
	}
	fmt.Println("this responce", response)
	return "reesee", nil

}

func (s *EtcdStore) Create(ctx context.Context, cancel context.CancelFunc, log *slog.Logger, request user.Request) (bool, error) {
	const op = "storage.etcd.Create"

	kv := clientv3.NewKV(s.Cli)

	userId, err := uuid.NewUUID()
	if err != nil {
		//Тут так возвращаем ошибку так как ее в логи выводим в хандлере
		return false, fmt.Errorf("%s: %w", op, err)
	}
	us := userId.String()

	log = log.With(
		slog.String("op", op),
		slog.String("geerate user id", us),
		slog.String("request", request.Position),
		slog.String("request", request.Name),
	)

	firstUserPath := storage.GetUserPath(us)
	_, err = kv.Put(ctx, firstUserPath, "i'm user")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	cancel()
	return true, nil
}
