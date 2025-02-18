package OpenAI

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"net/http"
	"os"
	"time"
)

type ExtractedInfo struct {
	Start   time.Time `json:"start" bson:"start"`
	End     time.Time `json:"end" bson:"end"`
	Context string    `json:"context" bson:"context"`
}

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
			openai.UserMessage("明日10時から12時まで新宿で会議がある"),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})

	var extracted ExtractedInfo
	rawResponse := chatCompletion.Choices[0].Message.Content

	err = json.Unmarshal([]byte(rawResponse), &extracted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
		println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"start":   extracted.Start,
		"end":     extracted.End,
		"context": extracted.Context})
}
