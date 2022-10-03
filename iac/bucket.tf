resource "aws_s3_bucket" "images" {
  bucket = "user-images-lambda-bucket-01"

  tags = {
    Name        = "user images"
    Environment = "dev"
  }
}

resource "aws_s3_bucket_acl" "images" {
  bucket = aws_s3_bucket.images.id
  acl    = "private"
}