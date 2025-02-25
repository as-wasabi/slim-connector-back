package OpenAI

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"net/http"
	"os"
	"slim-connector-back/model"
	"time"
)

type ExtractedInfo struct {
	Start   time.Time `json:"start" bson:"start"`
	End     time.Time `json:"end" bson:"end"`
	Context string    `json:"context" bson:"context"`
}

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func CreateNewClient() (*openai.Client, error) {
	ApiKey, found := os.LookupEnv("OPENAI_API_KEY")
	if !found {
		return nil, fmt.Errorf("API key not found in environment variables")
	}
	client := openai.NewClient(
		option.WithAPIKey(ApiKey),
	)
	return client, nil
}

func GetAIResponse(client *openai.Client, userPrompt string) (*openai.ChatCompletion, error) {
	// UTCに変える...?
	systemPrompt := `あなたはスケジュール管理AIです。
[開始日時]「終了日時」は次のレイアウトに従ってください。
layout = "2006-01-02T15:04:05Z07:00"
ユーザーの入力から予定の「開始日時」「終了日時」「コンテキスト」を抽出し、以下の JSON 形式で出力してください。
{
  "start": "layout",
  "end": "layout",
  "context": "予定の内容"
}

見つからない場合は "start" または "end" を null にしてください。
{}の中身以外の内容は出力しないようにしてください
今日の日付は2025年2月18日です。
`
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})

	return chatCompletion, err

}

func ExtractedAIResp(chatCompletion *openai.ChatCompletion) (ExtractedInfo, error) {
	var extracted ExtractedInfo
	rawResponse := chatCompletion.Choices[0].Message.Content

	err := json.Unmarshal([]byte(rawResponse), &extracted)

	return extracted, err
}

func (h *OpenAIHandler) ExtractedTask(c *gin.Context) {
	err := LoadEnv()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Not Found .ENV FILE": "error"})
	}
	client, err := CreateNewClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//userPrompt := c.Param("prompt")
	userPrompt := ""

	chatCompletion, err := GetAIResponse(client, userPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to GET AI response"})
	}

	extracted, err := ExtractedAIResp(chatCompletion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
		//println(err.Error())
		return
	}

	taskData := &model.Task{
		Start:   extracted.Start,
		End:     extracted.End,
		Context: extracted.Context,
	}

	err = h.TaskHandler.CreateTaskFromAIResponse(taskData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}
