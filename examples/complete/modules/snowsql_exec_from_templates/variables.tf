variable "name" {
  description = "The name of the directory containing the create and delete SnowSQL templates."
  type        = string
}

variable "sql_path" {
  description = "The relative path to the `sql` directory containing the SnowSQL create and delete SnowSQL templates. Defaults to the `sql` directory in the root of your Terraform workspace."
  default     = null
  type        = string
}

variable "template_vars" {
  description = "The variables used to render the create and delete SnowSQL templates. This map will be passed to the [template_file](https://registry.terraform.io/providers/hashicorp/template/latest/docs/data-sources/file#vars) data resource."
  type        = map
}
