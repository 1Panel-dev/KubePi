package user

import (
	"fmt"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	customStorm "github.com/KubeOperator/ekko/pkg/storm"
	"github.com/asdine/storm/v3"
	"testing"
)

func TestUserService_Create(t *testing.T) {
	db, err := storm.Open("/var/lib/ekko/db/ekko.db")
	if err != nil {
		t.Error(err)
	}

	//rolebinding := v1Role.Banding{
	//	BaseModel: v1.BaseModel{
	//		ApiVersion: "v1",
	//		Kind:       "RoleBinding",
	//		CreateAt:   time.Now(),
	//		UpdateAt:   time.Now(),
	//		CreatedBy:  "system",
	//	},
	//	Metadata: v1.Metadata{
	//		UUID: uuid.New().String(),
	//		Name: "test",
	//	},
	//	Subjects: []v1Role.Subject{
	//		{
	//			Kind: "User",
	//			Name: "zhangsan",
	//		},
	//		{
	//			Kind: "Group",
	//			Name: "group1",
	//		},
	//	},
	//	RoleRef: "administrator",
	//}
	//if err := db.Save(&rolebinding); err != nil {
	//	t.Error(err)
	//}
	var u []v1Role.Banding

	query := db.Select(customStorm.Containers("Subjects", v1Role.Subject{Name: "lisi", Kind: "User"}))

	if err := query.Find(&u); err != nil {
		t.Error(err)
	}
	//if err := db.Find("Subjects", []v1Role.Subject{{Kind: "User", Name: "zhangsan"}, {Kind: "Group", Name: "group1"}}, &u); err != nil {
	//	t.Error(err)
	//}
	fmt.Println(u)

}
