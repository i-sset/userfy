DROP TABLE IF EXISTS users;


CREATE TABLE users (
    ID    INTEGER PRIMARY KEY AUTOINCREMENT,
	Name  TEXT,
	Email TEXT,
	Age   INT
);

INSERT INTO users VALUES 
    (1, "Josset Garcia", "isset.joset@gmail.com", 26),
    (2, "Silvana Ferreiro", "silvanaf@thoughtworks.com", 42),
    (3, "Adriana Ortega", "adriortega@gmail.com", 30),
    (4, "Javiera Lasus", "javivu@thoughtworks.com", 27),
    (5, "Nicolas Bedregal", "nicobe@gmail.com", 32);
    