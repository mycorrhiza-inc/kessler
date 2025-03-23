
group "default" {
  targets = ["frontend", "backend"]
}

target "frontend" {
  context = "./frontend"
  dockerfile = "dev.Dockerfile"
  tags = ["kessler/frontend:latest"]
}


current_go_version = "1.24"
target "backend" {
  context = "./backend"
  dockerfile = "dev.server.Dockerfile"
  args = {
    GO_VERSION = current_go_version
  }
  tags = ["kessler/backend:latest"]
}

target "ingest" {
  context = "./backend"
  dockerfile = "dev.ingest.Dockerfile"
  args = {
    GO_VERSION = current_go_version
  }
  tags = ["kessler/ingest:latest"]
}

