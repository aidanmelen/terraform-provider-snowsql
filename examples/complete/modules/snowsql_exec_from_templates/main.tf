locals {
  sql_path = var.sql_path == null ? "${path.root}/sql" : var.sql_path
}

data "template_file" "create" {
  template = file(join("/", [local.sql_path, var.name, "create.tpl"]))
  vars     = var.template_vars
}

data "template_file" "delete" {
  template = file(join("/", [local.sql_path, var.name, "delete.tpl"]))
  vars     = var.template_vars
}

resource "snowsql_exec" "stmts" {
  name = var.name

  create {
    statements = data.template_file.create.rendered
  }

  delete {
    statements = data.template_file.delete.rendered
  }
}
