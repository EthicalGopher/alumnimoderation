package main

import (
	"fmt"
	"os"

	"strings"
// "github.com/joho/godotenv" 
	"github.com/EthicalGopher/rag/groq"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const model = "deepseek-r1-distill-llama-70b"

func main() {
	// err:=godotenv.Load()
	// if err!=nil{
	// 	fmt.Println(err)
	// }
	// var api = os.Getenv("APIKEY")
	api := "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"
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "POST,GET",
	}))
	slangs:=fmt.Sprint(Decodefile())
	app.Post("/add",func(c*fiber.Ctx)error{
		input:=c.Query("input")
		Addtext(input)

		return c.SendString("Added to the list")
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
`+slangs
		response := groq.Ragfromgroq(api, input, about, model)
		if strings.Contains(strings.ToLower(response), "flag") {
			return c.SendString("")
		}
		return c.SendString(input)

	})
	
	app.Listen("0.0.0.0:8090")
}




func Decodefile() []string {
	file, err := os.ReadFile("slang.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		
	}

	alltext := string(file)

	alltext = strings.ReplaceAll(alltext, "\r\n", "\n")

	decoded := strings.Split(alltext, "\n")

	for len(decoded) > 0 && decoded[len(decoded)-1] == "" {
		decoded = decoded[:len(decoded)-1]
	}

	return decoded
}


func Addtext(input string)error{
	file,err:=os.OpenFile("slang.txt",os.O_APPEND|os.O_WRONLY,0644)
	if err!=nil{
		fmt.Println(err)
		return err
	}
	defer file.Close()
	input=fmt.Sprintln(input)
	_,err=file.WriteString(input)
	if err!=nil{
		fmt.Println(err)
		return err
	}
	return nil
}
