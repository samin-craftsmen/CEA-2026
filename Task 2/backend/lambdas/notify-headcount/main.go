package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/config"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/database"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/service"
)

func handler(ctx context.Context) error {
	ist := time.FixedZone("IST", 5*60*60+30*60)
	tomorrow := time.Now().In(ist).AddDate(0, 0, 1).Format("2006-01-02")

	summary, err := service.GetHeadcountSummaryForDate(tomorrow)
	if err != nil {
		return fmt.Errorf("failed to get headcount: %w", err)
	}

	content := buildMessage(summary)
	return postToDiscord(content)
}

func buildMessage(s *service.HeadcountSummaryResponse) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 **Meal Headcount for %s**\n", s.Date))
	sb.WriteString(fmt.Sprintf("🏢 Office: **%d** | 🏠 WFH: **%d**\n\n", s.WorkLocation.Office, s.WorkLocation.WFH))

	// Sort meal types for consistent ordering
	mealTypes := make([]string, 0, len(s.Summary))
	for mt := range s.Summary {
		mealTypes = append(mealTypes, mt)
	}
	sort.Strings(mealTypes)

	for _, mt := range mealTypes {
		entry := s.Summary[mt]
		sb.WriteString(fmt.Sprintf("**%s** — ✅ Yes: **%d** | ❌ No: **%d**\n",
			strings.ToUpper(mt[:1])+mt[1:], entry.Yes, entry.No))
	}

	return sb.String()
}

func postToDiscord(content string) error {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	if token == "" || channelID == "" {
		return fmt.Errorf("DISCORD_BOT_TOKEN and DISCORD_CHANNEL_ID env vars must be set")
	}

	body, _ := json.Marshal(map[string]string{"content": content})
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", channelID),
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bot "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("discord API returned status %d", resp.StatusCode)
	}
	return nil
}

func main() {
	cfg := config.LoadConfig()
	if err := database.InitDynamoDB(cfg); err != nil {
		panic(err)
	}
	lambda.Start(handler)
}
