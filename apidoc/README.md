# Project: gitsast

## End-point: add repository

### Method: POST

> ```
> http://127.0.0.1:8080/api/v1/repository
> ```

### Body (**raw**)

```json
{
  "name": "",
  "remote_url": "https://github.com/mdlayher/wireguard_exporter"
}
```

### Response: 200

```json
{
  "id": "bb42beea-4118-4885-8f88-f291aa0fe790",
  "name": "atque",
  "remote_url": "https://github.com/mdlayher/wireguard_exporter.git",
  "created_at": "2023-02-22T04:59:45.723021Z",
  "updated_at": "2023-02-22T04:59:45.723021Z"
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: list repository

### Method: GET

> ```
> http://127.0.0.1:8080/api/v1/repository
> ```

### Body (**raw**)

```json

```

### Query Params

| Param      | value                                             |
| ---------- | ------------------------------------------------- |
| name       | quo                                               |
| remote_url | github.com/Alexandrine_Waters/dolor-aliquam-minus |
| limit      | 2                                                 |
| offset     | 1                                                 |

### Response: 200

```json
{
  "repositories": [
    {
      "id": "7e898d91-0cf8-4c03-a7ec-67f2f0a0f034",
      "name": "voluptatibus",
      "remote_url": "https://github.com/mdlayher/wireguard_exporter.git",
      "created_at": "2023-02-22T04:58:59.794583Z",
      "updated_at": "2023-02-22T04:58:59.794583Z"
    },
    {
      "id": "235025e4-73d7-4294-8ecd-06d9a71d2015",
      "name": "expedita",
      "remote_url": "https://github.com/mdlayher/wireguard_exporter.git",
      "created_at": "2023-02-22T04:59:31.941957Z",
      "updated_at": "2023-02-22T04:59:31.941957Z"
    },
    {
      "id": "bb42beea-4118-4885-8f88-f291aa0fe790",
      "name": "atque",
      "remote_url": "https://github.com/mdlayher/wireguard_exporter.git",
      "created_at": "2023-02-22T04:59:45.723021Z",
      "updated_at": "2023-02-22T04:59:45.723021Z"
    }
  ],
  "total": 3
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get repository by ID

### Method: GET

> ```
> http://127.0.0.1:8080/api/v1/repository/:id
> ```

### Response: 200

```json
{
  "id": "bb42beea-4118-4885-8f88-f291aa0fe790",
  "name": "atque",
  "remote_url": "https://github.com/mdlayher/wireguard_exporter.git",
  "created_at": "2023-02-22T04:59:45.723021Z",
  "updated_at": "2023-02-22T04:59:45.723021Z"
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: remove repository by ID

### Method: DELETE

> ```
> http://127.0.0.1:8080/api/v1/repository/:id
> ```

### Body (**raw**)

```json

```

### Response: 200

```json
null
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: update repository by ID

### Method: PUT

> ```
> http://127.0.0.1:8080/api/v1/repository/:id
> ```

### Body (**raw**)

```json
{
  "name": "",
  "remote_url": "github.com//"
}
```

### Response: 200

```json
null
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: trigger repository scan by ID

### Method: POST

> ```
> http://127.0.0.1:8080/api/v1/repository/:id/scan
> ```

### Body (**raw**)

```json

```

### Response: 200

```json
{
  "id": "62b8b624-4588-4f83-8d61-87bc3a560276",
  "repository_id": "bb42beea-4118-4885-8f88-f291aa0fe790",
  "status": "enqueued",
  "created_at": "2023-02-22T11:59:58.935452+07:00",
  "updated_at": "2023-02-22T11:59:58.935452+07:00",
  "enqueue_at": "2023-02-22T11:59:58.938824+07:00",
  "started_at": "0001-01-01T00:00:00Z",
  "finished_at": "0001-01-01T00:00:00Z"
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: get report by repository ID

### Method: GET

> ```
> http://127.0.0.1:8080/api/v1/repository/:id/report
> ```

### Response: 200

```json
{
  "id": "78589fe4-8365-4c86-9b78-3401ab7f35a3",
  "repository_id": "1336634e-73f4-438a-ab3a-152d52f071d9",
  "status": "success",
  "created_at": "2023-02-23T15:14:05.673514Z",
  "updated_at": "2023-02-23T15:14:07.121661Z",
  "enqueue_at": "2023-02-23T15:14:05.6767Z",
  "started_at": "2023-02-23T15:14:05.691266Z",
  "finished_at": "2023-02-23T15:14:07.121661Z",
  "findings": [
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse.go",
        "position": {
          "begin": {
            "line": 13
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 66
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 67
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 68
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 69
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 70
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 71
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 72
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 73
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 74
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 75
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 76
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 77
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 78
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 37
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 49
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 32
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 39
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 47
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 51
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 32
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 39
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 47
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    },
    {
      "type": "sast",
      "ruleId": "G001",
      "location": {
        "path": "/parse_test.go",
        "position": {
          "begin": {
            "line": 51
          }
        }
      },
      "metadata": {
        "description": "A secret starts with the prefix public_key",
        "severity": "LOW"
      }
    }
  ]
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: health check

### Method: GET

> ```
> http://127.0.0.1:8080/health
> ```

### Response: 200

```json
{
  "status": "ok"
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
