package dbtools

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ====Users====
type User struct {
	User_id  int64  `field:"user_id"`
	Username string `field:"username"`
	Email    string `field:"email"`
	Password string `field:"password"`
}

func NewUser(username, email, password string) *User {
	return &User{
		User_id:  -1,
		Username: username,
		Email:    email,
		Password: password,
	}
}

func EmptyUser() *User {
	return &User{User_id: -1}
}

func (u *User) AsMap() fiber.Map {
	return fiber.Map{
		"User_id":  u.User_id,
		"Username": u.Username,
		"Email":    u.Email,
	}
}

func (u *User) AddToDB(db *sql.DB) error {
	// Make statement
	prompt := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	statement, err := db.Prepare(prompt)
	if err != nil {
		fmt.Println("Error at statement")
		return err
	}

	// Execute statement
	result, err := statement.Exec(u.Username, u.Email, u.Password)
	if err != nil {
		fmt.Println("Error at execution")
		return err
	}

	// Give user there id
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error at id")
		return err
	}
	fmt.Println(id)
	u.User_id = id

	return nil
}

func GetUserFromEmail(email string, db *sql.DB) (*User, error) {
	prompt := "SELECT * FROM users WHERE email=?"
	row := db.QueryRow(prompt, email)

	user := EmptyUser()
	err := row.Scan(&user.User_id, &user.Username, &user.Email, &user.Password)

	return user, err
}

func GetUserFromId(user_id int64, db *sql.DB) (*User, error) {
	prompt := "SELECT * FROM users WHERE user_id=?"
	row := db.QueryRow(prompt, user_id)

	user := EmptyUser()
	err := row.Scan(&user.User_id, &user.Username, &user.Email, &user.Password)

	return user, err
}

func GetUserFromUsername(username string, db *sql.DB) (*User, error) {
	prompt := "SELECT * FROM users WHERE username=?"
	row := db.QueryRow(prompt, username)

	user := EmptyUser()
	err := row.Scan(&user.User_id, &user.Username, &user.Email, &user.Password)

	return user, err
}

// ====Offers====
type Offer struct {
	Offer_id    int64     `field:"offer_id"`
	User_id     int64     `field:"user_id"`
	Title       string    `field:"title"`
	Description string    `field:"description"`
	Uploaded    time.Time `field:"uploaded"`
	PrettyTime  string
	Tags        []Tag
	User        *User
}

type Tag struct {
	Offer_id int64  `field:"offer_id"`
	Tag      string `field:"tag"`
}

func NewOffer(user_id int64, title, description string, stringTags []string) *Offer {
	tags := make([]Tag, len(stringTags))
	for i, tag := range stringTags {
		tags[i].Offer_id = -1
		tags[i].Tag = tag
	}
	return &Offer{
		Offer_id:    -1,
		User_id:     user_id,
		Title:       title,
		Description: description,
		Tags:        tags,
	}
}

func EmptyOffer() *Offer {
	return &Offer{
		Offer_id: -1,
		User_id:  -1,
		User:     EmptyUser(),
	}
}

func EmptyTag() *Tag {
	return &Tag{}
}

func (o *Offer) AsMap() fiber.Map {
	return fiber.Map{
		"Offer_id":    o.Offer_id,
		"Title":       o.Title,
		"Description": o.Description,
		"Time":        o.PrettyTime,
		"User":        o.User.AsMap(),
		"Tags":        TagsAsMap(o.Tags),
	}
}

func TagsAsMap(tags []Tag) []string {
	sTags := make([]string, len(tags))
	for i, tag := range tags {
		sTags[i] = tag.Tag
	}
	return sTags
}
func (o *Offer) AddToDB(db *sql.DB) error {
	// Make statement
	prompt := "INSERT INTO offers (user_id, title, description) VALUES (?, ?, ?)"
	statement, err := db.Prepare(prompt)
	if err != nil {
		fmt.Println("Error at statement")
		return err
	}

	// Execute statement
	result, err := statement.Exec(o.User_id, o.Title, o.Description)
	if err != nil {
		fmt.Println("Error at execution")
		return err
	}

	// Give offer it's id
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error at id")
		return err
	}
	fmt.Println(id)
	o.Offer_id = id

	for i := range o.Tags {
		o.Tags[i].Offer_id = o.Offer_id
		err = o.Tags[i].AddToDB(db)
		if err != nil {
			fmt.Println("Error trying to add tag")
			return err
		}
	}

	return nil
}

func (t *Tag) AddToDB(db *sql.DB) error {
	// Make statement
	prompt := "INSERT INTO offer_tags (offer_id, tag) VALUES (?, ?)"
	statement, err := db.Prepare(prompt)
	if err != nil {
		fmt.Println("Error at statement")
		return err
	}

	// Execute statement
	_, err = statement.Exec(t.Offer_id, t.Tag)
	if err != nil {
		fmt.Println("Error at execution")
		return err
	}

	return nil
}

func GetOfferFromId(offer_id int64, db *sql.DB) (*Offer, error) {
	prompt := "SELECT * FROM offers WHERE offer_id=?"
	row := db.QueryRow(prompt, offer_id)

	offer := EmptyOffer()
	err := row.Scan(&offer.Offer_id, &offer.User_id, &offer.Title, &offer.Description, &offer.Uploaded)
	if err != nil {
		return offer, err
	}

	offer.PrettyTime = formatTimeDiff(time.Now().UTC().Sub(offer.Uploaded))

	offer.User, err = GetUserFromId(offer.User_id, db)
	if err != nil {
		return offer, err
	}

	offer.Tags, err = GetTagsFromId(offer_id, db)
	if err != nil {
		return offer, err
	}

	return offer, nil
}

func GetOffers(db *sql.DB) ([]Offer, error) {
	offers := make([]Offer, 0)
	prompt := "SELECT * FROM offers"
	rows, err := db.Query(prompt)
	defer rows.Close()

	if err != nil {
		return offers, err
	}

	for rows.Next() {
		offer := EmptyOffer()
		err = rows.Scan(&offer.Offer_id, &offer.User_id, &offer.Title, &offer.Description, &offer.Uploaded)
		if err != nil {
			return offers, err
		}

		offer.PrettyTime = formatTimeDiff(time.Now().UTC().Sub(offer.Uploaded))

		offer.User, err = GetUserFromId(offer.User_id, db)
		if err != nil {
			return offers, err
		}

		offer.Tags, err = GetTagsFromId(offer.Offer_id, db)
		if err != nil {
			return offers, err
		}
		offers = append(offers, *offer)
	}

	// Sort based on time
	sort.Slice(offers, func(i, j int) bool {
		return !offers[i].Uploaded.Before(offers[j].Uploaded)
	})

	return offers, nil
}

func OffersAsMaps(offers []Offer) []fiber.Map {
	maps := make([]fiber.Map, len(offers))

	for i, offer := range offers {
		maps[i] = offer.AsMap()
	}

	return maps
}

func GetTagsFromId(offer_id int64, db *sql.DB) ([]Tag, error) {
	tags := make([]Tag, 0)
	prompt := "SELECT * FROM offer_tags WHERE offer_id=?"
	rows, err := db.Query(prompt, offer_id)
	defer rows.Close()

	if err != nil {
		return tags, err
	}

	for rows.Next() {
		tag := EmptyTag()
		err = rows.Scan(&tag.Offer_id, &tag.Tag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, *tag)
	}

	return tags, nil
}

func formatTimeDiff(d time.Duration) string {
	hours := d.Hours()
	minutes := d.Minutes() - float64(math.Floor(hours))*60
	seconds := d.Seconds() - float64(math.Floor(minutes))*60

	days := int(hours) / 24
	hours = float64(int64(hours) % 24.0)

	if days > 0 {
		return fmt.Sprintf("%d days %dh", days, int(hours))
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", int(hours), int(minutes))
	} else if minutes > 0 {
		return fmt.Sprintf("%dm", int(minutes))
	} else {
		return fmt.Sprintf("%ds", int(seconds))
	}
}
