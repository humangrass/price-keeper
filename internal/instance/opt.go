package instance

import (
	"github.com/humangrass/gommon/database"
	"github.com/humangrass/price-keeper/pgk/x/xhttp"
)

type Opt struct {
	Database     *database.Opt
	IsProduction bool
	XHttpOpt     xhttp.Opt
}
