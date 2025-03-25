package main

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// Improved Leaky Bucket Rate Limiter
type LeakyBucket struct {
	capacity     int64
	leakRate     time.Duration
	current      int64
	lastLeakTime int64
}

func NewLeakyBucket(capacity int64, leakRate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity:     capacity,
		leakRate:     leakRate,
		current:      0,
		lastLeakTime: time.Now().UnixNano(),
	}
}

func (lb *LeakyBucket) Allow() bool {
	now := time.Now().UnixNano()
	elapsed := time.Duration(now - atomic.LoadInt64(&lb.lastLeakTime))
	
	// Atomic operations for thread-safety
	leaked := int64(elapsed / lb.leakRate)
	current := atomic.LoadInt64(&lb.current)
	
	newCurrent := max(0, current - leaked)
	atomic.StoreInt64(&lb.current, newCurrent)
	atomic.StoreInt64(&lb.lastLeakTime, now)

	if newCurrent < lb.capacity {
		atomic.AddInt64(&lb.current, 1)
		return true
	}
	return false
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

var leakybucket *LeakyBucket

func init() {
	leakybucket = NewLeakyBucket(10, 200*time.Millisecond)
}

func rateLimiterMiddleware(c *fiber.Ctx) error {
	if !leakybucket.Allow() {
		return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
	}
	return c.Next()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	var api = os.Getenv("APIKEY")

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "POST,GET",
	}))
	app.Use(rateLimiterMiddleware)
	slangs := fmt.Sprint(Decodefile())
	
	app.Get("/html", func(c *fiber.Ctx) error {
		component := Show()
		c.Set("Content-Type", "text/html")
		return component.Render(c.Context(), c)
	})

	app.Post("/add", func(c *fiber.Ctx) error {
		input := c.Query("input")
		Addtext(input)
		return c.SendString("Added to the list")
	})

	app.Get("/list-slangs", func(c *fiber.Ctx) error {
		slangs := Decodefile()
		return c.JSON(slangs)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		input := c.Query("input")

		about := `
	You are a slang detection AI. Your only task is to detect slang words in a given text based on the provided slang list.

✅ Rules (DOs):
✔ Check the text against the slang list provided in the context.
✔ If the text contains any slang from the list, respond with only the word:
flag
✔ If the text does not contain any slang, return the text exactly as it is.
✔ Keep the original script and language unchanged.

❌ Restrictions (DON'Ts):
❌ Do not return anything other than flag when slang is found.
❌ Do not modify, replace, or translate any words.
❌ Do not explain, censor, or change the sentence structure.

📝 Example Inputs & Outputs
✅ No Slang Found → Return Text As-Is
Input:
"hello bhai kesa hai tu"
Output:
"hello bhai kesa hai tu"

✅ Slang Found → Return Only flag
Input:
"Bhai tu ekdum haramjade jaisa behave kar raha hai."
Output:
flag

✅ Slang Found → Return Only flag
Input:
"kutte mera peer daba ke de."
Output:
flag

✅ Slang Found → Return Only flag
Input:
"Tumi ekdom fokora manuh!"
Output:
flag

🔹 Slang List (To Be Provided in Context)
(Example:)
` + slangs
		response, err := Verify(api, input, about)
		if err != nil {
			fmt.Println(err)
		}
		if strings.Contains(strings.ToLower(response), "flag") {
			return c.SendString("")
		}
		return c.SendString(input)
	})

	fmt.Println("Server starting on port 8090")
	err = app.Listen("0.0.0.0:8090")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func Decodefile() []string {
	file, err := os.ReadFile("slang.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return []string{}
	}

	alltext := string(file)
	alltext = strings.ReplaceAll(alltext, "\r\n", "\n")
	decoded := strings.Split(alltext, "\n")

	// Remove trailing empty strings
	for len(decoded) > 0 && decoded[len(decoded)-1] == "" {
		decoded = decoded[:len(decoded)-1]
	}

	return decoded
}

func Addtext(input string) error {
	file, err := os.OpenFile("slang.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()
	
	input = fmt.Sprintln(input)
	_, err = file.WriteString(input)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}

func Verify(api, input, about string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(api))
	if err != nil {
		return "", fmt.Errorf("error creating Gemini client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(about)},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(input))
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(text), nil
		}
	}
	return "", fmt.Errorf("unexpected response format from Gemini")
}
