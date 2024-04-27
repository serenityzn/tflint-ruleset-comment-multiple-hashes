resource "aws_instance" "test" {
  ami = var.asb
  instance_type = var.test
  associate_public_ip_address =  data.aws_ec2_host.test
}




