locals {
  s3_origin_id  = "s3-videos"
  alb_origin_id = "alb-frontend"
}

# ── WAF Web ACL (must be in us-east-1 for CloudFront) ───────────

resource "aws_wafv2_web_acl" "main" {
  provider = aws.us_east_1

  name        = "sravas-waf-${var.environment}"
  description = "WAF rules for CloudFront distribution ${var.environment}"
  scope       = "CLOUDFRONT"

  default_action {
    allow {}
  }

  rule {
    name     = "AWSManagedRulesCommonRuleSet"
    priority = 1
    override_action {
      none {}
    }
    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesCommonRuleSet"
        vendor_name = "AWS"
      }
    }
    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name               = "AWSCommonRuleSet-${var.environment}"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "AWSManagedRulesKnownBadInputsRuleSet"
    priority = 10
    override_action {
      none {}
    }
    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesKnownBadInputsRuleSet"
        vendor_name = "AWS"
      }
    }
    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name               = "KnownBadInputs-${var.environment}"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "AWSManagedRulesAmazonIpReputationList"
    priority = 20
    override_action {
      none {}
    }
    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesAmazonIpReputationList"
        vendor_name = "AWS"
      }
    }
    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name               = "IpReputation-${var.environment}"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "RateBasedRule"
    priority = 30
    action {
      block {}
    }
    statement {
      rate_based_statement {
        limit              = 2000
        aggregate_key_type = "IP"
      }
    }
    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name               = "RateBasedRule-${var.environment}"
      sampled_requests_enabled   = true
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name               = "sravas-waf-${var.environment}"
    sampled_requests_enabled   = true
  }

  tags = merge({
    Name        = "sravas-waf-${var.environment}"
    Environment = var.environment
  }, var.tags)
}

# ── CloudFront Distribution ──────────────────────────────────────

resource "aws_cloudfront_distribution" "main" {
  enabled             = var.enabled
  is_ipv6_enabled     = true
  default_root_object = ""
  price_class         = var.price_class
  web_acl_id          = aws_wafv2_web_acl.main.arn

  aliases = []
  comment = "sravas-${var.environment}"

  logging_config {
    bucket          = "sravas-logs-${var.environment}.s3.amazonaws.com"
    prefix          = "cloudfront/"
    include_cookies = false
  }

  # ── S3 origin (video segments) ────────────────────────────────
  origin {
    domain_name = var.s3_bucket_regional_domain_name
    origin_id   = local.s3_origin_id

    s3_origin_config {
      origin_access_identity = var.cloudfront_oai_path
    }
  }

  # ── ALB origin (frontend, optional) ───────────────────────────
  dynamic "origin" {
    for_each = var.alb_domain_name != null ? [1] : []
    content {
      domain_name = var.alb_domain_name
      origin_id   = local.alb_origin_id

      custom_origin_config {
        http_port              = 80
        https_port             = 443
        origin_protocol_policy = "https-only"
        origin_ssl_protocols   = ["TLSv1.2"]
      }
    }
  }

  # ── Default behavior ──────────────────────────────────────────
  default_cache_behavior {
    target_origin_id       = var.alb_domain_name != null ? local.alb_origin_id : local.s3_origin_id
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = var.alb_domain_name != null ? ["GET", "HEAD", "OPTIONS"] : ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true
    min_ttl                = var.min_ttl
    default_ttl            = var.default_ttl
    max_ttl                = var.max_ttl

    forwarded_values {
      query_string = var.alb_domain_name != null ? true : false
      cookies {
        forward = var.alb_domain_name != null ? "all" : "none"
      }
    }

    lambda_function_association {
      event_type   = "origin-request"
      lambda_arn   = ""  # TODO: add URL rewrite function in Week 5
      include_body = false
    }
  }

  # ── S3 /videos/* behavior (only when ALB is present) ─────────
  dynamic "ordered_cache_behavior" {
    for_each = var.alb_domain_name != null ? [1] : []
    content {
      path_pattern           = "/videos/*"
      target_origin_id       = local.s3_origin_id
      viewer_protocol_policy = "redirect-to-https"
      allowed_methods        = ["GET", "HEAD"]
      cached_methods         = ["GET", "HEAD"]
      compress               = true
      min_ttl                = var.min_ttl
      default_ttl            = var.default_ttl
      max_ttl                = var.max_ttl

      forwarded_values {
        query_string = false
        cookies {
          forward = "none"
        }
      }
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = merge({
    Name        = "sravas-cloudfront-${var.environment}"
    Environment = var.environment
  }, var.tags)
}
