resource "aws_s3_bucket" "website" {
  bucket = var.bucket_name
}

resource "aws_s3_bucket_website_configuration" "website" {
  bucket = aws_s3_bucket.website.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}

resource "aws_s3_bucket_public_access_block" "website" {
  bucket = aws_s3_bucket.website.id

  block_public_acls       = false
  ignore_public_acls      = false
  block_public_policy     = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "website" {
  bucket = aws_s3_bucket.website.id

  depends_on = [
    aws_s3_bucket_public_access_block.website
  ]

  policy = jsonencode({
    Version = "2012-10-17"

    Statement = [{
      Sid = "PublicRead"

      Effect = "Allow"

      Principal = "*"

      Action = [
        "s3:GetObject"
      ]

      Resource = [
        "${aws_s3_bucket.website.arn}/*"
      ]
    }]
  })
}

locals {
  mime_types = {
    ".html" = "text/html"
    ".css"  = "text/css"
    ".js"   = "application/javascript"
    ".wasm" = "application/wasm"
    ".json" = "application/json"
    ".png"  = "image/png"
    ".jpg"  = "image/jpeg"
    ".jpeg" = "image/jpeg"
    ".svg"  = "image/svg+xml"
    ".ico"  = "image/x-icon"
    ".txt"  = "text/plain"
  }
}

resource "aws_s3_object" "website_files" {
  for_each = fileset("${path.module}/../dist", "**")

  bucket = aws_s3_bucket.website.id

  key    = each.value
  source = "${path.module}/../dist/${each.value}"

  etag = filemd5("${path.module}/../dist/${each.value}")

  content_type = lookup(
    local.mime_types,
    regex("\\.[^.]+$", each.value),
    "application/octet-stream"
  )

  depends_on = [
    aws_s3_bucket_policy.website
  ]
}