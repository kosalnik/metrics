package backup

type Config struct {
	FileStoragePath string
	StoreInterval   int
	Restore         bool
}
