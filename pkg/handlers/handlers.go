package handlers

import (
	"canonicalTestAPI/pkg/models"
	"canonicalTestAPI/pkg/service"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// handles getting a specific book by its ID.
func GetBook(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	bookID := c.Param("id")
	book, err := sv.GetBook(bookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book doesn't exist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "sucess", "data": book})
}

// handles the creation of a new book.
func CreateBook(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if book.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book must have a title"})
		return
	}

	id, err := sv.InsertBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created", "data": id})
}

// handles the creation of a new collection.
func CreateCollection(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	var collection models.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if collection.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection must have a name"})
		return
	}

	if err := sv.InsertCollection(collection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

// handles adding a book to a collection.
func CollectBook(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	collectionName := c.Param("name")
	bookID := c.Param("id")

	if _, err := sv.GetBook(bookID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book doesn't exist"})
		return
	}

	if _, err := sv.GetCollection(collectionName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Collection doesn't exist"})
		return
	}

	if err := sv.Collect(collectionName, bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// handles removing a book from a collection.
func DiscardBook(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	collectionName := c.Param("name")
	bookID := c.Param("id")

	if _, err := sv.GetBook(bookID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book doesn't exist"})
		return
	}

	if _, err := sv.GetCollection(collectionName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Collection doesn't exist"})
		return
	}

	if !sv.BelongsTo(collectionName, bookID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book isn't a part of this collection"})
		return
	}

	if err := sv.Discard(collectionName, bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// handles the deletion of a book.
func DeleteBook(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	bookID := c.Param("id")
	if _, err := sv.GetBook(bookID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book doesn't exist"})
		return
	}

	if err := sv.DeleteBook(bookID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// handles the deletion of a collection.
func DeleteCollection(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	collectionName := c.Param("name")
	if _, err := sv.GetCollection(collectionName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection doesn't exist"})
		return
	}

	if err := sv.DeleteCollection(collectionName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// handles getting all collections.
func GetAllCollections(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	collections, err := sv.GetAllCollections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": collections})
}

func parseQuery(values url.Values) (models.QueryFilter, error) {
	var filter models.QueryFilter

	flattened := make(map[string]interface{})
	for k, v := range values {
		flattened[k] = v[0]
	}

	if flattened["max"] != nil {
		max, err := strconv.Atoi(flattened["max"].(string))
		if err != nil {
			return filter, err
		}
		flattened["max"] = max
	}

	config := &mapstructure.DecoderConfig{
		DecodeHook:  mapstructure.TextUnmarshallerHookFunc(),
		ErrorUnused: true,
		ErrorUnset:  false,
		Result:      &filter,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return models.QueryFilter{}, err
	}

	if err := decoder.Decode(flattened); err != nil {
		return models.QueryFilter{}, err
	}

	return filter, nil
}

// handles querying books with various filters.
func Query(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	queryParams := c.Request.URL.Query()

	filter, err := parseQuery(queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	books, err := sv.Query(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": books})
}

// handles updating book information.
func UpdateBook(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	bookID := c.Param("id")

	if _, err := sv.GetBook(bookID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book doesn't exist"})
		return
	}

	var updates models.Book
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updates.ID != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot update book ID"})
		return
	}

	if updates == (models.Book{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided"})
		return
	}

	err := sv.UpdateBook(bookID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": updates})
}

// handles updating collection information.
func UpdateCollection(c *gin.Context) {
	sv_raw, _ := c.Get("service")
	sv, ok := sv_raw.(service.Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't connect to database"})
		return
	}

	collectionName := c.Param("name")

	if _, err := sv.GetCollection(collectionName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection doesn't exist"})
		return
	}

	var updates models.Collection
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updates.Id != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot update collection name"})
		return
	}

	if updates == (models.Collection{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided"})
		return
	}

	if err := sv.UpdateCollection(collectionName, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": updates})
}
