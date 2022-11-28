package db

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"

	_ "image/jpeg"

	"github.com/reiver/go-porterstemmer"
)

//go:embed stockimage.json
var stockImageFile embed.FS

const stockImageDBJsonFileName string = "stockimage.json"

const defaultNewsImageUrl string = "https://photos.google.com/share/AF1QipMoj8Yhuzm3b-b7XguhmS0hSUA8zd9lkQRjW8rzEtz88Wxs9ETgju0Mp9VtU2giBA/photo/AF1QipNuCGWjwFs3Km73Em-8X7BFR2ECQseonYtnakef?key=ejREcHZRcGgtN0l4V2NrNmM3eWU0LWtPaWtCQk93"

type StockImageInfo struct {
	Title  string   `json:"title"`
	Tags   []string `json:"tags"`
	Source string   `json:"source"`
	Image  string   `json:"image"`
}

type ImagesDB struct {
	Images []StockImageInfo `json:"images"`
}

var stockImageDB ImagesDB

var stockImageTagsList []string

// "imageseach" package consists of
// 1) a JSON file "stockimage.json" which contains all the images and their metadata for matching
// 2) a set of APIs for finding the images based on the set of tags.
// this init function will be called to load the metadata into the memory form an JSON file of image metadata
// when this package is imported from other package.
func init() {
	byteValue, _ := stockImageFile.ReadFile(stockImageDBJsonFileName)

	//fmt.Println("Loading Image Metadata")
	json.Unmarshal(byteValue, &stockImageDB)

	for _, imageData := range stockImageDB.Images {
		//tokens := NormalizeStringTokens(imageData.Title)
		tokens := []string{}
		for _, tag := range imageData.Tags {
			stem := porterstemmer.StemWithoutLowerCasing([]rune(strings.ToLower(tag)))
			tokens = append(tokens, string(stem))
		}
		sort.Strings(tokens)
		stockImageTagsList = append(stockImageTagsList, strings.Join(tokens, " "))
	}
	//fmt.Println("Normalize Image Metadata")
	//for index, tags := range stockImageTagsList {
	//	fmt.Println(index, tags)
	//}
}

// A list of stopwords that would be removed
var stopWords = []string{
	"about", "all", "alone", "also", "am", "and", "as", "at", "is", "at",
	"because", "before", "beside", "besides", "between", "but", "he", "him",
	"she", "her", "the", "by", "etc", "for", "i", "of", "on", "other", "others", "so",
	"than", "that", "though", "to", "too", "through", "until", "&"}

// This function tests if a string is in the array
func listContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// This function does following
// 1) tokenize the string
// 2) eliminate the stop words and redundant words
// 3) find stemming word of the original work in the string
// 4) lowercase the stemming word
// 5) sort all the normalized words
// it will return a string concatinating all the normalize words
// This function should be called before passing a string to Fuzzy match function
// to improve the matching efficiency and reduce the distance
func NormalizeTokens(words []string) []string {
	wordsAndCounter := make(map[string]int)
	var tokens []string

	for _, word := range words {
		// keyword is not a stopword and not empty
		wordInLowerCase := strings.ToLower(word)
		if !listContainsString(stopWords, wordInLowerCase) {
			_, hasKey := wordsAndCounter[wordInLowerCase]
			if !hasKey {
				wordsAndCounter[wordInLowerCase] = 1
				stem := porterstemmer.StemWithoutLowerCasing([]rune(wordInLowerCase))
				tokens = append(tokens, string(stem))
			}
		}
	}

	sort.Strings(tokens)
	return tokens
}

func NormalizeString(str string) string {
	words := strings.Fields(str)
	normalizedWords := NormalizeTokens(words)
	return strings.Join(normalizedWords, " ")
}

func LoadImageDb(db_file_name string) {
	jsonFile, err := os.Open(db_file_name)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		fmt.Println(err)
	}

	//var data ImagesDB
	json.Unmarshal(byteValue, &stockImageDB)
	fmt.Println("Image Metadata")
	fmt.Println(stockImageDB)

	fmt.Println(NormalizeString("i am manoj and dheeraj and dheeraj is going"))
}

func PrintRank(ranks fuzzy.Ranks) {
	for _, value := range ranks {
		fmt.Printf("source: %s target: %s Distance: %d OriginalIndex: %d \n",
			value.Source, value.Target, value.Distance, value.OriginalIndex)
	}
}

func searchStockImageDB(title string) []string {
	matches := fuzzy.RankFindNormalizedFold(title, stockImageTagsList)
	//PrintRank(matches)
	strlist := []string{}
	for _, value := range matches {
		strlist = append(strlist, stockImageDB.Images[value.OriginalIndex].Image)
	}
	return strlist
}

// We must return a default image if we can not find any matching image
func SearchImagesByTitle(title string) []string {
	images := searchStockImageDB(NormalizeString(title))

	if len(images) == 0 {
		return []string{defaultNewsImageUrl}
	} else {
		return images
	}
}

func SearchImagesByTags(tags []string) []string {
	normalizedTags := NormalizeTokens(tags)
	images := searchStockImageDB(strings.Join(normalizedTags, " "))

	if len(images) == 0 {
		return []string{defaultNewsImageUrl}
	} else {
		return images
	}
}
