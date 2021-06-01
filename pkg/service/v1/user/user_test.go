package user

import (
	"fmt"
	"github.com/KubeOperator/ekko/pkg/model/v1/user"
	"github.com/asdine/storm/v3"
	"testing"
)

func TestUserService_Create(t *testing.T) {
	db, err := storm.Open("/var/lib/ekko/db/ekko.db")
	if err != nil {
		t.Error(err)
	}
	//u := user.User{
	//	BaseModel: v1.BaseModel{
	//		ApiVersion: "v1",
	//		Kind:       "User",
	//		CreateAt:   time.Now(),
	//		UpdateAt:   time.Now(),
	//	},
	//	Metadata: v1.Metadata{
	//		Name: "test",
	//		UUID: uuid.New().String(),
	//	},
	//	Spec: user.Spec{
	//		Info: user.Info{
	//			NickName: "iamtest",
	//			Email:    "chenyang@fit2cloud.com",
	//		},
	//		Authenticate: user.Authenticate{
	//			Password: "Calong@2015",
	//			Token:    "",
	//		},
	//	},
	//}
	//if err := db.Save(&u); err != nil {
	//	t.Error(err)
	//}
	u := user.User{}
	if err := db.One("Name", "test", &u); err != nil {
		t.Error(err)
	}
	fmt.Println(u)

}
