package consts

import (
	"time"

	"github.com/uncomonq/FinalCalc/pckg/types"
)

// Статусы задач и выражений.
const (
	WaitStatus        types.Status = "Wait"        // Задача ждёт выполнения
	OKStatus          types.Status = "OK"          // Задача готова
	CalculationStatus types.Status = "Calculation" // Задача в процессе выполнения
	ErrorStatus       types.Status = "Error"       // Ошибка
)

// Относительные директории файла data.db
const (
	AppStoragePath  = "./db/data.db"     // workdir=/cmd
	TestStoragePath = "../../db/data.db" // workdir=/internal/application
)

// Время запроса агента
var AgentReqestDelay = time.Millisecond * 1
