CREATE SCHEMA IF NOT EXISTS `auto_reply_bot` DEFAULT CHARACTER SET utf8mb4;

-- サーバーのSQL
CREATE TABLE IF NOT EXISTS `auto_reply_bot`.`servers`
(
    `id` VARCHAR(100) NOT NULL,
    `content` JSON NOT NULL,
    `created` DATETIME NOT NULL,
    `updated` DATETIME NOT NULL,

    PRIMARY KEY(`id`),
    UNIQUE INDEX `id_UNIQUE`(`id` ASC)
)
    ENGINE = InnoDB;

-- 他のコンテキストがある場合は、この下にコピペしてください。
