schema "public" {
  name = "public"
}

table "users" {
  schema = schema.public

  column "id" {
    type = uuid
    null = false
  }

  column "role_id" {
    type = uuid
    null = true
  }
  
  column "first_name" {
    type = varchar(32)
    null = false
  }

  column "last_name" {
    type = varchar(32)
    null = false
  }

  column "email" {
    type = text
    null = false
  }

  column "password" {
    type = text
    null = false
  }

  column "phone" {
    type = text
    null = false
    default = ""
  }

  column "is_active" {
    type = bool
    null = false
    default = true
  }

  column "created_at" {
    type = timestamp
    null = false
    default = sql("now()")
  }

  column "updated_at" {
    type = timestamp
    null = true
  }

  primary_key {
    columns = [column.id]
  }
  
  foreign_key "role_id" {
    columns     = [column.role_id]
    ref_columns = [table.roles.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  
  unique "users_email_key" {
      columns = [column.email]
  }

  unique "users_phone_key" {
      columns = [column.phone]
  }
}

table "roles" {
  schema = schema.public

  column "id" {
    type = uuid
    null = false
  }

  column "name" {
    type = varchar(32)
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  unique "roles_name_key" {
      columns = [column.name]
  }
}

table "posts" {
  schema = schema.public

  column "id" {
    type = uuid
    null = false
  }

  column "title" {
    type = varchar(255)
    null = false
  }

  column "content" {
    type = text
    null = true
  }

  column "published" {
    type = bool
    null = false
    default = false
  }

  column "user_id" {
    type = uuid
    null = false
  }

  column "created_at" {
    type = timestamp
    null = false
    default = sql("now()")
  }

  column "updated_at" {
    type = timestamp
    null = true
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
