package store

type Store struct {
	UserStore                *UserStore
	NotificationStore        *NotificationStore
	NotificationDataStore    *NotificationDataStore
	NotificationImgUrlsStore *NotificationImgUrlsStore
	NotificationPackStore    *NotificationPackStore
}

func NewStore() *Store {
	return &Store{
		UserStore:                NewUserStore(),
		NotificationStore:        NewNotificationStore(),
		NotificationDataStore:    NewDataStore(),
		NotificationImgUrlsStore: NewImgUrlsStore(),
		NotificationPackStore:    NewPacksStore(),
	}
}
