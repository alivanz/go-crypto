package bitcoin

import (
	"fmt"
)

type Satoshi uint64

const Bitcoin = Satoshi(100000000)

func (sat Satoshi) String() string {
	return fmt.Sprintf("%d:%08d", sat/Bitcoin, sat%Bitcoin)
}
