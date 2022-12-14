START TRANSACTION;

CREATE TABLE `candidate` (
    `id` varchar(50) NOT NULL,
    `name` varchar(100) DEFAULT '',
    `address` varchar(250) DEFAULT '',
    `experience` tinyint(2) DEFAULT '0',
    `willing_to_relocate` ENUM('yes', 'no') DEFAULT 'no',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `job` (
    `id` varchar(50) NOT NULL,
    `position` varchar(100) DEFAULT '',
    `department` varchar(100) DEFAULT '',
    `requester` varchar(100) DEFAULT '',
    `job_description` text,
    `criteria` text,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `recruitment` (
    `id` varchar(50) NOT NULL,
    `job_id` varchar(50) DEFAULT NULL,
    `status` ENUM('open', 'close') DEFAULT 'open',
    `deadline` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_rec_1` FOREIGN KEY (`job_id`) REFERENCES `job` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `candidate_score` (
    `id` varchar(50) NOT NULL,
    `candidate_id` varchar(50) DEFAULT NULL,
    `recruitment_id` varchar(50) DEFAULT NULL,
    `willing_to_relocate_score` DECIMAL(5,2) DEFAULT '0.00',
    `attitude_score` DECIMAL(5,2) DEFAULT '0.00',
    `skill_score` DECIMAL(5,2) DEFAULT '0.00',
    `experience_score` DECIMAL(5,2) DEFAULT '0.00',
    `overall_score` DECIMAL(5,2) DEFAULT '0.00',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_cs_1` FOREIGN KEY (`candidate_id`) REFERENCES `candidate` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT `fk_cs_2` FOREIGN KEY (`recruitment_id`) REFERENCES `recruitment` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

COMMIT;