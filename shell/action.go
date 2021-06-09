package shell

import (
	"bytes"
	"errors"
	"fmt"
	"m2np_cli/model"
	"os"
	"strings"
	"time"
)

type action func(s *shell, inputs []byte) (string, error)

func distribute(s *shell, inputs []byte, actions ...map[string]action) (string, error){
	args:= strings.Split(string(inputs), " ")
	for _, as := range actions{
		if action, ok := as[args[0]]; ok {
			res, err := action(s, inputs)
			return res, err
		}
	}
	keysm := map[string]bool{}
	for _, as := range actions{
		for k := range as {
			keysm[k] = true
		}
	}
	keys := make([]string, 0)
	for k, _ := range keysm{
		keys = append(keys, k)
	}
	return fmt.Sprintf("unknow command %s, reference %v\n", args[0], keys), nil
}

var baseActions = map[string]action{
	"cd": func(s *shell, inputs []byte) (string, error) {
		args := strings.Split(string(inputs), " ")
		if len(args) != 2 {
			return "", errors.New("cd $dirname")
		}

		dirs := strings.Split(args[1], "/")
		for _, dir := range dirs {
			if dir == ".." {
				if s.position.parent != nil {
					s.position = s.position.parent
				}
			}else if dir == "."{
				// no things to do!
			}else{
				for _, v := range s.position.child {
					if v.name == dir {
						s.position = v
					}
				}
			}
		}
		return "", nil
	},
	"ls": func(s *shell, inputs []byte) (string, error) {
		buffer := bytes.Buffer{}
		for _, v := range s.position.child {
			buffer.WriteString(v.name + "\n")
		}
		return buffer.String(), nil
	},
	"os": func(s *shell, inputs []byte) (string, error) {
		if s.ctx.User != "" {
			return fmt.Sprintf("username: %s\n", s.ctx.User), nil
		}
		return "no login.\n", nil
	},
	"exit": func(s *shell, inputs []byte) (string, error) {
		os.Exit(0)
		return "", nil
	},
}

var rootActions = map[string]action{
	"login": func(s *shell, inputs []byte) (string, error) {
		args := strings.Split(string(inputs), " ")
		if len(args) != 3 {
			return "", errors.New("login $user $pwd")
		}
		err := s.ctx.Login(args[1], args[2])
		if err != nil{
			return "", errors.New("login fail")
		}
		return "login successful!\n", nil
	},
}

var postsActions = map[string]action{
	"ls": func(s *shell, inputs []byte) (string, error) {
		inbox, err := s.ctx.GetInBox()
		if err != nil{
			return "", errors.New(fmt.Sprintf("get inbox fail: %v", err))
		}
		users := map[int64]model.User{}
		for _, user := range inbox.Users {
			users[user.ID] = user
		}
		buffer := bytes.Buffer{}
		for _, v := range inbox.Posts {
			buffer.WriteString(fmt.Sprintf("@%-10s %s %s\n", users[v.UserID].Username, time.Unix(v.CreatedAt, 0).Format("Jan 02 15:04"), v.Content))
		}
		return buffer.String(), nil
	},
	"touch": func(s *shell, inputs []byte) (string, error) {
		args := strings.Split(string(inputs), " ")
		if len(args) != 2 {
			return "", errors.New("touch $content")
		}
		err := s.ctx.PostPost(args[1])
		if err != nil{
			return "", errors.New(fmt.Sprintf("post fail: %v", err))
		}
		return "", nil
	},
}

var followersActions = map[string]action{
	"ls": func(s *shell, inputs []byte) (string, error) {
		return "output all followers of login user.\n", nil
	},
}