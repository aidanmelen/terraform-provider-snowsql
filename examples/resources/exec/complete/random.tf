resource "random_pet" "name" {
  prefix    = upper(basename(path.cwd))
  separator = "_"
}

resource "random_password" "password" {
  length  = 16
  special = false
}

locals {
  name = upper(random_pet.name.id)
}
