# FecharChats

A Go application that automatically closes old chats in the Huggy platform by fetching users from a database, retrieving their chats via API, and logging the closure operations.

## Features

- 🔄 **Automated Chat Closure**: Automatically closes old chats based on time criteria
- 📊 **Database Integration**: Fetches users from PostgreSQL database
- 🔌 **API Integration**: Integrates with Huggy API for chat management
- 📝 **Comprehensive Logging**: Detailed logging for all operations and errors
- 🛡️ **Error Handling**: Robust error handling with retry mechanisms
- 🐳 **Docker Support**: Containerized application with multi-stage build

## Architecture

```
FecharChats/
├── cmd/
│   └── main.go          # Main application entry point
├── internal/
│   ├── api/
│   │   └── api.go       # Huggy API integration
│   ├── config/
│   │   └── config.go    # Environment configuration
│   └── database/
│       ├── connect.go   # Database connection
│       ├── fetch.go     # User fetching
│       └── insert.go    # Log insertion
├── Dockerfile           # Docker configuration
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
└── README.md           # This file
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Huggy API access
- Docker (optional)

## Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# Database Configuration
DB_URL=postgresql://username:password@localhost:5432/database_name

# Huggy API Configuration
API_KEY=your_huggy_api_key_here
```

### Database Schema

The application expects the following database tables:

#### `usuarios` table
```sql
CREATE TABLE usuarios (
    userid VARCHAR(255) PRIMARY KEY,
    cargo VARCHAR(255)
);
```

#### `logs_fechamento_chamados` table
```sql
CREATE TABLE logs_fechamento_chamados (
    id SERIAL PRIMARY KEY,
    chatid INTEGER NOT NULL,
    userid INTEGER NOT NULL,
    tabulation_id INTEGER NOT NULL,
    chat_last_message TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Installation

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd FecharChats
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run cmd/main.go
   ```

### Docker Deployment

1. **Build the Docker image**
   ```bash
   docker build -t fechar-chats .
   ```

2. **Run the container**
   ```bash
   docker run -d \
     --name fechar-chats \
     --env-file .env \
     fechar-chats
   ```

## Usage

The application runs continuously and performs the following operations:

1. **Database Connection**: Connects to PostgreSQL database
2. **User Fetching**: Retrieves users with `cargo = 'vendedor'`
3. **Chat Processing**: For each user:
   - Fetches chats from Huggy API
   - Filters chats older than 3 days
   - Logs chat data to database
   - Attempts to close chats via API
4. **Logging**: Comprehensive logging of all operations
5. **Retry Logic**: Handles errors with 5-minute retry intervals

### Log Output

The application provides detailed logging:

```
2025/08/01 15:50:01 === Starting FecharChats Application ===
2025/08/01 15:50:01 Loading environment configuration...
2025/08/01 15:50:01 SUCCESS: .env file loaded successfully
2025/08/01 15:50:01 SUCCESS: DB_URL found (length: 45 characters)
2025/08/01 15:50:01 SUCCESS: API_KEY found (length: 32 characters)
2025/08/01 15:50:01 --- Starting new iteration ---
2025/08/01 15:50:01 Attempting to connect to database...
2025/08/01 15:50:01 SUCCESS: Database connection established
2025/08/01 15:50:01 SUCCESS: Retrieved 5 users from database
2025/08/01 15:50:01 Processing user 1/5: 157100
2025/08/01 15:50:01 Starting to fetch chats for user ID: 157100
2025/08/01 15:50:01 SUCCESS: Completed fetching chats for user 157100
2025/08/01 15:50:01 Found 116 chats for user 157100
```

## API Integration

### Huggy API Endpoints

- **GET** `/v3/chats?agent={id}&situation=in_chat&page={page}` - Fetch user chats
- **PUT** `/v3/chats/{id}/close` - Close a chat

### Chat Filtering

Chats are filtered based on:
- **Time criteria**: Only chats older than 3 days are processed
- **Status**: Only chats with `situation=in_chat` are fetched
- **User role**: Only users with `cargo='vendedor'` are processed

## Error Handling

The application includes comprehensive error handling:

- **Database connection failures**: 5-minute retry intervals
- **API request failures**: Logged but don't stop processing
- **Invalid data**: Graceful handling of null/empty values
- **Network issues**: Automatic retry mechanisms

## Security Features

- **Non-root Docker container**: Runs as unprivileged user
- **Masked logging**: Sensitive data is masked in logs
- **Environment-based configuration**: No hardcoded secrets
- **Input validation**: All inputs are validated before processing

## Monitoring

### Key Metrics to Monitor

- **Database connection success rate**
- **API request success rate**
- **Chat processing throughput**
- **Error rates by type**
- **Processing time per iteration**

### Log Analysis

The application logs can be analyzed for:
- Performance metrics
- Error patterns
- Processing efficiency
- API response times

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check `DB_URL` environment variable
   - Verify database is running and accessible
   - Check network connectivity

2. **API Authentication Failed**
   - Verify `API_KEY` is correct
   - Check API key permissions
   - Ensure API endpoint is accessible

3. **Chat Processing Errors**
   - Check API response status codes
   - Verify chat data format
   - Review tabulation ID parsing

### Debug Mode

For detailed debugging, check the logs for:
- `ERROR:` messages for failures
- `WARNING:` messages for potential issues
- `SUCCESS:` messages for completed operations

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license information here]

## Support

For issues and questions:
- Check the logs for error details
- Review the configuration
- Verify API and database connectivity
- Contact the development team 