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

	sess.Set("user_id", user.User_id)
	sess.Set("username", user.Username)
	sess.Set("email", user.Email)

	err = sess.Save()

	if err != nil {
		return err
	}
	return nil
}

func GetUser(c *fiber.Ctx, s *session.Store) fiber.Map {
	sess, err := s.Get(c)
	if err != nil {
		return fiber.Map{}
	}

	return fiber.Map{
		"User_id":  sess.Get("user_id"),
		"Username": sess.Get("username"),
		"Email":    sess.Get("email"),
	}
}

func HasSess(c *fiber.Ctx, s *session.Store) bool {
	mp := GetUser(c, s)
	if len(mp) == 0 {
		return false
	}

	if mp == nil {
		return false
	}

	return mp["User_id"] != nil
}
