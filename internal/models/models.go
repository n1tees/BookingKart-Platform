package models

import (
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
	RaceCanceled RaceStatus = "Отменена"
	RaceReserved RaceStatus = "Забронирована"
	RaceClosed   RaceStatus = "Завершена"
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
	BookingCanceled BookingStatus = "Отменен"
	BookingReserved BookingStatus = "Забронирован"
	BookingClosed   BookingStatus = "Завершен"
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
	Availible  KartStatus = "Доступен"
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
	AuthID     uint           `gorm:"not null"`
	Auth       AuthCredential `gorm:"foreignKey:AuthID;references:ID"`
	UserInfoID uint           `gorm:"not null"`
	UserInfo   UserInfo       `gorm:"foreignKey:UserInfoID;references:ID"`
	UserType   UserType       `gorm:"type:varchar(25);not null"`
}

type UserInfo struct {
	gorm.Model
	FName    string    `gorm:"type:varchar(50);not null"`
	SName    string    `gorm:"type:varchar(50);not null"`
	Phone    string    `gorm:"type:varchar(11);not null;unique"`
	Email    string    `gorm:"type:varchar(50);unique"`
	Weight   float64   `gorm:"type:numeric(12, 2)"`
	Height   float64   `gorm:"type:numeric(12, 2)"`
	BirthDay time.Time `gorm:"type:date"`
	Balance  float64   `gorm:"type:numeric(12, 2)"`
}

type Payment struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uint      `gorm:"not null"`
	User   UserInfo  `gorm:"foreignKey:UserID;references:ID"`
	Amount float64   `gorm:"type:numeric(12, 2);not null"`
	Date   time.Time `gorm:"not null"`
}

type RaceRider struct {
	ID             uint       `gorm:"primaryKey"`
	RiderID        uint       `gorm:"not null"`
	Rider          User       `gorm:"foreignKey:RiderID;references:ID"`
	RaceID         uint       `gorm:"not null"`
	Race           Race       `gorm:"foreignKey:RaceID;references:ID"`
	ResultTypeID   uint       `gorm:"not null"`
	ResultType     ResultType `gorm:"foreignKey:ResultTypeID;references:ID"`
	PersonalResult uint       `gorm:"not null"`
}

type Booking struct {
	ID          uint          `gorm:"primaryKey"`
	TrackID     uint          `gorm:"not null"`
	Track       Track         `gorm:"foreignKey:TrackID;references:ID"`
	CustomerID  uint          `gorm:"not null"`
	Customer    User          `gorm:"foreignKey:CustomerID;references:ID"`
	Date        time.Time     `gorm:"type:date;not null"`
	TimeStart   time.Time     `gorm:"type:time;not null"`
	Duration    uint          `gorm:"not null"`
	TotalPrice  float64       `gorm:"type:numeric(12, 2);not null"`
	BookingType BookingType   `gorm:"type:varchar(50);not null"`
	Status      BookingStatus `gorm:"type:varchar(25);not null"`
}

type KartBooking struct {
	ID         uint          `gorm:"primaryKey"`
	BookingID  uint          `gorm:"not null"`
	Booking    Booking       `gorm:"foreignKey:BookingID;references:ID"`
	CustomerID uint          `gorm:"not null"`
	Customer   User          `gorm:"foreignKey:CustomerID;references:ID"`
	KartID     uint          `gorm:"not null"`
	Kart       Kart          `gorm:"foreignKey:KartID;references:ID"`
	Status     BookingStatus `gorm:"type:varchar(25);not null"`
}

type ResultType struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(50);not null"`
	Unit Unit   `gorm:"type:varchar(25);not null"`
}

type Race struct {
	ID         uint       `gorm:"primaryKey"`
	TrackID    uint       `gorm:"not null"`
	Track      Track      `gorm:"foreignKey:TrackID;references:ID"`
	Date       time.Time  `gorm:"type:date;not null"`
	TimeStart  time.Time  `gorm:"type:time;not null"`
	Duration   uint       `gorm:"not null"`
	Laps       uint       `gorm:"not null"`
	TotalPrice float64    `gorm:"type:numeric(12, 2);not null"`
	RaceType   RaceType   `gorm:"type:varchar(50);not null"`
	Status     RaceStatus `gorm:"type:varchar(25);not null"`
}

type Track struct {
	ID          uint      `gorm:"primaryKey"`
	KartodromID uint      `gorm:"not null"`
	Kartodrom   Kartodrom `gorm:"foreignKey:KartodromID;references:ID"`
	Name        string    `gorm:"type:varchar(50);not null"`
	Length      uint      `gorm:"not null"`
	DifLevel    DifLevel  `gorm:"type:varchar(20);not null"`
	PricePerMin float64   `gorm:"type:numeric(12, 2);not null"`
	MaxKarts    uint      `gorm:"not null"`
}

type Kart struct {
	ID          uint       `gorm:"primaryKey"`
	KartodromID uint       `gorm:"not null"`
	Kartodrom   Kartodrom  `gorm:"foreignKey:KartodromID;references:ID"`
	KartModelID uint       `gorm:"not null"`
	KartModel   KartModel  `gorm:"foreignKey:KartModelID;references:ID"`
	Status      KartStatus `gorm:"type:varchar(20);not null"`
}

type TrackStat struct {
	ID      uint      `gorm:"primaryKey"`
	TrackID uint      `gorm:"not null"`
	Track   Track     `gorm:"foreignKey:TrackID;references:ID"`
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
	ID            uint      `gorm:"primaryKey"`
	KartodromID   uint      `gorm:"not null"`
	Kartodrom     Kartodrom `gorm:"foreignKey:KartodromID;references:ID"`
	Date          time.Time `gorm:"type:date;not null"`
	TotalCustomer uint      `gorm:"not null"`
	TotalTime     uint      `gorm:"not null"`
	TotalRevenue  float64   `gorm:"type:numeric(12, 2);not null"`
}

type KartodromSchedule struct {
	ID          uint      `gorm:"primaryKey"`
	KartodromID uint      `gorm:"not null"`
	Kartodrom   Kartodrom `gorm:"foreignKey:KartodromID;references:ID"`
	DayOfWeek   uint      `gorm:"not null"`
	OpenTime    time.Time `gorm:"type:time;not null"`
	CloseTime   time.Time `gorm:"type:time;not null"`
}

type Kartodrom struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(100);not null"`
	Location string `gorm:"type:varchar(100);not null"`
	Phone    string `gorm:"type:varchar(11);not null;unique"`
	Email    string `gorm:"type:varchar(50);not null;unique"`
}

type KartStat struct {
	ID            uint      `gorm:"primaryKey"`
	KartID        uint      `gorm:"not null"`
	Kart          Kart      `gorm:"foreignKey:KartID;references:ID"`
	Date          time.Time `gorm:"type:date;not null"`
	TotalCustomer uint      `gorm:"not null"`
	TotalTime     uint      `gorm:"not null"`
	TotalRevenue  float64   `gorm:"type:numeric(12, 2);not null"`
}
