DROP TABLE IF EXISTS "users";
CREATE TABLE "users"
(
    "id"       serial       NOT NULL,
    "username" VARCHAR(255) NOT NULL,
    "email"    VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "is_admin" SMALLINT DEFAULT 0,
    "uid"      VARCHAR(255) NOT NULL,
    "score"    NUMERIC,
    "avatar"   VARCHAR(255) NOT NULL,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "datasets";
CREATE TABLE "datasets"
(
    "id"         serial       NOT NULL,
    "name"       VARCHAR(255) NOT NULL,
    "creator_id" INT          NOT NULL,
    "type_id"    INT          NOT NULL,
    "total"      INT DEFAULT 0,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "dataset_types";
CREATE TABLE "dataset_types"
(
    "id"   serial       NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "dataset_categories";
CREATE TABLE "dataset_categories"
(
    "id"   serial       NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "dataset_category_relations";
CREATE TABLE "dataset_category_relations"
(
    "id"          serial NOT NULL,
    "category_id" INT    NOT NULL,
    "dataset_id"  INT    NOT NULL,
    "count"       INT DEFAULT 0,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "annotations";
CREATE TABLE "annotations"
(
    "id"              serial       NOT NULL,
    "status"          SMALLINT DEFAULT 0,
    "url"             VARCHAR(255) NOT NULL,
    "dataset_id"      INT          NOT NULL,
    "replica_count"   SMALLINT,
    "qualified_count" SMALLINT,
    "delivered_count" SMALLINT,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "annotation_users";
CREATE TABLE "annotations_users"
(
    "id"            serial NOT NULL,
    "annotation_id" INT    NOT NULL,
    "user_id"       INT    NOT NULL,
    "status"        SMALLINT DEFAULT 0,
    "result"        json,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "datasets_users";
CREATE TABLE "datasets_users"
(
    "id"         serial NOT NULL,
    "user_id"    INT    NOT NULL,
    "dataset_id" INT    NOT NULL,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "dataset_tags";
CREATE TABLE "dataset_tags"
(
    "id"  serial       NOT NULL,
    "tag" VARCHAR(255) NOT NULL,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "datasets_tags";
CREATE TABLE "datasets_tags"
(
    "id"         serial NOT NULL,
    "dataset_id" INT    NOT NULL,
    "tag_id"     INT    NOT NULL,
    PRIMARY KEY ("id")
);
