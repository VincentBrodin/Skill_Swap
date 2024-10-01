package sesstools

import (
	"example.com/SkillSwap/tools/dbtools"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AddUser(c *fiber.Ctx, s *session.Store, user *dbtools.User) error {
	//Add user to session
	sess, err := s.Get(c)

	if err != nil {
		return err
	}

	sess.Set("User_id", user.User_id)

	err = sess.Save()

	if err != nil {
		return err
	}
	return nil
}

// Returns -1 if no user was found and the user id otherwise
func GetUser(c *fiber.Ctx, s *session.Store) int64 {
	sess, err := s.Get(c)
	if err != nil {
		return -1
	}
	user_interface := sess.Get("User_id")
	if user_interface == nil {
		return -1
	}
	user_id, ok := user_interface.(int64)
	if !ok {
		return -1
	}
	return user_id
}

func HasSess(c *fiber.Ctx, s *session.Store) bool {
	user_id := GetUser(c, s)
	return user_id != -1
}
