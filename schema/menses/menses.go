package menses

import (
	"time"

	"github.com/jamesmukumu/backup/schema/admin"
)

//schema for menses
type Menses struct{
Normalcycleday int `json:"normalcycledays"`
Lastcycledate string `json:"lastcycledate" bson:"lastcycledate"`
Nextexpectedperioddate time.Time `json:"nextexpectedperiod" bson:"nextexpectedperiod"`
Email admin.User `json:"Email" bson:"Email"`
Safedays time.Time `json:"safedays" bson:"safedays"`
Lastcycledatetime time.Time `json:"lastcycletime" bson:"lastcycletime"`
} 