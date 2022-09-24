package users

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uuid        string   `toml:"uuid"`
	Username    string   `toml:"username"`
	Password    string   `toml:"password"`
	Permissions []string `toml:"permissions"`
}

var Permissions = []string{
	"admin",
	"project:create",
	"project:read",
	"assets:image:read",
	"assets:image:write",
	"assets:model:read",
	"assets:model:read-licensed",
	"assets:model:write",
	"assets:file:read",
	"assets:file:write",
	"assets:slice:read",
	"assets:slice:write",
}

var users map[string]*User
var userFile string

func Register(protected *echo.Group, public *echo.Group) {
	userFile = fmt.Sprintf("%s/users/users.toml", runtime.Cfg.SystemPath)
	initFs()
	initUsers()
	protected.GET("/users", nil)
	public.POST("/login", login)

}

func initFs() {
	err := os.MkdirAll(path.Dir(userFile), os.ModePerm)

	if err != nil {
		log.Println(err)
	}

	_, err = os.Stat(userFile)
	if os.IsNotExist(err) {
		file, err := os.Create(userFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

}

func initUsers() {

	var data struct {
		Users map[string]*User `toml:"users"`
	}
	_, err := toml.DecodeFile(userFile, &data)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(data)

	_ = toml.NewEncoder(os.Stdout).Encode(data)

	changed := false
	users = make(map[string]*User)
	for _, user := range data.Users {
		if c, err := bcrypt.Cost([]byte(user.Password)); err != nil || c == 0 {
			bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
			if err != nil {
				log.Println(err)
			}
			user.Password = string(bytes)
			changed = true
		}
		if user.Uuid == "" {
			user.Uuid = uuid.New().String()
			changed = true
		}
		users[user.Username] = user
	}

	if changed {
		f, err := os.OpenFile(userFile, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		err = toml.NewEncoder(f).Encode(data)

		if err != nil {
			log.Println(err)
		}
	}
}
