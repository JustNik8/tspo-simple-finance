CREATE TABLE "users"(
                        "id" UUID NOT NULL,
                        "email" TEXT NOT NULL,
                        "username" TEXT NOT NULL,
                        "hash_pass" TEXT NOT NULL,
                        "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "users" ADD PRIMARY KEY("id");
CREATE TABLE "transactions"(
                               "id" UUID NOT NULL,
                               "user_id" UUID NOT NULL,
                               "amout" DOUBLE PRECISION NOT NULL,
                               "category_id" UUID NOT NULL,
                               "comment" TEXT NOT NULL,
                               "date" DATE NOT NULL,
                               "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "transactions" ADD PRIMARY KEY("id");
CREATE TABLE "categories"(
                             "id" UUID NOT NULL,
                             "user_id" UUID NOT NULL,
                             "name" TEXT NOT NULL
);
ALTER TABLE
    "categories" ADD PRIMARY KEY("id");
CREATE TABLE "tags"(
                       "id" UUID NOT NULL,
                       "user_id" BIGINT NOT NULL,
                       "name" TEXT NOT NULL
);
ALTER TABLE
    "tags" ADD PRIMARY KEY("id");
CREATE TABLE "transaction_tags"(
                                   "transaction_id" UUID NOT NULL,
                                   "tag_id" UUID NOT NULL
);
CREATE TABLE "incomes"(
                          "id" UUID NOT NULL,
                          "user_id" UUID NOT NULL,
                          "amount" BIGINT NOT NULL,
                          "comment" TEXT NOT NULL,
                          "date" DATE NOT NULL,
                          "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "incomes" ADD PRIMARY KEY("id");
ALTER TABLE
    "incomes" ADD CONSTRAINT "incomes_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "transaction_tags" ADD CONSTRAINT "transaction_tags_transaction_id_foreign" FOREIGN KEY("transaction_id") REFERENCES "transactions"("id");
ALTER TABLE
    "transactions" ADD CONSTRAINT "transactions_category_id_foreign" FOREIGN KEY("category_id") REFERENCES "categories"("id");
ALTER TABLE
    "transactions" ADD CONSTRAINT "transactions_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "transaction_tags" ADD CONSTRAINT "transaction_tags_tag_id_foreign" FOREIGN KEY("tag_id") REFERENCES "tags"("id");