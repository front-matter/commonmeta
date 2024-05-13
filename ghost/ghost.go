package ghost

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/jsonfeed"
	"github.com/golang-jwt/jwt"
)

// generate a short-lived JWT for the Ghost Admin API.
// adapted for Go from https://ghost.org/docs/admin-api/#token-authentication
func GenerateGhostToken(key string) (string, error) {
	// Split the key into ID and SECRET
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return "", errors.New("invalid key format")
	}
	id := parts[0]
	secret, _ := hex.DecodeString(parts[1])

	// Create the claims for the token, expires in 5 minutes
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		Audience:  "/admin/",
	}

	// Create and return the token (including decoding secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = id
	singnedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return singnedToken, nil
}

func UpdateGhostPost(id string, apiKey string, apiURL string) (string, error) {
	// get post doi and url from Rogue Scholar API
	// post url is needed to find post via Ghost API
	type Post struct {
		Slug         string `json:"slug"`
		ID           string `json:"id"`
		UUID         string `json:"uuid"`
		CanonicalURL string `json:"canonical_url"`
		UpdatedAt    string `json:"updated_at"`
	}

	type Ghost struct {
		Posts []Post `json:"posts"`
	}

	var content Ghost
	post, err := jsonfeed.Get(id)
	if err != nil {
		return "", err
	}
	doi := doiutils.NormalizeDOI(post.DOI)

	urlString := post.URL
	if doi == "" || urlString == "" {
		return "", errors.New("DOI or URL not found")
	}

	// get post_id and updated_at from ghost api
	token, err := GenerateGhostToken(apiKey)
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	u, _ := url.Parse(urlString)
	path := strings.Split(u.Path, "/")
	slug := path[len(path)-1]
	ghostURL := apiURL + "/ghost/api/admin/posts/slug/" + slug
	req, err := http.NewRequest(http.MethodGet, ghostURL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	addGhostHeaders(req, token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching post", err)
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", errors.New("Error fetching post: " + resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &content)
	if err != nil {
		fmt.Println("error:", err)
	}
	ghostPost := content.Posts[0]
	guid := ghostPost.ID
	updatedAt := ghostPost.UpdatedAt
	if guid == "" || updatedAt == "" {
		return "", errors.New("guid or updatedAt not found")
	}

	// update post canonical_url using the DOI. This requires sending
	// the updated_at timestamp to avoid conflicts, and must use guid
	// rather than url for put requests
	posts := append([]Post{}, Post{
		CanonicalURL: doi,
		UpdatedAt:    updatedAt,
	})
	p := Ghost{Posts: posts}
	payload, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	ghostPutURL := apiURL + "/ghost/api/admin/posts/" + guid
	req, err = http.NewRequest(http.MethodPut, ghostPutURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	addGhostHeaders(req, token)
	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("error updating post: " + resp.Status)
	}
	message := fmt.Sprintf("Canonical URL %s updated for GUID %s at %s", doi, guid, updatedAt)
	return message, nil
}

func addGhostHeaders(r *http.Request, token string) {
	r.Header.Set("Authorization", "Ghost "+token)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept-Version", "v5")
}
