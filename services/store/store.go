package store

type Store struct {
	NotificationStore        *NotificationStore
	NotificationDataStore    *NotificationDataStore
	NotificationImgUrlsStore *NotificationImgUrlsStore
	NotificationGifUrlsStore *NotificationGifUrlsStore
	NotificationPackStore    *NotificationPackStore
}

func NewStore() *Store {
	return &Store{
		NotificationStore:        NewNotificationStore(),
		NotificationDataStore:    NewDataStore(),
		NotificationImgUrlsStore: NewImgUrlsStore(),
		NotificationGifUrlsStore: NewGifUrlsStore(),
		NotificationPackStore:    NewPacksStore(),
	}
}
