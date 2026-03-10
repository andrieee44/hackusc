CREATE TABLE files (
	id          BINARY(16)      NOT NULL,
	checksum    BINARY(32)      NOT NULL,
	storage_key VARCHAR(512)    NOT NULL,
	mime_type   VARCHAR(255)    NOT NULL,
	file_size   BIGINT UNSIGNED NOT NULL,

	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT uq_files_checksum UNIQUE (checksum)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE addresses (
	id          BINARY(16)    NOT NULL,
	street      VARCHAR(255),
	locality    VARCHAR(255),
	city        VARCHAR(255),
	province    VARCHAR(255),
	postal_code VARCHAR(20),
	country     CHAR(2),
	location    POINT,

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT chk_addresses_location CHECK (
    	location IS NULL OR (
        	ST_X(location) BETWEEN -180 AND 180 AND
        	ST_Y(location) BETWEEN -90 AND 90
    	)
	)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE timestamps (
	id         BINARY(16)  NOT NULL,
	created_at DATETIME(6) NOT NULL,
	updated_at DATETIME(6) NOT NULL,

	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT chk_timestamps_updated_at CHECK (updated_at >= created_at)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE moderation_states (
	id          BINARY(16)   NOT NULL,
	name        VARCHAR(255) NOT NULL,
	description TEXT         NOT NULL,

	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT uq_moderation_states_name UNIQUE (name)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE users (
	id            BINARY(16)     NOT NULL,
	timestamp_id  BINARY(16)     NOT NULL,
	address_id    BINARY(16)     NOT NULL,
	email         VARCHAR(320)   NOT NULL,
	first_name    VARCHAR(255)   NOT NULL,
	last_name     VARCHAR(255)   NOT NULL,
	password_hash VARBINARY(255) NOT NULL,
	middle_name   VARCHAR(255),

	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT uq_users_email UNIQUE (email),

	CONSTRAINT fk_users_timestamp_id
		FOREIGN KEY (timestamp_id)
		REFERENCES timestamps(id),

	CONSTRAINT fk_users_address_id
		FOREIGN KEY (address_id)
		REFERENCES addresses(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE authorships (
	id          BINARY(16) NOT NULL,
	creator_id  BINARY(16),
	modifier_id BINARY(16),

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT fk_authorships_creator_id
		FOREIGN KEY (creator_id)
		REFERENCES users(id)
		ON DELETE SET NULL,

	CONSTRAINT fk_authorships_modifier_id
		FOREIGN KEY (modifier_id)
		REFERENCES users(id)
		ON DELETE SET NULL
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE moderations (
	id           BINARY(16)  NOT NULL,
	state_id     BINARY(16)  NOT NULL,
	timestamp_id BINARY(16)  NOT NULL,
	moderator_id BINARY(16),

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT fk_moderations_state_id
		FOREIGN KEY (state_id)
		REFERENCES moderation_states(id),

	CONSTRAINT fk_moderations_timestamp_id
		FOREIGN KEY (timestamp_id)
		REFERENCES timestamps(id),

	CONSTRAINT fk_moderations_moderator_id
		FOREIGN KEY (moderator_id)
		REFERENCES users(id)
		ON DELETE SET NULL
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

ALTER TABLE users
	ADD COLUMN moderation_id BINARY(16) NULL AFTER id,
	ADD CONSTRAINT fk_users_moderation_id
		FOREIGN KEY (moderation_id)
		REFERENCES moderations(id);

CREATE TABLE avatars (
	user_id      BINARY(16)    NOT NULL,
	timestamp_id BINARY(16)    NOT NULL,
	file_id      BINARY(16),
	width        INT UNSIGNED,
	height       INT UNSIGNED,

	CONSTRAINT PRIMARY KEY (user_id),

	CONSTRAINT fk_avatars_user_id
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_avatars_timestamp_id
		FOREIGN KEY (timestamp_id)
		REFERENCES timestamps(id),

	CONSTRAINT fk_avatars_file_id
		FOREIGN KEY (file_id)
		REFERENCES files(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE comments (
	id            BINARY(16) NOT NULL,
	moderation_id BINARY(16) NOT NULL,
	timestamp_id  BINARY(16) NOT NULL,
	address_id    BINARY(16) NOT NULL,
	creator_id    BINARY(16),
	content       TEXT       NOT NULL,
	pinned        BOOLEAN    NOT NULL,

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT fk_comments_moderation_id
		FOREIGN KEY (moderation_id)
		REFERENCES moderations(id),

	CONSTRAINT fk_comments_timestamp_id
		FOREIGN KEY (timestamp_id)
		REFERENCES timestamps(id),

	CONSTRAINT fk_comments_address_id
		FOREIGN KEY (address_id)
		REFERENCES addresses(id),

	CONSTRAINT fk_comments_creator_id
		FOREIGN KEY (creator_id)
		REFERENCES users(id)
		ON DELETE SET NULL
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE community_types (
	id          BINARY(16)   NOT NULL,
	name        VARCHAR(255) NOT NULL,
	description TEXT         NOT NULL,

	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT uq_community_types_name UNIQUE (name)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE communities (
	id            BINARY(16)   NOT NULL,
	type_id       BINARY(16)   NOT NULL,
	root_id       BINARY(16)   NOT NULL,
	moderation_id BINARY(16)   NOT NULL,
	timestamp_id  BINARY(16)   NOT NULL,
	address_id    BINARY(16)   NOT NULL,
	creator_id    BINARY(16),
	name          VARCHAR(255) NOT NULL,
	description   TEXT,

	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT uq_communities_type_id_name UNIQUE (type_id, name),

	CONSTRAINT fk_communities_type_id
		FOREIGN KEY (type_id)
		REFERENCES community_types(id),

	CONSTRAINT fk_communities_root_id
		FOREIGN KEY (root_id)
		REFERENCES users(id),

	CONSTRAINT fk_communities_moderation_id
		FOREIGN KEY (moderation_id)
		REFERENCES moderations(id),

	CONSTRAINT fk_communities_timestamp_id
		FOREIGN KEY (timestamp_id)
		REFERENCES timestamps(id),

	CONSTRAINT fk_communities_address_id
		FOREIGN KEY (address_id)
		REFERENCES addresses(id),

	CONSTRAINT fk_communities_creator_id
		FOREIGN KEY (creator_id)
		REFERENCES users(id)
		ON DELETE SET NULL
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE users_communities (
	user_id      BINARY(16)   NOT NULL,
	community_id BINARY(16)   NOT NULL,
	created_by   TIMESTAMP(6) NOT NULL,

	CONSTRAINT PRIMARY KEY (user_id, community_id),

	CONSTRAINT fk_users_communities_user_id
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_users_communities_community_id
		FOREIGN KEY (community_id)
		REFERENCES communities(id)
		ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE ticket_states (
	id   BINARY(16)   NOT NULL,
	name VARCHAR(255) NOT NULL,

	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT uq_ticket_states_name UNIQUE (name)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tickets (
	id            BINARY(16)   NOT NULL,
	community_id  BINARY(16)   NOT NULL,
	moderation_id BINARY(16)   NOT NULL,
	timestamp_id  BINARY(16)   NOT NULL,
	address_id    BINARY(16)   NOT NULL,
	state_id      BINARY(16)   NOT NULL,
	creator_id    BINARY(16),
	title         VARCHAR(255) NOT NULL,
	description   TEXT,

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT fk_tickets_community_id
		FOREIGN KEY (community_id)
		REFERENCES communities(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_tickets_moderation_id
		FOREIGN KEY (moderation_id)
		REFERENCES moderations(id),

	CONSTRAINT fk_tickets_timestamp_id
		FOREIGN KEY (timestamp_id)
		REFERENCES timestamps(id),

	CONSTRAINT fk_tickets_address_id
		FOREIGN KEY (address_id)
		REFERENCES addresses(id),

	CONSTRAINT fk_tickets_state_id
		FOREIGN KEY (state_id)
		REFERENCES ticket_states(id),

	CONSTRAINT fk_tickets_creator_id
		FOREIGN KEY (creator_id)
		REFERENCES users(id)
		ON DELETE SET NULL
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tickets_states_history (
	id         BINARY(16)  NOT NULL,
	state_id   BINARY(16)  NOT NULL,
	ticket_id  BINARY(16)  NOT NULL,
	creator_id BINARY(16),
	created_at DATETIME(6) NOT NULL,

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT fk_tickets_states_history_state_id
		FOREIGN KEY (state_id)
		REFERENCES ticket_states(id),

	CONSTRAINT fk_tickets_states_history_ticket_id
		FOREIGN KEY (ticket_id)
		REFERENCES tickets(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_tickets_states_history_creator_id
		FOREIGN KEY (creator_id)
		REFERENCES users(id)
		ON DELETE SET NULL
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tickets_assignees (
	ticket_id BINARY(16) NOT NULL,
	user_id   BINARY(16) NOT NULL,

	CONSTRAINT PRIMARY KEY (ticket_id, user_id),

	CONSTRAINT fk_tickets_assignees_ticket_id
		FOREIGN KEY (ticket_id)
		REFERENCES tickets(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_tickets_assignees_user_id
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tickets_assignment_history (
	id          BINARY(16)  NOT NULL,
	ticket_id   BINARY(16)  NOT NULL,
	creator_id  BINARY(16),
	assignee_id BINARY(16),
	created_at  DATETIME(6) NOT NULL,
	added       BOOLEAN     NOT NULL,

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT fk_tickets_assignment_history_ticket_id
		FOREIGN KEY (ticket_id)
		REFERENCES tickets(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_tickets_assignment_history_creator_id
		FOREIGN KEY (creator_id)
		REFERENCES users(id)
		ON DELETE SET NULL,

	CONSTRAINT fk_tickets_assignment_history_assignee_id
		FOREIGN KEY (assignee_id)
		REFERENCES users(id)
		ON DELETE SET NULL
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE global_roles (
	id          BINARY(16)      NOT NULL,
	name        VARCHAR(255)    NOT NULL,
	permissions BIGINT UNSIGNED NOT NULL,
	description TEXT            NOT NULL,

	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT uq_global_roles_name UNIQUE (name)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE users_global_roles (
	user_id    BINARY(16)  NOT NULL,
	role_id    BINARY(16)  NOT NULL,
	created_at DATETIME(6) NOT NULL,

	CONSTRAINT PRIMARY KEY (user_id, role_id),

	CONSTRAINT fk_users_global_roles_user_id
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_users_global_roles_role_id
		FOREIGN KEY (role_id)
		REFERENCES global_roles(id)
		ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE communities_roles (
	id            BINARY(16)      NOT NULL,
	community_id  BINARY(16)      NOT NULL,
	timestamp_id  BINARY(16)      NOT NULL,
	authorship_id BINARY(16)      NOT NULL,
	name          VARCHAR(255)    NOT NULL,
	permissions   BIGINT UNSIGNED NOT NULL,
	description   TEXT,

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT uq_communities_roles_community_id_name
		UNIQUE (community_id, name),

	CONSTRAINT fk_communities_roles_community_id
		FOREIGN KEY (community_id)
		REFERENCES communities(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_communities_roles_timestamp_id
		FOREIGN KEY (timestamp_id)
		REFERENCES timestamps(id),

	CONSTRAINT fk_communities_roles_authorship_id
		FOREIGN KEY (authorship_id)
		REFERENCES authorships(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE users_communities_roles (
	user_id    BINARY(16)  NOT NULL,
	role_id    BINARY(16)  NOT NULL,
	created_at DATETIME(6) NOT NULL,

	CONSTRAINT PRIMARY KEY (user_id, role_id),

	CONSTRAINT fk_users_communities_roles_user_id
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_users_communities_roles_role_id
		FOREIGN KEY (role_id)
		REFERENCES communities_roles(id)
		ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE labels (
	id            BINARY(16)   NOT NULL,
	community_id  BINARY(16)   NOT NULL,
	timestamp_id  BINARY(16)   NOT NULL,
	authorship_id BINARY(16)   NOT NULL,
	name          VARCHAR(255) NOT NULL,
	color         CHAR(6)      NOT NULL,
	description   TEXT,

	CONSTRAINT PRIMARY KEY (id),
	CONSTRAINT uq_labels_community_id_name UNIQUE (community_id, name),

	CONSTRAINT fk_labels_community_id
		FOREIGN KEY (community_id)
		REFERENCES communities(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_labels_timestamp_id
		FOREIGN KEY (timestamp_id)
		REFERENCES timestamps(id),

	CONSTRAINT fk_labels_authorship_id
		FOREIGN KEY (authorship_id)
		REFERENCES authorships(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tickets_labels (
	ticket_id  BINARY(16)  NOT NULL,
	label_id   BINARY(16)  NOT NULL,
	created_at DATETIME(6) NOT NULL,

	CONSTRAINT PRIMARY KEY (ticket_id, label_id),

	CONSTRAINT fk_tickets_labels_ticket_id
		FOREIGN KEY (ticket_id)
		REFERENCES tickets(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_tickets_labels_label_id
		FOREIGN KEY (label_id)
		REFERENCES labels(id)
		ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tickets_files (
	ticket_id  BINARY(16)   NOT NULL,
	file_id    BINARY(16)   NOT NULL,
	name       VARCHAR(255) NOT NULL,
	created_at DATETIME(6)  NOT NULL,

	CONSTRAINT PRIMARY KEY (ticket_id, file_id),

	CONSTRAINT fk_tickets_files_ticket_id
		FOREIGN KEY (ticket_id)
		REFERENCES tickets(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_tickets_files_file_id
		FOREIGN KEY (file_id)
		REFERENCES files(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tickets_comments (
	id         BINARY(16) NOT NULL,
	ticket_id  BINARY(16) NOT NULL,
	comment_id BINARY(16) NOT NULL,

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT fk_tickets_comments_ticket_id
		FOREIGN KEY (ticket_id)
		REFERENCES tickets(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_tickets_comments_comment_id
		FOREIGN KEY (comment_id)
		REFERENCES comments(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE tickets_comments_files (
	tickets_comment_id BINARY(16)   NOT NULL,
	file_id            BINARY(16)   NOT NULL,
	name               VARCHAR(255) NOT NULL,
	created_at         DATETIME(6)  NOT NULL,

	CONSTRAINT PRIMARY KEY (tickets_comment_id, file_id),

	CONSTRAINT fk_tickets_comments_files_ticket_tickets_comment_id
		FOREIGN KEY (tickets_comment_id)
		REFERENCES tickets_comments(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_comments_files_file_id
		FOREIGN KEY (file_id)
		REFERENCES files(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE posts (
	id            BINARY(16)   NOT NULL,
	community_id  BINARY(16)   NOT NULL,
	moderation_id BINARY(16)   NOT NULL,
	timestamp_id  BINARY(16)   NOT NULL,
	address_id    BINARY(16)   NOT NULL,
	creator_id    BINARY(16),
	title         VARCHAR(255) NOT NULL,
	content       TEXT         NOT NULL,

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT fk_posts_community_id
		FOREIGN KEY (community_id)
		REFERENCES communities(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_posts_moderation_id
		FOREIGN KEY (moderation_id)
		REFERENCES moderations(id),

	CONSTRAINT fk_posts_timestamp_id
		FOREIGN KEY (timestamp_id)
		REFERENCES timestamps(id),

	CONSTRAINT fk_posts_address_id
		FOREIGN KEY (address_id)
		REFERENCES addresses(id),

	CONSTRAINT fk_posts_creator_id
		FOREIGN KEY (creator_id)
		REFERENCES users(id)
		ON DELETE SET NULL
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE posts_files (
	post_id    BINARY(16)   NOT NULL,
	file_id    BINARY(16)   NOT NULL,
	name       VARCHAR(255) NOT NULL,
	created_at DATETIME(6)  NOT NULL,

	CONSTRAINT PRIMARY KEY (post_id, file_id),

	CONSTRAINT fk_posts_files_post_id
		FOREIGN KEY (post_id)
		REFERENCES posts(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_posts_files_file_id
		FOREIGN KEY (file_id)
		REFERENCES files(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE posts_comments (
	id         BINARY(16) NOT NULL,
	post_id    BINARY(16) NOT NULL,
	comment_id BINARY(16) NOT NULL,

	CONSTRAINT PRIMARY KEY (id),

	CONSTRAINT fk_posts_comments_post_id
		FOREIGN KEY (post_id)
		REFERENCES posts(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_posts_comments_comment_id
		FOREIGN KEY (comment_id)
		REFERENCES comments(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE posts_comments_files (
	posts_comment_id BINARY(16)   NOT NULL,
	file_id          BINARY(16)   NOT NULL,
	name             VARCHAR(255) NOT NULL,
	created_at       DATETIME(6)  NOT NULL,

	CONSTRAINT PRIMARY KEY (posts_comment_id, file_id),

	CONSTRAINT fk_posts_comments_files_post_posts_comment_id
		FOREIGN KEY (posts_comment_id)
		REFERENCES posts_comments(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_posts_comments_files_file_id
		FOREIGN KEY (file_id)
		REFERENCES files(id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
