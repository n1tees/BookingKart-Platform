package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Unit string

const (
	IndividualScore  Unit = "Индивидуальные баллы"
	FinalTime        Unit = "Итоговое время"
	Placement        Unit = "Место"
	PerformanceScore Unit = "Лучший результат"
	ResultValue      Unit = "Резульатат"
)

type DifLevel string

const (
	Kids   DifLevel = "Детский"
	Light  DifLevel = "Легкий"
	Medium DifLevel = "Средний"
	Hard   DifLevel = "Сложный"
)

type RaceStatus string

const (
	RaceCreate   RaceStatus = "Создана"
	RaceStart    RaceStatus = "Начата"
	RaceFinish   RaceStatus = "Завершена"
	RaceCanceled RaceStatus = "Отменена"
)

type RaceType string

const (
	FreeRide      RaceType = "Свободный заезд"
	TimeAttack    RaceType = "Одиночный заезд на время"
	Duo           RaceType = "Парный заезд"
	SprintRace    RaceType = "Заезд на короткую дистанцию"
	EnduranceRace RaceType = "Заезд на длинную дистанцию"
)

type BookingStatus string

const (
	BookingActive  BookingStatus = "Активный"
	BookingReserve BookingStatus = "Забронирован"
	BookingClose   BookingStatus = "Завершен"
	BookingCancel  BookingStatus = "Отменен"
)

type BookingType string

const (
	CommonBooking   BookingType = "Обычное бронирование"
	OneTrackBooking BookingType = "Бронирование всего трека"
	AllTrackBooking BookingType = "Бронирование всего картодрома"
)

type KartModelStatus string

const (
	KidKart    KartModelStatus = "детский"
	CommonKart KartModelStatus = "обычный"
	SportKart  KartModelStatus = "спортивный"
	RaceKart   KartModelStatus = "гоночный"
	ElectoKart KartModelStatus = "электрический"
)

type KartStatus string

const (
	Available  KartStatus = "Доступен"
	InUse      KartStatus = "В использовании"
	Broken     KartStatus = "Сломан"
	InStopList KartStatus = "Недоступен"
)

type UserType string

const (
	Customer UserType = "customer"
	Admin    UserType = "admin"
)

type AuthCredential struct {
	gorm.Model
	Login        string `gorm:"type:varchar(50);not null;unique"`
	PasswordHash []byte `gorm:"not null"`
}

type User struct {
	gorm.Model
	AuthID    uint
	Auth      AuthCredential
	ProfileID uint
	Profile   Profile
	UserType  UserType `gorm:"type:varchar(25);not null"`
}

type Profile struct {
	gorm.Model
	FName    string    `gorm:"type:varchar(50);not null"`
	SName    string    `gorm:"type:varchar(50)"`
	Phone    string    `gorm:"type:varchar(11);not null;unique"`
	Email    string    `gorm:"type:varchar(50)"`
	Weight   float64   `gorm:"type:numeric(12, 2)"`
	Height   float64   `gorm:"type:numeric(12, 2)"`
	BirthDay time.Time `gorm:"type:date"`
	Balance  float64   `gorm:"type:numeric(12, 2)"`
}

type Payment struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	User   User
	Amount float64   `gorm:"type:numeric(12, 2);not null"`
	Date   time.Time `gorm:"not null"`
}

type RaceRider struct {
	ID             uint `gorm:"primaryKey"`
	RiderID        uint
	Rider          User
	RaceID         uint
	Race           Race
	ResultTypeID   uint
	ResultType     ResultType
	PersonalResult uint
}

type Booking struct {
	ID          uint `gorm:"primaryKey"`
	TrackID     uint
	Track       Track
	CustomerID  uint
	Customer    User
	RiderCount  uint          `gorm:"not null"`
	Date        time.Time     `gorm:"type:date;not null"`
	StartTime   LocalTime     `gorm:"type:time;not null"`
	EndTime     LocalTime     `gorm:"type:time;not null"`
	Duration    uint          `gorm:"not null"`
	TotalPrice  float64       `gorm:"type:numeric(12,2);not null"`
	BookingType BookingType   `gorm:"type:varchar(50);not null"`
	Status      BookingStatus `gorm:"type:varchar(25);not null"`
}

type KartBooking struct {
	ID         uint `gorm:"primaryKey"`
	BookingID  uint
	Booking    Booking
	CustomerID uint
	Customer   User
	KartID     uint
	Kart       Kart
	Status     BookingStatus `gorm:"type:varchar(25);not null"`
}

type ResultType struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(50);not null"`
	Unit Unit   `gorm:"type:varchar(25);not null"`
}

type Race struct {
	ID         uint `gorm:"primaryKey"`
	TrackID    uint
	Track      Track
	Date       time.Time `gorm:"type:date;not null"`
	TimeStart  LocalTime `gorm:"type:time;not null"`
	Duration   uint      `gorm:"not null"`
	Laps       uint
	TotalPrice float64    `gorm:"type:numeric(12, 2);not null"`
	RaceType   RaceType   `gorm:"type:varchar(50);not null"`
	Status     RaceStatus `gorm:"type:varchar(25);not null"`
}

type Track struct {
	ID          uint `gorm:"primaryKey"`
	KartodromID uint
	Kartodrom   Kartodrom
	Name        string   `gorm:"type:varchar(50);not null"`
	Length      uint     `gorm:"not null"`
	DifLevel    DifLevel `gorm:"type:varchar(20);not null"`
	PricePerMin float64  `gorm:"type:numeric(12, 2);not null"`
	MaxKarts    uint     `gorm:"not null"`
}

type Kart struct {
	ID          uint `gorm:"primaryKey"`
	KartodromID uint
	Kartodrom   Kartodrom
	KartModelID uint
	KartModel   KartModel
	Status      KartStatus `gorm:"type:varchar(20);not null"`
}

type TrackStat struct {
	ID      uint `gorm:"primaryKey"`
	TrackID uint
	Track   Track
	Date    time.Time `gorm:"type:date;not null"`

	FreeRideCount   uint    `gorm:"not null"`
	FreeRideTime    uint    `gorm:"not null"`
	TimeAttackCount uint    `gorm:"not null"`
	TimeAttackTime  uint    `gorm:"not null"`
	DuoCount        uint    `gorm:"not null"`
	DuoTime         uint    `gorm:"not null"`
	SprintCount     uint    `gorm:"not null"`
	SprintTime      uint    `gorm:"not null"`
	EnduranceCount  uint    `gorm:"not null"`
	EnduranceTime   uint    `gorm:"not null"`
	TotalCustomer   uint    `gorm:"not null"`
	TotalTime       uint    `gorm:"not null"`
	TotalRevenue    float64 `gorm:"type:numeric(12, 2);not null"`
}

type KartModel struct {
	ID          uint            `gorm:"primaryKey"`
	Name        string          `gorm:"type:varchar(50);not null"`
	Category    KartModelStatus `gorm:"type:varchar(50);not null"`
	MaxSpeed    uint            `gorm:"not null"`
	MaxWeight   uint            `gorm:"not null"`
	MaxHeight   uint            `gorm:"not null"`
	PricePerMin float64         `gorm:"type:numeric(12, 2);not null"`
	Desc        string          `gorm:"type:varchar(255);not null"`
}

type CommonStat struct {
	ID            uint `gorm:"primaryKey"`
	KartodromID   uint
	Kartodrom     Kartodrom
	Date          time.Time `gorm:"type:date;not null"`
	TotalCustomer uint      `gorm:"not null"`
	TotalTime     uint      `gorm:"not null"`
	TotalRevenue  float64   `gorm:"type:numeric(12, 2);not null"`
}

type KartodromSchedule struct {
	ID          uint `gorm:"primaryKey"`
	KartodromID uint
	Kartodrom   Kartodrom
	DayOfWeek   uint   `gorm:"not null"`
	OpenTime    string `gorm:"type:time;not null"`
	CloseTime   string `gorm:"type:time;not null"`
}

type Kartodrom struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"type:varchar(100);not null"`
	City      string  `gorm:"type:varchar(100);not null"`
	Location  string  `gorm:"type:varchar(100);not null"`
	Latitude  float64 `gorm:"not null"`
	Longitude float64 `gorm:"not null"`
	Phone     string  `gorm:"type:varchar(11);not null;unique"`
	Email     string  `gorm:"type:varchar(50);not null;unique"`

	Schedules []KartodromSchedule
}

type KartStat struct {
	ID            uint `gorm:"primaryKey"`
	KartID        uint
	Kart          Kart
	Date          time.Time `gorm:"type:date;not null"`
	TotalCustomer uint      `gorm:"not null"`
	TotalTime     uint      `gorm:"not null"`
	TotalRevenue  float64   `gorm:"type:numeric(12, 2);not null"`
}

type LocalTime struct {
	time.Time
}

func (lt *LocalTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		lt.Time = v
		return nil
	case string:
		t, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		lt.Time = t
		return nil
	default:
		return fmt.Errorf("unsupported type %T for LocalTime", v)
	}
}

func (lt LocalTime) Value() (driver.Value, error) {
	return lt.Format("15:04:05"), nil
}
