package modules

import "time"

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
	Employee UserType = "Работник"
	Customer UserType = "Клиент"
)

type Kartodrom struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(100);not null"`
	Location string `gorm:"type:varchar(100);not null"`
	Phone    string `gorm:"type:varchar(11);not null;unique"`
	Email    string `gorm:"type:varchar(50);not null;unique"`
}

type Kart struct {
	ID          uint       `gorm:"primaryKey"`
	KartodromID uint       `gorm:"not null"`
	Kartodrom   Kartodrom  `gorm:"foreignKey:KartodromID;references:ID"`
	KartModelID uint       `gorm:"not null"`
	KartModel   KartModel  `gorm:"foreignKey:KartModelID;references:ID"`
	KartStatus  KartStatus `gorm:"type:varchar(20);not null"`
}

type KartModel struct {
	ID             uint            `gorm:"primaryKey"`
	Name           string          `gorm:"type:varchar(50);not null"`
	Category       KartModelStatus `gorm:"type:varchar(20);not null"`
	MaxSpeed       int             `gorm:"not null"`
	MaxWeight      int             `gorm:"not null"`
	MaxHeight      int             `gorm:"not null"`
	PricePerMinute float64         `gorm:"type:decimal;not null"`
	Desc           string          `gorm:"type:varchar(255);not null"`
}

type UserInfo struct {
	ID        uint      `gorm:"primaryKey"`
	FName     string    `gorm:"type:varchar(50);not null"`
	SName     string    `gorm:"type:varchar(50);not null"`
	Phone     string    `gorm:"type:varchar(11);not null;unique"`
	Email     string    `gorm:"type:varchar(50);unique"`
	Weight    float64   `gorm:"type:decimal"`
	Height    float64   `gorm:"type:decimal"`
	BirthDay  time.Time `gorm:"type:date"`
	Balance   float64   `gorm:"type:decimal"`
	CreatedAt time.Time
}

type EmployeeInfo struct {
	ID       uint   `gorm:"primaryKey"`
	Position string `gorm:"type:varchar(50);not null"`
}

type AuthCredential struct {
	ID               uint   `gorm:"primaryKey"`
	Login            string `gorm:"type:varchar(50);not null;unique"`
	PassWordHashbyte []byte `gorm:"not null"`
}

type User struct {
	ID             uint           `gorm:"primaryKey"`
	UserType       UserType       `gorm:"type:varchar(20);not null"`
	AuthID         uint           `gorm:"not null"`
	Auth           AuthCredential `gorm:"foreignKey:AuthID;references:ID"`
	UserInfoID     uint           `gorm:"not null"`
	UserInfo       UserInfo       `gorm:"foreignKey:UserInfoID;references:ID"`
	EmployeeInfoID *uint
	EmployeeInfo   EmployeeInfo `gorm:"foreignKey:EmployeeInfoID;references:ID"`
}

type RaceRider struct {
	ID      uint `gorm:"primaryKey"`
	RiderID uint `gorm:"not null"`
	Rider   User `gorm:"foreignKey:RiderID;references:ID"`
	RaceID  uint `gorm:"not null"`
	Race    Race `gorm:"foreignKey:RaceID;references:ID"`
}

type Booking struct {
	ID          uint          `gorm:"primaryKey"`
	TrackID     uint          `gorm:"not null"`
	Track       Track         `gorm:"foreignKey:TrackID;references:ID"`
	CustomerID  uint          `gorm:"not null"`
	Customer    User          `gorm:"foreignKey:CustomerID;references:ID"`
	Date        time.Time     `gorm:"type:date;not null"`
	BookingType BookingType   `gorm:"type:varchar(50);not null"`
	Status      BookingStatus `gorm:"type:varchar(25);not null"`
	Duration    uint          `gorm:"not null"`
	TimeStart   time.Time     `gorm:"not null"`
	TimeEnd     time.Time     `gorm:"not null"`
	TotalPrice  float64       `gorm:"type:decimal;not null"`
}

type KartBooking struct {
	ID                uint          `gorm:"primaryKey"`
	BookingID         uint          `gorm:"not null"`
	Booking           Booking       `gorm:"foreignKey:BookingID;references:ID"`
	CustomerID        uint          `gorm:"not null"`
	Customer          User          `gorm:"foreignKey:CustomerID;references:ID"`
	KartID            uint          `gorm:"not null"`
	Kart              Kart          `gorm:"foreignKey:KartID;references:ID"`
	KartBookingStatus BookingStatus `gorm:"type:varchar(25);not null"`
}

type Race struct {
	ID         uint       `gorm:"primaryKey"`
	TrackID    uint       `gorm:"not null"`
	Track      Track      `gorm:"foreignKey:TrackID;references:ID"`
	Date       time.Time  `gorm:"type:date;not null"`
	RaceType   RaceType   `gorm:"type:varchar(50);not null"`
	TimeStart  time.Time  `gorm:"not null"`
	TimeEnd    time.Time  `gorm:"not null"`
	Duration   uint       `gorm:"not null"`
	Laps       uint       `gorm:"not null"`
	TotalPrice float64    `gorm:"type:decimal;not null"`
	RaceStatus RaceStatus `gorm:"type:varchar(25);not null"`
}

type Track struct {
	ID             uint      `gorm:"primaryKey"`
	KartodromID    uint      `gorm:"not null"`
	Kartodrom      Kartodrom `gorm:"foreignKey:KartodromID;references:ID"`
	Name           string    `gorm:"type:varchar(50);not null"`
	Length         uint      `gorm:"not null"`
	DifLevel       DifLevel  `gorm:"type:varchar(20);not null"`
	PricePerMinute float64   `gorm:"type:decimal;not null"`
	MaxKart        uint      `gorm:"not null"`
}

type KartStat struct {
	ID            uint      `gorm:"primaryKey"`
	KartID        uint      `gorm:"not null"`
	Kart          Kart      `gorm:"foreignKey:KartID;references:ID"`
	Date          time.Time `gorm:"type:date;not null"`
	TotalCustomer int       `gorm:"not null"`
	TotalTime     int       `gorm:"not null"`
	TotalRevenue  float64   `gorm:"type:decimal;not null"`
}

type TrackStat struct {
	ID      uint      `gorm:"primaryKey"`
	TrackID uint      `gorm:"not null"`
	Track   Track     `gorm:"foreignKey:TrackID;references:ID"`
	Date    time.Time `gorm:"type:date;not null"`

	FreeRideCount int `gorm:"not null"`
	FreeRideTime  int `gorm:"not null"`

	TimeAttackCount int `gorm:"not null"`
	TimeAttackTime  int `gorm:"not null"`

	DuoCount int `gorm:"not null"`
	DuoTime  int `gorm:"not null"`

	SprintCount int `gorm:"not null"`
	SprintTime  int `gorm:"not null"`

	EnduranceCount int `gorm:"not null"`
	EnduranceTime  int `gorm:"not null"`

	TotalCustomer int     `gorm:"not null"`
	TotalTime     int     `gorm:"not null"`
	TotalRevenue  float64 `gorm:"type:decimal;not null"`
}

type CommonStat struct {
	ID            uint      `gorm:"primaryKey"`
	KartodromID   uint      `gorm:"not null"`
	Kartodrom     Kartodrom `gorm:"foreignKey:KartodromID;references:ID"`
	Date          time.Time `gorm:"type:date;not null"`
	TotalCustomer int       `gorm:"not null"`
	TotalTime     int       `gorm:"not null"`
	TotalRevenue  float64   `gorm:"type:decimal;not null"`
}
