// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/pairs": {
            "get": {
                "description": "Get paginated list of token pairs with sorting options",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "pairs"
                ],
                "summary": "List token pairs",
                "parameters": [
                    {
                        "minimum": 0,
                        "type": "integer",
                        "default": 0,
                        "description": "Pagination offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "maximum": 100,
                        "minimum": 1,
                        "type": "integer",
                        "default": 10,
                        "description": "Items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "asc",
                            "desc"
                        ],
                        "type": "string",
                        "default": "asc",
                        "description": "Sort order",
                        "name": "orderBy",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful operation",
                        "schema": {
                            "$ref": "#/definitions/internal_usecases_keeper.PairsResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid parameters",
                        "schema": {
                            "$ref": "#/definitions/github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Register new token pair in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "pairs"
                ],
                "summary": "Create token pair",
                "parameters": [
                    {
                        "description": "Token pair data",
                        "name": "pair",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_usecases_keeper.NewPairRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully created",
                        "schema": {
                            "$ref": "#/definitions/internal_usecases_keeper.Pair"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "$ref": "#/definitions/github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Token not found",
                        "schema": {
                            "$ref": "#/definitions/github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/tokens": {
            "get": {
                "description": "Get paginated list of tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tokens"
                ],
                "summary": "Get tokens list",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Offset for pagination",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Limit for pagination",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "asc",
                            "desc"
                        ],
                        "type": "string",
                        "default": "asc",
                        "description": "Order by field",
                        "name": "orderBy",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/internal_usecases_keeper.TokensResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a new token to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tokens"
                ],
                "summary": "Create new token",
                "parameters": [
                    {
                        "description": "Token data",
                        "name": "token",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_usecases_keeper.NewTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_humangrass_price-keeper_pgk_x_xhttp.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "internal_usecases_keeper.NewPairRequest": {
            "type": "object",
            "required": [
                "denominator",
                "numerator",
                "timeframe"
            ],
            "properties": {
                "denominator": {
                    "type": "string",
                    "maxLength": 10
                },
                "numerator": {
                    "type": "string",
                    "maxLength": 10
                },
                "timeframe": {
                    "maximum": 5,
                    "allOf": [
                        {
                            "$ref": "#/definitions/time.Duration"
                        }
                    ]
                }
            }
        },
        "internal_usecases_keeper.NewTokenRequest": {
            "type": "object",
            "required": [
                "name",
                "network",
                "network_id",
                "symbol"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 100
                },
                "network": {
                    "type": "string",
                    "maxLength": 100
                },
                "network_id": {
                    "type": "string",
                    "maxLength": 100
                },
                "symbol": {
                    "type": "string",
                    "maxLength": 10
                }
            }
        },
        "internal_usecases_keeper.Pair": {
            "type": "object",
            "properties": {
                "ticket": {
                    "type": "string"
                },
                "timeframe": {
                    "type": "string"
                }
            }
        },
        "internal_usecases_keeper.PairsResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal_usecases_keeper.Pair"
                    }
                },
                "page": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "internal_usecases_keeper.Token": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "network": {
                    "type": "string"
                },
                "network_id": {
                    "type": "string"
                },
                "symbol": {
                    "type": "string"
                }
            }
        },
        "internal_usecases_keeper.TokensResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal_usecases_keeper.Token"
                    }
                },
                "page": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "time.Duration": {
            "type": "integer",
            "enum": [
                -9223372036854775808,
                9223372036854775807,
                1,
                1000,
                1000000,
                1000000000,
                60000000000,
                3600000000000,
                -9223372036854775808,
                9223372036854775807,
                1,
                1000,
                1000000,
                1000000000,
                60000000000,
                3600000000000
            ],
            "x-enum-varnames": [
                "minDuration",
                "maxDuration",
                "Nanosecond",
                "Microsecond",
                "Millisecond",
                "Second",
                "Minute",
                "Hour",
                "minDuration",
                "maxDuration",
                "Nanosecond",
                "Microsecond",
                "Millisecond",
                "Second",
                "Minute",
                "Hour"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8888",
	BasePath:         "/api/",
	Schemes:          []string{},
	Title:            "Price Keeper API",
	Description:      "API for cryptocurrency tokens",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
