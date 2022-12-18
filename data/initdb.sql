CREATE SEQUENCE articles_seq;

-- Tambah uuid (GET /articles/:uuid, PUT /articles/:uuid, DELETE /articles/:uuid)
CREATE TABLE articles
(
	id INT NOT NULL DEFAULT NEXTVAL ('articles_seq'),
	uuid CHAR(36) NOT NULL UNIQUE,
	author TEXT,
	title TEXT,
    body TEXT,
	created_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);