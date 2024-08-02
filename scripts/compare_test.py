#!/usr/bin/python3
"""
    Script that checks the compatibility of all APIs for given endpoints with multiple base URLs.
"""

import requests
from requests.exceptions import RequestException
import uuid

BASE_URLS = [
    "https://deployment.api-csharp.boilerplate.hng.tech/api/v1",
    "https://deployment.api-golang.boilerplate.hng.tech/api/v1",
    "https://deployment.api-php.boilerplate.hng.tech/api/v1"
]
# Dynamic Endpoints and sample payloads
ENDPOINTS = {
    "register_csharp": {
        "url": "/auth/register",
        "method": "POST",
        "payload_template": {
            "firstName": "test",
            "lastName": "kim",
            "email": "test.kim@example.com",
            "password": "Password123",
            "phoneNumber": "123-456-7890"
        },
        "expected_status_code": 201
    },
    "register_golang": {
        "url": "/auth/register",
        "method": "POST",
        "payload_template": {
            "FirstName": "go",
            "LastName": "lang",  
            "email": "go.lang@example.com",
            "password": "Password123"
        },
        "expected_status_code": 201
    },
    "register_php": {
        "url": "/auth/register",
        "method": "POST",
        "payload_template": {
            "first_name": "Php",  # Use correct field names
            "last_name": "Diana",
            "email": "john.php@example.com",
            "password": "Janphp163!"  # Use a new, secure password
        },
        "expected_status_code": 201
    },
    "login": {
        "url": "/auth/login",
        "method": "POST",
        "templates": {
            "csharp": {
                "payload": {
                    "email": "test.kim@example.com",
                    "password": "Password123"
                },
                "expected_status_code": 200
            },
            "golang": {
                "payload": {
                    "email": "test.kim@example.com",
                    "password": "Password123"
                },
                "expected_status_code": 200
            },
            "php": {
                "payload": {
                    "email": "test.kim@example.com",
                    "password": "Password123"
                },
                "expected_status_code": 200
            }
        }
    },
    "products": {
        "url": "/products",
        "method": "POST",
        "templates": {
            "csharp": {
                "payload": {
                    "name": "product",
                    "description": "product",
                    "category": "category",
                    "price": 0.01
                },
                "expected_status_code": 201
            },
            "golang": {
                "payload": {
                    "name": "product",
                    "description": "product",
                    "category": "category",
                    "price": 0.01
                },
                "expected_status_code": 201
            },
            "php": {
                "payload": {
                   "name": "product",
                    "description": "product",
                    "category": "category",
                    "price": 0.01
                },
                "expected_status_code": 201
            }
        }
    },

 "organizations": {
        "url": "/products",
        "method": "POST",
        "templates": {
            "csharp": {
                "payload": {
                    "name": "string",
                    "description": "string",
                    "email": "string",
                    "industry": "string",
                    "type": "string",
                    "country": "string",
                    "address": "string",
                    "state": "string"
                },
                "expected_status_code": 201
            },
            "golang": {
                "payload": {
                    "name": "string",
                    "description": "string",
                    "email": "string",
                    "industry": "string",
                    "type": "string",
                    "country": "string",
                    "address": "string",
                    "state": "string"
                },
                "expected_status_code": 201
            },
            "php": {
                "payload": {
                    "name": "string",
                    "description": "string",
                    "email": "string",
                    "industry": "string",
                    "type": "string",
                    "country": "string",
                    "address": "string",
                    "state": "string"
                },
                "expected_status_code": 201
            }
        }
    },

    "subscriptions": {
        "url": "/subscriptions/free",
        "method": "POST",
        "templates": {
            "csharp": {
                "payload": {
                    "userId": "string",
                    "organizationId": "string"
                },
                "expected_status_code": 201
            },
            "golang": {
                "payload": {
                    "userId": "string",
                    "organizationId": "string"
                },
                "expected_status_code": 201
            },
            "php": {
                "payload": {
                    "userId": "string",
                    "organizationId": "string"
                },
                "expected_status_code": 201
            }
        }
    },

    "jobs": {
        "url": "/jobs",
        "method": "POST",
        "templates": {
            "csharp": {
                "payload": {
                    "title": "string",
                    "description": "string",
                    "location": "string",
                    "salary": 0,
                    "level": 0,
                    "company": "string"
                },
                "expected_status_code": 201
            },
            "golang": {
                "payload": {
                    "title": "string",
                    "salary": "5000-7000",
                    "job_type": "string",
                    "location": "string",
                    "deadline": "2024-12-31T23:59:59Z",
                    "work_mode": "string",
                    "experience": "string",
                    "how_to_apply": "string",
                    "job_benefits": "string",
                    "company_name": "string",
                    "description": "string",
                    "key_responsibilities": "string",
                    "qualifications": "string"
                },
                "expected_status_code": 201
            },
            "php": {
                "payload": {
                    "title": "string",
                    "description": "string",
                    "location": "string",
                    "salary": 0,
                    "level": 0,
                    "company": "string"
                },
                "expected_status_code": 201
            }
        }
    }

}

def is_logged_in(base_url, api_type):
    check_login_endpoint = ENDPOINTS["check_login"]
    url = f"{base_url}{check_login_endpoint['url']}"
    method = check_login_endpoint["method"]
    expected_status_code = check_login_endpoint["templates"][api_type]["expected_status_code"]
    
    try:
        if method == "GET":
            response = requests.get(url)

        status_code = response.status_code
        return status_code == expected_status_code
    except RequestException as e:
        print(f"Request to {url} failed: {e}")
        return False


def generate_unique_payload(template, api_type):
    """ Generate a unique payload for testing to avoid conflicts. """
    payload = template.copy()
    payload['email'] = f"{uuid.uuid4().hex[:8]}@example.com"
    return payload

def test_endpoint(base_url, endpoint_key, api_type):
    endpoint = ENDPOINTS.get(endpoint_key, {})
    url = f"{base_url}{endpoint.get('url', '')}"
    method = endpoint.get("method", "GET")
    
    # Generate a unique payload for registration endpoints
    payload_template = endpoint.get('payload_template', {})
    payload = generate_unique_payload(payload_template, api_type)
    
    expected_status_code = endpoint.get("expected_status_code", 200)
    try:
        if method == "POST":
            response = requests.post(url, json=payload)
        status_code = response.status_code
        # Handle non-JSON responses
        try:
            response_json = response.json()
        except ValueError:
            response_json = response.text
        print(f"Testing {method} {url}")
        print(f"Expected status code: {expected_status_code}, Got: {status_code}")
        if status_code == expected_status_code:
            print("Test passed!")
        else:
            print("Test failed.")
        print("Response:", response_json)
        print("=" * 50)
    except RequestException as e:
        print(f"Request to {url} failed: {e}")
        print("=" * 50)

if __name__ == "__main__":
    for base_url in BASE_URLS:
        for endpoint_key in ENDPOINTS:
            if "login" in endpoint_key:
                api_type = "csharp" if "csharp" in base_url else "golang" if "golang" in base_url else "php"
            elif "products":
                api_type = "csharp" if "csharp" in base_url else "golang" if "golang" in base_url else "php"
            elif "organizations":
                api_type = "csharp" if "csharp" in base_url else "golang" if "golang" in base_url else "php"
            elif "subscriptions":
                api_type = "csharp" if "csharp" in base_url else "golang" if "golang" in base_url else "php"
            elif "jobs":
                api_type = "csharp" if "csharp" in base_url else "golang" if "golang" in base_url else "php"
            else:
                api_type = "csharp" if "csharp" in base_url else "golang" if "golang" in base_url else "php"
                if "register" in endpoint_key:
                    endpoint_key = f"register_{api_type}"
            
            test_endpoint(base_url, endpoint_key, api_type)