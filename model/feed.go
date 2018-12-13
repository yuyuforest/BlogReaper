package model

import (
	"github.com/boltdb/bolt"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/core/errors"
)

type FeedModel struct {
	*Model
}

type Feed struct {
	ID         bson.ObjectId   `bson:"id"`         // 订阅源的ID
	URL        string          `bson:"url"`        // 订阅源的URL
	Title      string          `bson:"title"`      // 订阅源的标题
	Categories []bson.ObjectId `bson:"categories"` // 订阅源的分类
	Articles   []Article       `bson:"articles"`   // 订阅源包括的文章
}

type Article struct {
	URL   string `bson:"url"`
	Read  bool   `bson:"read"`
	Later bool   `bson:"later"`
}

func (m *FeedModel) AddFeed(userID, categoryID string, publicFeed PublicFeed) (feed Feed, err error) {
	return feed, m.Update(func(b *bolt.Bucket) error {
		ub, err := b.CreateBucketIfNotExists([]byte(userID))
		if err != nil {
			return err
		}
		uub, err := ub.CreateBucketIfNotExists([]byte("key_url_value_userId"))
		if err != nil {
			return err
		}
		if uub.Get([]byte(publicFeed.URL)) == nil {
			return errors.New("repeat_url")
		}
		var articles []Article
		for _, a := range publicFeed.Articles {
			articles = append(articles, Article{
				URL:   a,
				Read:  false,
				Later: false,
			})
		}
		feed = Feed{
			ID:         bson.NewObjectId(),
			URL:        publicFeed.URL,
			Title:      publicFeed.Title,
			Categories: []bson.ObjectId{bson.ObjectIdHex(categoryID)},
			Articles:   articles,
		}
		bytes, err := bson.Marshal(&feed)
		if err != nil {
			return err
		}
		err = ub.Put([]byte(feed.ID), bytes)
		if err != nil {
			return err
		}
		return uub.Put([]byte(feed.URL), []byte(feed.URL))
	})
}

func (m *FeedModel) EditFeed(userID, categoryID, url, title string) (feed Feed, err error) {
	panic("not implement")
}

func (m *FeedModel) EditArticle(userID, categoryID, url, articleURL string, read, later bool) (err error) {
	panic("not implement")
}
