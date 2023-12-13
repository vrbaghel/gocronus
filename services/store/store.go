package store

type Store struct {
	NotificationStore        *NotificationStore
	NotificationDataStore    *NotificationDataStore
	NotificationImgUrlsStore *NotificationImgUrlsStore
	NotificationPackStore    *NotificationPackStore
}

func NewStore() *Store {
	return &Store{
		NotificationStore:        NewNotificationStore(),
		NotificationDataStore:    NewDataStore(),
		NotificationImgUrlsStore: NewImgUrlsStore(),
		NotificationPackStore:    NewPacksStore(),
	}
}
