# Likha Configuration Examples

This directory contains various configuration examples demonstrating the capabilities of Likha. Each example shows different generator types and output formats.

## Available Examples

### 1. Basic User Data (`basic_users.yaml`)
Simple user profile generation with basic field types.

**Configuration Features:**
- Sequential ID generation
- Name selection from predefined list
- Email composition using expressions
- Registration timestamp generation
- Status selection

**Sample Output (JSON):**
```json
{
  "id": 1001,
  "name": "Alice",
  "email": "alice-42@example.com",
  "registered_at": "2022-08-15T14:30:22Z",
  "status": "active"
}
```

### 2. E-commerce Products (`ecommerce_products.yaml`)
Product catalog generation with pricing and categories.

**Configuration Features:**
- Product name from predefined list
- Price range generation
- Category-based discount mapping using foreign keys
- Stock quantity generation
- Simple boolean availability generation

**Sample Output (JSON):**
```json
{
  "product_id": "PROD-7842",
  "name": "Wireless Headphones",
  "price": 129.99,
  "category": "Electronics",
  "stock_quantity": 45,
  "discount_percentage": 12.5,
  "shipping_weight": 0.85
}
```

### 3. Log File Simulation (`server_logs.yaml`)
Server log entry generation with realistic timestamps and status codes.

**Configuration Features:**
- Timestamp generation within specific range
- IP address generation
- HTTP method and status code correlation
- Response time simulation
- User agent string generation

**Sample Output (JSON):**
```json
{
  "timestamp": "2023-10-15T09:23:45Z",
  "ip_address": "192.168.1.45",
  "method": "GET",
  "endpoint": "/api/users",
  "status_code": 200,
  "response_time_ms": 145,
  "user_agent": "Mozilla/5.0 (compatible; bot/1.0)"
}
```

### 4. Financial Transactions (`financial_transactions.yaml`)
Banking transaction simulation with account relationships.

**Configuration Features:**
- Account number generation
- Transaction type-dependent amount ranges
- Currency code selection
- Balance calculation expressions
- Merchant name generation

**Sample Output (JSON):**
```json
{
  "transaction_id": "TXN-9876543210",
  "account_number": "ACC-4532",
  "transaction_type": "debit",
  "amount": -45.67,
  "currency": "USD",
  "merchant": "Coffee Shop Downtown",
  "timestamp": "2023-09-20T16:42:13Z",
  "balance_after": 1234.56
}
```

### 5. IoT Sensor Data (`iot_sensors.yaml`)
Internet of Things sensor reading generation.

**Configuration Features:**
- Device ID generation
- Sensor type-specific value ranges
- Timestamp generation with intervals
- Location coordinate generation
- Status monitoring

**Sample Output (JSON):**
```json
{
  "device_id": "SENSOR-A001",
  "sensor_type": "temperature",
  "value": 23.5,
  "unit": "celsius",
  "timestamp": "2023-11-01T12:15:30Z",
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060
  },
  "status": "online"
}
```

### 6. Social Media Posts (`social_posts.yaml`)
Social media post generation with user interactions.

**Configuration Features:**
- User handle generation
- Post content from templates
- Like/share count generation
- Hashtag insertion
- Comment relationship modeling

**Sample Output (JSON):**
```json
{
  "post_id": "POST-1234567890",
  "user_handle": "@user_alice_42",
  "content": "Just had an amazing coffee! #MorningFuel #CoffeeLovers",
  "timestamp": "2023-10-25T08:30:15Z",
  "likes": 23,
  "shares": 5,
  "comments": 7,
  "hashtags": ["MorningFuel", "CoffeeLovers"]
}
```

### 7. Complex CSV with Headers (`complex_csv.yaml`)
Demonstrates CSV-specific configuration options.

**Configuration Features:**
- Custom delimiter configuration
- Header row inclusion
- Null value handling
- Date formatting for CSV
- Numeric precision control

**Sample Output (CSV):**
```csv
id,name,email,signup_date,last_login,status
1001,"John Doe","john.doe@email.com","2023-01-15","2023-10-20T14:30:00Z","active"
```

### 8. XML Document Generation (`xml_records.yaml`)
XML format with nested structure and custom root node.

**Configuration Features:**
- Custom root node configuration
- Nested element generation
- Attribute value generation
- XML-safe string generation
- Namespace handling

**Sample Output (XML):**
```xml
<users>
  <user>
    <id>1001</id>
    <profile>
      <name>Alice Johnson</name>
      <email>alice.johnson@example.com</email>
      <created_at>2023-05-10T09:15:30Z</created_at>
    </profile>
    <status>active</status>
  </user>
</users>
```

### 9. YAML Configuration (`yaml_config.yaml`)
YAML output with proper indentation and structure.

**Configuration Features:**
- Custom indentation settings
- Nested object generation
- Array field generation
- Multi-line string handling
- Boolean and null value representation

**Sample Output (YAML):**
```yaml
user:
  id: 1001
  name: "Alice Johnson"
  email: "alice.johnson@example.com"
  preferences:
    theme: "dark"
    notifications: true
    languages: ["en", "es"]
  created_at: "2023-05-10T09:15:30Z"
```

### 10. Custom Generator Example (`simple_custom_generator.yaml`)
Demonstrates integration with external binary generators.

**Configuration Features:**
- External binary execution

**Sample Output (JSON):**
```json
{
  "custom_execution": "hello",
  "id": 2201,
  "username": "user_bK2XnA7V"
}
```

## Running Examples

To run any example:

```bash
# Generate 1000 records using basic users example
likha --config examples/basic_users.yaml --count 1k

# Generate 100k records using e-commerce products example
likha --config examples/ecommerce_products.yaml --count 100k

# Generate 1 million log entries
likha --config examples/server_logs.yaml --count 1m
```

## Customization Tips

1. **Field Dependencies**: Use foreign key generators to create realistic relationships between fields
2. **Expression Power**: Combine multiple builtin functions in expressions for complex data patterns
3. **Performance Tuning**: For large datasets, consider simpler generators over complex custom ones
4. **Format-Specific Settings**: Leverage format-specific settings for optimal output structure
5. **Validation**: Test configurations with small counts before generating large datasets

## Contributing Examples

When adding new examples:
1. Include comprehensive comments in the YAML configuration
2. Provide sample output for clarity
3. Document any special requirements or dependencies
4. Test with various count sizes to ensure performance
5. Add explanation to this README following the existing format
