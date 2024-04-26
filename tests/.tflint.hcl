config {
  format = "compact"
  //plugin_dir = "~/.tflint.d/plugins"

  force = false
  disabled_by_default = false

  ignore_module = {
    "terraform-aws-modules/vpc/aws"            = true
    "terraform-aws-modules/security-group/aws" = true
  }

}

#rule "turyachka_detected" {
#  enabled = true
#}

plugin "template" {
  enabled = true
}

#plugin "comment-multiple-hashes" {
#  enabled = false
#  version = "0.0.1"
#  source = "github.com/serenityzn/tflint-ruleset-comment-multiple-hashes"
#}