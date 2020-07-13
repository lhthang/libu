package repository

import (
	"libu/app/form"
	"libu/app/model"
	"libu/my_db"
	"libu/utils/firebase"
	"net/http"
	"sync"
	"time"

	"github.com/araddon/dateparse"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var BookEntity IBook

type bookEntity struct {
	resource *my_db.Resource
	repo     *mongo.Collection
}

type BookChan struct {
	Ids []string
	Err error
}


type IBook interface {
	GetRecommendBooks(skip, limit int64,categoryIds []string) ([]form.BookResponse, int, error)
	GetSimilarBooks(skip, limit int64,id string) ([]form.BookResponse, int, error)
	GetAll(skip, limit int64) ([]form.BookResponse, int, error)
	GetHighRatedBooks(skips, limit int64) ([]form.BookResponse, int, error)
	GetPopularBooks(skips, limit int64) ([]form.BookResponse, int, error)
	GetNewBooks(skip, limit int64) ([]form.BookResponse, int, error)
	Search(keyword string) ([]form.BookResponse, int, error)
	Create(bookForm form.BookForm) (form.BookResponse, int, error)
	GetOneByID(id string) (form.BookResponse, int, error)
	Delete(id string) (form.BookResponse, int, error)
	Update(id string, bookForm form.UpdateBookForm) (form.BookResponse, int, error)
}

func NewBookEntity(resource *my_db.Resource) IBook {
	bookRepo := resource.DB.Collection("book")
	BookEntity = &bookEntity{
		resource: resource,
		repo:     bookRepo,
	}
	return BookEntity
}

func getCategoryOfBook(book *model.Book) []model.Category {
	var categories []model.Category
	var validCategoryIds []string
	var wg sync.WaitGroup
	categoryChan := make(chan *model.Category)
	for i, _ := range book.CategoryIds {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			chann := new(model.Category)
			category, _, err := CategoryEntity.GetOneByID(book.CategoryIds[i])
			if err != nil {
				return
			}
			validCategoryIds = append(validCategoryIds, book.CategoryIds[i])
			chann = category
			categoryChan <- chann
		}(i)
	}
	go func() {
		wg.Wait()
		close(categoryChan)
	}()

	for category := range categoryChan {
		categories = append(categories, *category)
	}
	book.CategoryIds = validCategoryIds
	return categories
}

func getAuthorsOfBook(book *model.Book) []model.Author {
	var authors []model.Author
	var validAuthorIds []string

	var wg sync.WaitGroup
	authorChan := make(chan *model.Author)
	for i, _ := range book.AuthorIds {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			chann := new(model.Author)
			author, _, err := AuthorEntity.GetOneByID(book.AuthorIds[i])
			if err != nil {
				return
			}
			validAuthorIds = append(validAuthorIds, book.AuthorIds[i])
			chann = author
			authorChan <- chann
		}(i)
	}

	go func() {
		wg.Wait()
		close(authorChan)
	}()

	for author := range authorChan {
		authors = append(authors, *author)
	}
	book.AuthorIds = validAuthorIds
	return authors
}

func getReviewsOfBook(book model.Book) *form.ReviewResponse {
	reviewResp, _, err := ReviewEntity.GetByBookId(book.Id.Hex())
	if err != nil {
		return nil
	}
	return reviewResp
}

func (entity bookEntity) GetAll(skip, limit int64) ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var booksResp []form.BookResponse
	cursor, err := entity.repo.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return booksResp, getHTTPCode(err), err
	}

	for cursor.Next(ctx) {
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Print(err)
		}
		reviewResp := getReviewsOfBook(book)
		booksResp = append(booksResp, form.BookResponse{
			Book: &book,
			//Reviews:    reviewResp.Reviews,
			Rating:     reviewResp.AvgRating,
			Categories: getCategoryOfBook(&book),
			Authors:    getAuthorsOfBook(&book),
		})
	}
	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) GetRecommendBooks(skip, limit int64, categoryIds []string) ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()

	var booksResp []form.BookResponse

	logrus.Println(categoryIds)
	pipeline := []bson.M{{"$addFields": bson.M{"categoryMatch": bson.M{"$size": bson.M{"$setIntersection": []interface{}{"$categoryIds", categoryIds}}}}},
		{"$sort": bson.M{"categoryMatch": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}



	cursor, err := entity.repo.Aggregate(ctx, pipeline)
	if err != nil {
		logrus.Println(err)
	}
	for cursor.Next(ctx) {
		//logrus.Printf("%v",cursor)
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Println(err)
			continue
		}
		reviewResp := getReviewsOfBook(book)
		booksResp = append(booksResp, form.BookResponse{
			Book: &book,
			//Reviews:    reviewResp.Reviews,
			Rating:     reviewResp.AvgRating,
			Categories: getCategoryOfBook(&book),
			Authors:    getAuthorsOfBook(&book),
		})
	}


	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) GetNewBooks(skip, limit int64) ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var booksResp []form.BookResponse
	cursor, err := entity.repo.Find(ctx, bson.M{}, &options.FindOptions{
		Sort:  bson.D{{"createAt", -1}},
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return booksResp, getHTTPCode(err), err
	}

	for cursor.Next(ctx) {
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Print(err)
		}
		reviewResp := getReviewsOfBook(book)
		booksResp = append(booksResp, form.BookResponse{
			Book: &book,
			//Reviews:    reviewResp.Reviews,
			Rating:     reviewResp.AvgRating,
			Categories: getCategoryOfBook(&book),
			Authors:    getAuthorsOfBook(&book),
		})
	}
	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) GetHighRatedBooks(skip, limit int64) ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()

	var booksResp []form.BookResponse

	//var bookIds []primitive.ObjectID
	pipeline := []bson.M{{"$project": bson.M{"_id": bson.M{"$toString": "$_id"},"title":1,"description":1,"publisher":1,"releaseAt":1,"createAt":1,"authorIds":1,"categoryIds":1,"link":1,"image":1}},
		{"$lookup": bson.M{"from": "review", "localField": "_id", "foreignField": "bookId", "as": "reviews"}},
		{"$project": bson.M{"_id": 1, "rating": bson.M{"$avg": "$reviews.rating"},"title":1,"description":1,"publisher":1,"releaseAt":1,"createAt":1,"authorIds":1,"categoryIds":1,"link":1,"image":1}},
		{"$sort": bson.M{"rating": -1}},
		{"$skip": skip},
		{"$limit": limit},
		{"$project": bson.M{"_id": bson.M{"$toObjectId": "$_id"},"title":1,"description":1,"publisher":1,"releaseAt":1,"createAt":1,"authorIds":1,"categoryIds":1,"link":1,"image":1}},
	}

	cursor, err := entity.repo.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	for cursor.Next(ctx) {
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Println(err)
			continue
		}
		reviewResp := getReviewsOfBook(book)
		booksResp = append(booksResp, form.BookResponse{
			Book: &book,
			//Reviews:    reviewResp.Reviews,
			Rating:     reviewResp.AvgRating,
			Categories: getCategoryOfBook(&book),
			Authors:    getAuthorsOfBook(&book),
		})
	}


	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) GetPopularBooks(skip, limit int64) ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()

	var booksResp []form.BookResponse

	pipeline := []bson.M{{"$project": bson.M{"_id": bson.M{"$toString": "$_id"},"title":1,"description":1,"publisher":1,"releaseAt":1,"createAt":1,"authorIds":1,"categoryIds":1,"link":1,"image":1}},
		{"$lookup": bson.M{"from": "review", "localField": "_id", "foreignField": "bookId", "as": "reviews"}},
		{"$project": bson.M{"_id": 1, "totalReviews": bson.M{"$cond": bson.M{"if": bson.M{"$isArray": "$reviews"}, "then": bson.M{"$size": "$reviews"}, "else": 0}},"title":1,"description":1,"publisher":1,"releaseAt":1,"createAt":1,"authorIds":1,"categoryIds":1,"link":1,"image":1}},
		{"$sort": bson.M{"totalReviews": -1}},
		{"$skip": skip},
		{"$limit": limit},
		{"$project": bson.M{"_id": bson.M{"$toObjectId": "$_id"},"title":1,"description":1,"publisher":1,"releaseAt":1,"createAt":1,"authorIds":1,"categoryIds":1,"link":1,"image":1}},
	}

	cursor, err := entity.repo.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	for cursor.Next(ctx) {
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Println(err)
			continue
		}
		reviewResp := getReviewsOfBook(book)
		booksResp = append(booksResp, form.BookResponse{
			Book: &book,
			//Reviews:    reviewResp.Reviews,
			Rating:     reviewResp.AvgRating,
			Categories: getCategoryOfBook(&book),
			Authors:    getAuthorsOfBook(&book),
		})
	}

	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) GetSimilarBooks(skip, limit int64,id string) ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()

	var booksResp []form.BookResponse

	book, _, err := entity.GetOneByID(id)
	if err != nil {
		return []form.BookResponse{}, http.StatusBadRequest, err
	}

	pipeline := []bson.M{{"$addFields": bson.M{"categoryMatch": bson.M{"$size": bson.M{"$setIntersection": []interface{}{"$categoryIds",book.CategoryIds}}}}},
		{"$sort": bson.M{"categoryMatch": -1}},
		{"$addFields": bson.M{"authorMatch": bson.M{"$size": bson.M{"$setIntersection": []interface{}{"$authorIds",book.AuthorIds}}}}},
		{"$sort": bson.M{"authorMatch": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}



	cursor, err := entity.repo.Aggregate(ctx, pipeline)
	if err != nil {
		logrus.Println(err)
	}
	for cursor.Next(ctx) {
		//logrus.Printf("%v",cursor)
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Println(err)
			continue
		}
		reviewResp := getReviewsOfBook(book)
		booksResp = append(booksResp, form.BookResponse{
			Book: &book,
			//Reviews:    reviewResp.Reviews,
			Rating:     reviewResp.AvgRating,
			Categories: getCategoryOfBook(&book),
			Authors:    getAuthorsOfBook(&book),
		})
	}

	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) Create(bookForm form.BookForm) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()

	releaseAt, err := dateparse.ParseAny(bookForm.ReleaseAt)
	if err != nil {
		logrus.Println(err)
	}
	var categoryIds []string
	for _, id := range bookForm.CategoryIds {
		category, _, _ := CategoryEntity.GetOneByID(id)
		if category != nil {
			categoryIds = append(categoryIds, id)
		}
	}
	var authorIds []string
	for _, id := range bookForm.AuthorIds {
		author, _, _ := AuthorEntity.GetOneByID(id)
		if author != nil {
			authorIds = append(authorIds, id)
		}
	}
	book := model.Book{
		Id:          primitive.NewObjectID(),
		ReleaseAt:   releaseAt,
		Title:       bookForm.Title,
		CreateAt:    time.Now(),
		AuthorIds:   authorIds,
		Publisher:   bookForm.Publisher,
		CategoryIds: categoryIds,
		Image:       bookForm.Image,
		Description: bookForm.Description,
		Link:        bookForm.Link,
	}
	channel := make(chan string)
	if bookForm.File != nil {
		go func() {

			url, _, err := firebase.UploadFile(*bookForm.File)
			if err != nil {
				logrus.Println(err)
			}
			channel <- url
		}()
		book.Link = <-channel
	}
	_, err = entity.repo.InsertOne(ctx, book)
	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}
	bookResp := form.BookResponse{
		Book:       &book,
		Reviews:    nil,
		Categories: getCategoryOfBook(&book),
		Authors:    getAuthorsOfBook(&book),
	}
	return bookResp, http.StatusOK, nil
}

func (entity bookEntity) GetOneByID(id string) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var book model.Book
	objID, _ := primitive.ObjectIDFromHex(id)

	err := entity.repo.FindOne(ctx, bson.M{"_id": objID}).Decode(&book)

	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}

	reviewResp := getReviewsOfBook(book)

	bookResp := form.BookResponse{
		Book:       &book,
		Reviews:    reviewResp.ReviewResp,
		Rating:     reviewResp.AvgRating,
		Categories: getCategoryOfBook(&book),
		Authors:    getAuthorsOfBook(&book),
	}
	return bookResp, http.StatusOK, nil
}

func (entity bookEntity) Search(keyword string) ([]form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	var booksResp []form.BookResponse

	title := map[string]interface{}{}
	description := map[string]interface{}{}
	// author := map[string]interface{}{}
	publisher := map[string]interface{}{}
	title["title"] = bson.M{"$regex": keyword, "$options": "i"}
	description["description"] = bson.M{"$regex": keyword, "$options": "i"}
	// author["authors"] = bson.M{"$regex": keyword, "$options": "i"}
	publisher["publisher"] = bson.M{"$regex": keyword, "$options": "i"}
	query := map[string]interface{}{}
	// query["$or"] = []interface{}{title, description, author, publisher}
	query["$or"] = []interface{}{title, description, publisher}

	cursor, err := entity.repo.Find(ctx, query)
	if err != nil {
		return booksResp, getHTTPCode(err), err
	}

	for cursor.Next(ctx) {
		var book model.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.Print(err)
		}

		reviewResp := getReviewsOfBook(book)
		booksResp = append(booksResp, form.BookResponse{
			Book: &book,
			//Reviews:    reviewResp.Reviews,
			Rating:     reviewResp.AvgRating,
			Categories: getCategoryOfBook(&book),
			Authors:    getAuthorsOfBook(&book),
		})
	}
	return booksResp, http.StatusOK, nil
}

func (entity bookEntity) Delete(id string) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	bookResp, _, err := entity.GetOneByID(id)

	if err != nil || bookResp.Book == nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}

	_, err = entity.repo.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}
	return bookResp, http.StatusOK, nil
}

func (entity bookEntity) Update(id string, bookForm form.UpdateBookForm) (form.BookResponse, int, error) {
	ctx, cancel := initContext()

	defer cancel()
	objID, _ := primitive.ObjectIDFromHex(id)
	bookResp, _, err := entity.GetOneByID(id)
	if err != nil || bookResp.Book == nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}

	var categoryIds []string
	for _, id := range bookForm.CategoryIds {
		category, _, _ := CategoryEntity.GetOneByID(id)
		if category != nil {
			categoryIds = append(categoryIds, id)
		}
	}
	bookForm.CategoryIds = categoryIds
	var authorIds []string
	for _, id := range bookForm.AuthorIds {
		author, _, _ := AuthorEntity.GetOneByID(id)
		if author != nil {
			authorIds = append(authorIds, id)
		}
	}
	bookForm.AuthorIds = authorIds
	err = copier.Copy(bookResp.Book, bookForm)
	if err != nil {
		return form.BookResponse{}, getHTTPCode(err), err
	}

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	var updatedBook model.Book

	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": bookResp.Book}, opts).Decode(&updatedBook)
	newBookResp := form.BookResponse{
		Book:       &updatedBook,
		Reviews:    bookResp.Reviews,
		Categories: getCategoryOfBook(&updatedBook),
		Authors:    getAuthorsOfBook(&updatedBook),
		Rating:     bookResp.Rating,
	}
	return newBookResp, http.StatusOK, nil
}
