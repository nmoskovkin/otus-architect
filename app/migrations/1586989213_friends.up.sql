CREATE TABLE friends
(
    `from_user` varchar(36) NOT NULL,
    `to_user` varchar(36) NOT NULL,
    FOREIGN KEY (`from_user`) REFERENCES users(id),
    FOREIGN KEY (`to_user`) REFERENCES users(id),
    PRIMARY KEY(`from_user`, `to_user`)
);
