package adidasv2

import (
	"time"
	"umbrella/internal/modules"
	globaltypes "umbrella/internal/utils/global_types"
)

const (
	USER_AGENT = "Screaming Frog SEO Spider/10.0"
)

const (
	ITEM_DATA_TYPE_BACKEND  = "backend"
	ITEM_DATA_TYPE_FRONTEND = "frontend"
)

type TaskStates struct {
	Login     globaltypes.TaskState
	CartClear globaltypes.TaskState
	AddToCart globaltypes.TaskState
	Address   globaltypes.TaskState
	Monitor   globaltypes.TaskState
	Order     globaltypes.TaskState
	Akamai    globaltypes.TaskState
}

type Resources struct {
	SessionTimeout time.Time
	BasketID       string
	CheckoutID     string
}

type Region struct {
	RegionCode string
}

type Config struct {
	DefaultConfig modules.DefaultConfig
	TaskStates    TaskStates
	Resources     Resources
	Region        *Region
	Cosmetics     Cosmetics
}

type Variant struct {
	Pid       string
	SizePid   string
	SizeValue string
}

type Cosmetics struct {
	Price string
	Pid   string
	Name  string
	Image string
	Size  string
}
