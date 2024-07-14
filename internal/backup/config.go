package backup

type Config struct {
	FileStoragePath string `json:"store_file"`
	StoreInterval   int    `json:"store_interval"`
	Restore         bool   `json:"restore"`
}
