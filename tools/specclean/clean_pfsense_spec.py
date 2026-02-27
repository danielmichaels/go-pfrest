#!/usr/bin/env python3
"""
pfSense OpenAPI Specification Cleaner for Fern

Cleans the pfSense OpenAPI specification to make it compatible with Fern code
generation. Addresses invalid enum values, duplicate GraphQL schemas, security
scheme format, and problematic escape sequences.

Usage:
    python clean_pfsense_spec.py input_spec.json output_spec.json
"""

import json
import sys
import re
from typing import Any, Dict

sys.setrecursionlimit(50000)


def fix_enum_values(obj: Any) -> Any:
    if isinstance(obj, dict):
        if "enum" in obj and isinstance(obj["enum"], list):
            fixed_enum = []
            for i, value in enumerate(obj["enum"]):
                if isinstance(value, str):
                    if value.isdigit():
                        fixed_enum.append(f"value_{value}")
                    elif value == "":
                        fixed_enum.append(f"empty_value_{i}")
                    elif value and value[0].isdigit():
                        fixed_enum.append(f"value_{value}")
                    elif value and value.startswith("-") and value[1:].isdigit():
                        fixed_enum.append(f"negative_{value[1:]}")
                    elif len(value) == 1 and not value.isalpha():
                        fixed_enum.append(f"char_{ord(value)}")
                    else:
                        fixed_enum.append(value)
                else:
                    fixed_enum.append(value)
            obj["enum"] = fixed_enum

        return {k: fix_enum_values(v) for k, v in obj.items()}
    elif isinstance(obj, list):
        return [fix_enum_values(item) for item in obj]
    return obj


def remove_graphql_duplicates(spec: Dict[str, Any]) -> Dict[str, Any]:
    if "components" not in spec or "schemas" not in spec["components"]:
        return spec

    schemas = spec["components"]["schemas"]
    graphql_patterns = [
        r"^__.*",
        r".*Query$",
        r".*Mutation$",
        r".*Subscription$",
        r"^GraphQL.*",
        r"^GraphQl.*",
    ]

    graphql_keys = []
    for schema_name in schemas.keys():
        for pattern in graphql_patterns:
            if re.match(pattern, schema_name):
                graphql_keys.append(schema_name)
                break

    for key in graphql_keys:
        del schemas[key]

    print(f"  Removed {len(graphql_keys)} GraphQL schemas")
    return spec


def fix_security_schemes(spec: Dict[str, Any]) -> Dict[str, Any]:
    if "security" in spec and isinstance(spec["security"], list):
        current_security = spec["security"]
        if len(current_security) == 1 and len(current_security[0]) > 1:
            new_security = []
            for auth_method in current_security[0]:
                new_security.append({auth_method: []})
            spec["security"] = new_security
            print(f"  Split combined auth into individual options: {list(current_security[0].keys())}")

    return spec


def fix_escape_sequences(obj: Any) -> Any:
    if isinstance(obj, dict):
        fixed_obj = {}
        for k, v in obj.items():
            if isinstance(v, str):
                v = v.replace("\\/", "/")
                v = v.replace("\\n", " ")
                v = v.replace("\\t", " ")
                v = v.replace("\\r", " ")
                fixed_obj[k] = v
            else:
                fixed_obj[k] = fix_escape_sequences(v)
        return fixed_obj
    elif isinstance(obj, list):
        return [fix_escape_sequences(item) for item in obj]
    elif isinstance(obj, str):
        obj = obj.replace("\\/", "/")
        obj = obj.replace("\\n", " ")
        obj = obj.replace("\\t", " ")
        obj = obj.replace("\\r", " ")
        return obj
    return obj


def fix_schema_references(spec: Dict[str, Any]) -> Dict[str, Any]:
    spec_str = json.dumps(spec)

    graphql_refs = [
        "#/components/schemas/GraphQL",
        "#/components/schemas/GraphQLResponse",
    ]

    for ref in graphql_refs:
        if ref in spec_str:
            print(f"  Replacing reference: {ref}")
            spec_str = spec_str.replace(f'"{ref}"', '"#/components/schemas/Success"')

    return json.loads(spec_str)


def clean_pfsense_spec(input_file: str, output_file: str) -> None:
    print(f"Loading spec from: {input_file}")

    with open(input_file, "r", encoding="utf-8") as f:
        spec = json.load(f)

    schema_count = len(spec.get("components", {}).get("schemas", {}))
    print(f"Original spec: {schema_count} schemas")

    steps = [
        ("Fixing enum values", fix_enum_values),
        ("Removing GraphQL duplicates", remove_graphql_duplicates),
        ("Fixing security schemes", fix_security_schemes),
        ("Fixing escape sequences", fix_escape_sequences),
        ("Fixing schema references", fix_schema_references),
    ]

    for i, (name, fn) in enumerate(steps, 1):
        print(f"Step {i}: {name}...")
        spec = fn(spec)

    schema_count = len(spec.get("components", {}).get("schemas", {}))
    print(f"Cleaned spec: {schema_count} schemas")

    with open(output_file, "w", encoding="utf-8") as f:
        json.dump(spec, f, indent=2, ensure_ascii=False)

    print(f"Saved to: {output_file}")


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python clean_pfsense_spec.py <input_file> <output_file>")
        sys.exit(1)

    clean_pfsense_spec(sys.argv[1], sys.argv[2])
