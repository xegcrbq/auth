CREATE TABLE users(
                      userid SERIAL PRIMARY KEY ,
                      username varchar(50) NOT NULL,
                      password varchar(50) NOT NULL
);

CREATE TABLE refreshSessions(
                                id SERIAL PRIMARY KEY,
                                userid integer,
                                refreshtoken varchar(300) NOT NULL,
                                useragent character varying(200) NOT NULL, /* user-agent */
                                fingerprint varchar(300) NOT NULL,
                                ip character varying(15) NOT NULL,
                                expiresin bigint NOT NULL,
                                createdat timestamp with time zone NOT NULL DEFAULT now()
);

INSERT INTO users (username, password) VALUES
    ('admin', 'admin'),
    ('user', 'password');