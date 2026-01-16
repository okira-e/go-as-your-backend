env "env" {
  src = "file://schema.hcl"
  
  dev = "docker://postgres/15/dev?search_path=public"
  
  migration {
    dir = "file://migrations"
  }
}