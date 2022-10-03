resource "aws_iam_role" "s3_read_write" {
  name = "s3_read_write-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

    managed_policy_arns = [aws_iam_policy.s3_read_write.arn]
}

resource "aws_iam_policy" "s3_read_write" {
  name = "s3_read_write-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["s3:GetObject", "s3:PutObject"]
        Effect   = "Allow"
        Resource = "${aws_s3_bucket.images.arn}/*"
      },
    ]
  })
}

resource "aws_lambda_function" "img_zipping" {
  # If the file is not in the current working directory you will need to include a
  # path.module in the filename.
  filename      = "${path.module}/../img_zipping_lambda.zip"
  function_name = "img_zipping_lambda"
  role          = aws_iam_role.s3_read_write.arn
  handler       = "main"

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = filebase64sha256("../img_zipping_lambda.zip")

  runtime = "go1.x"

  environment {
    variables = {
      Environment = "dev"
    }
  }
}