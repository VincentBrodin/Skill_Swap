package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"example.com/SkillSwap/tools/dbtools"
	"example.com/SkillSwap/tools/pswdhash"
	"example.com/SkillSwap/tools/sesstools"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sujit-baniya/flash"
)

func main() {
	// Load db
	db, err := sql.Open("sqlite3", "data.db")

	if err != nil {
		fmt.Println("db error", err.Error())
		os.Exit(1)
	}

	// Start new session
	store := session.New()

	// Start app
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Static
	app.Static("/js/", "src")

	// Routes

	//====Index====
	app.Get("/", func(c *fiber.Ctx) error {
		data := flash.Get(c)

		return c.Render("index", fiber.Map{
			"Title": "Landing",
			"Flash": data,
			"User":  sesstools.GetUser(c, store),
		}, "layouts/main")
	})

	//====Home====
	app.Get("/home", func(c *fiber.Ctx) error {
		data := flash.Get(c)

		return c.Render("home", fiber.Map{
			"Title": "Home",
			"Flash": data,
			"User":  sesstools.GetUser(c, store),
		}, "layouts/main")
	})

	//====Login====
	app.Get("/login", func(c *fiber.Ctx) error {
		// Redirect user if already logged in
		if sesstools.HasSess(c, store) {
			return c.Redirect("/")
		}

		data := flash.Get(c)
		return c.Render("login", fiber.Map{
			"Title": "Login",
			"Flash": data,
			"User":  sesstools.GetUser(c, store),
		}, "layouts/main")
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		// Redirect user if already logged in
		if sesstools.HasSess(c, store) {
			return c.Redirect("/")
		}

		email := c.FormValue("email")
		password := c.FormValue("password")

		user, err := dbtools.GetUserFromEmail(email, db)

		if err != nil {
			mp := fiber.Map{
				"message": "Incorrect email or password!",
			}
			fmt.Println(err.Error())
			return flash.WithError(c, mp).Redirect("/login")
		}

		if !pswdhash.VerifyPassword(password, user.Password) {
			mp := fiber.Map{
				"message": "Incorrect email or password!",
			}
			fmt.Println("Wrong password")
			return flash.WithError(c, mp).Redirect("/login")
		}

		// Add user to session
		err = sesstools.AddUser(c, store, user)
		if err != nil {

			mp := fiber.Map{
				"message": "Could not get session!",
			}
			fmt.Println(err.Error())
			return flash.WithError(c, mp).Redirect("/register")

		}

		mp := fiber.Map{
			"message": "Logged in :)!",
		}
		return flash.WithSuccess(c, mp).Redirect("/")
	})

	//====Logout====
	app.Get("/logout", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)

		if err != nil {
			return err
		}
		sess.Reset()

		err = sess.Save()

		if err != nil {
			return c.Redirect("/")
		}

		return c.Redirect("/login")
	})

	//====Register====
	app.Get("/register", func(c *fiber.Ctx) error {
		// Redirect user if already logged in
		if sesstools.HasSess(c, store) {
			return c.Redirect("/")
		}

		data := flash.Get(c)
		return c.Render("register", fiber.Map{
			"Title": "Register",
			"Flash": data,
			"User":  sesstools.GetUser(c, store),
		}, "layouts/main")
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		// Redirect user if already logged in
		if sesstools.HasSess(c, store) {
			return c.Redirect("/")
		}

		// Get user information
		username := c.FormValue("username")
		email := c.FormValue("email")
		password := c.FormValue("password")
		repeatPassword := c.FormValue("password-repeat")

		//Check that email is not taken
		_, err := dbtools.GetUserFromEmail(email, db)
		if err == nil {
			mp := fiber.Map{
				"message": "An account is already using that email!",
			}
			return flash.WithError(c, mp).Redirect("/register")
		}

		//Check that username is not taken
		_, err = dbtools.GetUserFromUsername(username, db)
		if err == nil {
			mp := fiber.Map{
				"message": "An account is already using that username!",
			}
			return flash.WithError(c, mp).Redirect("/register")
		}

		// Check that passwords matches
		if password != repeatPassword {
			mp := fiber.Map{
				"message": "Password does not match!",
			}
			fmt.Println(err.Error())
			return flash.WithError(c, mp).Redirect("/register")
		}

		//Hash password
		hashedPassword, err := pswdhash.HashPassword(password)

		if err != nil {
			mp := fiber.Map{
				"message": "Could not hash password!",
			}
			fmt.Println(err.Error())
			return flash.WithError(c, mp).Redirect("/register")
		}

		//Create user and add user
		user := dbtools.NewUser(username, email, hashedPassword)
		err = user.AddToDB(db)
		fmt.Println(user)

		if err != nil {
			mp := fiber.Map{
				"message": "Could not add user to db!",
			}
			fmt.Println(err.Error())
			return flash.WithError(c, mp).Redirect("/register")
		}

		// Add user to session
		err = sesstools.AddUser(c, store, user)
		if err != nil {

			mp := fiber.Map{
				"message": "Could not get session!",
			}
			fmt.Println(err.Error())
			return flash.WithError(c, mp).Redirect("/register")

		}

		mp := fiber.Map{
			"message": "Accout created!",
		}
		return flash.WithSuccess(c, mp).Redirect("/")
	})

	//===Offers====
	app.Get("/offer/*", func(c *fiber.Ctx) error {
		param := c.Params("*")
		offer_id, err := strconv.ParseInt(param, 10, 64)

		if err != nil {
			data := flash.Get(c)
			fmt.Println(err.Error())
			return c.Render("error", fiber.Map{
				"Title":        "Error",
				"ErrorCode":    "404",
				"ErrorMessage": "Page not fount.",
				"Flash":        data,
				"User":         sesstools.GetUser(c, store),
			}, "layouts/main")

		}

		offer, err := dbtools.GetOfferFromId(offer_id, db)

		if err != nil {
			data := flash.Get(c)
			fmt.Println(err.Error())
			return c.Render("error", fiber.Map{
				"Title":        "Error",
				"ErrorCode":    "404",
				"ErrorMessage": "Page not fount.",
				"Flash":        data,
				"User":         sesstools.GetUser(c, store),
			}, "layouts/main")

		}
		fmt.Println(offer)

		data := flash.Get(c)
		return c.Render("offer", fiber.Map{
			"Title": "Error",
			"Flash": data,
			"User":  sesstools.GetUser(c, store),
		}, "layouts/main")
	})

	app.Post("/offer", func(c *fiber.Ctx) error {
		// Redirect user if already logged in
		if !sesstools.HasSess(c, store) {
			mp := fiber.Map{
				"message": "You are not logged in!",
			}
			return flash.WithError(c, mp).Redirect("/home")
		}

		fmt.Println(sesstools.GetUser(c, store))
		user_id, ok := sesstools.GetUser(c, store)["User_id"].(int64)
		if !ok {
			mp := fiber.Map{
				"message": "Could not find user!",
			}
			return flash.WithError(c, mp).Redirect("/home")

		}

		title := c.FormValue("title")
		description := c.FormValue("description")
		tags := c.FormValue("tags")
		tagsArr := strings.Split(tags, ",")

		offer := dbtools.NewOffer(user_id, title, description, tagsArr)
		err := offer.AddToDB(db)

		if err != nil {
			mp := fiber.Map{
				"message": "Could not add offer to db!",
			}
			fmt.Println(err.Error())
			return flash.WithError(c, mp).Redirect("/home")
		}

		mp := fiber.Map{
			"message": "Created offer!",
		}
		return flash.WithSuccess(c, mp).Redirect("/home")
	})

	log.Fatal(app.Listen(":3000"))
}
