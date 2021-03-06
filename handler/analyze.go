package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/IBM/go-sdk-core/core"
	"github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
)

var (
	nlu, nluErr = naturallanguageunderstandingv1.NewNaturalLanguageUnderstandingV1(&naturallanguageunderstandingv1.NaturalLanguageUnderstandingV1Options{
		Version: os.Getenv("NATURALLANGUAGEUNDERSTANDING_VERSION"),
		Authenticator: &core.IamAuthenticator{
			ApiKey: os.Getenv("NATURALLANGUAGEUNDERSTANDING_APIKEY"),
		},
		URL: os.Getenv("NATURALLANGUAGEUNDERSTANDING_URL"),
	})
)

// AnalyzeText and return all entities found by Natural Language Understanding service
func AnalyzeText(w http.ResponseWriter, h *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(h.Body)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Fail to retrieve the body. Try again later."})
		return
	}

	var parsedBody map[string]string
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Fail to parse the body. Try again later."})
		return
	}
	text := parsedBody["text"]

	customModel, language := loadEnvironmentVariables()

	nluResult, _, err := nlu.Analyze(&naturallanguageunderstandingv1.AnalyzeOptions{
		Text: &text,
		Features: &naturallanguageunderstandingv1.Features{
			Entities: &naturallanguageunderstandingv1.EntitiesOptions{
				Model: &customModel,
			},
			Relations: &naturallanguageunderstandingv1.RelationsOptions{
				Model: &customModel,
			},
		},
		Language: &language,
	})

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Fail to extract the entities and relations from your text. Try again later.", "description": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(nluResult)
}

func loadEnvironmentVariables() (customModel string, language string) {
	customModel = os.Getenv("KNOWLEDGESTUDIO_CUSTOMMODEL")
	language = os.Getenv("NATURALLANGUAGEUNDERSTANDING_LANGUAGE")
	return
}
