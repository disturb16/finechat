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
  `user_id` INT UNSIGNED NOT NULL,
  `name` VARCHAR(45) NULL,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `fk_CHATROOMS_USERS1_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_CHATROOMS_USERS1`
    FOREIGN KEY (`user_id`)
    REFERENCES `finechat`.`USERS` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `finechat`.`CHATROOM_GUESTS`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `finechat`.`CHATROOM_GUESTS` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `chatroom_id` INT UNSIGNED NOT NULL,
  `user_id` INT UNSIGNED NOT NULL,
  `created_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `fk_CHATROOM_USERS_USERS1_idx` (`user_id` ASC) VISIBLE,
  INDEX `fk_CHATROOM_USERS_CHATROOMS1_idx` (`chatroom_id` ASC) VISIBLE,
  CONSTRAINT `fk_CHATROOM_USERS_USERS1`
    FOREIGN KEY (`user_id`)
    REFERENCES `finechat`.`USERS` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_CHATROOM_USERS_CHATROOMS1`
    FOREIGN KEY (`chatroom_id`)
    REFERENCES `finechat`.`CHATROOMS` (`id`)
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
  `message` VARCHAR(500) NOT NULL,
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
	in pi_email varchar(50)
)
BEGIN
	SELECT
		cr.id,
		cr.name
	FROM CHATROOMS cr
	LEFT JOIN CHATROOM_GUESTS crg
			ON crg.chatroom_id = cr.id
	LEFT JOIN USERS u
			ON u.id = cr.user_id
		    OR u.id = crg.user_id
	WHERE u.email = pi_email
	  AND u.enabled = 1;
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
		u.email,
		CONCAT(u.first_name, ' ', u.last_name) as name
    FROM CHATROOM_GUESTS cru
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
-- procedure getUserByEmail
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `getUserByEmail` (
	in pi_email varchar(50)
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

-- -----------------------------------------------------
-- procedure saveChatRoom
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `saveChatRoom` (
	in pi_name varchar(45),
    in pi_userId int
)
BEGIN
	INSERT INTO CHATROOMS(name, user_id) VALUES (pi_name, pi_userId);
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure getChatRoomMessages
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `getChatRoomMessages` (
	in pi_chatRoomId int
)
BEGIN
	SELECT
		crm.message,
        u.id as user_id,
        CONCAT(u.first_name, ' ', u.last_name) as user,
        crm.created_date
	FROM CHATROOM_MESSAGES crm
    INNER JOIN USERS u
			ON crm.user_id = u.id
	WHERE crm.chatroom_id = pi_chatRoomId
    ORDER BY crm.created_date DESC
    LIMIT 50;
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure saveChatRoomMessage
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `saveChatRoomMessage` (
	in pi_chatRoomId int,
    in pi_email varchar(50),
    in pi_message varchar(500),
    in pi_createdDate datetime
)
BEGIN
	SELECT @userId := id FROM USERS WHERE email = pi_email and enabled = 1;

	INSERT INTO CHATROOM_MESSAGES (chatroom_id, user_id, message, created_date)
		VALUES(pi_chatRoomId, @userId, pi_message, pi_createdDate);
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure saveChatRoomUser
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `saveChatRoomUser` (
	in pi_chatRoomId int,
    in pi_userId int
)
BEGIN
	INSERT INTO CHATROOM_GUESTS(chatroom_id, user_id)
		VALUES (pi_chatRoomId, pi_userId);
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure saveFriend
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `saveFriend` (
	in pi_email varchar(50),
    in pi_friendEmail varchar(50)
)
BEGIN
	select @userId := id from USERS where email = pi_email and enabled =1;
    select @friendId := id from USERS where email = pi_friendEmail and enabled =1;

	INSERT INTO USER_FRIENDS (user_id, friend_id)
		VALUES(@userId, @friendId);
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure getUserFriends
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `getUserFriends` (
	in pi_email varchar(50)
)
BEGIN
	SELECT
		 u.email,
		 u.first_name,
		 u.last_name
	FROM USER_FRIENDS uf
	INNER JOIN USERS u
			ON u.id  = uf.friend_id
		   AND u.enabled = 1
	WHERE uf.user_id = (select id from USERS where email = pi_email and enabled =1);
END$$

DELIMITER ;

-- -----------------------------------------------------
-- procedure removeChatRoomUser
-- -----------------------------------------------------

DELIMITER $$
USE `finechat`$$
CREATE PROCEDURE `removeChatRoomUser` (
	in pi_chatRoomId int,
	in pi_email varchar(50)
)
BEGIN
	DELETE FROM CHATROOM_GUESTS
    WHERE chatroom_id = pi_chatRoomId
      AND user_id = (SELECT id FROM USERS WHERE email = pi_email);
END$$

DELIMITER ;

CREATE USER
IF NOT EXISTS 'appuser'@'%' IDENTIFIED BY '1234';
GRANT EXECUTE ON `finechat`.* TO 'appuser'@'%';

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
