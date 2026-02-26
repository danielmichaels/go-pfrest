package pfrest

//go:generate go run ./tools/specprep specs/v2.7/openapi.json specs/v2.7/openapi-processed.json
//go:generate oapi-codegen --config api/oapi-types.cfg.yaml specs/v2.7/openapi-processed.json
//go:generate oapi-codegen --config api/oapi-client.cfg.yaml specs/v2.7/openapi-processed.json
