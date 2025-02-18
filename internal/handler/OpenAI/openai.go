package OpenAI

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"net/http"
	"os"
)

func (h *OpenAIHandler) FetchApi(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not Found .ENV FILE": "error"})
	}

	ApiKey, found := os.LookupEnv("OPENAI_API_KEY")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"Not Found API KEY": "error"})
	}

	client := openai.NewClient(
		option.WithAPIKey(ApiKey),
	)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Say this is a test"),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})

	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
}
