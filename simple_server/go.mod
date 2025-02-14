module example.com/simple_server

go 1.22.0

require example.com/kv_store v0.0.0-00010101000000-000000000000

require github.com/google/uuid v1.6.0

replace example.com/kv_store => ../kv_store
