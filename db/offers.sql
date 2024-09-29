CREATE TABLE offers(
    offer_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    uploaded TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE offer_tags(
    offer_id INTEGER NOT NULL,
    tag TEXT NOT NULL,

    FOREIGN KEY (offer_id) REFERENCES offers(offer_id)
);
