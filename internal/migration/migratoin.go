package migration

import (
	"encoding/json"
	"fmt"
	"github.com/AlexandrLitkevich/qwery/internal/http-server/handlers/user"
	"github.com/AlexandrLitkevich/qwery/internal/storage"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func KeyMissing(key string) clientv3.Cmp {
	return clientv3.Compare(clientv3.Version(key), "=", 0)
}

func UpgradeTo0001() ([]clientv3.Cmp, []clientv3.Op, error) {

	defaultAdmin := user.User{
		ID:       "1000",
		Name:     "Default Admin",
		Age:      33,
		Position: "clerc",
	}

	userPath := storage.GetUserPath(defaultAdmin.ID)

	jsonUser, err := json.Marshal(defaultAdmin)
	if err != nil {
		fmt.Println("err json marshal")
	}

	cmps := []clientv3.Cmp{
		KeyMissing(userPath),
	}

	ops := []clientv3.Op{
		clientv3.OpPut(userPath, string(jsonUser)),
	}

	return cmps, ops, nil

}
