package backup

type Config struct {
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}
