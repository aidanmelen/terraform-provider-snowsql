resource "random_pet" "name" {
  prefix    = upper(basename(path.cwd))
  separator = "_"
}

locals {
  name = upper(random_pet.name.id)
}

resource "snowsql_exec" "role" {
  name = local.name

  create {
    statements = "CREATE ROLE IF NOT EXISTS ${local.name};"
  }

  # uncomment optional update statements to alter the user in-place after creation
  # update {
  #   statements = "ALTER ROLE IF EXISTS ${local.name} SET COMMENT = 'updated with terraform';"
  # }

  delete {
    statements = "DROP ROLE IF EXISTS ${local.name};"
  }
}
