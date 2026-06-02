terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "cope-hope-terraform-state"
    key            = "api/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}

provider "aws" {
  region = "us-east-1"
}

variable "owm_api_key" {
  description = "OpenWeatherMap API Key"
  type        = string
  sensitive   = true
}
