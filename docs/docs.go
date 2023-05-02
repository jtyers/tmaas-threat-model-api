// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "ThreatPlane",
            "url": "http://www.threatplane.io",
            "email": "support@threatplane.io"
        },
        "license": {
            "name": "Commercial licence",
            "url": "https://threatplane.io/terms"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/threatModel": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Get all Threat Models",
                "operationId": "get-all-threat-models",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.ThreatModel"
                            }
                        }
                    }
                }
            },
            "put": {
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new Threat Model",
                "operationId": "create-threat-model",
                "parameters": [
                    {
                        "description": "todo data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ThreatModel"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.ThreatModel"
                        }
                    },
                    "400": {
                        "description": "Nope",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "patch": {
                "produces": [
                    "application/json"
                ],
                "summary": "Update a new Threat Model",
                "operationId": "update-threat-model",
                "parameters": [
                    {
                        "description": "todo data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ThreatModel"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.ThreatModel"
                        }
                    },
                    "400": {
                        "description": "Nope",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/threatModel/{id}": {
            "delete": {
                "produces": [
                    "application/json"
                ],
                "summary": "delete a todo item by ID",
                "operationId": "delete-todo-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "todo ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.ThreatModel"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/{threatModelID}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Retrieves a threatModel by ID.",
                "operationId": "get-threat-model-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "threat model ID",
                        "name": "threatModelID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.ThreatModel"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.ThreatModel": {
            "type": "object",
            "properties": {
                "dfdID": {
                    "description": "The DFD of the threat model.",
                    "type": "string"
                },
                "id": {
                    "description": "A randomised string e.g. th_abcdef123456",
                    "type": "string"
                },
                "title": {
                    "description": "The title of the threat model, must be at least two characters long. Required.",
                    "type": "string",
                    "minLength": 2
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "ThreatPlane Threat Model API",
	Description:      "The API used to interact with Threat Models",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
