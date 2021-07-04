package apiutil

import (
    "github.com/SevereCloud/vksdk/v2/api"
    "github.com/SevereCloud/vksdk/v2/api/params"
    "strconv"
)

func GetName (session *api.VK, ID int) string {
    b := params.NewUsersGetBuilder()

    i := []string{strconv.Itoa(ID)}

    b.UserIDs(i)
    u, err := session.UsersGet(b.Params)

    if err != nil {
    }

    return u[0].FirstName
}
