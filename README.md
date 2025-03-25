# 🚨 Slang Detector Application

## 📝 Overview

The Slang Detector is a powerful Go-based web application that leverages AI to detect and flag potentially offensive or inappropriate slang words in text. Built with Fiber and Google's Gemini AI, this tool provides a robust solution for content moderation and language filtering.

## ✨ Features

- 🤖 AI-powered slang detection
- 📝 Dynamic slang word management
- 🛡️ Rate limiting for API protection
- 🌐 RESTful API endpoints
- 🖥️ Simple web interface
- 🔒 Cross-Origin Resource Sharing (CORS) support

## 🚀 Prerequisites

- Go 1.20+
- Gemini API Key
- Internet connection

## 🛠️ Installation

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/slang-detector.git
cd slang-detector
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Set Up Environment
Create a `.env` file in the project root:
```
APIKEY=your_gemini_api_key_here
```

## 🔧 Configuration

### Slang List
- Create a `slang.txt` file in the project root
- Add slang words line by line
- Example:
  ```
  haramjade
  kutte
  fokora
  ```

### Rate Limiting
- Currently configured with:
  - Capacity: 10 requests
  - Leak Rate: 200 milliseconds
- Modify in `init()` function if needed

## 🌈 API Endpoints

### 1. Slang Check `GET /`
- Query Parameter: `input`
- Returns:
  - Original text if no slang detected
  - Empty response if slang found

### 2. Add Slang `POST /add`
- Query Parameter: `input`
- Adds new slang word to `slang.txt`

### 3. Web Interface `GET /html`
- Provides interactive UI for testing

## 🏃 Running the Application

```bash
go run main.go
```
- Server starts on `0.0.0.0:8090`

## 📋 API Usage Examples

### Check Text for Slang
```bash
curl "http://localhost:8090/?input=Bhai%20tu%20haramjade%20jaisa%20behave%20kar%20raha%20hai"
# Returns empty response (slang detected)

curl "http://localhost:8090/?input=Hello%20world"
# Returns "Hello world"
```

### Add Slang Word
```bash
curl -X POST "http://localhost:8090/add?input=newslangword"
```

## 🛡️ Rate Limiting

- Maximum of 10 requests per 200ms
- Excess requests receive a "Too many requests" error

## 📦 Dependencies

- Fiber: Web framework
- Godotenv: Environment variable management
- Google Generative AI: AI-powered text analysis

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit changes
4. Push to the branch
5. Create a Pull Request

## ⚠️ Limitations

- Requires active Gemini API connection
- Slang detection accuracy depends on AI model
- Limited to provided slang list

## 📄 License

[Your License Here - e.g., MIT]

## 🐛 Reporting Issues

Report issues on the GitHub repository's issue tracker.

## 🌟 Future Roadmap

- [ ] Machine learning model for dynamic slang detection
- [ ] Multi-language support
- [ ] Advanced rate limiting configurations
- [ ] Persistent slang word storage