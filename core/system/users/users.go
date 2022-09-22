package users

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Role     string `toml:"role"`
}

var data struct {
	Users map[string]*User `toml:"users"`
}
var userFile string

func Register(g echo.Group) {
	userFile = fmt.Sprintf("%s/users/users.toml", runtime.Cfg.SystemPath)
	initFs()
	initUsers()
	g.GET("/users", nil)

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

	_, err = toml.DecodeFile(userFile, &data)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(data)

	_ = toml.NewEncoder(os.Stdout).Encode(data)

}

func initUsers() {

	changed := false

	for _, user := range data.Users {
		if c, err := bcrypt.Cost([]byte(user.Password)); err != nil || c == 0 {
			bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
			if err != nil {
				log.Println(err)
			}
			user.Password = string(bytes)
			changed = true
		}
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
