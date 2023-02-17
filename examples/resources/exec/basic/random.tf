resource "random_pet" "name" {
  prefix    = upper(basename(path.cwd))
  separator = "_"
}

locals {
  name = upper(random_pet.name.id)
}
