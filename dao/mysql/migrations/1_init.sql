-- +migrate Up
CREATE TABLE `tasks`
(
    `tsk_id`             int(11)      NOT NULL AUTO_INCREMENT,
    `tsk_active`         tinyint(1)   NOT NULL,
    `tsk_title`          varchar(255) NOT NULL,
    `tsk_start_height`   int(11)      NOT NULL,
    `tsk_current_height` int(11)      NOT NULL,
    `tsk_end_height`     int(11)      NOT NULL,
    `tsk_batch`          int(11)      NOT NULL,
    PRIMARY KEY (`tsk_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci;

-- +migrate Down
DROP TABLE `tasks`;
