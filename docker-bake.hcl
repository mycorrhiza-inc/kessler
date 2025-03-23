group "default" {
  targets = ["frontend", "backend"]
  }

target "frontend" {
    context = "./frontend"
    dockerfile = "Dockerfile"
  }
