-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema finechat
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema finechat
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `finechat` DEFAULT CHARACTER SET utf8 ;
USE `finechat` ;

-- -----------------------------------------------------
-- Table `finechat`.`USERS`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `finechat`.`USERS` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `first_name` VARCHAR(50) NOT NULL,
  `last_name` VARCHAR(50) NOT NULL,
  `email` VARCHAR(100) NOT NULL,
  `password` VARCHAR(350) NOT NULL,
  `enabled` TINYINT NOT NULL DEFAULT 1,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `finechat`.`CHATROOMS`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `finechat`.`CHATROOMS` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NULL,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `finechat`.`CHATROOM_PERMISSIONS`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `finechat`.`CHATROOM_PERMISSIONS` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `permission` VARCHAR(45) NOT NULL,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_date` DATETIME NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `finechat`.`CHATROOM_USERS`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `finechat`.`CHATROOM_USERS` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `chatroom_id` INT UNSIGNED NOT NULL,
  `user_id` INT UNSIGNED NOT NULL,
  `permission_id` INT UNSIGNED NOT NULL,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `fk_CHATROOM_USERS_USERS1_idx` (`user_id` ASC) VISIBLE,
  INDEX `fk_CHATROOM_USERS_CHATROOMS1_idx` (`chatroom_id` ASC) VISIBLE,
  INDEX `fk_CHATROOM_USERS_CHATROOM_PERMISSIONS1_idx` (`permission_id` ASC) VISIBLE,
  CONSTRAINT `fk_CHATROOM_USERS_USERS1`
    FOREIGN KEY (`user_id`)
    REFERENCES `finechat`.`USERS` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_CHATROOM_USERS_CHATROOMS1`
    FOREIGN KEY (`chatroom_id`)
    REFERENCES `finechat`.`CHATROOMS` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_CHATROOM_USERS_CHATROOM_PERMISSIONS1`
    FOREIGN KEY (`permission_id`)
    REFERENCES `finechat`.`CHATROOM_PERMISSIONS` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `finechat`.`CHATROOM_MESSAGES`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `finechat`.`CHATROOM_MESSAGES` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `chatroom_id` INT UNSIGNED NOT NULL,
  `user_id` INT UNSIGNED NOT NULL,
  `message` LONGTEXT NOT NULL,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `fk_CHATROOM_MESSAGES_CHATROOMS1_idx` (`chatroom_id` ASC) VISIBLE,
  INDEX `fk_CHATROOM_MESSAGES_USERS1_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_CHATROOM_MESSAGES_CHATROOMS1`
    FOREIGN KEY (`chatroom_id`)
    REFERENCES `finechat`.`CHATROOMS` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_CHATROOM_MESSAGES_USERS1`
    FOREIGN KEY (`user_id`)
    REFERENCES `finechat`.`USERS` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `finechat`.`USER_FRIENDS`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `finechat`.`USER_FRIENDS` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `friend_id` INT UNSIGNED NOT NULL,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_date` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_USER_FRIENDS_USERS_idx` (`user_id` ASC) VISIBLE,
  INDEX `fk_USER_FRIENDS_USERS1_idx` (`friend_id` ASC) VISIBLE,
  CONSTRAINT `fk_USER_FRIENDS_USERS`
    FOREIGN KEY (`user_id`)
    REFERENCES `finechat`.`USERS` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_USER_FRIENDS_USERS1`
    FOREIGN KEY (`friend_id`)
    REFERENCES `finechat`.`USERS` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

USE `finechat` ;

-- -----------------------------------------------------
-- procedure getChatRoomsByUser
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `getChatRoomsByUser` (
	in pi_userId int
)
BEGIN
	SELECT
		cr.name
	FROM CHATROOMS cr
	INNER JOIN CHATROOM_USERS cru
			ON cru.chatroom_id = cr.id
	WHERE cru.user_id = pi_userId
      AND cru.permission_id = 1;
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure getChatroomUsers
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `getChatroomUsers` (
	in pi_chatroomId int
)
BEGIN
	SELECT
		u.name
    FROM CHATROOM_USERS cru
	INNER JOIN USERS u
			ON u.id = cru.user_id
	WHERE cru.chatroom_id = pi_chatroomId;
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure saveUser
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `saveUser` (
	in pi_fname varchar(50),
    in pi_lname varchar(50),
    in pi_email varchar(100),
    in pi_password varchar(350)
)
BEGIN
	INSERT INTO USERS(first_name, last_name, email, password)
		VALUES (pi_fname, pi_lname, pi_email, pi_password);
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure getUserByEmail_SYNTAX_ERROR
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `getUserByEmail` (
	in pi_email varchar(50),
    in pi_password varchar(350)
)
BEGIN
	SELECT
		u.id,
		u.first_name,
        u.last_name,
        u.email,
        u.password
	FROM USERS u
    WHERE u.email = pi_email;
END$$

DELIMITER ;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
