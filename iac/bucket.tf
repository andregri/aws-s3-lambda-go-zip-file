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

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.images.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.img_zipping.arn
    events              = ["s3:ObjectCreated:*"]
    #filter_prefix       = "AWSLogs/"
    #filter_suffix       = ".log"
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}