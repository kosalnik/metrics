package backup

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kosalnik/metrics/internal/models"
)

type MyStorage struct {
	s []models.Metrics
	u time.Time
}

func (m *MyStorage) UpsertAll(_ context.Context, list []models.Metrics) error {
	m.s = make([]models.Metrics, len(list))
	for k := range list {
		v := list[k]
		m.s[k] = v
	}
	m.u = time.Now()
	return nil
}

func (m *MyStorage) UpdatedAt() time.Time {
	return m.u
}

func (m *MyStorage) GetAll(_ context.Context) ([]models.Metrics, error) {
	return m.s, nil
}

var _ Storage = &MyStorage{}

func ExampleDump_Store() {
	tmp, err := os.CreateTemp(os.TempDir(), "example")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := os.Remove(tmp.Name()); err != nil {
			fmt.Println(err.Error())
		}
	}()
	s := &MyStorage{s: []models.Metrics{
		{ID: "pi", MType: models.MGauge, Value: 3.14},
		{ID: "cnt", MType: models.MCounter, Value: 1},
	}}
	b := NewDump(s, tmp.Name())
	if err := b.Store(context.Background()); err != nil {
		panic(err)
	}
	got, err := os.ReadFile(tmp.Name())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", got)
	// output:
	// {"Data":[{"id":"pi","type":"gauge","value":3.14},{"id":"cnt","type":"counter","delta":0}]}
}

func ExampleRecover_Recover() {
	tmp, err := os.CreateTemp(os.TempDir(), "example")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := os.Remove(tmp.Name()); err != nil {
			fmt.Println(err.Error())
		}
	}()
	_, err = fmt.Fprint(tmp, `{"Data":[{"id":"pi","type":"gauge","value":3.14},{"id":"cnt","type":"counter","delta":0}]}`)
	if err != nil {
		panic(err)
	}
	s := &MyStorage{}
	b := NewRecover(s, tmp.Name())
	if err := b.Recover(context.Background()); err != nil {
		panic(err)
	}
	m, err := s.GetAll(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", m)
	// output:
	// [{ID:pi MType:gauge Delta:0 Value:3.14} {ID:cnt MType:counter Delta:0 Value:0}]
}

func Example() {
	tmp, err := os.CreateTemp(os.TempDir(), "example")
	if err != nil {
		panic(err)
	}
	fname := tmp.Name()
	defer func() {
		if err := os.Remove(tmp.Name()); err != nil {
			fmt.Println(err.Error())
		}
	}()
	_, err = fmt.Fprint(tmp, `{"Data":[{"id":"e","type":"gauge","value":2.71828}]}`)
	if err != nil {
		panic(err)
	}
	if err := tmp.Close(); err != nil {
		panic(err)
	}

	s := &MyStorage{}
	b, err := NewBackupManager(s, Config{
		StoreInterval:   1,
		FileStoragePath: fname,
		Restore:         true,
	})
	if err != nil {
		panic(err)
	}

	m, err := s.GetAll(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Init: %+v\n", m)

	if err := b.Recover(context.Background()); err != nil {
		panic(err)
	}
	m, err = s.GetAll(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Recovered: %+v\n", m)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()

	go b.BackupLoop(ctx)

	s.s = []models.Metrics{
		{ID: "pi", MType: models.MGauge, Value: 3.14},
		{ID: "cnt", MType: models.MCounter, Value: 1},
	}
	s.u = time.Now().Add(time.Second)
	<-time.After(time.Second * 3)

	got, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Dumped: %s", got)
	// output:
	// Init: []
	// Recovered: [{ID:e MType:gauge Delta:0 Value:2.71828}]
	// Dumped: {"Data":[{"id":"pi","type":"gauge","value":3.14},{"id":"cnt","type":"counter","delta":0}]}
}
