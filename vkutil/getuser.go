package vkutil

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"strconv"
)

/*
Получает инфу об ОДНОМ пользователе.
Если тебе нужно получить инфу о нескольких, то воспользуйся обычным интерфейсом из api
*/

func GetUser(session *api.VK, ID int) (*object.UsersUser, error) {
	b := params.NewUsersGetBuilder()

	i := []string{strconv.Itoa(ID)}

	b.UserIDs(i)
	b.Lang(0)
	u, err := session.UsersGet(b.Params)

	if err != nil {
		return nil, err
	}

	return &u[0], nil
}
