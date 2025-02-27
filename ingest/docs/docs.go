// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/add-task/ingest": {
            "post": {
                "description": "Creates a new default ingestion task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Add Default Ingest Task",
                "parameters": [
                    {
                        "description": "Scraper information",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tasks.ScraperInfoPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tasks.KesslerTaskInfo"
                        }
                    },
                    "400": {
                        "description": "Error decoding request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Error adding task",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/add-task/ingest/nypuc": {
            "post": {
                "description": "Creates a new NYPUC-specific ingestion task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Add NYPUC Ingest Task",
                "parameters": [
                    {
                        "description": "NYPUC document information",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tasks.NYPUCDocInfo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tasks.KesslerTaskInfo"
                        }
                    },
                    "400": {
                        "description": "Error decoding request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Error adding task",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/task/{id}": {
            "get": {
                "description": "Retrieves information about a specific task by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Get Task Information",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tasks.KesslerTaskInfo"
                        }
                    },
                    "404": {
                        "description": "Error retrieving task info",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "501": {
                        "description": "Not implemented",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/version_hash": {
            "get": {
                "description": "Returns the current version hash of the application",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "system"
                ],
                "summary": "Get Version Hash",
                "responses": {
                    "200": {
                        "description": "Version hash string",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "tasks.KesslerTaskInfo": {
            "type": "object",
            "properties": {
                "queue": {
                    "type": "string"
                },
                "state": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "task_id": {
                    "type": "string"
                }
            }
        },
        "tasks.NYPUCDocInfo": {
            "type": "object",
            "properties": {
                "date_filed": {
                    "type": "string"
                },
                "docket_id": {
                    "type": "string"
                },
                "file_name": {
                    "type": "string"
                },
                "item_no": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "nypuc_doctype": {
                    "type": "string"
                },
                "organization": {
                    "type": "string"
                },
                "serial": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "tasks.ScraperInfoPayload": {
            "type": "object",
            "properties": {
                "author_individual": {
                    "type": "string"
                },
                "author_individual_email": {
                    "type": "string"
                },
                "author_organisation": {
                    "type": "string"
                },
                "docket_id": {
                    "type": "string"
                },
                "file_class": {
                    "type": "string"
                },
                "file_type": {
                    "type": "string"
                },
                "file_url": {
                    "type": "string"
                },
                "internal_source_name": {
                    "type": "string"
                },
                "item_number": {
                    "type": "string"
                },
                "lang": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "published_date": {
                    "type": "string"
                },
                "state": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "petstore.swagger.io",
	BasePath:         "/ingest_v1",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server Petstore server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
