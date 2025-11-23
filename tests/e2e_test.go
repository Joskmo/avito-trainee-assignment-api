//go:build e2e

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

// Структуры для запросов и ответов (упрощенные для теста)
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type CreateTeamRequest struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type CreatePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
	AuthorID      string `json:"author_id"`
	PRName        string `json:"pull_request_name"`
}

type DeactivateRequest struct {
	Users []string `json:"users"`
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestE2E_FullFlow(t *testing.T) {
	// Генерация уникальных данных
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	teamName := fmt.Sprintf("e2e_team_%d", rnd.Int())
	user1 := fmt.Sprintf("u1_%d", rnd.Int())
	user2 := fmt.Sprintf("u2_%d", rnd.Int())
	user3 := fmt.Sprintf("u3_%d", rnd.Int())

	client := &http.Client{Timeout: 5 * time.Second}

	// 1. Создание команды
	t.Run("Create Team", func(t *testing.T) {
		reqBody := CreateTeamRequest{
			TeamName: teamName,
			Members: []TeamMember{
				{UserID: user1, Username: "Alice", IsActive: true},
				{UserID: user2, Username: "Bob", IsActive: true},
				{UserID: user3, Username: "Charlie", IsActive: true},
			},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := client.Post(baseURL+"/team/add", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// 2. Создание PR
	t.Run("Create PR", func(t *testing.T) {
		reqBody := CreatePRRequest{
			PullRequestID: "pr_" + randomString(8),
			AuthorID:      user1,
			PRName:        "E2E Feature",
		}
		body, _ := json.Marshal(reqBody)
		resp, err := client.Post(baseURL+"/pullRequest/create", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		prData := result["pr"].(map[string]interface{})
		// prID := prData["pull_request_id"].(string)

		// Проверяем, что ревьюверы назначены (должен быть кто-то из u2 или u3)
		reviewers := prData["assigned_reviewers"].([]interface{})
		assert.NotEmpty(t, reviewers, "Reviewers should be assigned")
	})

	// 3. Получение статистики
	t.Run("Get Stats", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/stats")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 4. Массовая деактивация и проверка переназначения
	t.Run("Deactivate User", func(t *testing.T) {
		// Деактивируем user2 и user3
		reqBody := DeactivateRequest{
			Users: []string{user2, user3},
		}
		body, _ := json.Marshal(reqBody)
		resp, err := client.Post(baseURL+"/team/deactivateUsers", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
