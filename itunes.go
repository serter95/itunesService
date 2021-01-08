package itunesservicego

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Result struct ...
type Result struct {
	Kind           string `json:"kind"`
	WrapperType    string `json:"wrapperType"`
	ArtistName     string `json:"artistName"`
	TrackName      string `json:"trackName"`
	CollectionName string `json:"collectionName"`
	PreviewURL     string `json:"previewUrl"`
	TrackViewURL   string `json:"trackViewUrl"`
}

// AppleResponse struct ...
type AppleResponse struct {
	ResultCount int `json:"resultCount"`
	Results     []Result
}

// StandardResponse struct ...
type StandardResponse struct {
	Category   string `json:"category"`
	Name       string `json:"name"`
	Author     string `json:"author"`
	PreviewURL string `json:"previewUrl"`
	Origin     string `json:"origin"`
}

// FormatText Return ...
func FormatText(text string) string {
	// se quitan los espacios del principio y el final
	// se quitan los dobles espacios y se deja un solo espacio
	// luego de eso se reemplaza el espacio por un '+'
	// todas las letras se colocan en minusculas
	filteredText := strings.ToLower(strings.Replace(regexp.MustCompile(`\s\s*`).ReplaceAllString(strings.TrimSpace(text), " "), " ", "+", -1))
	// fmt.Println(filteredText)
	return filteredText
}

// FindResults Return ...
func FindResults(textToFind string) ([]StandardResponse, error) {
	// creo el array que tendra la respuesta unificada
	responseArray := make([]StandardResponse, 0)

	// url de la api
	apiURL := "https://itunes.apple.com"
	// recurso o ruta
	resource := "/search/"
	// variables de la peticion
	data := url.Values{}
	data.Set("term", FormatText(textToFind))
	data.Set("limit", "25")
	data.Set("media", "all")
	// se parsea la url
	u, _ := url.ParseRequestURI(apiURL)
	// se iguala el path al recurso solicitado
	u.Path = resource
	// se codifican las variables en la query
	u.RawQuery = data.Encode()
	// se asigna la url completa a la variable final
	urlStr := u.String() // "https://itunes.apple.com/search/?limit=25&media=all&term=foo"
	// se crea el cliente que hara la peticion (se usa un puntero)
	client := &http.Client{}
	// se crea la request
	r, _ := http.NewRequest(http.MethodGet, urlStr, nil) // URL-encoded payload\
	// se asigna el header
	r.Header.Add("Content-Type", "application/json")

	// se ejecuta el request
	resp, err := client.Do(r)
	// handle error, la peticion fallo
	if err != nil {
		// fmt.Println("PETICION FALLIDA")
		// fmt.Println(err)
		return responseArray, err
	}

	// despues del flujo cierra la respuesta
	defer resp.Body.Close()
	// leo todo el body
	body, _ := ioutil.ReadAll(resp.Body)
	// creo una variable de tipo AppleResponse
	var appleResponse AppleResponse
	// formateo el body de json a go y se asigna a appleResponse
	json.Unmarshal(body, &appleResponse)
	// asigno los valores de la query a las variables
	resultCount := appleResponse.ResultCount
	arrayResult := appleResponse.Results

	// recorro los resultados
	for i := 0; i < resultCount; i++ {
		// creo un iterador que contiene el valor de cada posicion
		iterator := arrayResult[i]
		// variables para validar
		var category string
		var name string
		var author string
		var previewURL string
		var origin string
		// validacion para category
		if iterator.Kind != "" {
			category = iterator.Kind
		} else if iterator.WrapperType != "" {
			category = iterator.WrapperType
		}
		// validacion para name
		if iterator.TrackName != "" {
			name = iterator.TrackName
		} else if iterator.CollectionName != "" {
			name = iterator.CollectionName
		}
		// se asigna el autor
		author = iterator.ArtistName
		// validacion para previewURL
		if iterator.PreviewURL != "" {
			previewURL = iterator.PreviewURL
		} else if iterator.TrackViewURL != "" {
			previewURL = iterator.TrackViewURL
		}
		// se asigna el origen
		origin = "apple"
		responseArray = append(responseArray, StandardResponse{Category: category, Name: name, Author: author, PreviewURL: previewURL, Origin: origin})
	}

	// responseArrayJSON, _ := json.MarshalIndent(responseArray, "", "  ")
	// fmt.Println(string(responseArrayJSON))
	return responseArray, nil
}
