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
type promptRequest struct {
	Prompt string `json:"prompt"`
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
	today := time.Now().UTC().Format("YYYY-MM-DD")
	// UTCに変える...?
	systemPrompt := fmt.Sprintf(`あなたはスケジュール管理AIです。
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
今日の日付は %s です。
`, today)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
		//	ここで呼び出し
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

	var req promptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Json"})
		return
	}

	//userPrompt := "以下の条件で予定を作成してください。\n\n- 開始日: 2025年3月10日\n- 終了日: 2025年3月12日\n- " +
	//	"内容: 新しいタスク管理ツールの設計と実装\n\n" +
	//	"詳細:  \nこの期間中に、新しいタスク管理ツールの基本設計と初期実装を行います。" +
	//	"要件定義、データベース設計、UIワイヤーフレームの作成を含めて、効率的に進めるための計画を立ててください。\n"

	chatCompletion, err := GetAIResponse(client, req.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to GET AI response"})
	}

	extracted, err := ExtractedAIResp(chatCompletion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
		//println(err.Error())
		return
	}

	//　これをJson形式で返してフロントでチェックした後にCreateTaskでもたたけばいいんじゃね？
	taskData := &model.Task{
		Start:   extracted.Start,
		End:     extracted.End,
		Context: extracted.Context,
	}

	// 今は自動的にタスク登録されちゃうよーん
	//err = h.TaskHandler.CreateTaskFromAIResponse(taskData)

	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	//}

	c.JSON(http.StatusOK, taskData)
}
