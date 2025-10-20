package analysis

import "github.com/BerylCAtieno/string-analyzer/utils"

type Properties map[string]any

type AnalysisPayload struct {
	Id         string     `json:"id"`
	Value      string     `json:"value"`
	Properties Properties `json:"properties"`
	CreatedAt  string     `json:"created_at"`
}

func AnalyzeString(s string) Properties {

	properties := make(Properties)

	length := utils.Length(s)
	is_palindrome := utils.IsPalindrome(s)
	unique_characters := utils.UniqueCharacters(s)
	word_count := utils.NumberOfWords(s)
	sha256_hash := utils.GenerateId(s)
	character_frequency_map := utils.CharFrequency(s)

	properties["length"] = length
	properties["is_palindrome"] = is_palindrome
	properties["unique_characters"] = unique_characters
	properties["word_count"] = word_count
	properties["sha256_hash"] = sha256_hash
	properties["character_frequency_map"] = character_frequency_map

	return properties

}
