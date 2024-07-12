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
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/health/": {
            "get": {
                "description": "Health Check",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health Check"
                ],
                "summary": "Health Check",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    },
                    "400": {
                        "description": "Failed",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/members/": {
            "post": {
                "description": "Registration a member",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Members"
                ],
                "summary": "Registration a member",
                "parameters": [
                    {
                        "description": "member",
                        "name": "Request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.MemberCreate"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    },
                    "400": {
                        "description": "Failed",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    },
                    "409": {
                        "description": "Failed",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/members/login": {
            "post": {
                "description": "Login a member",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Members"
                ],
                "summary": "Login a member",
                "parameters": [
                    {
                        "description": "member",
                        "name": "Request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.MemberAuth"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    },
                    "400": {
                        "description": "Failed",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    },
                    "409": {
                        "description": "Failed",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/roles/create": {
            "post": {
                "description": "Создание роли",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Создание роли",
                "parameters": [
                    {
                        "description": "role",
                        "name": "Request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Role"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    },
                    "400": {
                        "description": "Failed",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    },
                    "409": {
                        "description": "Failed",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/roles/list": {
            "get": {
                "description": "Вывод всех ролей",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Вывод всех ролей",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.RoleList"
                            }
                        }
                    },
                    "400": {
                        "description": "Failed",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    },
                    "409": {
                        "description": "Failed",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/vehicles/": {
            "post": {
                "security": [
                    {
                        "AuthBearer": []
                    }
                ],
                "description": "Create a vehicle",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Vehicles"
                ],
                "summary": "Create a vehicle",
                "parameters": [
                    {
                        "description": "Create a vehicle model",
                        "name": "Request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateVehicleRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created response",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/components.BaseHttpResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "$ref": "#/definitions/dto.DBVehicleDTO"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/components.BaseHttpResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "components.BaseHttpResponse": {
            "type": "object",
            "properties": {
                "error": {},
                "result": {},
                "resultCode": {
                    "$ref": "#/definitions/components.ResultCode"
                },
                "success": {
                    "type": "boolean"
                },
                "validationErrors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/validations.ValidationError"
                    }
                }
            }
        },
        "components.ResultCode": {
            "type": "integer",
            "enum": [
                0,
                40001,
                40101,
                40301,
                40401,
                42901,
                442902,
                50001,
                50002
            ],
            "x-enum-varnames": [
                "Success",
                "ValidationError",
                "AuthError",
                "ForbiddenError",
                "NotFoundError",
                "LimiterError",
                "OtpLimiterError",
                "CustomRecovery",
                "InternalError"
            ]
        },
        "dto.CreateVehicleRequest": {
            "type": "object",
            "properties": {
                "location": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "dto.DBVehicleDTO": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "dto.MemberAuth": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@comecord.com"
                },
                "password": {
                    "type": "string",
                    "example": "calista78Batista"
                }
            }
        },
        "dto.MemberCreate": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 6,
                    "example": "user@comecord.com"
                },
                "password": {
                    "type": "string",
                    "minLength": 6,
                    "example": "calista78Batista"
                },
                "phone": {
                    "type": "string",
                    "example": "+7 (999) 999-99-99"
                }
            }
        },
        "dto.Role": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "example": "Admin"
                },
                "permissions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "dto.RoleList": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "permissions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "validations.ValidationError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "property": {
                    "type": "string"
                },
                "tag": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
